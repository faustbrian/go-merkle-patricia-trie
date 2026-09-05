# Documentation

`merkle-patricia-trie` is stable at v1 and requires Go 1.26.6. Its root public
package uses the default import identifier `mpt`; the public `memory` and
`filesystem` packages provide process-local and durable stores. The
implementation is portable Go with no platform-specific production source
files or cgo. The optional filesystem store is caller-opened and caller-closed;
the other public packages have no close lifecycle or background work. See the
[compatibility policy](../COMPATIBILITY.md) for the supported public contract.

- [Detailed package reference](reference.md)

## Getting started

- [Adoption](adoption.md)
- [Executable examples](../example_test.go)
- [FAQ](adoption.md#faq)

## Concepts and design

- [Architecture](architecture.md)
- [Encoding](encoding.md)
- [Filesystem Store](filesystem-store.md)
- [Profiles And Proofs](profiles-and-proofs.md)
- [Source Provenance](source-provenance.md)
- [Specification Decisions](specification-decisions.md)

## Operations and security

- [Operations](operations.md)
- [Filesystem store](filesystem-store.md)
- [Security model](security.md)
- [Security reporting](../SECURITY.md)
- [Support](../SUPPORT.md)

## Reference and maintenance

- [Benchmarks](benchmarks.md)
- [Compatibility Decisions](compatibility-decisions.md)
- [Compatibility policy](../COMPATIBILITY.md)
- [Deprecation policy](../DEPRECATION.md)

## Contributing

- [Contribution guide](../CONTRIBUTING.md)
- [Release history](../CHANGELOG.md)
- [License](../LICENSE)
