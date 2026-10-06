# MPT family threat model

## Scope and assets

This model covers the root `mpt` package, canonical internal RLP codec, memory
adapter, and filesystem adapter. It includes raw/secure and Ethereum state,
storage, envelope, single/multi/range/EIP-1186 proof, construction, traversal,
rebuild, recovery, commit, retention, and pruning paths. The `_interop` sources
and `go.interop.mod` are test-oracle inputs, not production dependencies.

Protect exact root commitments, canonical bytes, immutable snapshots, complete
publication, retained historical roots, bounded owned work, and confidential
keys, values, proofs, preimages, paths, and collaborator diagnostics. Root/node
hashes and error categories are permitted bounded operational metadata.

Untrusted inputs include keys, values, encoded nodes, proofs, storage responses,
and imported fixtures. Stores may omit, substitute, corrupt, replay, or partially
persist data. Caller-selected authoritative roots, finite resource policies,
request admission, collaborator implementations, and an exclusively owned
filesystem directory are separate trust boundaries. No proof authenticates the
authority, freshness, finality, or chain membership of its supplied root.

## Owned work and bounds

| Boundary | Owned admission and work bounds | Owner |
| --- | --- | --- |
| Trie mutations and traversal | Positive `Limits`: mutation key/value admission, traversal keys/depth/nodes, encoding/hash/read work, batch operations, iterator results and rebuild work; stored-value reads use the canonical node ceiling | Root maintainers |
| Storage response resolution | Existing canonical RLP ceiling of 16 MiB before copy/hash; read/hash/traversal counts; context recheck after reader return | Root maintainers |
| RLP and compact paths | Canonical prefix/length/arity rules, checked length arithmetic, 16 MiB bytes, depth 1024 and 1,048,576 items; bounded compact paths and 31/32-byte reference rules | Codec maintainers |
| Proof decoding/verification | Key, node, aggregate byte, depth and hash limits; exact claim/path/witness consumption; no surplus-node acceptance | Root maintainers |
| Recovery overlays | Recovery node/byte caps before copying, hash verification and canonical decode; immutable source-store identity | Root maintainers |
| Pending snapshots and materialized cache | Positive `MaxPendingNodes`/`MaxPendingBytes`; incremental admitted-union accounting, conservative shared-encoding retention and bounded reachable filtering at compaction, encoded-child preflight and immutable decoded-buffer summaries; cache admission also obeys encoding/depth ceilings without rescanning shared subtrees | Root maintainers |
| Reachability and prune mark | Root/node/read/hash/depth limits, remaining aggregate byte budget and canonical node ceiling before copy/hash; cycle detection and deduplication | Root maintainers |
| Commit handles | Owned handle slice, private immutable encodings, owned `Encoded` bytes; `NodeCount`/`EncodedLen` permit preflight without payload copies | Root and adapter maintainers |
| Memory admission | Default cumulative 4,194,304 unique nodes and 256 MiB encoded bytes, including old roots; positive explicit `NewWithLimits`; zero-value uses defaults; preflight before payload or historical-map copies | Memory maintainers |
| Filesystem adapter | Positive node/commit-byte/commit-count/stored-node/retention limits; commit count before handle copying, cancellable per-node validation before payload copying/hashing; bounded regular-file reads and directory inventory, checked products and counters | Filesystem maintainers |

The limits bound library work, not total application memory. Allocator/map
overhead, caller-owned input buffers, caller-retained snapshots/returned bytes,
and caller-controlled concurrency need application admission and lifetime
policies. Copying convenience values does not authorize unbounded transport
input; validate lengths before constructing them. Use explicit smaller limits
for application-specific budgets rather than treating defaults as sizing advice.
The encoded pending chain and materialized cache are separately bounded by
the snapshot ceilings; working buffers and externally retained old snapshots
are additional resources, not an aggregate heap guarantee.

## Integrity, ownership, and lifecycle

Legacy Keccak-256 from `golang.org/x/crypto/sha3` implements Ethereum's required
fixed commitment; filesystem metadata checksums use standard SHA-256. There is
no custom cryptographic primitive or secret-key management. Secure keys are
hashed exactly once; profile-specific paths and typed envelope activation remain
explicit. These helpers validate trie framing, not full transaction semantics
or signatures. Selected interoperability differences remain in the
[decision register](specification-decisions.md).

