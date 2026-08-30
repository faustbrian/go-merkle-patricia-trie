# Changelog

All notable changes follow Keep a Changelog and Semantic Versioning.

## [Unreleased]

### Changed

- Advance the shared workflow to `go-library-tools` commit
  `3367371735ac409ba101778aafbdaf8e2560064c` so CI validates specification
  decisions and monitored authorities while retaining the released v1.0.13
  local tool contract.

### Removed

- Remove copied repository-local Golib scripts and tool-version definitions;
  package-specific conformance, interoperability, and benchmark commands now
  remain in the source-owned verification Makefile.

### Documentation

- Add machine-validated MPT and RLP decision, conformance, authority-monitoring,
  and history records with bounded maintained-peer evidence and no broad
  compliance claim.

- Align maintained documentation with the stable v1 release contract and
  remove obsolete candidate or verdict wording.

- Replace obsolete portal links and completed execution artifacts with a
  standalone, human-oriented documentation structure.

Decision records:

- MPT-DEC-001 sha256:2ba5a763a1942a391bc46ef49af29e6a0dc1a949ecf9f7cd8ac07aeeca9d7418
- MPT-DEC-002 sha256:0b6d356faacc3df0bc841b6c5e01b7361f1cc0641f1ce4ffbe960394235e6785
- MPT-DEC-003 sha256:98926b20816783f5814d61b612cfabc7215d82ecc5b80891b86c223ee4e9084e
- MPT-DEC-004 sha256:3492de50569a3a40a425d17bfcf70a56c8fff7708224883881f989e6b2fff4a5
- MPT-DEC-005 sha256:460277c8bb85a3c75b430a190a99bd6fb4f8484216df5736d2c02686ceb9fd6d
- MPT-DEC-006 sha256:7bfbd351ec14d2498f3bed5b6519b2e4d48810bdeba7a481b695100c931ab0b6
- MPT-DEC-007 sha256:51cd3a66aefb2d5f155b45c0f93504aabbd6317ffe4d75a53c04e52884a122f4
- MPT-DEC-008 sha256:a77e9e6397f4953c5c4a9abddce85ca2000338515dace510a62b0a16d39676cc
- MPT-DEC-009 sha256:cfeda2751d8c15f1b80973efd5ac5c75913440fdfeecb2f2f0a2bba156ef9967
- MPT-DEC-010 sha256:3d5b790cc5c951e801de2d1c17369da1a01bf863a6766ca85938f5d578873b8b

[Decision register](docs/specification-decisions.md)

## [1.0.0] - 2026-08-25

### Fixed

- Preserve the pinned EthereumJS MPT and RLP oracle dependencies alongside
  documentation tooling in the standalone Node manifest.

### Changed

- Upgrade the cryptographic dependency set to its current secure releases.

- Exclude intentional nested modules from root local-proxy archives so local,
  bootstrap, CI, and public module checksums describe the same source
  boundary.

- Track the pinned documentation-tool lockfile so clean CI checkouts install
  the exact validated cspell dependency.

- Reconcile standalone dependency checksums against deterministic current
  module archives so CI, local verification, and release consumers resolve
  identical content.

- Harden standalone documentation validation with deterministic spelling and
  link checks, package-specific documentation gates, and repository-local
  contributor guidance.

### Documentation

- Link the package README to package-owned documentation.

### Changed

- Publish the module from its standalone `github.com/faustbrian/go-merkle-patricia-trie` identity while preserving its documented API and behavior.
- Complete the Linux-specific Geth interoperability module graph so the pinned
  reference suite resolves read-only on CI.
- Link compatibility and source-provenance guidance directly to the canonical
  specification decision register.
- Convert consensus and client-divergence choices into an auditable
  specification register and bind every pinned source to an immutable URL.

### Added

- Added the stable-release hardening report with compatibility, fixture,
  evidence, and unsupported-scope matrices.
- Added a checksummed local EIP-1186 account and storage proof regression
  fixture alongside pinned Geth and EthereumJS proof interoperability.
- Established the Ethereum MPT module boundary, authoritative-source pins,
  compatibility decisions, and hostile-input threat model.
- Added canonical hex-prefix paths, bounded canonical RLP, legacy Keccak root
  commitments, canonical node encoding and validation, and exact embedded
  versus hashed child references.
- Added immutable raw and secure trie snapshots with bounded lookup, update,
  replacement, deletion, canonical compaction, and history-independent roots.
- Added atomic batch mutation, ordered raw and hashed-key iteration, lazy
  hash-addressed loading, integrity-checked atomic commits, and a concurrent
  in-memory store with stale-root protection.
- Added optional deterministic bounded node-store iteration for audits and
  rebuild tooling.
