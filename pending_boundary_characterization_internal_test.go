package mpt

import (
	"bytes"
	"context"
	"errors"
	"os"
	"testing"
)

func TestHostedPendingEncodingExactChildByteBoundary(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("canonical encoding boundary executes only in hosted CI")
	}
	first, err := newLeaf([]byte{0}, bytes.Repeat([]byte{1}, 29))
	if err != nil {
		t.Fatal(err)
	}
	second, err := newLeaf([]byte{0}, bytes.Repeat([]byte{2}, 30))
	if err != nil {
		t.Fatal(err)
	}
	var children [16]node
	children[0], children[1] = first, second
	branch, err := newBranch(children, nil)
	if err != nil {
		t.Fatal(err)
	}
	// One compact-path byte, a short-string prefix plus 29/30 value bytes,
	// and a short-list prefix produce canonical 32/33-byte child encodings.
	for _, tc := range []struct {
		name  string
		bytes int
		limit bool
	}{
		{"exact", 65, false},
		{"one byte short", 64, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			limits := DefaultLimits()
			limits.MaxPendingNodes, limits.MaxPendingBytes = 3, tc.bytes
			budget := &workBudget{hashesLeft: limits.MaxHashOperations}
			encoded, pending, err := encodeNodeBounded(context.Background(), branch, 3, budget, limits)
			if tc.limit {
				if !errors.Is(err, ErrResourceLimit) || encoded != nil || pending != nil {
					t.Fatalf("short byte budget returned encoding %t, pending %t, error %v", encoded != nil, pending != nil, err)
				}
				return
			}
			if err != nil || len(encoded) == 0 || len(pending) != 2 {
				t.Fatalf("exact byte budget returned encoding %t, %d pending children, error %v", len(encoded) != 0, len(pending), err)
			}
			var total int
			counts := map[int]int{}
			for _, child := range pending {
				total += len(child)
				counts[len(child)]++
			}
			if total != 65 || counts[32] != 1 || counts[33] != 1 {
				t.Fatalf("persisted child bytes = %d, sizes = %v, want 65 and one each of 32/33", total, counts)
			}
		})
	}
}

func TestMaterializedPendingCompactionTraversalDepthBoundary(t *testing.T) {
	leaf, err := newLeaf(nil, []byte{1})
	if err != nil {
		t.Fatal(err)
	}
	var children [16]node
	children[0], children[1] = leaf, leaf
	inner, err := newBranch(children, nil)
	if err != nil {
		t.Fatal(err)
	}
	sibling, err := newLeaf([]byte{0}, []byte{1})
	if err != nil {
		t.Fatal(err)
	}
	children[0], children[1] = inner, sibling
	outer, err := newBranch(children, nil)
	if err != nil {
		t.Fatal(err)
	}
	extension, err := newExtension([]byte{0}, inner)
	if err != nil {
		t.Fatal(err)
	}
	for _, graph := range []struct {
		name string
		root node
	}{
		{"nested branch", outer},
		{"extension", extension},
	} {
		for _, depth := range []int{1, 2} {
			name := graph.name + "/exact"
			if depth == 1 {
				name = graph.name + "/one short"
			}
			t.Run(name, func(t *testing.T) {
				limits := DefaultLimits()
				limits.MaxTraversalDepth = depth
				retained, err := retainReferencedPending(context.Background(), graph.root, Root{}, nil, limits)
				if depth == 1 {
					if !errors.Is(err, ErrResourceLimit) || retained != nil {
						t.Fatalf("short depth returned retained %v, error %v", retained, err)
					}
					return
				}
				if err != nil || retained == nil || len(retained) != 0 {
					t.Fatalf("exact depth returned retained %v, error %v", retained, err)
				}
			})
		}
	}
}
