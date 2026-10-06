package mpt_test

import (
	"bytes"
	"context"
	"os"
	"testing"

	mpt "github.com/faustbrian/go-merkle-patricia-trie/v2"
)

func TestHostedCanonicalRootEncodedByteLimitIsInclusive(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("finite canonical encoded-root boundary executes only in hosted CI")
	}
	ctx := context.Background()
	const valueBytes = (16 << 20) - 56
	limits := mpt.DefaultLimits()
	limits.MaxValueBytes = valueBytes
	trie, err := mpt.NewRawTrie(limits)
	if err != nil {
		t.Fatal(err)
	}
	value := bytes.Repeat([]byte{1}, valueBytes)
	trie, err = trie.Update(ctx, nil, value)
	if err != nil {
		t.Fatal(err)
	}
	trie, err = trie.Update(ctx, []byte{0x00}, bytes.Repeat([]byte{2}, 29))
	if err != nil {
		t.Fatal(err)
	}
	store := newTestNodeStore()
	committed, err := trie.Commit(ctx, store)
	if err != nil {
		t.Fatal(err)
	}
	root, err := committed.Root()
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := store.GetNode(ctx, root)
	if err != nil {
		t.Fatal(err)
	}
	// A 33-byte hashed child reference, 15 empty references, four bytes
	// of value prefix, and four of list prefix add exactly 56 bytes.
	if len(encoded) != 16<<20 {
		t.Fatalf("canonical root encoded bytes = %d, want 16777216", len(encoded))
	}
	// Reload prevents the committed materialized cache from bypassing the
	// backing-reader encoded-byte check before copy, hash, and decode.
	loaded, err := mpt.LoadRawTrie(root, store, limits)
	if err != nil {
		t.Fatal(err)
	}
	got, err := loaded.Get(ctx, nil)
	if err != nil || !bytes.Equal(got, value) {
		t.Fatalf("reloaded exact-cap value bytes = %d, error %v, want %d identical bytes", len(got), err, valueBytes)
	}
}
