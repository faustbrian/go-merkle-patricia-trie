package memory_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	mpt "github.com/faustbrian/go-merkle-patricia-trie/v2"
	"github.com/faustbrian/go-merkle-patricia-trie/v2/memory"
)

func TestMemoryAdmissionRequiresEachPositiveDimension(t *testing.T) {
	for _, limits := range []memory.Limits{{MaxStoredNodes: 0, MaxStoredBytes: 1}, {MaxStoredNodes: 1, MaxStoredBytes: 0}} {
		if store, err := memory.NewWithLimits(limits); store != nil || !errors.Is(err, mpt.ErrResourceLimit) {
			t.Fatal("zero admission dimension was accepted")
		}
	}
	if store, err := memory.NewWithLimits(memory.Limits{MaxStoredNodes: 1, MaxStoredBytes: 1}); store == nil || err != nil {
		t.Fatal("positive admission dimensions were refused")
	}
}

func TestMemoryAdmissionAcceptsExactOneNodeCapacity(t *testing.T) {
	ctx := context.Background()
	trie := mustMemoryTrie(t, map[string]string{"key": "value"})
	_, commit := captureTrieCommit(t, trie, newCaptureStore())
	if commit.NodeCount() != 1 {
		t.Fatal("fixture is not a one-node canonical commit")
	}
	store := boundedMemoryStore(t, memory.Limits{MaxStoredNodes: 1, MaxStoredBytes: commit.Nodes()[0].EncodedLen()})
	if err := store.CommitTrie(ctx, commit); err != nil {
		t.Fatalf("exact capacity was refused: %v", err)
	}
	assertMemoryValue(t, store, "key", "value")
}

func TestMemoryAdmissionCountsEveryFreshNodeCumulatively(t *testing.T) {
	ctx := context.Background()
	capture := newCaptureStore()
	trie, first := captureTrieCommit(t, mustMemoryTrie(t, map[string]string{"alpha": strings.Repeat("a", 64)}), capture)
	var err error
	trie, err = trie.Update(ctx, []byte("alpha"), []byte(strings.Repeat("c", 64)))
	if err != nil {
		t.Fatal(err)
	}
	trie, err = trie.Update(ctx, []byte("beta"), []byte(strings.Repeat("b", 64)))
	if err != nil {
		t.Fatal(err)
	}
	_, next := captureTrieCommit(t, trie, capture)
	if first.NodeCount() != 1 || next.NodeCount() < 2 {
		t.Fatal("fixture lacks a historical node and multiple fresh nodes")
	}
	for _, stored := range next.Nodes() {
		if stored.Hash() == first.Nodes()[0].Hash() {
			t.Fatal("fixture unexpectedly reuses the historical node")
		}
	}
	limits := memory.DefaultLimits()
	limits.MaxStoredNodes = next.NodeCount()
	store := boundedMemoryStore(t, limits)
	if err := store.CommitTrie(ctx, first); err != nil {
		t.Fatal(err)
	}
	assertMemoryCommitRefused(t, store, next)
	assertMemoryValue(t, store, "alpha", strings.Repeat("a", 64))
}

func TestMemoryPrunePreservesCumulativeByteAdmission(t *testing.T) {
	ctx := context.Background()
	capture := newCaptureStore()
	trie, first := captureTrieCommit(t, mustMemoryTrie(t, map[string]string{
		"alpha": strings.Repeat("a", 64), "beta": strings.Repeat("b", 64),
	}), capture)
	initialBytes := 0
	for _, stored := range first.Nodes() {
		initialBytes += stored.EncodedLen()
	}
	probe := memory.New()
	if err := probe.CommitTrie(ctx, first); err != nil {
		t.Fatal(err)
	}
	if _, err := probe.Prune(ctx, mpt.DefaultReachabilityLimits()); err != nil {
		t.Fatal(err)
	}
	retained := memoryNodeSnapshot(t, probe)
	retainedBytes, largest := 0, 0
	for _, encoded := range retained {
		retainedBytes += len(encoded)
		largest = max(largest, len(encoded))
	}
	var err error
	trie, err = trie.Delete(ctx, []byte("alpha"))
	if err != nil {
		t.Fatal(err)
	}
	trie, err = trie.Update(ctx, []byte("beta"), []byte(strings.Repeat("c", 80)))
	if err != nil {
		t.Fatal(err)
	}
	_, next := captureTrieCommit(t, trie, capture)
	addedBytes := 0
	for _, stored := range next.Nodes() {
		if _, exists := retained[stored.Hash()]; !exists {
			addedBytes += stored.EncodedLen()
		}
	}
	// Every possible final map entry is too small to represent the sum. The
	// candidate fits that incorrect slack, but not the true remaining budget.
	if len(retained) < 2 || addedBytes <= initialBytes-retainedBytes || addedBytes > initialBytes-largest {
		t.Fatal("fixture does not distinguish cumulative accounting for every map order")
	}
	limits := memory.DefaultLimits()
	limits.MaxStoredBytes = initialBytes
	store := boundedMemoryStore(t, limits)
	if err := store.CommitTrie(ctx, first); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Prune(ctx, mpt.DefaultReachabilityLimits()); err != nil {
		t.Fatal(err)
	}
	assertMemoryCommitRefused(t, store, next)
	assertMemoryValue(t, store, "alpha", strings.Repeat("a", 64))
	assertMemoryValue(t, store, "beta", strings.Repeat("b", 64))
}

