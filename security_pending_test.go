package mpt_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	mpt "github.com/faustbrian/go-merkle-patricia-trie/v2"
	"github.com/faustbrian/go-merkle-patricia-trie/v2/memory"
)

type boundedSnapshot[T any] interface {
	Root() (mpt.Root, error)
	Get(context.Context, []byte) ([]byte, error)
	Update(context.Context, []byte, []byte) (T, error)
	Commit(context.Context, mpt.NodeStore) (T, error)
}

func boundedLoadedSnapshot[T boundedSnapshot[T]](limits mpt.Limits, create func(mpt.Limits) (T, error), load func(mpt.Root, mpt.NodeReader, mpt.Limits) (T, error)) (T, error) {
	var zero T
	ctx := context.Background()
	trie, err := create(mpt.DefaultLimits())
	if err != nil {
		return zero, err
	}
	trie, err = trie.Update(ctx, []byte{'a'}, bytes.Repeat([]byte{'v'}, 64))
	if err != nil {
		return zero, err
	}
	store := memory.New()
	committed, err := trie.Commit(ctx, store)
	if err != nil {
		return zero, err
	}
	root, err := committed.Root()
	if err != nil {
		return zero, err
	}
	return load(root, store, limits)
}

func exercisePendingAdmission[T boundedSnapshot[T]](t *testing.T, create func(mpt.Limits) (T, error), limits mpt.Limits) {
	t.Helper()
	ctx := context.Background()
	trie, err := create(limits)
	if err != nil {
		t.Fatal(err)
	}
	value := bytes.Repeat([]byte{'v'}, 64)
	trie, err = trie.Update(ctx, []byte{'a'}, value)
	if err != nil {
		t.Fatal(err)
	}
	root, err := trie.Root()
	if err != nil {
		t.Fatal(err)
	}
	candidate, err := trie.Update(ctx, []byte{'b'}, value)
	if !errors.Is(err, mpt.ErrResourceLimit) {
		newRoot, rootErr := candidate.Root()
		t.Errorf("over-limit snapshot accepted: error=%v, new root=%t", err, rootErr == nil && newRoot != root)
	}
	unchanged, err := trie.Root()
	if err != nil || unchanged != root {
		t.Fatal("rejected update changed prior commitment")
	}
	value[0] ^= 0xff
	got, err := trie.Get(ctx, []byte{'a'})
	if err != nil || len(got) != 64 || got[0] != 'v' {
		t.Fatal("rejected update or caller mutation changed prior value")
	}
	if _, err := trie.Get(ctx, []byte{'b'}); !errors.Is(err, mpt.ErrAbsentKey) {
		t.Fatal("rejected update changed prior key set")
	}
}

func TestPendingAdmissionBoundsRawAndSecureSnapshots(t *testing.T) {
	for _, dimension := range []string{"nodes", "bytes"} {
		limits := mpt.DefaultLimits()
		if dimension == "nodes" {
			limits.MaxPendingNodes = 1
		} else {
			limits.MaxPendingBytes = 160
		}
		t.Run("raw "+dimension, func(t *testing.T) { exercisePendingAdmission(t, mpt.NewRawTrie, limits) })
		t.Run("secure "+dimension, func(t *testing.T) { exercisePendingAdmission(t, mpt.NewSecureTrie, limits) })
		t.Run("loaded raw "+dimension, func(t *testing.T) {
			exercisePendingAdmission(t, func(limits mpt.Limits) (mpt.RawTrie, error) {
				return boundedLoadedSnapshot(limits, mpt.NewRawTrie, mpt.LoadRawTrie)
			}, limits)
		})
		t.Run("loaded secure "+dimension, func(t *testing.T) {
			exercisePendingAdmission(t, func(limits mpt.Limits) (mpt.SecureTrie, error) {
				return boundedLoadedSnapshot(limits, mpt.NewSecureTrie, mpt.LoadSecureTrie)
			}, limits)
		})
	}
}

func TestMaterializedCacheAdmissionCountsDistinctBuffersWithSharedHashes(t *testing.T) {
	ctx := context.Background()
	limits := mpt.DefaultLimits()
	limits.MaxPendingBytes = 3 << 19
	trie, err := mpt.NewRawTrie(limits)
	if err != nil {
		t.Fatal(err)
	}
	value := bytes.Repeat([]byte{'v'}, 1<<20)
	trie, err = trie.Update(ctx, []byte{0x10}, value)
	if err != nil {
		t.Fatal(err)
	}
	// After splitting the last nibble both children have the same encoded
	// leaf hash, but the materialized fast path owns two value buffers.
	if _, err := trie.Update(ctx, []byte{0x11}, value); !errors.Is(err, mpt.ErrResourceLimit) {
		t.Error("hash deduplication hid materialized cache growth")
	}
	store := memory.New()
	committed, err := trie.Commit(ctx, store)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := committed.Update(ctx, []byte{0x11}, value); !errors.Is(err, mpt.ErrResourceLimit) {
		t.Error("commit reset materialized cache admission")
	}
}

func TestPendingPressureCompactsReplacementAndDeletion(t *testing.T) {
	ctx := context.Background()
	limits := mpt.DefaultLimits()
	limits.MaxPendingNodes = 1
	limits.MaxPendingBytes = 160
	trie, err := mpt.NewRawTrie(limits)
	if err != nil {
		t.Fatal(err)
	}
	for i := range 40 {
		value := bytes.Repeat([]byte{byte(i + 1)}, 64)
		trie, err = trie.Update(ctx, []byte{'a'}, value)
		if err != nil {
			t.Fatalf("replacement did not reclaim pending ancestry: %v", err)
		}
	}
	store := memory.New()
	committed, err := trie.Commit(ctx, store)
	if err != nil {
		t.Fatal(err)
	}
	root := store.Root()
	encoded, err := store.GetNode(ctx, root)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := mpt.LoadRawTrie(root, store, limits)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err = loaded.RecoverNode(ctx, root, encoded)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err = loaded.RecoverNode(ctx, root, encoded)
	if err != nil {
		t.Fatal("duplicate recovery consumed pending admission")
	}
	if _, err := loaded.Rebuild(ctx); err != nil {
		t.Fatal(err)
	}
	deleted, err := committed.Delete(ctx, []byte{'a'})
	if err != nil {
		t.Fatal(err)
	}
	deleted, err = deleted.Update(ctx, []byte{'b'}, bytes.Repeat([]byte{'v'}, 64))
	if err != nil {
		t.Fatal("deletion did not reclaim pending admission")
	}
	if _, err := deleted.Commit(ctx, store); err != nil {
		t.Fatal(err)
	}
}
