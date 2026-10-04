package mpt

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestRangeItemByteAccountingBoundaries(t *testing.T) {
	for _, tc := range []struct {
		name      string
		key       []byte
		value     []byte
		budget    int
		wantLimit bool
	}{
		{"exact", []byte{1, 2}, []byte{3, 4}, 4, false},
		{"key fills capacity", []byte{1, 2}, []byte{3}, 2, true},
		{"key exceeds capacity", []byte{1, 2}, []byte{3}, 1, true},
		{"value exceeds remainder", []byte{1}, []byte{2, 3, 4}, 3, true},
		{"spare capacity", []byte{1, 2}, []byte{3}, 4, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			limits := DefaultLimits()
			limits.MaxProofBytes = tc.budget
			state := rangeGenerationState{builder: newMultiProofBuilder(context.Background(), &trieSnapshot{limits: limits})}
			err := state.emit(bytesToNibbles(tc.key), tc.value)
			if errors.Is(err, ErrResourceLimit) != tc.wantLimit || (err != nil && !tc.wantLimit) {
				t.Fatalf("emit = %v, want limit %t", err, tc.wantLimit)
			}
			if tc.wantLimit {
				if len(state.items) != 0 || state.itemBytes != 0 {
					t.Fatal("rejection changed output accounting")
				}
			} else if len(state.items) != 1 || state.itemBytes != len(tc.key)+len(tc.value) || !bytes.Equal(state.items[0].Key(), tc.key) || !bytes.Equal(state.items[0].Value(), tc.value) {
				t.Fatal("admitted output or byte accounting differs")
			}
			err = validateRangeItems([]RangeItem{NewRangeItem(tc.key, tc.value)}, rangeBounds{}, limits)
			if errors.Is(err, ErrResourceLimit) != tc.wantLimit || (err != nil && !tc.wantLimit) {
				t.Fatalf("validate = %v, want limit %t", err, tc.wantLimit)
			}
		})
	}
	limits := DefaultLimits()
	limits.MaxProofBytes = 5
	state := rangeGenerationState{builder: newMultiProofBuilder(context.Background(), &trieSnapshot{limits: limits})}
	if err := state.emit([]byte{0, 1}, []byte{2, 3}); err != nil {
		t.Fatal(err)
	}
	if err := state.emit([]byte{0, 2}, []byte{4, 5}); !errors.Is(err, ErrResourceLimit) || len(state.items) != 1 || state.itemBytes != 3 {
		t.Fatal("second item ignored remaining capacity")
	}
	if err := validateRangeItems([]RangeItem{NewRangeItem([]byte{1}, []byte{2, 3}), NewRangeItem([]byte{2}, []byte{4, 5})}, rangeBounds{}, limits); !errors.Is(err, ErrResourceLimit) {
		t.Fatal("claim ignored cumulative item bytes")
	}
}

