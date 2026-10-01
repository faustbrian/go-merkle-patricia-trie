package mpt_test

import (
	"bytes"
	"context"
	"runtime"
	"testing"

	mpt "github.com/faustbrian/go-merkle-patricia-trie/v2"
)

type securityCommitCapture struct {
	commit mpt.StoreCommit
}

func (store *securityCommitCapture) GetNode(_ context.Context, hash mpt.Root) ([]byte, error) {
	for _, node := range store.commit.Nodes() {
		if node.Hash() == hash {
			return node.Encoded(), nil
		}
	}
	return nil, mpt.ErrMissingNode
}

func (store *securityCommitCapture) CommitTrie(_ context.Context, commit mpt.StoreCommit) error {
	store.commit = commit
	return nil
}

func TestStoreCommitMetadataDoesNotCopyEncodedPayloads(t *testing.T) {
	ctx := context.Background()
	value := bytes.Repeat([]byte{'v'}, 1<<20)
	trie, err := mpt.NewRawTrie(mpt.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	trie, err = trie.Update(ctx, []byte("key"), value)
	if err != nil {
		t.Fatal(err)
	}
	store := &securityCommitCapture{}
	if _, err := trie.Commit(ctx, store); err != nil {
		t.Fatal(err)
	}
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	nodes := store.commit.Nodes()
	runtime.ReadMemStats(&after)
	if len(nodes) != 1 {
		t.Fatal("expected one canonical root node")
	}
	if after.TotalAlloc-before.TotalAlloc > 64<<10 {
		t.Fatalf("metadata enumeration copied payload: allocated %d bytes", after.TotalAlloc-before.TotalAlloc)
	}
	runtime.KeepAlive(nodes)
}

func TestStoreCommitHandlesPreservePublicOwnershipAndReuse(t *testing.T) {
	ctx := context.Background()
	key := []byte("key")
	value := []byte("private caller-owned value")
	trie, err := mpt.NewRawTrie(mpt.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	trie, err = trie.Update(ctx, key, value)
	if err != nil {
		t.Fatal(err)
	}
	key[0], value[0] = 'x', 'x'
	first := &securityCommitCapture{}
	committed, err := trie.Commit(ctx, first)
	if err != nil {
		t.Fatal(err)
	}
	nodes := first.commit.Nodes()
	if len(nodes) != 1 {
		t.Fatal("expected one canonical root node")
	}
	hash := nodes[0].Hash()
	expected := nodes[0].Encoded()
	encoded := nodes[0].Encoded()
	encoded[0] ^= 0xff
	nodes[0] = mpt.StoredNode{}
	if !bytes.Equal(first.commit.Nodes()[0].Encoded(), expected) {
		t.Fatal("public node handles exposed mutable commit storage")
	}
	read, err := first.GetNode(ctx, hash)
	if err != nil {
		t.Fatal(err)
	}
	read[0] ^= 0xff
	second := &securityCommitCapture{}
	if err := second.CommitTrie(ctx, first.commit); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(second.commit.Nodes()[0].Encoded(), expected) {
		t.Fatal("cross-store reuse observed a mutated commit")
	}
	got, err := committed.Get(ctx, []byte("key"))
	if err != nil || string(got) != "private caller-owned value" {
		t.Fatal("caller mutation changed committed snapshot")
	}
	got[0] ^= 0xff
	got, err = committed.Get(ctx, []byte("key"))
	if err != nil || string(got) != "private caller-owned value" {
		t.Fatal("returned value alias changed committed snapshot")
	}
}