Constructors and updates own retained bytes. Commit handles cannot expose their
private encoding; `Encoded`, `GetNode`, and returned values copy bytes. Readers
and stores must not concurrently mutate buffers they return while the library
is copying them. Root and memory packages perform no hidden network, filesystem,
goroutine, retry, or background work. User callbacks run synchronously. Context
checks bound traversal and publication seams, including after a reader returns;
they cannot forcibly stop an arbitrary collaborator that ignores its context.
Pending materialization, compaction, and commit-handle preparation check each
retained entry and discard partial preparation on cancellation before store
publication. Recovery overlay merging and inheritance check each entry as well.
Immutable summaries keep cache admission path-local rather than
rescanning the complete decoded graph after each mutation.

Resolved hash encodings remain available because other paths can still share
them. Only existing pressure/depth compaction seams filter provably unreferenced
owned nodes. An unresolved backing-store frontier conservatively preserves the
bounded admitted union without additional backend reads; admission can reject
when that frontier prevents safe reclamation. Compaction has bounded traversal
cost, while ordinary cache admission remains incremental.

Memory publication compares the immutable base state after bounded preparation;
rejection or stale publication preserves the old state. Admission counts unique
hashes, and successful pruning recomputes retained bytes. Historical-root leases
must outlive their readers. Filesystem node writes are synced before root
replacement; root rename is the visible publication point. Errors after that
point can leave the new complete root visible with uncertain sync durability:
reopen and reconcile before retrying. Retention and prune journals have explicit
recovery boundaries; refer to [filesystem operations](filesystem-store.md).

Default storage errors render only the fixed category. `errors.Is` and
`errors.As` retain deliberate original-cause inspection. Missing/corrupt-node
errors expose only their hash metadata. Cause inspection is not safe-to-log
authorization; application loggers must preserve confidentiality.

## Residual risks and review triggers

| Risk and rationale | Mitigation and owner | Revisit when |
| --- | --- | --- |
| Caller supplies an unauthorized or stale root; commitment verification cannot establish root authority | Application owner authenticates root provenance/finality before verification | Root authority or chain trust changes |
| A reader/callback or operating-system file operation blocks despite cancellation; synchronous collaborators and regular-file syscalls are not forcibly preemptible | Application owner uses context-aware readers, bounded request concurrency/timeouts, and trusted local storage; maintainers check owned seams | Collaborator or runtime I/O boundary changes |
| Large finite policies, many concurrent requests, or retained snapshots consume additional process resources outside one store's encoded-byte accounting | Application owner bounds admission/lifetimes and budgets overhead; memory maintainer enforces cumulative stored bytes/nodes and pruning | Service concurrency, admission policy, or store state representation changes |
| Filesystem path replacement by another writer defeats check/open assumptions; this is outside the exclusive-directory contract | Application owner restricts directory/ancestor permissions and ensures one owning Store; adapter rejects observed symlinks and files that are not regular | Shared/untrusted directory access is introduced |
| Crash or sync failure after a durable publication point requires reconciliation, not blind retry | Application owner reopens/inventories; filesystem maintainer preserves checksummed journals and complete-root publication | Filesystem publication, journal, or supported platform changes |
| Backend diagnostics remain recoverable through explicit causes and can leak through application logging | Application owner redacts explicit cause output; maintainers keep default strings metadata-only | Error wrappers, formatting, logging, or backend integration changes |
| Dependency, fixture or peer movement may affect trust without changing pinned source | Maintainers review checksum-pinned authorities/provenance, vulnerability and secret scans, required CI and release consumer evidence | Authority drift, advisories, dependency changes or release boundary |

## Scoped scanner evidence

The 2026-09-30 production scan used `govulncheck` v1.1.4 rebuilt with Go 1.27.1,
`golang.org/x/crypto` v0.56.0 and `golang.org/x/sys` v0.47.0. It reported zero
reachable-symbol and imported-package findings. Module-only GO-2026-5932 covers
the unmaintained OpenPGP package, has no fixed version, and is not imported;
production imports SHA-3 only. Maintainers must revisit that disposition if
crypto imports or the dependency graph change. Gitleaks found no leaks in the
scanned working tree; that result does not certify external backend diagnostics,
all historical Git objects, or future artifacts. Required hosted workflow,
CodeQL, release rehearsal and public-consumer evidence remain separate gates.

These are explicit ownership boundaries, not a blanket security certification.
The v2 upgrade is required because finite memory admission changes previously
accepted operations; canonical data formats remain compatible. Release
readiness requires current review, required CI, rehearsal and an actual published
v2 clean consumer. Local source/tests or scanner success alone do not establish
that delivery state.