func TestPendingAccountingExactReplacementAndRetention(t *testing.T) {
	ctx := context.Background()
	a, b := Root{1}, Root{2}
	previous := &trieSnapshot{pending: map[Root][]byte{a: {1, 2}}, pendingStats: pendingAccounting{liveNodes: 1, liveBytes: 2, retainedNodes: 1, retainedBytes: 2}}
	for _, tc := range []struct {
		name           string
		added          map[Root][]byte
		removed        map[Root]struct{}
		nodes, bytes   int
		want           pendingAccounting
		compact, limit bool
	}{
		{"add exact", map[Root][]byte{b: {3}}, nil, 2, 3, pendingAccounting{2, 3, 2, 3}, false, false},
		{"replace spare", map[Root][]byte{a: {3, 4, 5}}, nil, 3, 8, pendingAccounting{1, 3, 2, 5}, false, false},
		{"replace byte pressure", map[Root][]byte{a: {3, 4, 5}}, nil, 3, 3, pendingAccounting{1, 3, 1, 3}, true, false},
		{"replace node pressure", map[Root][]byte{a: {3}}, nil, 1, 8, pendingAccounting{1, 1, 1, 1}, true, false},
		{"remove and add", map[Root][]byte{b: {3}}, map[Root]struct{}{a: {}}, 2, 3, pendingAccounting{1, 1, 2, 3}, false, false},
		{"remove replacement", map[Root][]byte{a: {3}}, map[Root]struct{}{a: {}}, 2, 3, pendingAccounting{1, 1, 2, 3}, false, false},
		{"node excess", map[Root][]byte{b: {3}}, nil, 1, 8, pendingAccounting{}, false, true},
		{"byte excess", map[Root][]byte{b: {3, 4}}, nil, 3, 3, pendingAccounting{}, false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			limits := DefaultLimits()
			limits.MaxPendingNodes, limits.MaxPendingBytes = tc.nodes, tc.bytes
			got, compact, err := admitPending(ctx, previous, tc.added, tc.removed, limits)
			if errors.Is(err, ErrResourceLimit) != tc.limit || (err != nil && !tc.limit) || got != tc.want || compact != tc.compact {
				t.Fatalf("accounting = %+v, compact %t, error %v; want %+v, %t, limit %t", got, compact, err, tc.want, tc.compact, tc.limit)
			}
			if !reflect.DeepEqual(previous.pending[a], []byte{1, 2}) || previous.pendingStats != (pendingAccounting{1, 2, 1, 2}) {
				t.Fatal("admission changed prior snapshot")
			}
		})
	}
}

func TestDecodedFootprintSparseChildrenAndCompleteness(t *testing.T) {
	leaf, err := newLeaf([]byte{1}, []byte{2, 3})
	if err != nil {
		t.Fatal(err)
	}
	var children [16]node
	children[1], children[15] = leaf, hashNode(Root{1})
	branch, err := newBranch(children, []byte{4})
	if err != nil {
		t.Fatal(err)
	}
	if got := decodedNodeFootprint(branch); got != (decodedFootprint{nodes: 2, bytes: 4, depth: 1, complete: false}) {
		t.Fatalf("sparse unresolved footprint = %+v", got)
	}
	children[15] = leaf
	branch, err = newBranch(children, nil)
	if err != nil {
		t.Fatal(err)
	}
	extension, err := newExtension([]byte{1, 2}, branch)
	if err != nil {
		t.Fatal(err)
	}
	if got := decodedNodeFootprint(extension); got != (decodedFootprint{nodes: 4, bytes: 8, depth: 2, complete: true}) {
		t.Fatalf("complete extension footprint = %+v", got)
	}
}

func TestEncodingCountsDistinctPendingChildren(t *testing.T) {
	first, err := newLeaf(nil, bytes.Repeat([]byte{1}, 64))
	if err != nil {
		t.Fatal(err)
	}
	second, err := newLeaf(nil, bytes.Repeat([]byte{2}, 64))
	if err != nil {
		t.Fatal(err)
	}
	var children [16]node
	children[0], children[1] = first, second
	branch, err := newBranch(children, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name              string
		encoding, pending int
		limit             bool
	}{
		{"exact", 3, 3, false}, {"encoding short", 2, 3, true}, {"pending short", 3, 2, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			limits := DefaultLimits()
			limits.MaxPendingNodes = tc.pending
			budget := &workBudget{hashesLeft: limits.MaxHashOperations}
			encoded, pending, err := encodeNodeBounded(context.Background(), branch, tc.encoding, budget, limits)
			if errors.Is(err, ErrResourceLimit) != tc.limit || (err != nil && !tc.limit) {
				t.Fatalf("encoding = %v, want limit %t", err, tc.limit)
			}
			if tc.limit {
				if encoded != nil || pending != nil {
					t.Fatal("failed encoding returned partial result")
				}
			} else if len(encoded) == 0 || len(pending) != 2 {
				t.Fatal("successful encoding lost distinct children")
			}
		})
	}
}
