package mpt_test

import (
	"context"
	"errors"
	"runtime"
	"strings"
	"testing"

	mpt "github.com/faustbrian/go-merkle-patricia-trie/v2"
)

// The reader is deliberately a transport boundary: its returned storage bytes
// and diagnostics are untrusted, unlike the caller's root and resource policy.
type hostileStorageReader struct {
	encoded []byte
	failure error
	cancel  context.CancelFunc
}

func (reader hostileStorageReader) GetNode(context.Context, mpt.Root) ([]byte, error) {
	if reader.cancel != nil {
		reader.cancel()
	}
	return reader.encoded, reader.failure
}

type confidentialStorageFailure struct{}

func (*confidentialStorageFailure) Error() string {
	return "synthetic-private-backend-diagnostic"
}

func TestMultiProofInputAdmissionPrecedesKeyCopies(t *testing.T) {
	limits := mpt.DefaultLimits()
	limits.MaxKeyBytes = 4
	raw, err := mpt.NewRawTrie(limits)
	if err != nil {
		t.Fatal(err)
	}
	secure, err := mpt.NewSecureTrie(limits)
	if err != nil {
		t.Fatal(err)
	}
	// Caller transport owns a modest oversized payload. Neither rejected
	// key bytes nor rejected key counts authorize a library-owned copy.
	oversized := make([]byte, 1<<20)
	tooMany := make([][]byte, limits.MaxProofKeys+1)
	for index := range tooMany {
		tooMany[index] = oversized
	}
	for name, prove := range map[string]func(context.Context, [][]byte) (mpt.MultiProof, error){
		"raw": raw.ProveMany, "secure": secure.ProveMany,
	} {
		for _, input := range []struct {
			name string
			keys [][]byte
			want error
		}{
			{"key bytes", [][]byte{[]byte("ok"), oversized}, mpt.ErrInvalidKey},
			{"key count", tooMany, mpt.ErrInvalidProofClaim},
		} {
			t.Run(name+" "+input.name, func(t *testing.T) {
				var before, after runtime.MemStats
				runtime.ReadMemStats(&before)
				_, err := prove(context.Background(), input.keys)
				runtime.ReadMemStats(&after)
				if !errors.Is(err, input.want) {
					t.Fatalf("input admission returned %v, want %v", err, input.want)
				}
				if allocated := after.TotalAlloc - before.TotalAlloc; allocated > 64<<10 {
					t.Fatalf("rejected proof input copied payload: allocated %d bytes", allocated)
				}
			})
		}
	}
}

func TestStorageResponsesAreBoundedBeforeOwnedWork(t *testing.T) {
	// The existing node decoder admits at most 16 MiB. The hostile transport
	// owns this one bounded test payload; operations must not duplicate it.
	reader := hostileStorageReader{encoded: make([]byte, (16<<20)+1)}
	root := mpt.Root{1}
	raw, err := mpt.LoadRawTrie(root, reader, mpt.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	secure, err := mpt.LoadSecureTrie(root, reader, mpt.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	for name, operation := range map[string]func() error{
		"raw lookup": func() error {
			_, err := raw.Get(context.Background(), []byte("key"))
			return err
		},
		"secure lookup": func() error {
			_, err := secure.Get(context.Background(), []byte("key"))
			return err
		},
		"proof generation": func() error {
			_, err := raw.Prove(context.Background(), []byte("key"))
			return err
		},
		"reachability": func() error {
			limits := mpt.DefaultReachabilityLimits()
			limits.MaxBytes = 32
			_, err := mpt.CollectReachableNodes(
				context.Background(), []mpt.Root{root}, reader, limits,
			)
			return err
		},
		"reachability canonical ceiling": func() error {
			_, err := mpt.CollectReachableNodes(context.Background(), []mpt.Root{root}, reader, mpt.DefaultReachabilityLimits())
			return err
		},
	} {
		t.Run(name, func(t *testing.T) {
			var before, after runtime.MemStats
			runtime.ReadMemStats(&before)
			err := operation()
			runtime.ReadMemStats(&after)
			if !errors.Is(err, mpt.ErrResourceLimit) {
				t.Errorf("oversized storage response did not return resource limit")
			}
			// A generous allowance for incidental runtime work still excludes
			// a copy of the 16 MiB response. No elapsed-time assertion is used.
			if after.TotalAlloc-before.TotalAlloc > 1<<20 {
				t.Errorf("rejected response allocated %d bytes", after.TotalAlloc-before.TotalAlloc)
			}
		})
	}
}

func TestStorageResponseHonorsCancellationBeforeOwnedWork(t *testing.T) {
	for _, operation := range []string{"lookup", "reachability"} {
		t.Run(operation, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			reader := hostileStorageReader{encoded: make([]byte, 1<<20), cancel: cancel}
			root := mpt.Root{1}
			var err error
			if operation == "lookup" {
				var trie mpt.RawTrie
				trie, err = mpt.LoadRawTrie(root, reader, mpt.DefaultLimits())
				if err != nil {
					t.Fatal(err)
				}
				_, err = trie.Get(ctx, nil)
			} else {
				_, err = mpt.CollectReachableNodes(ctx, []mpt.Root{root}, reader, mpt.DefaultReachabilityLimits())
			}
			if !errors.Is(err, mpt.ErrCanceled) || !errors.Is(err, context.Canceled) {
				t.Fatal("reader-boundary cancellation did not stop owned work")
			}
		})
	}
}

func TestStorageFailuresRedactDiagnosticsAndPreserveCauses(t *testing.T) {
	failure := &confidentialStorageFailure{}
	reader := hostileStorageReader{failure: failure}
	root := mpt.Root{1}
	raw, err := mpt.LoadRawTrie(root, reader, mpt.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	secure, err := mpt.LoadSecureTrie(root, reader, mpt.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	for name, operation := range map[string]func() error{
		"raw lookup": func() error {
			_, err := raw.Get(context.Background(), nil)
			return err
		},
		"secure lookup": func() error {
			_, err := secure.Get(context.Background(), nil)
			return err
		},
		"rebuild": func() error {
			_, err := raw.Rebuild(context.Background())
			return err
		},
		"reachability": func() error {
			_, err := mpt.CollectReachableNodes(
				context.Background(), []mpt.Root{root}, reader,
				mpt.DefaultReachabilityLimits(),
			)
			return err
		},
	} {
		t.Run(name, func(t *testing.T) {
			assertRedactedStorageFailure(t, operation(), mpt.ErrStorageRead, failure)
		})
	}

	t.Run("commit", func(t *testing.T) {
		trie := mustRawTrie(t, map[string]string{"key": "value"})
		store := newTestNodeStore()
		store.commitErr = failure
		_, err := trie.Commit(context.Background(), store)
		assertRedactedStorageFailure(t, err, mpt.ErrStorageCommit, failure)
		if store.root != mpt.EmptyRoot() || len(store.nodes) != 0 {
			t.Fatal("failed commit changed published storage")
		}
	})
}

func assertRedactedStorageFailure(t *testing.T, err, category error, failure *confidentialStorageFailure) {
	t.Helper()
	var cause *confidentialStorageFailure
	if !errors.Is(err, category) || !errors.Is(err, failure) ||
		!errors.As(err, &cause) || cause != failure {
		t.Fatal("storage failure lost its classification or inspectable cause")
	}
	if strings.Contains(err.Error(), failure.Error()) {
		t.Fatal("storage failure rendered confidential backend diagnostics")
	}
}
