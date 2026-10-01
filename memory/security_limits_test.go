package memory_test

import (
	"context"
	"errors"
	"runtime"
	"strings"
	"testing"

	mpt "github.com/faustbrian/go-merkle-patricia-trie/v2"
	"github.com/faustbrian/go-merkle-patricia-trie/v2/memory"
)

func TestStoreRetentionHonorsLoweredAdmissionBound(t *testing.T) {
	ctx := context.Background()
	store := memory.New()
	limits := mpt.DefaultReachabilityLimits()
	limits.MaxRetentions = 2
	for range 2 {
		if _, err := store.RetainRoot(ctx, mpt.EmptyRoot(), limits); err != nil {
			t.Fatal(err)
		}
	}
	limits.MaxRetentions = 1
	if retention, err := store.RetainRoot(ctx, mpt.EmptyRoot(), limits); retention != nil || !errors.Is(err, mpt.ErrResourceLimit) {
		t.Fatal("lowered retention bound admitted another lease")
	}
}

func TestStoreAdmissionCancellationDoesNotPublish(t *testing.T) {
	trie := mustMemoryTrie(t, map[string]string{"key": "value"})
	_, commit := captureTrieCommit(t, trie, newCaptureStore())
	for cancellationAt := 1; cancellationAt <= 5; cancellationAt++ {
		store := memory.New()
		ctx := &stepContext{cancelAt: cancellationAt}
		if err := store.CommitTrie(ctx, commit); !errors.Is(err, mpt.ErrCanceled) || !errors.Is(err, context.Canceled) {
			t.Fatalf("canceled admission at seam %d returned %v", cancellationAt, err)
		}
		if store.Root() != mpt.EmptyRoot() || countStoreNodes(t, store) != 0 {
			t.Fatal("canceled admission changed published state")
		}
	}
}

func TestStorePruneAccountingCancellationPreservesSnapshot(t *testing.T) {
	ctx := context.Background()
	trie := mustMemoryTrie(t, map[string]string{"key": "value"})
	probe := memory.New()
	if _, err := trie.Commit(ctx, probe); err != nil {
		t.Fatal(err)
	}
	counting := &stepContext{}
	if _, err := probe.Prune(counting, mpt.DefaultReachabilityLimits()); err != nil {
		t.Fatal(err)
	}
	// Exercise each final accounting/publication seam with the same one-node
	// graph, rather than depending on elapsed scheduling or goroutine races.
	for cancellationAt := counting.calls - 2; cancellationAt <= counting.calls; cancellationAt++ {
		store := memory.New()
		committed, err := trie.Commit(ctx, store)
		if err != nil {
			t.Fatal(err)
		}
		root := store.Root()
		if _, err := store.Prune(&stepContext{cancelAt: cancellationAt}, mpt.DefaultReachabilityLimits()); !errors.Is(err, mpt.ErrCanceled) {
			t.Fatalf("canceled accounting returned %v", err)
		}
		if store.Root() != root || countStoreNodes(t, store) != 1 {
			t.Fatal("canceled accounting changed retained state")
		}
		if value, err := committed.Get(ctx, []byte("key")); err != nil || string(value) != "value" {
			t.Fatal("canceled accounting changed snapshot contents")
		}
	}
}

func TestStoreAdmissionRejectsPayloadBeforeAllocation(t *testing.T) {
	ctx := context.Background()
	trie := mustMemoryTrie(t, map[string]string{"key": strings.Repeat("v", 1<<20)})
	capture := newCaptureStore()
	_, commit := captureTrieCommit(t, trie, capture)
	limits := memory.DefaultLimits()
	limits.MaxStoredBytes = 32
	store, err := memory.NewWithLimits(limits)
	if err != nil {
		t.Fatal(err)
	}
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	err = store.CommitTrie(ctx, commit)
	runtime.ReadMemStats(&after)
	if !errors.Is(err, mpt.ErrResourceLimit) || store.Root() != mpt.EmptyRoot() || countStoreNodes(t, store) != 0 {
		t.Fatal("rejected oversized admission changed the store")
	}
	if after.TotalAlloc-before.TotalAlloc > 64<<10 {
		t.Fatal("rejected admission copied the encoded payload")
	}
}

