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
