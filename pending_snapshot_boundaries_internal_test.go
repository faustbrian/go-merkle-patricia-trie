package mpt

import (
	"bytes"
	"context"
	"errors"
	"os"
	"testing"
)

func requireHostedSnapshotBoundary(t *testing.T) {
	t.Helper()
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("snapshot encoding boundaries execute only in hosted CI")
	}
}

func assertSnapshotBoundaryValue(t *testing.T, trie RawTrie, key, want []byte) {
	t.Helper()
	got, err := trie.Get(context.Background(), key)
	if err != nil || !bytes.Equal(got, want) {
		t.Fatalf("Get() = %v, error %v, want %v", got, err, want)
	}
}

func TestHostedMaterializedSnapshotExactByteBoundary(t *testing.T) {
	requireHostedSnapshotBoundary(t)
	ctx := context.Background()
	limits := DefaultLimits()
	limits.MaxPendingBytes = 257
	trie, err := NewRawTrie(limits)
	if err != nil {
		t.Fatal(err)
	}
	value := bytes.Repeat([]byte{1}, 128)
	previous, err := trie.Update(ctx, []byte{0x00}, value)
	if err != nil {
		t.Fatal(err)
	}
	previousRoot, err := previous.Root()
	if err != nil {
		t.Fatal(err)
	}
	next, err := previous.Update(ctx, []byte{0x01}, value)
	if err != nil {
		t.Fatalf("exact materialized byte boundary rejected: %v", err)
	}
	if got := decodedNodeFootprint(next.snapshot.readRoot); got.bytes != 257 || got.depth != 2 || !got.complete {
		t.Fatalf("materialized footprint = %+v, want 257 bytes and depth 2", got)
	}
	// Equal leaf encodings deduplicate to one 133-byte record. The branch
	// encoding is 83 bytes and the extension root is 35: 251 live bytes.
	if next.snapshot.pendingStats.liveBytes != 251 {
		t.Fatalf("pending live bytes = %d, want 251", next.snapshot.pendingStats.liveBytes)
	}
	assertSnapshotBoundaryValue(t, next, []byte{0x00}, value)
	assertSnapshotBoundaryValue(t, next, []byte{0x01}, value)
	assertSnapshotBoundaryValue(t, previous, []byte{0x00}, value)
	if got, err := previous.Root(); err != nil || got != previousRoot {
		t.Fatal("update changed previous root")
	}
	if _, err := previous.Get(ctx, []byte{0x01}); !errors.Is(err, ErrAbsentKey) {
		t.Fatal("update changed previous key set")
	}
}

func TestHostedMaterializedSnapshotExactDepthBoundary(t *testing.T) {
	requireHostedSnapshotBoundary(t)
	for _, depth := range []int{1, 2} {
		name := "exact"
		if depth == 1 {
			name = "one short"
		}
		t.Run(name, func(t *testing.T) {
			limits := DefaultLimits()
			limits.MaxTraversalDepth = depth
			trie, err := NewRawTrie(limits)
			if err != nil {
				t.Fatal(err)
			}
			previous, err := trie.Update(context.Background(), []byte{0x00}, []byte{1})
			if err != nil {
				t.Fatal(err)
			}
			previousRoot, err := previous.Root()
			if err != nil {
				t.Fatal(err)
			}
			next, err := previous.Update(context.Background(), []byte{0x01}, []byte{2})
			if depth == 1 {
				if !errors.Is(err, ErrResourceLimit) {
					t.Fatalf("short depth error = %v, want ErrResourceLimit", err)
				}
				if _, err := next.Root(); !errors.Is(err, ErrUninitialized) {
					t.Fatal("rejected update returned initialized candidate")
				}
			} else {
				if err != nil {
					t.Fatalf("exact depth rejected: %v", err)
				}
				assertSnapshotBoundaryValue(t, next, []byte{0x00}, []byte{1})
				assertSnapshotBoundaryValue(t, next, []byte{0x01}, []byte{2})
			}
			assertSnapshotBoundaryValue(t, previous, []byte{0x00}, []byte{1})
			if got, err := previous.Root(); err != nil || got != previousRoot {
				t.Fatal("update changed previous root")
			}
		})
	}
}

func TestHostedPendingReplacementCompactsRetentionPressure(t *testing.T) {
	requireHostedSnapshotBoundary(t)
	limits := DefaultLimits()
	limits.MaxPendingNodes, limits.MaxPendingBytes = 1, 10
	trie, err := NewRawTrie(limits)
	if err != nil {
		t.Fatal(err)
	}
	previous, err := trie.Update(context.Background(), []byte{0x00}, []byte{1})
	if err != nil {
		t.Fatal(err)
	}
	next, err := previous.Update(context.Background(), []byte{0x00}, []byte{2})
	if err != nil {
		t.Fatal(err)
	}
	if next.snapshot.parent != nil || len(next.snapshot.removed) != 0 || len(next.snapshot.pending) != 1 || next.snapshot.pendingStats != (pendingAccounting{1, 5, 1, 5}) {
		t.Fatalf("replacement retained stale ancestry: parent %t, removed %d, pending %d, accounting %+v", next.snapshot.parent != nil, len(next.snapshot.removed), len(next.snapshot.pending), next.snapshot.pendingStats)
	}
	assertSnapshotBoundaryValue(t, next, []byte{0x00}, []byte{2})
	assertSnapshotBoundaryValue(t, previous, []byte{0x00}, []byte{1})
}