func TestMemoryAdmissionContinuesAfterDuplicateBeforeFreshNode(t *testing.T) {
	ctx := context.Background()
	for _, recoveredKey := range []string{"alpha", "beta"} {
		for _, replacement := range []string{strings.Repeat("c", 64), strings.Repeat("d", 64)} {
			capture := newCaptureStore()
			trie, first := captureTrieCommit(t, mustMemoryTrie(t, map[string]string{
				"alpha": strings.Repeat("a", 64), "beta": strings.Repeat("b", 64),
			}), capture)
			proof, err := trie.Prove(ctx, []byte(recoveredKey))
			if err != nil {
				t.Fatal(err)
			}
			path := proof.Nodes()
			if len(path) < 2 {
				t.Fatal("fixture lacks a hashed child")
			}
			encoded := path[len(path)-1]
			var hash mpt.Root
			found := false
			base := make(map[mpt.Root][]byte)
			for _, stored := range first.Nodes() {
				base[stored.Hash()] = stored.Encoded()
				if bytes.Equal(stored.Encoded(), encoded) {
					hash, found = stored.Hash(), true
				}
			}
			if !found {
				t.Fatal("proof child is not in the canonical base commit")
			}
			trie, err = trie.RecoverNode(ctx, hash, encoded)
			if err != nil {
				t.Fatal(err)
			}
			changedKey := "alpha"
			if recoveredKey == changedKey {
				changedKey = "beta"
			}
			trie, err = trie.Update(ctx, []byte(changedKey), []byte(replacement))
			if err != nil {
				t.Fatal(err)
			}
			_, next := captureTrieCommit(t, trie, capture)
			duplicateSeen, freshAfterDuplicate, freshCount := false, false, 0
			nodes := next.Nodes()
			for index, stored := range nodes {
				if index > 0 {
					previous, current := nodes[index-1].Hash(), stored.Hash()
					if bytes.Compare(previous[:], current[:]) >= 0 {
						t.Fatal("canonical commit does not have strictly ascending hashes")
					}
				}
				if _, exists := base[stored.Hash()]; exists {
					duplicateSeen = true
				} else {
					freshCount++
					freshAfterDuplicate = freshAfterDuplicate || duplicateSeen
				}
			}
			if !freshAfterDuplicate {
				continue
			}
			limits := memory.DefaultLimits()
			limits.MaxStoredNodes = len(base) + freshCount - 1
			if next.NodeCount() > limits.MaxStoredNodes {
				t.Fatal("fixture would be rejected before duplicate accounting")
			}
			store := boundedMemoryStore(t, limits)
			if err := store.CommitTrie(ctx, first); err != nil {
				t.Fatal(err)
			}
			assertMemoryCommitRefused(t, store, next)
			assertMemoryValue(t, store, "alpha", strings.Repeat("a", 64))
			assertMemoryValue(t, store, "beta", strings.Repeat("b", 64))
			t.Logf("canonical mixed batch: recovered=%s changed=%s base=%d fresh=%d batch=%d", recoveredKey, changedKey, len(base), freshCount, next.NodeCount())
			return
		}
	}
	t.Fatal("finite canonical fixtures lack duplicate-before-fresh order")
}

func boundedMemoryStore(t *testing.T, limits memory.Limits) *memory.Store {
	t.Helper()
	store, err := memory.NewWithLimits(limits)
	if err != nil {
		t.Fatal(err)
	}
	return store
}

func memoryNodeSnapshot(t *testing.T, store *memory.Store) map[mpt.Root][]byte {
	t.Helper()
	nodes := make(map[mpt.Root][]byte)
	if err := store.IterateNodes(context.Background(), 32, func(hash mpt.Root, encoded []byte) error {
		nodes[hash] = encoded
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return nodes
}

func assertMemoryCommitRefused(t *testing.T, store *memory.Store, commit mpt.StoreCommit) {
	t.Helper()
	root, before := store.Root(), memoryNodeSnapshot(t, store)
	if err := store.CommitTrie(context.Background(), commit); !errors.Is(err, mpt.ErrResourceLimit) {
		t.Fatalf("cumulative admission error = %v", err)
	}
	after := memoryNodeSnapshot(t, store)
	if store.Root() != root || len(after) != len(before) {
		t.Fatal("refused admission changed publication or node count")
	}
	for hash, encoded := range before {
		if !bytes.Equal(encoded, after[hash]) {
			t.Fatal("refused admission changed stored bytes")
		}
	}
}

func assertMemoryValue(t *testing.T, store *memory.Store, key, expected string) {
	t.Helper()
	trie, err := mpt.LoadRawTrie(store.Root(), store, mpt.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	value, err := trie.Get(context.Background(), []byte(key))
	if err != nil || string(value) != expected {
		t.Fatalf("prior snapshot value changed: %v", err)
	}
}
