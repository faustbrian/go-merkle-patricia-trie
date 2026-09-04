# Ethereum MPT and RLP conformance matrix

The root module implements the bounded behavior recorded in the
[specification decision register](../docs/specification-decisions.md). This
matrix does not claim full Ethereum, EVM, JSON-RPC, execution-specification, or
client compliance. It binds only the named decisions to pinned authorities,
official fixtures, and maintained-peer evidence.

[`monitoring.json`](monitoring.json) pins immutable authority content and the
reviewed change feeds. [`provenance.json`](provenance.json) pins imported
fixture and peer versions. [`conformance.json`](conformance.json) is the
machine-readable evidence map.

## Upstream review history

### 2026-09-04

- The execution-specifications Amsterdam feed advanced from commit
  `1855bb169fdf8b29ff7fb1eb6396e855549c9d7e` to
  `132d1149a257c5174dfd2f38f8cf1cb521780f06`. The seven intervening commits
  changed tests, tooling, and unrelated gas accounting only; the exact MPT
  source retained Git blob
  `0dbf455ad215e7c8f25ae35cf5149e1fc957b2a1`. Pinned authority bindings,
  decisions, and runtime behavior remain unchanged.
- The EIPs feed advanced from commit
  `94f5a3e3c146c28625d9ab2f8a7c0a848530a13a` to
  `7243c92ba812437c64bae9fc6524ee269b29daa9`. The eleven intervening commits
  updated EIP-2780, EIP-7906, EIP-8037, EIP-8130, and EIP-8272 and added an
  unrelated proposal; EIP-2718 and EIP-1186 retained Git blobs
  `83a19b0fa865dc31b483cd97f35c417928792d7c` and
  `1a341c3f9b8094955d386ad63b5c04e3bcc491f3`. Their decisions and runtime
  behavior remain unchanged.

### 2026-09-03

- The execution-specifications Amsterdam feed advanced from commit
  `c4deda5b3cfc5c1c8429dcd9159a6fb5636d8486` to
  `1855bb169fdf8b29ff7fb1eb6396e855549c9d7e`; the exact MPT source retained
  Git blob `0dbf455ad215e7c8f25ae35cf5149e1fc957b2a1`.
- The EIPs feed advanced from commit
  `889f8c1e26e9b418f83721083098ca225b14fc0b` to
  `94f5a3e3c146c28625d9ab2f8a7c0a848530a13a`; EIP-2718 and EIP-1186 retained
  Git blobs `83a19b0fa865dc31b483cd97f35c417928792d7c` and
  `1a341c3f9b8094955d386ad63b5c04e3bcc491f3`.
- The Go release feed advanced through Go 1.26.8 without changing the pinned
  Go memory-model bytes. These source checks leave MPT decisions 005 through
  007 and 009 behavior-neutral; decision 010 records the reviewed, benign feed
  movement. Pinned normative source bindings and runtime behavior are
  unchanged.

| Decision | Primary authority | Evidence boundary | Differential classification |
| --- | --- | --- | --- |
| MPT-DEC-001 | `yellow-paper-mpt-source` | Root commitment and official roots | Maintained peer agreement |
| MPT-DEC-002 | `yellow-paper-mpt-source` | Empty value and empty raw key histories | Deliberate policy difference |
| MPT-DEC-003 | `yellow-paper-mpt-source` | Compact paths and 31/32-byte references | Maintained peer agreement |
| MPT-DEC-004 | `yellow-paper-rlp-source` | Canonical RLP and malformed input | Deliberate policy difference |
| MPT-DEC-005 | `eip-2718-source` | Activated typed transaction and receipt roots | Maintained peer agreement |
| MPT-DEC-006 | `execution-spec-mpt-source` | Covered account, state, and storage roots | Maintained peer agreement |
| MPT-DEC-007 | `eip-1186-source` | Package proof strictness and covered peer proofs | Maintained peer agreement |
| MPT-DEC-008 | `geth-proof-source` | Half-open range witnesses | Deliberate policy difference |
| MPT-DEC-009 | `go-memory-model-source` | Package snapshot and publication policy | Not assessed |
| MPT-DEC-010 | `execution-spec-mpt-source` | Authority precedence and disagreement handling | Deliberate policy difference |

Source movement requires review; it never changes runtime behavior
automatically. A decision digest change requires an explicit changelog entry,
and unresolved normative contradictions remain release-blocking.
