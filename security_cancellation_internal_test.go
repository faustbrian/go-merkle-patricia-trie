package mpt

import (
	"context"
	"errors"
	"runtime"
	"testing"
)

func TestProofCancellationPrecedesPendingMaterialization(t *testing.T) {
	ctx := context.Background()
	trie, err := NewRawTrie(DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	for index := range 1024 {
		key := []byte{byte(index >> 8), byte(index)}
		value := make([]byte, 64)
		value[0], value[1] = key[0], key[1]
		trie, err = trie.Update(ctx, key, value)
		if err != nil {
			t.Fatal(err)
		}
	}
	root, err := trie.Root()
	if err != nil {
		t.Fatal(err)
	}
	// Two existing checks admit the context and single key. The next check
	// must stop pending preparation before copying the whole unrelated map.
	canceled := &nthErrorContext{at: 3}
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	_, err = trie.ProveMany(canceled, [][]byte{{0, 0}})
	runtime.ReadMemStats(&after)
	if !errors.Is(err, ErrCanceled) || !errors.Is(err, context.Canceled) {
		t.Fatalf("pending preparation did not observe cancellation: %v", err)
	}
	if allocated := after.TotalAlloc - before.TotalAlloc; allocated > 64<<10 {
		t.Errorf("canceled proof copied unrelated pending state: allocated %d bytes", allocated)
	}
	if unchanged, err := trie.Root(); err != nil || unchanged != root {
		t.Fatal("canceled proof changed original commitment")
	}
}

func TestSecurityAdmissionCancellationPreservesSnapshots(t *testing.T) {
	ctx := context.Background()
	trie, err := NewRawTrie(DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	trie, err = trie.Update(ctx, []byte("key"), []byte("value"))
	if err != nil {
		t.Fatal(err)
	}
	root, err := trie.Root()
	if err != nil {
		t.Fatal(err)
	}
	proof, err := trie.Prove(ctx, []byte("key"))
	if err != nil {
		t.Fatal(err)
	}
	encoded := proof.Nodes()[0]
	loaded, err := LoadRawTrie(root, nodeReaderFunc(func(context.Context, Root) ([]byte, error) { return encoded, nil }), DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	for name, operation := range map[string]func(context.Context) error{
		"update": func(ctx context.Context) error {
			_, err := trie.Update(ctx, []byte("key"), []byte("replacement"))
			return err
		},
		"recovery":          func(ctx context.Context) error { _, err := loaded.RecoverNode(ctx, root, encoded); return err },
		"multi-proof input": func(ctx context.Context) error { _, err := trie.ProveMany(ctx, [][]byte{[]byte("key")}); return err },
		"commit preparation": func(ctx context.Context) error {
			store := &cancellationCaptureStore{}
			_, err := trie.Commit(ctx, store)
			if errors.Is(err, ErrCanceled) && store.published {
				t.Fatal("canceled preparation reached store publication")
			}
			return err
		},
	} {
		t.Run(name, func(t *testing.T) {
			probe := &nthErrorContext{at: int(^uint(0) >> 1)}
			if err := operation(probe); err != nil {
				t.Fatal(err)
			}
			// This tiny graph exercises all observed operation seams. The
			// deterministic context avoids timing or goroutine-based tests.
			for at := 1; at <= probe.calls; at++ {
				if err := operation(&nthErrorContext{at: at}); !errors.Is(err, ErrCanceled) || !errors.Is(err, context.Canceled) {
					t.Fatalf("cancellation at seam %d returned %v", at, err)
				}
				if unchanged, err := trie.Root(); err != nil || unchanged != root {
					t.Fatal("canceled operation changed original root")
				}
				if value, err := trie.Get(ctx, []byte("key")); err != nil || string(value) != "value" {
					t.Fatal("canceled operation changed original value")
				}
			}
		})
	}
}

type cancellationCaptureStore struct{ published bool }

func (*cancellationCaptureStore) GetNode(context.Context, Root) ([]byte, error) {
	return nil, ErrMissingNode
}

func (store *cancellationCaptureStore) CommitTrie(context.Context, StoreCommit) error {
	store.published = true
	return nil
}

func TestPendingPressureCancellationPreservesSnapshot(t *testing.T) {
	ctx := context.Background()
	limits := DefaultLimits()
	limits.MaxPendingNodes, limits.MaxPendingBytes = 1, 160
	trie, err := NewRawTrie(limits)
	if err != nil {
		t.Fatal(err)
	}
	trie, err = trie.Update(ctx, []byte("key"), make([]byte, 64))
	if err != nil {
		t.Fatal(err)
	}
	root, err := trie.Root()
	if err != nil {
		t.Fatal(err)
	}
	probe := &nthErrorContext{at: int(^uint(0) >> 1)}
	if _, err := trie.Update(probe, []byte("key"), []byte("replacement")); err != nil {
		t.Fatal(err)
	}
	for at := 1; at <= probe.calls; at++ {
		if _, err := trie.Update(&nthErrorContext{at: at}, []byte("key"), []byte("replacement")); !errors.Is(err, ErrCanceled) {
			t.Fatalf("canceled compaction at seam %d returned %v", at, err)
		}
		if unchanged, err := trie.Root(); err != nil || unchanged != root {
			t.Fatal("canceled compaction changed original commitment")
		}
	}
}

func TestRecoveryCompactionCancellationPreservesSnapshot(t *testing.T) {
	ctx := context.Background()
	reader := nodeReaderFunc(func(context.Context, Root) ([]byte, error) { return nil, ErrMissingNode })
	first, _, err := encodeNode(&leafNode{value: []byte("original")})
	if err != nil {
		t.Fatal(err)
	}
	root := keccakRoot(first)
	trie, err := LoadRawTrie(root, reader, DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	trie, err = trie.RecoverNode(ctx, root, first)
	if err != nil {
		t.Fatal(err)
	}
	for index := 1; index < maximumPendingLayerDepth; index++ {
		encoded, _, err := encodeNode(&leafNode{value: []byte{byte(index)}})
		if err != nil {
			t.Fatal(err)
		}
		trie, err = trie.RecoverNode(ctx, keccakRoot(encoded), encoded)
		if err != nil {
			t.Fatal(err)
		}
	}
	last, _, err := encodeNode(&leafNode{value: []byte("last")})
	if err != nil {
		t.Fatal(err)
	}
	probe := &nthErrorContext{at: int(^uint(0) >> 1)}
	if _, err := trie.RecoverNode(probe, keccakRoot(last), last); err != nil {
		t.Fatal(err)
	}
	for at := 1; at <= probe.calls; at++ {
		if _, err := trie.RecoverNode(&nthErrorContext{at: at}, keccakRoot(last), last); !errors.Is(err, ErrCanceled) {
			t.Fatalf("canceled recovery compaction at seam %d returned %v", at, err)
		}
		if value, err := trie.Get(ctx, nil); err != nil || string(value) != "original" {
			t.Fatal("canceled recovery changed original snapshot")
		}
	}
}

func TestOwnedPreparationRejectsNilContextWithoutPartialResult(t *testing.T) {
	var invalidContext context.Context
	if pending, err := materializeSnapshotPendingContext(invalidContext, nil); pending != nil || !errors.Is(err, ErrInvalidContext) {
		t.Fatal("nil-context materialization returned partial state")
	}
	if pending, err := materializePendingLayerContext(invalidContext, nil); pending != nil || !errors.Is(err, ErrInvalidContext) {
		t.Fatal("nil-context layer materialization returned partial state")
	}
	if commit, err := newStoreCommitContext(invalidContext, Root{}, Root{}, nil); commit.NodeCount() != 0 || !errors.Is(err, ErrInvalidContext) {
		t.Fatal("nil-context commit preparation returned handles")
	}
}
