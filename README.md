# merkle-patricia-trie

[![CI](https://github.com/faustbrian/go-merkle-patricia-trie/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/faustbrian/go-merkle-patricia-trie/actions/workflows/ci.yml)
[![CodeQL](https://img.shields.io/badge/CodeQL-required-blue)](https://github.com/faustbrian/go-merkle-patricia-trie/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Mutation](https://img.shields.io/badge/mutation-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Documentation](https://img.shields.io/badge/docs-checked_in_CI-blue)](docs/)
[![Go Reference](https://pkg.go.dev/badge/github.com/faustbrian/go-merkle-patricia-trie.svg)](https://pkg.go.dev/github.com/faustbrian/go-merkle-patricia-trie)
[![Release](https://img.shields.io/github/v/release/faustbrian/go-merkle-patricia-trie?sort=semver)](https://github.com/faustbrian/go-merkle-patricia-trie/releases)
[![Go](https://img.shields.io/badge/go-1.27.0-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`merkle-patricia-trie` implements Ethereum's execution-layer modified Merkle
Patricia trie. Its root package is `mpt`. It provides immutable raw and secure
tries, canonical RLP encoding, legacy Keccak-256 commitments, proofs,
caller-owned storage, recovery, retention, and pruning.

It does not implement an EVM, blockchain, network, JSON-RPC server, binary
Merkle tree, SSZ merkleization, or Verkle tree.

## Status and portability

The module is stable at v1 and requires Go 1.27.0. Its root package uses the
default import identifier `mpt`; the public `memory` and `filesystem` packages
provide process-local and durable stores. The implementation is portable Go:
it has no platform-specific production source files and requires no cgo. The
root and memory packages require no operating-system service or external
runtime backend. The optional filesystem package uses a caller-owned local
directory. These portability statements do not imply validation on every
`GOOS` and `GOARCH` combination.

Callers own configuration, authoritative-root selection, and runtime
resources. Trie snapshots and the memory store have no close lifecycle or
background work. A filesystem store must be opened and closed by its caller,
and its directory must have one exclusive owner.

## Installation

```sh
go get github.com/faustbrian/go-merkle-patricia-trie
```

## Quick start

```go
trie, err := mpt.NewRawTrie(mpt.DefaultLimits())
if err != nil {
	return err
}

trie, err = trie.Update(ctx, []byte("dog"), []byte("puppy"))
if err != nil {
	return err
}

value, err := trie.Get(ctx, []byte("dog"))
if err != nil {
	return err
}
root, err := trie.Root()
```

Updates return new immutable snapshots. Empty values delete keys. Use
`NewSecureTrie` when keys must be legacy-Keccak transformed exactly once.

The checked-in [`ExampleRawTrie`](example_test.go) is the executable
five-minute version of this flow. It is compiled and run by the Go example test
gate. The same file contains executable examples for secure, state, and storage
tries, proofs, and caller-owned persistence.

## Guarantees and limits

- Canonical nibble paths, hex-prefix encoding, RLP, and embedded-versus-hashed
  child references are enforced.
- Proofs bind roots, key transformation, value or absence, canonical nodes,
  path transitions, and explicit resource limits.
- Loaded nodes are hash checked and canonically decoded before use.
- Storage publication is compare-and-swap and caller-owned; pruning validates
  the complete retained graph before deletion.
- Compatibility is limited to the pinned Ethereum profiles and fixtures.

## Documentation

Use the [documentation index](docs/README.md) for the complete guide set. Start
with [adoption and FAQ](docs/adoption.md), the
[detailed reference](docs/reference.md), and the
[executable examples](example_test.go). Operational and project navigation is
available through [storage and recovery](docs/operations.md), the
[filesystem-store guide](docs/filesystem-store.md), [support](SUPPORT.md),
[security reporting](SECURITY.md), the [compatibility policy](COMPATIBILITY.md),
the [changelog](CHANGELOG.md), and the [license](LICENSE). The
[specification decision register](docs/specification-decisions.md) defines the
audited conformance boundary and known peer divergences.

For ecosystem-wide package selection and ownership conventions, see the
[versioned Golib ecosystem index](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/README.md)
and its [Domain utilities family](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/design-language.md#package-families-and-selection).

## Development

Run `make check` for the repository contract and the documented conformance
gates before changing encoding, commitment, proof, or storage behavior.

## License

MIT. See [LICENSE](LICENSE).