func TestStoreAdmissionRejectsInvalidPoliciesAndOversizedNodeBatch(t *testing.T) {
	for _, limits := range []memory.Limits{{}, {MaxStoredNodes: -1, MaxStoredBytes: 1}, {MaxStoredNodes: 1, MaxStoredBytes: -1}} {
		if store, err := memory.NewWithLimits(limits); store != nil || !errors.Is(err, mpt.ErrResourceLimit) {
			t.Fatal("invalid memory admission policy was accepted")
		}
	}
	ctx := context.Background()
	trie := mustMemoryTrie(t, map[string]string{"alpha": strings.Repeat("a", 64), "beta": strings.Repeat("b", 64)})
	_, commit := captureTrieCommit(t, trie, newCaptureStore())
	if commit.NodeCount() <= 1 {
		t.Fatal("fixture does not exercise a multi-node batch")
	}
	limits := memory.DefaultLimits()
	limits.MaxStoredNodes = 1
	store, err := memory.NewWithLimits(limits)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.CommitTrie(ctx, commit); !errors.Is(err, mpt.ErrResourceLimit) || store.Root() != mpt.EmptyRoot() {
		t.Fatal("oversized node batch changed published state")
	}
}

func TestStoreAdmissionBoundsHistoricalNodesAndReclaimsAfterPrune(t *testing.T) {
	ctx := context.Background()
	first := mustMemoryTrie(t, map[string]string{"key": "first"})
	probe := memory.New()
	if _, err := first.Commit(ctx, probe); err != nil {
		t.Fatal(err)
	}
	var nodeBytes int
	if err := probe.IterateNodes(ctx, 10, func(_ mpt.Root, encoded []byte) error {
		nodeBytes += len(encoded)
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	for _, dimension := range []string{"nodes", "bytes"} {
		t.Run(dimension, func(t *testing.T) {
			limits := memory.DefaultLimits()
			if dimension == "nodes" {
				limits.MaxStoredNodes = 2
			} else {
				limits.MaxStoredBytes = 2 * nodeBytes
			}
			store, err := memory.NewWithLimits(limits)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := first.Commit(ctx, store); err != nil {
				t.Fatal(err)
			}
			root := store.Root()
			loaded, err := mpt.LoadRawTrie(root, store, mpt.DefaultLimits())
			if err != nil {
				t.Fatal(err)
			}
			updated, err := loaded.Update(ctx, []byte("key"), []byte("other"))
			if err != nil {
				t.Fatal(err)
			}
			if _, err := updated.Commit(ctx, store); err != nil {
				t.Fatal(err)
			}
			root = store.Root()
			loaded, err = mpt.LoadRawTrie(root, store, mpt.DefaultLimits())
			if err != nil {
				t.Fatal(err)
			}
			// A repeated root must not consume admission capacity again.
			encoded, err := store.GetNode(ctx, root)
			if err != nil {
				t.Fatal(err)
			}
			// Recovery forces a real store batch even though the root already
			// exists, instead of the trie's no-op commit fast path.
			duplicate, err := loaded.RecoverNode(ctx, root, encoded)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := duplicate.Commit(ctx, store); err != nil {
				t.Fatalf("duplicate commit consumed capacity: %v", err)
			}
			candidate, err := loaded.Update(ctx, []byte("key"), []byte("third"))
			if err != nil {
				t.Fatal(err)
			}
			if _, err := candidate.Commit(ctx, store); !errors.Is(err, mpt.ErrResourceLimit) {
				t.Errorf("over-limit historical commit error = %v", err)
			}
			if store.Root() != root || countStoreNodes(t, store) != 2 {
				t.Fatalf("over-limit commit changed publication: new root = %t, retained nodes = %d", store.Root() != root, countStoreNodes(t, store))
			}
			if value, err := loaded.Get(ctx, []byte("key")); err != nil || string(value) != "other" {
				t.Fatal("rejected commit changed snapshot contents")
			}
			result, err := store.Prune(ctx, mpt.DefaultReachabilityLimits())
			if err != nil || result.RemovedNodes() != 1 {
				t.Fatalf("Prune did not reclaim historical admission: %v", err)
			}
			if _, err := candidate.Commit(ctx, store); err != nil {
				t.Fatalf("reclaimed admission remained unavailable: %v", err)
			}
		})
	}
}
