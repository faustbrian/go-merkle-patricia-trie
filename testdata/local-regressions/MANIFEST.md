# Local EIP-1186 regression fixture

- Origin: deterministic local regression corpus; not an official Ethereum
  fixture.
- Format: transport-decoded EIP-1186 account and storage proof material.
- Generator inputs: the exact state, account, storage, and decoy values in
  `TestEthereumJSEIP1186AccountAndStorageProofs`.
- Confirmed historical oracle source: `@ethereumjs/mpt` 10.1.2 at
  `3adf102baf8991f82feda860e0d3a3ec644d0802`. The original manifest named
  `@ethereumjs/trie` incorrectly; the original scripts and package lock already
  used `@ethereumjs/mpt`. This corrects attribution, not fixture bytes.
- Additional semantic oracle: go-ethereum 1.17.3 at
  `117e067f0f0bae1a17082321f224dedb6765b10f`.
- License: CC0-1.0.
- Applicability: account membership, account absence, storage membership, and
  storage absence under explicit supplied roots. JSON-RPC encoding is outside
  scope.
- Local coverage: `eip1186_fixture_test.go`,
  `ethereumjs_interop_test.go`, and `_interop/geth_interop_test.go`.

| File | SHA-256 |
| --- | --- |
| `eip1186.json` | `9ef6b4b3a4b740172a6372c717277de725ace3d8ac0af958d4087bcc13f6d000` |

## Historical and future producer identities

The fixture was introduced in commit `37ca0d8fb6a9d7ab44463f76c11d858e53db1544`.
Its original writer, serializer and producing execution were not recorded and
have not been recovered. The preceding oracle implementation was introduced in
`a1fe5f33936ef511fc58f27730fb172be51aa91d`; source and dependency identity do not
prove which writer serialized the historical corpus.

The authenticated npm MPT 10.1.2 tarball has SHA-256
`4a4902f3580c21bbed4bf1fc6e5fde04d5565e2278e04556915998f8505ade9f`
and matches the original lock's SHA-512 integrity value. Its declared source is
`28d6a165971964f239b77004196c6c0e87697835`. The MPT source tree and package metadata
are identical at that revision and the originally pinned revision above. The
oracle is MPL-2.0; the existing local fixture's CC0 dedication is separate.

The future explicit writer is
[`scripts/generate-eip1186-fixture.mjs`](../../scripts/generate-eip1186-fixture.mjs),
using the unchanged locked MPT and RLP 10.1.2 dependencies and fixed inputs above.
It produces account/storage membership and absence proofs, serializes keys in
fixed order with two-space indentation and a final newline, and emits stdout
only after confirming the preserved length and checksum. It never overwrites
this corpus. Two executions with authenticated Node 22.18.0 and its bundled npm
10.9.3 reproduced the exact existing 1,769 bytes. A changed-input counterfactual
was refused before output. Reproduction does not identify the historical writer.

## Update procedure

1. Review EIP-1186, the pinned client revisions, and the deterministic corpus.
2. Use an authenticated, pinned Node runtime and install the unchanged package
   lock with `npm ci --ignore-scripts`. Record the actual runtime and lock identity.
3. Run `node scripts/generate-eip1186-fixture.mjs`, redirecting stdout to a
   task-owned temporary file, and compare its complete bytes with `eip1186.json`.
   Do not overwrite the corpus to resolve a mismatch; review inputs and output.
4. Generate proofs from the local state and storage tries only for a separately
   reviewed corpus change, updating the producer's fixed identity guard together.
5. Verify the exact roots, values, and proof nodes with EthereumJS and Geth.
6. Serialize the decoded proof bytes without JSON-RPC quantity conversion.
7. Recompute the SHA-256 value and run conformance and interoperability gates.
