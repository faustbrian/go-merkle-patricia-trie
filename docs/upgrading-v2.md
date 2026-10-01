# Upgrading to v2

Use `github.com/faustbrian/go-merkle-patricia-trie/v2` and append `/v2` before
the `memory` or `filesystem` package path. The source remains on main; the
release tag is `v2.0.0`. Go 1.27.0 is required. Do not mix v1 and v2 root,
snapshot, store, proof, or commit types in one operation.

The breaking admission policy is intentional: `memory.New` and the zero-value
`memory.Store` now retain at most 4,194,304 unique nodes and 256 MiB of encoded
node bytes. Historical nodes count until `Prune` removes them. Map/allocator
overhead, caller-held snapshots, and concurrent caller operations are additional
process resources; applications must also bound their requests and lifetimes.
Choose positive finite `memory.Limits` through `memory.NewWithLimits` when a
different capacity is required. An over-limit commit returns
`mpt.ErrResourceLimit` without changing the root or retained state. Release
historical-root leases only after their readers finish, then prune to reclaim
admission capacity. Existing leases and compare-and-swap semantics are unchanged.

Trie `Limits` also gains positive `MaxPendingNodes` and `MaxPendingBytes`,
defaulting to 1,048,576 nodes and 256 MiB. They bound both the encoded pending
chain and its stale ancestry, and separately the materialized node/path/value
cache, including decoded buffers that share identical hashes. Cache admission
also uses the encoding-node and traversal-depth ceilings. Construct policies
from `DefaultLimits` and adjust fields; update any complete struct literals.
Replacement/deletion and early ancestry compaction reclaim provably unused
capacity. Shared encodings remain retained; an unresolved backing-store frontier
can prevent reclamation and cause admission to reject without changing the root.
Compaction performs bounded traversal only at pressure/depth seams. Recovery
and rebuild use the same owner; failed admission preserves the original
snapshot. Committing does not waive the retained materialized-cache bound.
The two bounded representations, working buffers, map/allocator overhead and
caller-held older snapshots still require application-level memory sizing.

Storage responses exceeding the existing 16 MiB canonical node ceiling now
return `ErrResourceLimit` before copying or hashing, rather than first returning
hash corruption. Reachability also enforces its remaining aggregate byte budget
before copying. Cancellation after a reader returns is checked before owned
work resumes. Arbitrary reader implementations must still honor their supplied
context themselves.
Pending materialization, compaction and commit-handle preparation also observe
cancellation between entries, including recovery overlay merging/inheritance;
canceled preparation does not reach publication.

`ProveMany` validates every key length before copying or sorting caller keys.
Range generation and verification use `MaxProofBytes` separately for encoded
witnesses and aggregate item key/value bytes. A range with repeated large values
can exceed item admission even when its shared encoded witnesses fit; narrow
the requested range or choose a reviewed larger finite proof budget.

Default storage-error strings no longer include underlying diagnostics or paths.
Use `errors.Is` for the MPT category or original cause, and `errors.As` for
explicit cause inspection. Treat recovered backend diagnostics as confidential;
do not log them by default. Missing/corrupt-node metadata remains unchanged.

`StoreCommit.Nodes` still returns an owned slice of immutable handles and
`StoredNode.Encoded` still returns owned bytes. Handle enumeration avoids
duplicating payloads; `NodeCount` and `EncodedLen` support adapter admission
before payload allocation. Canonical RLP, legacy Keccak commitments, profile
identities, proof decisions, filesystem records, and existing stored node bytes
are unchanged; there is no data-format migration.
