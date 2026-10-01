package filesystem_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	mpt "github.com/faustbrian/go-merkle-patricia-trie/v2"
	"github.com/faustbrian/go-merkle-patricia-trie/v2/filesystem"
)

func TestFilesystemStorageFailuresRedactPathsAndPreserveCauses(t *testing.T) {
	parent := filepath.Join(t.TempDir(), "synthetic-private-storage-parent")
	if err := os.WriteFile(parent, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(parent, "child")
	_, err := filesystem.Open(context.Background(), path, filesystem.DefaultLimits())
	var cause *os.PathError
	if !errors.Is(err, mpt.ErrStorageRead) || !errors.As(err, &cause) {
		t.Fatal("filesystem failure lost storage classification or path-error cause")
	}
	if !errors.Is(err, cause.Err) {
		t.Fatal("filesystem failure lost underlying operating-system cause")
	}
	if strings.Contains(err.Error(), "synthetic-private-storage-parent") {
		t.Fatal("filesystem failure rendered confidential storage path")
	}
}

func TestFilesystemReadAndCommitDiagnosticsRemainInspectableButRedacted(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "synthetic-private-storage-parent")
	store, err := filesystem.Open(ctx, path, filesystem.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			t.Errorf("close owned store: %v", err)
		}
	}()
	nodes := filepath.Join(path, "nodes")
	if err := os.Remove(nodes); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(nodes, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	_, readErr := store.GetNode(ctx, mpt.Root{1})
	trie, err := mpt.NewRawTrie(mpt.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	trie, err = trie.Update(ctx, []byte("key"), []byte("value"))
	if err != nil {
		t.Fatal(err)
	}
	_, commitErr := trie.Commit(ctx, store)
	for _, failure := range []struct{ err, category error }{
		{readErr, mpt.ErrStorageRead}, {commitErr, mpt.ErrStorageCommit},
	} {
		var cause *os.PathError
		if !errors.Is(failure.err, failure.category) || !errors.As(failure.err, &cause) ||
			!errors.Is(failure.err, cause.Err) {
			t.Fatal("filesystem failure lost its inspectable cause")
		}
		if strings.Contains(failure.err.Error(), "synthetic-private-storage-parent") {
			t.Fatal("filesystem failure rendered confidential storage path")
		}
	}
	if store.Root() != mpt.EmptyRoot() {
		t.Fatal("failed filesystem commit published a new root")
	}
}

type boundedCommitCapture struct{ commit mpt.StoreCommit }

func (store *boundedCommitCapture) GetNode(context.Context, mpt.Root) ([]byte, error) {
	return nil, mpt.ErrMissingNode
}

func (store *boundedCommitCapture) CommitTrie(_ context.Context, commit mpt.StoreCommit) error {
	store.commit = commit
	return nil
}

func TestFilesystemAdmissionRejectsBeforeCopyingEncodedPayload(t *testing.T) {
	ctx := context.Background()
	trie, err := mpt.NewRawTrie(mpt.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	trie, err = trie.Update(ctx, []byte("key"), make([]byte, 1<<20))
	if err != nil {
		t.Fatal(err)
	}
	capture := &boundedCommitCapture{}
	if _, err := trie.Commit(ctx, capture); err != nil {
		t.Fatal(err)
	}
	limits := filesystem.DefaultLimits()
	limits.MaxNodeBytes, limits.MaxCommitBytes = 32, 32
	store, err := filesystem.Open(ctx, t.TempDir(), limits)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			t.Errorf("close owned store: %v", err)
		}
	}()
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	err = store.CommitTrie(ctx, capture.commit)
	runtime.ReadMemStats(&after)
	if !errors.Is(err, mpt.ErrResourceLimit) || store.Root() != mpt.EmptyRoot() {
		t.Fatal("over-limit adapter admission did not preserve root")
	}
	if after.TotalAlloc-before.TotalAlloc > 64<<10 {
		t.Fatal("rejected adapter admission copied encoded payload")
	}
}

