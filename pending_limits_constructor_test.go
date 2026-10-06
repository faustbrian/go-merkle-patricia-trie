package mpt_test

import (
	"errors"
	"testing"

	mpt "github.com/faustbrian/go-merkle-patricia-trie/v2"
)

func TestConstructorsRequirePositivePendingLimits(t *testing.T) {
	for _, dimension := range []string{"nodes", "bytes"} {
		for _, bound := range []int{0, 1} {
			limits := mpt.DefaultLimits()
			if dimension == "nodes" {
				limits.MaxPendingNodes = bound
			} else {
				limits.MaxPendingBytes = bound
			}
			for _, constructor := range []struct {
				name string
				new  func(mpt.Limits) (mpt.Root, error, error)
			}{
				{"raw", func(limits mpt.Limits) (mpt.Root, error, error) {
					trie, err := mpt.NewRawTrie(limits)
					root, rootErr := trie.Root()
					return root, err, rootErr
				}},
				{"secure", func(limits mpt.Limits) (mpt.Root, error, error) {
					trie, err := mpt.NewSecureTrie(limits)
					root, rootErr := trie.Root()
					return root, err, rootErr
				}},
			} {
				name := constructor.name + "/" + dimension
				if bound == 0 {
					name += "/zero"
				} else {
					name += "/one"
				}
				t.Run(name, func(t *testing.T) {
					root, err, rootErr := constructor.new(limits)
					if bound == 0 {
						if !errors.Is(err, mpt.ErrResourceLimit) {
							t.Fatalf("constructor error = %v, want ErrResourceLimit", err)
						}
						if !errors.Is(rootErr, mpt.ErrUninitialized) {
							t.Fatalf("rejected constructor Root error = %v, want ErrUninitialized", rootErr)
						}
						return
					}
					if err != nil || rootErr != nil || root != mpt.EmptyRoot() {
						t.Fatalf("positive-limit constructor = root %v, error %v, Root error %v", root, err, rootErr)
					}
				})
			}
		}
	}
}