func TestHostedRecoveryCompactsPendingRetentionPressure(t *testing.T) {
	requireHostedSnapshotBoundary(t)
	ctx := context.Background()
	limits := DefaultLimits()
	limits.MaxPendingNodes, limits.MaxPendingBytes = 3, 20
	// EmptyRoot has no backing node; all subsequent reads resolve the actual
	// replacement leaf from owned pending records, not from this empty reader.
	reader := nodeReaderFunc(func(context.Context, Root) ([]byte, error) {
		return nil, ErrMissingNode
	})
	trie, err := LoadRawTrie(EmptyRoot(), reader, limits)
	if err != nil {
		t.Fatal(err)
	}
	for value := byte(1); value <= 3; value++ {
		trie, err = trie.Update(ctx, []byte{0x00}, []byte{value})
		if err != nil {
			t.Fatal(err)
		}
	}
	if trie.snapshot.parent == nil || trie.snapshot.parent.depth != 2 || trie.snapshot.pendingStats != (pendingAccounting{1, 5, 3, 15}) {
		t.Fatalf("pre-recovery ancestry: parent %v, accounting %+v", trie.snapshot.parent, trie.snapshot.pendingStats)
	}
	root, err := trie.Root()
	if err != nil {
		t.Fatal(err)
	}
	leaf, err := newLeaf([]byte{1, 0}, []byte{4})
	if err != nil {
		t.Fatal(err)
	}
	encoded, _, err := encodeNode(leaf)
	if err != nil || len(encoded) != 5 {
		t.Fatalf("recovery canonical leaf size = %d, error %v, want 5", len(encoded), err)
	}
	// Recovery accepts a canonical hash-bound record without requiring current
	// root reachability; it adds an overlay, never changes the commitment.
	next, err := trie.RecoverNode(ctx, keccakRoot(encoded), encoded)
	if err != nil {
		t.Fatal(err)
	}
	if next.snapshot.parent != nil || len(next.snapshot.pending) != 2 || next.snapshot.pendingStats != (pendingAccounting{2, 10, 2, 10}) {
		t.Fatalf("recovery retained stale ancestry: parent %t, pending %d, accounting %+v", next.snapshot.parent != nil, len(next.snapshot.pending), next.snapshot.pendingStats)
	}
	if got, err := next.Root(); err != nil || got != root {
		t.Fatal("recovery changed commitment")
	}
	assertSnapshotBoundaryValue(t, next, []byte{0x00}, []byte{3})
	assertSnapshotBoundaryValue(t, trie, []byte{0x00}, []byte{3})
	if trie.snapshot.pendingStats != (pendingAccounting{1, 5, 3, 15}) || trie.snapshot.parent.depth != 2 {
		t.Fatal("recovery changed previous accounting or ancestry")
	}
	if got, err := trie.Root(); err != nil || got != root {
		t.Fatal("recovery changed previous root")
	}
}

func TestHostedPendingAccountingCompactsAtLayerDepthBoundary(t *testing.T) {
	requireHostedSnapshotBoundary(t)
	ctx := context.Background()
	limits := DefaultLimits()
	trie, err := NewRawTrie(limits)
	if err != nil {
		t.Fatal(err)
	}
	for value := byte(1); value <= 32; value++ {
		trie, err = trie.Update(ctx, []byte{0x00}, []byte{value})
		if err != nil {
			t.Fatal(err)
		}
	}
	if layer := snapshotPendingLayer(trie.snapshot); layer == nil || layer.depth != 32 || trie.snapshot.pendingStats != (pendingAccounting{1, 5, 32, 160}) {
		t.Fatalf("real replacement boundary: layer %v, accounting %+v", layer, trie.snapshot.pendingStats)
	}
	leaf, err := newLeaf([]byte{0, 0}, []byte{33})
	if err != nil {
		t.Fatal(err)
	}
	encoded, _, err := encodeNode(leaf)
	if err != nil || len(encoded) != 5 {
		t.Fatalf("next canonical leaf size = %d, error %v, want 5", len(encoded), err)
	}
	root, err := trie.Root()
	if err != nil {
		t.Fatal(err)
	}
	stats, compact, err := admitPending(ctx, trie.snapshot, map[Root][]byte{keccakRoot(encoded): encoded}, map[Root]struct{}{root: {}}, limits)
	if err != nil || !compact || stats != (pendingAccounting{1, 5, 1, 5}) {
		t.Fatalf("depth-boundary admission = %+v, compact %t, error %v", stats, compact, err)
	}
	if trie.snapshot.pendingStats != (pendingAccounting{1, 5, 32, 160}) || snapshotPendingLayer(trie.snapshot).depth != 32 {
		t.Fatal("admission changed previous accounting or ancestry")
	}
	assertSnapshotBoundaryValue(t, trie, []byte{0x00}, []byte{32})
}