func TestFilesystemCountAdmissionPrecedesHandleCopy(t *testing.T) {
	ctx := context.Background()
	trie, err := mpt.NewRawTrie(mpt.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	for index := range 1024 {
		key := []byte{byte(index >> 8), byte(index)}
		value := make([]byte, 64)
		value[0], value[1] = key[0], key[1]
		trie, err = trie.Update(ctx, key, value)
		if err != nil {
			t.Fatal(err)
		}
	}
	capture := &boundedCommitCapture{}
	if _, err := trie.Commit(ctx, capture); err != nil {
		t.Fatal(err)
	}
	if capture.commit.NodeCount() < 1024 {
		t.Fatal("fixture did not retain the intended unique node population")
	}
	limits := filesystem.DefaultLimits()
	limits.MaxCommitNodes = 1
	store, err := filesystem.Open(ctx, t.TempDir(), limits)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			t.Errorf("close owned store: %v", err)
		}
	}()
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	err = store.CommitTrie(ctx, capture.commit)
	runtime.ReadMemStats(&after)
	if !errors.Is(err, mpt.ErrResourceLimit) || store.Root() != mpt.EmptyRoot() {
		t.Fatal("over-count admission did not reject atomically")
	}
	if allocated := after.TotalAlloc - before.TotalAlloc; allocated > 16<<10 {
		t.Errorf("over-count admission copied rejected handles: allocated %d bytes", allocated)
	}
}

type commitCancellationContext struct{ calls, at int }

func (*commitCancellationContext) Deadline() (time.Time, bool) { return time.Time{}, false }
func (*commitCancellationContext) Done() <-chan struct{}       { return nil }
func (*commitCancellationContext) Value(any) any               { return nil }
func (ctx *commitCancellationContext) Err() error {
	ctx.calls++
	if ctx.calls >= ctx.at {
		return context.Canceled
	}
	return nil
}

func TestFilesystemCancellationPrecedesValidationPayloadCopy(t *testing.T) {
	ctx := context.Background()
	trie, err := mpt.NewRawTrie(mpt.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	trie, err = trie.Update(ctx, []byte("key"), make([]byte, 1<<20))
	if err != nil {
		t.Fatal(err)
	}
	capture := &boundedCommitCapture{}
	if _, err := trie.Commit(ctx, capture); err != nil {
		t.Fatal(err)
	}
	store, err := filesystem.Open(ctx, t.TempDir(), filesystem.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			t.Errorf("close owned store: %v", err)
		}
	}()
	// Initial public admission succeeds; the next owned preparation seam
	// must observe cancellation before copying or hashing the payload.
	canceled := &commitCancellationContext{at: 2}
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	err = store.CommitTrie(canceled, capture.commit)
	runtime.ReadMemStats(&after)
	if !errors.Is(err, mpt.ErrCanceled) || !errors.Is(err, context.Canceled) || store.Root() != mpt.EmptyRoot() {
		t.Fatal("canceled adapter admission did not preserve classification and root")
	}
	if allocated := after.TotalAlloc - before.TotalAlloc; allocated > 64<<10 {
		t.Errorf("canceled adapter validation copied payload: allocated %d bytes", allocated)
	}
}

func TestFilesystemPreparationCancellationPreservesRetry(t *testing.T) {
	ctx := context.Background()
	trie, err := mpt.NewRawTrie(mpt.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"first", "second"} {
		trie, err = trie.Update(ctx, []byte(key), []byte("bounded value"))
		if err != nil {
			t.Fatal(err)
		}
	}
	capture := &boundedCommitCapture{}
	if _, err := trie.Commit(ctx, capture); err != nil {
		t.Fatal(err)
	}
	// Admission and validation check four fixed seams and each handle;
	// preparation then checks each handle before reading owned storage.
	// Every cancellation before the write phase must leave a retryable store.
	for at := 1; at <= 4+2*capture.commit.NodeCount(); at++ {
		store, err := filesystem.Open(ctx, t.TempDir(), filesystem.DefaultLimits())
		if err != nil {
			t.Fatal(err)
		}
		err = store.CommitTrie(&commitCancellationContext{at: at}, capture.commit)
		if !errors.Is(err, mpt.ErrCanceled) || !errors.Is(err, context.Canceled) || store.Root() != mpt.EmptyRoot() {
			t.Errorf("preparation cancellation at seam %d did not preserve root: %v", at, err)
		}
		if err := store.CommitTrie(ctx, capture.commit); err != nil || store.Root() != capture.commit.Root() {
			t.Errorf("preparation cancellation at seam %d prevented retry: %v", at, err)
		}
		if err := store.Close(); err != nil {
			t.Errorf("close owned store: %v", err)
		}
	}
}
