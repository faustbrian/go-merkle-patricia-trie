package mpt

import (
	"bytes"
	"context"
	"errors"
	"runtime"
	"testing"
)

func TestPendingCompactionPropagatesDecodeFailureWithoutPartialResult(t *testing.T) {
	encoded := []byte{0xc0} // An empty RLP list is not a trie node.
	_, expected := decodeNode(encoded)
	if expected == nil {
		t.Fatal("ordinary error fixture unexpectedly decoded")
	}
	hash := keccakRoot(encoded)
	root, err := newExtension([]byte{1}, hashNode(hash))
	if err != nil {
		t.Fatal(err)
	}
	rootEncoding, _, err := encodeNode(root)
	if err != nil {
		t.Fatal(err)
	}
	pending := map[Root][]byte{hash: encoded}
	original := append([]byte(nil), encoded...)
	result, err := retainReferencedPending(context.Background(), root, keccakRoot(rootEncoding), pending, DefaultLimits())
	if result != nil || err == nil || err.Error() != expected.Error() || !errors.Is(err, ErrMalformedNode) {
		t.Fatal("compaction did not propagate the decode failure with a nil result")
	}
	if len(pending) != 1 || !bytes.Equal(pending[hash], original) {
		t.Fatal("failed compaction changed its input")
	}
}

type cancelAfterGuardContext struct {
	context.Context
	cancel       context.CancelFunc
	calls, after int
}

func (ctx *cancelAfterGuardContext) Err() error {
	err := ctx.Context.Err()
	ctx.calls++
	if ctx.calls == ctx.after {
		ctx.cancel()
	}
	return err
}

