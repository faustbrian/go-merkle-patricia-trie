// Package mpt implements Ethereum's execution-layer modified Merkle Patricia
// trie. The core package is storage independent and distinguishes raw,
// Keccak-secured, and Ethereum protocol key profiles.
//
// Proof verification establishes a value or absence claim under a caller
// supplied root. It does not establish that the root is canonical, finalized,
// recent, or authorized. The pinned execution-specification and Yellow Paper
// revisions, fixture scope and defensive proof policies are recorded in
// docs/source-provenance.md and docs/specification-decisions.md.
//
// Start with NewRawTrie or NewSecureTrie and explicit Limits; StateTrie,
// StorageTrie and indexed-root helpers impose their distinct protocol profiles.
// Updates return owned immutable snapshots, not durable publication. Loaded
// snapshots borrow a NodeReader; concurrent use and historical availability
// depend on that provider and any required root retention. Commit remains bound
// to the source store until Rebuild detaches it.
//
// Operations start no workers. Context checks are cooperative and cannot preempt
// a blocking provider/callback or one RLP/hash/file operation. Input constructors
// without Limits copy caller-sized data before later operation admission; bound
// untrusted transport first. This package owns neither RPC transport, chain-root
// authorization nor execution/state-transition policy. Storage adapters are
// optional packages; the core introduces no filesystem I/O or cgo dependency.
package mpt
