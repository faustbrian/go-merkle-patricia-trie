package mpt_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	mpt "github.com/faustbrian/go-merkle-patricia-trie/v2"
	"github.com/faustbrian/go-merkle-patricia-trie/v2/memory"
)

func TestRangeItemAdmissionCountsSharedLeafCopies(t *testing.T) {
	for _, keyBytes := range []int{1, mpt.RootBytes} {
		name := "raw"
		if keyBytes == mpt.RootBytes {
			name = "secure hashed"
		}
		t.Run(name, func(t *testing.T) { exerciseSharedLeafRangeAdmission(t, keyBytes) })
	}
}

func TestRangeWitnessBudgetRemainsIndependentOfItemBudget(t *testing.T) {
	proof, err := mpt.RangeProofFromNodes([][]byte{make([]byte, 128)}, mpt.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	limits := mpt.DefaultLimits()
	limits.MaxProofBytes = 64
	if err := mpt.VerifyRawRange(context.Background(), mpt.Root{1}, nil, nil, nil, proof, limits); !errors.Is(err, mpt.ErrResourceLimit) {
		t.Fatalf("empty item list bypassed encoded witness byte budget: %v", err)
	}
}

func exerciseSharedLeafRangeAdmission(t *testing.T, keyBytes int) {
	t.Helper()
	ctx := context.Background()
	trie, err := mpt.NewRawTrie(mpt.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	value := bytes.Repeat([]byte{'v'}, 128)
	// Sixteen distinct positions share one canonical leaf hash after their
	// branch consumes the final nibble. Witness deduplication does not bound
	// the sixteen independently owned output values.
	for key := byte(0); key < 16; key++ {
		path := make([]byte, keyBytes)
		path[keyBytes-1] = key
		trie, err = trie.Update(ctx, path, value)
		if err != nil {
			t.Fatal(err)
		}
	}
	proof, items, err := trie.ProveRange(ctx, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	proofBytes := 0
	for _, encoded := range proof.Nodes() {
		proofBytes += len(encoded)
	}
	if len(items) != 16 || proofBytes > 1024 {
		t.Fatalf("fixture did not isolate shared-leaf output: items=%d witness bytes=%d", len(items), proofBytes)
	}
	store := memory.New()
	if _, err := trie.Commit(ctx, store); err != nil {
		t.Fatal(err)
	}
	limits := mpt.DefaultLimits()
	limits.MaxProofBytes = 1024
	loaded, err := mpt.LoadRawTrie(store.Root(), store, limits)
	if err != nil {
		t.Fatal(err)
	}
	prove := loaded.ProveRange
	verify := mpt.VerifyRawRange
	if keyBytes == mpt.RootBytes {
		secure, err := mpt.LoadSecureTrie(store.Root(), store, limits)
		if err != nil {
			t.Fatal(err)
		}
		prove = secure.ProveHashedRange
		verify = mpt.VerifySecureHashedRange
	}
	if _, _, err := prove(ctx, nil, nil); !errors.Is(err, mpt.ErrResourceLimit) {
		t.Errorf("shared witnesses hid cumulative range output: %v", err)
	}
	if err := verify(ctx, store.Root(), nil, nil, items, proof, limits); !errors.Is(err, mpt.ErrResourceLimit) {
		t.Errorf("shared witnesses hid cumulative range claim bytes: %v", err)
	}
	if got, err := loaded.Get(ctx, make([]byte, keyBytes)); err != nil || !bytes.Equal(got, value) {
		t.Fatal("range rejection changed loaded snapshot")
	}
	limits.MaxProofBytes = 16 * (keyBytes + len(value))
	if err := verify(ctx, store.Root(), nil, nil, items, proof, limits); err != nil {
		t.Fatalf("exact item-byte boundary rejected: %v", err)
	}
	fitting, err := mpt.LoadRawTrie(store.Root(), store, limits)
	if err != nil {
		t.Fatal(err)
	}
	prove = fitting.ProveRange
	if keyBytes == mpt.RootBytes {
		secure, err := mpt.LoadSecureTrie(store.Root(), store, limits)
		if err != nil {
			t.Fatal(err)
		}
		prove = secure.ProveHashedRange
	}
	if _, fittedItems, err := prove(ctx, nil, nil); err != nil || len(fittedItems) != 16 {
		t.Fatalf("exact generation boundary rejected: items=%d error=%v", len(fittedItems), err)
	}
}