func TestRecoveredMutationCancellationPrecedesOverlayInheritance(t *testing.T) {
	ctx := context.Background()
	source, err := NewRawTrie(DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range [][]byte{{0}, {1}} {
		source, err = source.Update(ctx, key, make([]byte, 64))
		if err != nil {
			t.Fatal(err)
		}
	}
	reader := nodeReaderFunc(func(context.Context, Root) ([]byte, error) { return nil, ErrMissingNode })
	recovered, err := LoadRawTrie(source.snapshot.hash, reader, DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	for hash, encoded := range materializeSnapshotPending(source.snapshot) {
		recovered, err = recovered.RecoverNode(ctx, hash, encoded)
		if err != nil {
			t.Fatal(err)
		}
	}
	for _, operation := range []struct {
		name string
		run  func(context.Context, RawTrie) (RawTrie, error)
	}{
		{"update", func(ctx context.Context, trie RawTrie) (RawTrie, error) {
			return trie.Update(ctx, []byte{0}, []byte("replacement"))
		}},
		{"delete", func(ctx context.Context, trie RawTrie) (RawTrie, error) { return trie.Delete(ctx, []byte{0}) }},
		{"batch", func(ctx context.Context, trie RawTrie) (RawTrie, error) {
			return trie.ApplyBatch(ctx, []Mutation{Put([]byte{0}, []byte("replacement"))})
		}},
	} {
		t.Run(operation.name, func(t *testing.T) {
			withoutOverlay := *recovered.snapshot
			withoutOverlay.recovered = nil
			probe := &nthErrorContext{at: int(^uint(0) >> 1)}
			if _, err := operation.run(probe, RawTrie{snapshot: &withoutOverlay}); err != nil {
				t.Fatal(err)
			}
			base, cancel := context.WithCancel(ctx)
			defer cancel()
			// Cancellation becomes visible immediately after the final existing
			// guard, before ownership of the recovered overlay is transferred.
			canceled := &cancelAfterGuardContext{Context: base, cancel: cancel, after: probe.calls}
			candidate, err := operation.run(canceled, recovered)
			if !errors.Is(err, ErrCanceled) || !errors.Is(err, context.Canceled) {
				t.Errorf("canceled overlay inheritance returned a snapshot: %v", err)
			}
			if _, err := candidate.Root(); err == nil {
				t.Error("canceled mutation exposed partial result")
			}
			if value, err := recovered.Get(ctx, []byte{1}); err != nil || len(value) != 64 {
				t.Fatal("cancellation changed original recovered snapshot")
			}
			allSeams := &nthErrorContext{at: int(^uint(0) >> 1)}
			if _, err := operation.run(allSeams, recovered); err != nil {
				t.Fatal(err)
			}
			for at := 1; at <= allSeams.calls; at++ {
				candidate, err := operation.run(&nthErrorContext{at: at}, recovered)
				if !errors.Is(err, ErrCanceled) {
					t.Fatalf("cancellation seam %d returned %v", at, err)
				}
				if _, err := candidate.Root(); err == nil {
					t.Fatal("canceled inheritance exposed partial result")
				}
			}
		})
	}
}

func TestOwnedEncodingAndProofCancellationReturnsNoPartialResult(t *testing.T) {
	ctx := context.Background()
	trie, err := NewRawTrie(DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range [][]byte{{0}, {1}} {
		trie, err = trie.Update(ctx, key, []byte("value"))
		if err != nil {
			t.Fatal(err)
		}
	}
	prepared := *trie.snapshot
	prepared.root, prepared.pending, prepared.parent, prepared.removed = prepared.readRoot, nil, nil, nil
	owned := RawTrie{snapshot: &prepared}
	for name, operation := range map[string]func(context.Context) error{
		"update": func(ctx context.Context) error { _, err := trie.Update(ctx, []byte{2}, []byte("value")); return err },
		"proof":  func(ctx context.Context) error { _, err := owned.ProveMany(ctx, [][]byte{{0}, {1}}); return err },
	} {
		t.Run(name, func(t *testing.T) {
			probe := &nthErrorContext{at: int(^uint(0) >> 1)}
			if err := operation(probe); err != nil {
				t.Fatal(err)
			}
			for at := 1; at <= probe.calls; at++ {
				if err := operation(&nthErrorContext{at: at}); !errors.Is(err, ErrCanceled) {
					t.Fatalf("owned preparation seam %d returned %v", at, err)
				}
			}
		})
	}
}

func TestPendingCompactionCancellationAndUnknownFrontier(t *testing.T) {
	ctx := context.Background()
	trie, err := NewRawTrie(DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range [][]byte{{0}, {1}} {
		trie, err = trie.Update(ctx, key, make([]byte, 64))
		if err != nil {
			t.Fatal(err)
		}
	}
	pending := materializeSnapshotPending(trie.snapshot)
	probe := &nthErrorContext{at: int(^uint(0) >> 1)}
	if _, err := retainReferencedPending(probe, trie.snapshot.root, trie.snapshot.hash, pending, DefaultLimits()); err != nil {
		t.Fatal(err)
	}
	for at := 1; at <= probe.calls; at++ {
		if result, err := retainReferencedPending(&nthErrorContext{at: at}, trie.snapshot.root, trie.snapshot.hash, pending, DefaultLimits()); result != nil || !errors.Is(err, ErrCanceled) {
			t.Fatalf("compaction seam %d returned partial result: %v", at, err)
		}
	}
	limits := DefaultLimits()
	limits.MaxTraversalNodes = 1
	if result, err := retainReferencedPending(ctx, trie.snapshot.root, trie.snapshot.hash, pending, limits); result != nil || !errors.Is(err, ErrResourceLimit) {
		t.Fatal("compaction traversal bound did not reject atomically")
	}
	unknown := Root{1}
	if result, err := retainReferencedPending(ctx, hashNode(unknown), trie.snapshot.hash, pending, DefaultLimits()); err != nil || len(result) != len(pending) {
		t.Fatal("unknown backing frontier discarded potentially needed owned nodes")
	}
}

func TestRecoveryCancellationPrecedesOverlayMerge(t *testing.T) {
	ctx := context.Background()
	reader := nodeReaderFunc(func(context.Context, Root) ([]byte, error) { return nil, ErrMissingNode })
	encoded, _, err := encodeNode(&leafNode{value: []byte("original")})
	if err != nil {
		t.Fatal(err)
	}
	trie, err := LoadRawTrie(keccakRoot(encoded), reader, DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	for index := range 64 {
		encoded, _, err := encodeNode(&leafNode{value: []byte{byte(index)}})
		if err != nil {
			t.Fatal(err)
		}
		trie, err = trie.RecoverNode(ctx, keccakRoot(encoded), encoded)
		if err != nil {
			t.Fatal(err)
		}
	}
	next, _, err := encodeNode(&leafNode{value: []byte("next")})
	if err != nil {
		t.Fatal(err)
	}
	base, cancel := context.WithCancel(ctx)
	defer cancel()
	// Initial context admission, three pending-admission passes, and the
	// validated-node guard precede transfer of the existing recovery overlay.
	canceled := &cancelAfterGuardContext{Context: base, cancel: cancel, after: 5}
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	candidate, err := trie.RecoverNode(canceled, keccakRoot(next), next)
	runtime.ReadMemStats(&after)
	if !errors.Is(err, ErrCanceled) || !errors.Is(err, context.Canceled) {
		t.Fatalf("recovery did not observe cancellation: %v", err)
	}
	if _, err := candidate.Root(); err == nil {
		t.Fatal("canceled recovery exposed partial result")
	}
	if allocated := after.TotalAlloc - before.TotalAlloc; allocated > 4<<10 {
		t.Errorf("canceled recovery copied the already-owned overlay: %d bytes", allocated)
	}
}
