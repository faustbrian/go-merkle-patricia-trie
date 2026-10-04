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

### 2026-10-04

- The execution-specifications Amsterdam feed advanced from commit
  `903b48f152c932f6e47a615f0f7f009c56f1d92b` to
  `a87891f7e69eab1f903233c61c5514d8c94bd5d1`. The reviewed MPT source retains
  Git blob `0dbf455ad215e7c8f25ae35cf5149e1fc957b2a1`. Upstream fork and
  storage transitions can change execution results; callers own those
  transitions. The covered root, encoding, and proof contracts remain
  unchanged; no new fork profile is adopted.
- The EIPs feed advanced from commit
  `a9031bdc85949321a9707dd59ba44cdcba4a0eb0` to
  `666fb5ebc712bef2179af03317fffabf15a2f563`. EIP-2718 and EIP-1186 retain
  Git blobs `83a19b0fa865dc31b483cd97f35c417928792d7c` and
  `1a341c3f9b8094955d386ad63b5c04e3bcc491f3`.
- The Geth release feed advanced through v1.17.7. Its `trie/proof.go` source
  remains byte-identical to pinned v1.17.3, with SHA-256
  `f1092b71ebdda4f11a54a7a1c27f6ccbc0599455f2bfab68615600442564deb8`.
  Peer pins and executed interoperability evidence remain unchanged; this
  review does not claim execution against the latest client.
- All pinned authority payloads remain unchanged. The proof decision corrects the
  documented Geth proof direction; its historical digest is retained.
  The other nine decision identities and all selected behavior are unchanged.
  The current MPT source is distinct from the immutable pinned revision;
  this review preserves those separate identities. Monitoring remains strict.

### 2026-09-30

- All twelve monitored authority responses were fetched and checked. Immutable
  normative sources, the Yellow Paper change feed, and the Go release feed
  retain their recorded SHA-256 digests. The 30-day review interval is unchanged.
- The execution-specifications feed now ends at
  `847cdbdb7e0130132dbbc987ffa3a9d31ca3ade5`; its SHA-256 changed from
  `311471dca9b1d8a4c7ceb667d53340b800efdf0cfdb762a8a9584fd6de671a53`
  to `9905bddafe50bce0b82df411b98d2bed866bcc5ea28c6f460085c0667c2e5c8b`.
  The exact MPT source retains blob
  `0dbf455ad215e7c8f25ae35cf5149e1fc957b2a1`.
- The EIPs feed now ends at `66daa41124581e4e839e89d71eb06b6cd4b1f9b8`;
  its SHA-256 changed from
  `2525eea7c1153152a1278011cb4b02b43ba3e8782ea4d49274f40f46550d3cc3`
  to `be8e6dea57dff6ef27006628b595e1d488e34732e6fa914c6ab7cbb60631fbf6`.
  EIP-2718 and EIP-1186 retain blobs
  `83a19b0fa865dc31b483cd97f35c417928792d7c` and
  `1a341c3f9b8094955d386ad63b5c04e3bcc491f3`.
- The final same-day recheck observed one further EIPs feed entry at
  `52593e4aa929247bae05463f190b9d00699e9395`. The feed SHA-256 moved from
  `be8e6dea57dff6ef27006628b595e1d488e34732e6fa914c6ab7cbb60631fbf6`
  to `fbc11507122f66224e4fef8f25cd9471f179af8c414bbd730b113a92be9969f3`.
  Both covered EIPs retain their exact source digests and Git blobs above;
  normative pins and selected behavior remain unchanged.
- The Geth release feed includes v1.17.7 at
  `3d858f858a458effb2a563788aedf1fe65e1f0d3`; its SHA-256 changed from
  `196f09de85c8a92c65d0c648db0e11553294815e76ad77aeffa5120ad6f6f552`
  to `b7409d7179b3a4ec85dec82b45245f2a8fbfa16e50a9d10cfc1802e31b7b8ae1`.
  Its `trie/proof.go` is byte-identical to the pinned v1.17.3 proof source
  (`f1092b71ebdda4f11a54a7a1c27f6ccbc0599455f2bfab68615600442564deb8`).
  The maintained-peer oracle remains v1.17.3; no normative pin, selected
  specification decision, wire encoding, or runtime rule changes from this
  authority review.

### 2026-09-06

- The EIPs feed advanced from commit
  `9207c6011f526bd40abd79649484a1a342585bd4` to
  `a9031bdc85949321a9707dd59ba44cdcba4a0eb0`. The sole intervening commit
  added EIP-8360 only; EIP-2718 and EIP-1186 retained Git blobs
  `83a19b0fa865dc31b483cd97f35c417928792d7c` and
  `1a341c3f9b8094955d386ad63b5c04e3bcc491f3`. Their decisions and runtime
  behavior remain unchanged.

### 2026-09-05

- The execution-specifications Amsterdam feed advanced from commit
  `132d1149a257c5174dfd2f38f8cf1cb521780f06` to
  `903b48f152c932f6e47a615f0f7f009c56f1d92b`. The two intervening commits
  changed transaction-receipt test validation and EIP-7778 tests only; the
  exact MPT source retained Git blob
  `0dbf455ad215e7c8f25ae35cf5149e1fc957b2a1`. Pinned authority bindings,
  decisions, and runtime behavior remain unchanged.
- The EIPs feed advanced from commit
  `7243c92ba812437c64bae9fc6524ee269b29daa9` to
  `9207c6011f526bd40abd79649484a1a342585bd4`. The two intervening commits
  changed EIP-8246 only; EIP-2718 and EIP-1186 retained Git blobs
  `83a19b0fa865dc31b483cd97f35c417928792d7c` and
  `1a341c3f9b8094955d386ad63b5c04e3bcc491f3`. Their decisions and runtime
  behavior remain unchanged.

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
