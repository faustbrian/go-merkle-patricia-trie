# merkle-patricia-trie

[![CI](https://github.com/faustbrian/go-merkle-patricia-trie/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/faustbrian/go-merkle-patricia-trie/actions/workflows/ci.yml)
[![CodeQL](https://img.shields.io/badge/CodeQL-required-blue)](https://github.com/faustbrian/go-merkle-patricia-trie/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Mutation](https://img.shields.io/badge/mutation-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Documentation](https://img.shields.io/badge/docs-checked_in_CI-blue)](docs/)
[![Go Reference](https://pkg.go.dev/badge/github.com/faustbrian/go-merkle-patricia-trie.svg)](https://pkg.go.dev/github.com/faustbrian/go-merkle-patricia-trie)
[![Release](https://img.shields.io/github/v/release/faustbrian/go-merkle-patricia-trie?sort=semver)](https://github.com/faustbrian/go-merkle-patricia-trie/releases)
[![Go](https://img.shields.io/badge/go-1.26.6-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`merkle-patricia-trie` implements Ethereum's execution-layer modified Merkle
Patricia trie. Its root package is `mpt`. It provides immutable raw and secure
tries, canonical RLP encoding, legacy Keccak-256 commitments, proofs,
caller-owned storage, recovery, retention, and pruning.

It does not implement an EVM, blockchain, network, JSON-RPC server, binary
Merkle tree, SSZ merkleization, or Verkle tree.

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

Use the [documentation index](docs/README.md) for profiles, proofs, storage,
security, compatibility, and operations. The [detailed reference](docs/reference.md)
contains the full storage, recovery, EIP-1186, and state-trie contracts. The
[specification decision register](docs/specification-decisions.md) defines the
audited conformance boundary and known peer divergences.

## Development

Run `make check` for the repository contract and the documented conformance
gates before changing encoding, commitment, proof, or storage behavior.

## License

MIT. See [LICENSE](LICENSE).
