package mpt

import (
	"context"
	"errors"
	"testing"
)

func TestDecodedFootprintAdmissionFailsClosedOnOverflow(t *testing.T) {
	// Actual overflow cannot be reproduced with real buffers safely. Inject
	// only the immutable child summary at this private arithmetic boundary.
	maximum := int(^uint(0) >> 1)
	leaf, err := newLeaf(nil, []byte("value"))
	if err != nil {
		t.Fatal(err)
	}
	leaf.footprint = decodedFootprint{nodes: maximum, bytes: maximum, depth: maximum, complete: true}
	var children [16]node
	children[0] = leaf
	children[1] = leaf
	branch, err := newBranch(children, []byte("branch value"))
	if err != nil {
		t.Fatal(err)
	}
	if branch.footprint.nodes != maximum || branch.footprint.bytes != maximum || branch.footprint.depth != maximum {
		t.Fatal("overflow did not saturate")
	}
	limits := DefaultLimits()
	limits.MaxPendingNodes = maximum
	limits.MaxPendingBytes = maximum
	limits.MaxEncodingNodes = maximum
	limits.MaxTraversalDepth = maximum
	budget := &workBudget{hashesLeft: limits.MaxHashOperations}
	if _, err := finishSnapshotWithPending(context.Background(), branch, branch, limits, EmptyRoot(), nil, nil, nil, true, budget); !errors.Is(err, ErrResourceLimit) {
		t.Fatalf("saturated summary admitted by maximum policy: %v", err)
	}
}

func TestDecodedFootprintRejectsUnresolvedMaterializedGraph(t *testing.T) {
	root, err := newExtension([]byte{1}, hashNode(Root{1}))
	if err != nil {
		t.Fatal(err)
	}
	limits := DefaultLimits()
	budget := &workBudget{hashesLeft: limits.MaxHashOperations}
	if _, err := finishSnapshotWithPending(context.Background(), root, root, limits, EmptyRoot(), nil, nil, nil, true, budget); !errors.Is(err, ErrMalformedNode) {
		t.Fatalf("unresolved materialized graph admitted: %v", err)
	}
}
