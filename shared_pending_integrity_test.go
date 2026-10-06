package mpt_test

import (
	"bytes"
	"context"
	"testing"

	mpt "github.com/faustbrian/go-merkle-patricia-trie/v2"
	"github.com/faustbrian/go-merkle-patricia-trie/v2/memory"
)

func TestSharedPendingLeafSurvivesSiblingReplacement(t *testing.T) {
	ctx := context.Background()
	limits := mpt.DefaultLimits()
	trie, err := mpt.NewRawTrie(limits)
	if err != nil {
		t.Fatal(err)
	}
	value := bytes.Repeat([]byte{'v'}, 64)
	for _, key := range [][]byte{{0}, {1}} {
		trie, err = trie.Update(ctx, key, value)
		if err != nil {
			t.Fatal(err)
		}
	}
	trie, err = trie.Update(ctx, []byte{0}, bytes.Repeat([]byte{'w'}, 64))
	if err != nil {
		t.Fatal(err)
	}
	root, err := trie.Root()
	if err != nil {
		t.Fatal(err)
	}
	if got, err := trie.Get(ctx, []byte{1}); err != nil || !bytes.Equal(got, value) {
		t.Fatalf("cached sibling changed: %v", err)
	}
	if proof, err := trie.Prove(ctx, []byte{1}); err != nil {
		t.Errorf("sibling proof lost shared leaf: %v", err)
	} else if err := mpt.VerifyRawMembership(ctx, root, []byte{1}, value, proof, limits); err != nil {
		t.Errorf("sibling proof did not verify: %v", err)
	}
	if _, err := trie.ProveMany(ctx, [][]byte{{0}, {1}}); err != nil {
		t.Errorf("multi-proof lost shared leaf: %v", err)
	}
	store := memory.New()
	if _, err := trie.Commit(ctx, store); err != nil {
		t.Fatalf("commit shared sibling: %v", err)
	}
	loaded, err := mpt.LoadRawTrie(root, store, limits)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := loaded.Get(ctx, []byte{1}); err != nil || !bytes.Equal(got, value) {
		t.Errorf("committed sibling did not reload: %v", err)
	}
}