- Added root-verified raw and secure rebuilds that fully materialize snapshots
  for safe migration between stores.
- Added bounded immutable missing-node recovery overlays that validate fetched
  nodes, resume every traversal surface, and atomically repair the source store.
- Added bounded canonical reachability audits plus explicit historical-root
  leases and atomic mark-and-sweep pruning in the concurrent memory store.
- Rejected non-canonical hashed references to child encodings shorter than 32
  bytes across stored traversal, proofs, rebuilds, and reachability audits.
- Added bounded Ethereum-style membership and non-membership proof generation
  and verification with strict root, key, value, profile, ordering, and surplus
  node binding.
- Added deterministic raw and secure multi-key proofs with shared-node
  deduplication and mixed membership/absence verification.
- Added bounded raw and transformed secure-key range proofs for explicit
  `[start,end)` intervals, with consecutive-leaf completeness, strict witness
  ordering, and pinned Geth and EthereumJS verification interoperability.
- Made multi-key and range witness indexing observe context cancellation
  between proof nodes.
- Imported the pinned legacy Ethereum raw and secure trie fixture corpus
  byte-for-byte with checksum, license, update, applicability, and local
  coverage records.
- Imported pinned execution-spec-tests v5.4.0 blockchain fixtures byte-for-byte
  and verified official pre/post allocation roots plus legacy and type-1
  through type-4 transaction roots.
- Imported pinned Geth v1.17.3 transition fixtures byte-for-byte and verified
  legacy and typed receipt roots from exact receipt values.
- Added canonical RLP integer key derivation for raw transaction and receipt
  trie indexes.
- Added bounded fuzz harnesses for compact paths, canonical RLP, node decoding,
  proof verification, mutation sequences, rebuilds, EIP-1186, streaming
  construction, commit/recovery/pruning transitions, cancellation, and
  ordered iteration callback failures.
- Added transport-independent EIP-1186 account membership, account absence,
  canonical account decoding, and bounded one-or-many storage-slot proof
  verification with duplicate and conflicting slot rejection.
- Verified EIP-1186 account membership, account absence, storage membership,
  and storage absence bidirectionally against proofs generated by pinned Geth
  and EthereumJS.
- Added encoding, profile, proof, operation, storage, recovery, pruning,
  adoption, comparison, migration, and FAQ guides.
- Added reproducible local benchmark workloads and an ownership-equivalent
  populated raw lookup comparison against pinned Geth.
- Reworked immutable updates to structurally share canonical pending-node
  layers, bound overlay compaction, and discard superseded or unreachable
  hashed nodes without changing roots or snapshot semantics.
- Added immutable state and storage trie profiles with exact address and slot
  types, canonical account construction, unsigned word handling, zero-slot
  deletion, persistence, rebuild, recovery, proofs, and pinned Geth and
  EthereumJS root interoperability.
- Added separate validated transaction and receipt value types, explicit
  Berlin-through-Osaka EIP-2718 activation profiles, matching receipt-type
  enforcement, and canonical root construction from RLP indexes, with pinned
  Geth and EthereumJS interoperability.
- Added deterministic mutation-trace differential tests against pinned Geth
  and EthereumJS implementations for raw and secure trie profiles.
- Added direct canonical RLP encoding and hostile-decoding differential
  coverage against pinned Geth and EthereumJS implementations.
- Added a bounded sorted-input raw-trie root builder with strict ordering,
  transactional rejection, single finalization, and ordinary-insertion parity.
- Added a dependency-free durable filesystem node store with bounded reads and
  commits, immutable integrity-checked node files, atomic root publication,
  interrupted-write cleanup, process-termination recovery evidence, and
  ordinary/streaming/rebuild root parity.
- Added durable historical-root leases and crash-recoverable atomic
  mark-and-sweep pruning to the filesystem store, with bounded retention
  inventory and exclusive mutation ownership.

### Changed

- Filesystem retention and pruning now reject unopened stores, bound recovery
  without integer overflow, and document indeterminate post-publication
  storage outcomes.
- State-account nonces now use `uint64`, and balances use exact 32-byte words,
  matching the execution-spec account types instead of exposing ambiguous
  minimally encoded integer bytes.
- Replaced the ambiguous shared `EncodedTrieValue`, `LegacyTrieValue`, and
  `TypedTrieValue` pre-v1 API with profile-bound transaction and receipt types.
  Receipt-root callers must now provide the corresponding transaction values so
  EIP-2718 type equality is enforced.

[Unreleased]: https://github.com/faustbrian/go-merkle-patricia-trie/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/faustbrian/go-merkle-patricia-trie/releases/tag/v1.0.0
