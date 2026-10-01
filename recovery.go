package mpt

import (
	"bytes"
	"context"
	"fmt"
)

// RecoverNode copies, hash-checks, and canonically decodes one node retrieved
// after a MissingNodeError, then returns a new raw snapshot that consults the
// bounded recovery overlay before its backing reader. The receiver is
// unchanged. Commit atomically persists recovered nodes to the source store.
func (trie RawTrie) RecoverNode(
	ctx context.Context,
	hash Root,
	encoded []byte,
) (RawTrie, error) {
	snapshot, err := recoverSnapshot(ctx, trie.snapshot, hash, encoded)
	if err != nil {
		return RawTrie{}, err
	}
	return RawTrie{snapshot: snapshot}, nil
}

// RecoverNode copies, hash-checks, and canonically decodes one node retrieved
// after a MissingNodeError, then returns a new secure snapshot that consults
// the bounded recovery overlay before its backing reader. The receiver is
// unchanged. Commit atomically persists recovered nodes to the source store.
func (trie SecureTrie) RecoverNode(
	ctx context.Context,
	hash Root,
	encoded []byte,
) (SecureTrie, error) {
	snapshot, err := recoverSnapshot(ctx, trie.snapshot, hash, encoded)
	if err != nil {
		return SecureTrie{}, err
	}
	return SecureTrie{snapshot: snapshot}, nil
}

func recoverSnapshot(
	ctx context.Context,
	snapshot *trieSnapshot,
	hash Root,
	encoded []byte,
) (*trieSnapshot, error) {
	if snapshot == nil {
		return nil, ErrUninitialized
	}
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if !validStore(snapshot.reader) {
		return nil, fmt.Errorf("%w: recovery requires a backing reader", ErrInvalidStore)
	}
	if existing, exists := lookupSnapshotPending(snapshot, hash); exists {
		if bytes.Equal(existing, encoded) {
			return snapshot, nil
		}
		return nil, &CorruptNodeError{
			Hash:  hash,
			Cause: fmt.Errorf("%w: conflicting recovered bytes", ErrCorruptNode),
		}
	}
	if snapshot.recoveryNodes == snapshot.limits.MaxRecoveryNodes {
		return nil, fmt.Errorf("%w: recovery node bound exceeded", ErrResourceLimit)
	}
	if len(encoded) > snapshot.limits.MaxRecoveryBytes-snapshot.recoveryBytes {
		return nil, fmt.Errorf("%w: recovery byte bound exceeded", ErrResourceLimit)
	}
	stats, compact, err := admitPending(ctx, snapshot, map[Root][]byte{hash: encoded}, nil, snapshot.limits)
	if err != nil {
		return nil, err
	}

	owned := append([]byte(nil), encoded...)
	budget := workBudget{hashesLeft: snapshot.limits.MaxHashOperations}
	actual, err := budget.hash(owned)
	if err != nil {
		return nil, err
	}
	if actual != hash {
		return nil, &CorruptNodeError{
			Hash:  hash,
			Cause: fmt.Errorf("%w: recovered node hash mismatch", ErrCorruptNode),
		}
	}
	decoded, err := decodeNode(owned)
	if err != nil || decoded == nil {
		return nil, fmt.Errorf("%w: invalid recovered node", ErrMalformedNode)
	}
	if err := checkContext(ctx); err != nil {
		return nil, err
	}

	recoveredNodes := make(map[Root][]byte)
	if err := mergePersisted(ctx, recoveredNodes, snapshot.recovered); err != nil {
		return nil, err
	}
	recoveredNodes[hash] = owned
	recovered := *snapshot
	recovered.pending = map[Root][]byte{hash: owned}
	recovered.parent = snapshotPendingLayer(snapshot)
	recovered.removed = nil
	recovered.pendingStats = stats
	if recovered.parent != nil &&
		(compact || recovered.parent.depth >= maximumPendingLayerDepth) {
		compacted, err := materializePendingLayerContext(ctx, recovered.parent)
		if err != nil {
			return nil, err
		}
		compacted[hash] = owned
		recovered.pending = compacted
		recovered.parent = nil
	}
	recovered.recovered = recoveredNodes
	recovered.recoveryNodes++
	recovered.recoveryBytes += len(owned)
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	return &recovered, nil
}

func inheritRecovery(ctx context.Context, next, previous *trieSnapshot) (*trieSnapshot, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if next.root == nil || len(previous.recovered) == 0 {
		return next, nil
	}
	recovered := make(map[Root][]byte)
	recoveryBytes := 0
	for hash := range previous.recovered {
		if err := checkContext(ctx); err != nil {
			return nil, err
		}
		if encoded, reachable := lookupSnapshotPending(next, hash); reachable {
			recovered[hash] = encoded
			recoveryBytes += len(encoded)
		}
	}
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	next.recovered = recovered
	next.recoveryNodes = len(recovered)
	next.recoveryBytes = recoveryBytes
	return next, nil
}
