# Source documentation audit

[Documentation map](index.md) · [Repository README](../README.md)
[Next: Detailed reference](reference.md)

## Initial scope and bounded plan

Tier A comments/current contract prose only from
`5101b8b1e3cc20c0272ca3d0fb923a950ed280f9`. This worker owns source documentation;
the coordinator owns publication and exact-source CI. The unrelated original
checkout is preserved. Reproduce with `git ls-files "*.go"`, the module/package
catalogs, `golib source-docs check --json` and
`rg --no-ignore "^func (Test|Fuzz|Benchmark|Example)" --glob "*_test.go"`.

All 107 tracked Go files (27 implementation, 80 assurance) belong to the root
`github.com/faustbrian/go-merkle-patricia-trie/v2` module. Production packages are
root `mpt`, `filesystem`, `memory` and `internal/rlp`; `_interop` is tagged test-only
support using the explicit alternate `go.interop.mod`. Four assurance sources
carry the interoperability tag; no production platform variants or generated Go
were found. Existing nine security annotations are preserved. Inventory contains
343 test functions, 21 fuzz targets, 39 benchmarks and eight Examples; these are
coverage partitions to inspect, not passing runtime or conformance claims.

Outcome: truthful commitment/proof encodings and key profiles, immutable snapshots,
borrowed store/callback obligations, cooperative cancellation/resource admission,
root publication/recovery/pruning windows and provider differences. Direct
consumers are raw/secure/protocol trie users, proof verification, recovery and
reachability walkers, builders/batches, memory/filesystem stores and fixtures.
Named risks: legacy Keccak versus SHA3/profile confusion, RLP/path canonicality,
31/32-byte references, hostile proof bounds, pending-overlay ownership, partial
persistence/stale-root publication, retain/prune safety, callback reentrancy and
trusted filesystem ownership. No runtime/dependency/fixture repair or desired
specification-policy rewrite is scoped.

Smallest checks: pinned Go1.27.2 all-file formatting, executable AST/token/literal
and annotation identity; four production package renders; native docs/Examples
recipe and configured API mode discovery; current structural/source/local-link
and offline specification checks; pinned changed-file/added/staged secret scans;
one ordinary complete Tier A final-diff review and explicit-path Intent commit.
The only API baseline mode is apidiff; immutable baselines remain unchanged.
All Go/native execution awaits an explicit serialized coordinator grant, with
three owned disposable caches and exact cleanup on success/failure.
The native root documentation route has no configured example operation; existing
Examples are therefore selected directly, rather than reported as its output.

## Complete tracked inventory and initial gaps

The initial checker reports 46 declaration gaps in documented packages.
This distinguishes missing comments from ownership/security/completeness questions
that must be investigated against source and pinned decisions before writing.

| Go file | Initial structural gaps | Audited boundary |
| --- | --- | --- |
| `_interop/benchmark_test.go` | 0 | BenchmarkComparableRawGetOwned |
| `_interop/geth_interop_test.go` | 0 | TestGethStateAndStorageTrieProfiles; TestGethEIP1186AccountAndStorageProofs; additional cases/helpers |
| `_interop/rlp_interop_test.go` | 0 | TestGethAndEthereumJSRLPEncodingDifferential; TestGethAndEthereumJSRLPMalformedInputInventory |
| `accounting_confidence_internal_test.go` | 0 | TestRangeItemByteAccountingBoundaries; TestPendingAccountingExactReplacementAndRetention; additional cases/helpers |
| `api_error_test.go` | 0 | TestZeroTrieValuesRejectEveryOperation; TestSecureTrieCompleteOperationSurface; additional cases/helpers |
| `batch.go` | 0 | implementation contracts and private invariants |
| `batch_test.go` | 0 | TestRawTrieBatchIsAtomicAndMatchesSequentialUpdates; TestBatchRejectsInvalidDuplicateAbsentAndOversizedChanges; additional cases/helpers |
| `benchmark_extended_test.go` | 0 | BenchmarkSnapshotAdmissionUpdateScaling; BenchmarkComparableRawGetOwned; additional cases/helpers |
| `benchmark_test.go` | 0 | BenchmarkGetEmpty; BenchmarkGetPopulated; additional cases/helpers |
| `builder.go` | 0 | implementation contracts and private invariants |
| `builder_internal_test.go` | 0 | TestSortedBuilderRetainsOnlyBoundedFrontier; TestSortedBuilderInternalFailurePaths; additional cases/helpers |
| `builder_test.go` | 0 | TestSortedBuilderMatchesOrdinaryInsertion; TestSortedBuilderOrderingAndLifecycle; additional cases/helpers |
| `compact.go` | 0 | implementation contracts and private invariants |
| `compact_test.go` | 0 | TestCompactPathKnownVectors; TestCompactPathExhaustiveShortRoundTrip; additional cases/helpers |
| `concurrency_architecture_test.go` | 0 | TestProductionPackagesOwnNoGoroutines; TestProductionGoroutineScopeExcludesExternalTooling |
| `consumer_test.go` | 0 | TestCleanExternalConsumer; TestPublicPackages |
| `coverage_internal_test.go` | 0 | TestInternalCompactAndNodeFailureContracts; TestInternalNodeDecodeFailureContracts; additional cases/helpers |
| `doc.go` | 0 | implementation contracts and private invariants |
| `eip1186.go` | 0 | implementation contracts and private invariants |
| `eip1186_fixture_test.go` | 0 | TestEIP1186RegressionFixture |
| `eip1186_internal_test.go` | 0 | TestStorageProofSetValidation |
| `eip1186_test.go` | 0 | TestEIP1186AccountAndStorageProofs; TestEIP1186StorageProofSet; additional cases/helpers |
| `envelope.go` | 0 | implementation contracts and private invariants |
| `envelope_profile_test.go` | 0 | TestTypedEnvelopeProfilesBindForkActivationAndPayloads; TestReceiptRootBindsMatchingTransactionTypes |
| `envelope_test.go` | 0 | TestTransactionAndReceiptRootsUseRLPIndexesAndExactEnvelopes; TestEncodedEnvelopeValuesValidateAndOwnBytes; additional cases/helpers |
| `errors.go` | 0 | implementation contracts and private invariants |
| `ethereumjs_interop_test.go` | 0 | TestEthereumJSTransactionAndReceiptRoots; TestEthereumJSStateAndStorageTrieRoots; additional cases/helpers |
| `example_test.go` | 0 | ExampleRawTrie; ExampleSecureTrie; additional cases/helpers |
| `execution_spec_fixture_test.go` | 0 | TestExecutionSpecFixtureChecksums; TestExecutionSpecStateAndTransactionRoots |
| `filesystem/benchmark_test.go` | 0 | BenchmarkFilesystemWarmGet; BenchmarkFilesystemOpenAndGet; additional cases/helpers |
| `filesystem/crash_test.go` | 0 | TestCommitSurvivesProcessTerminationAtPublicationBoundaries; TestCommitCrashHelper; additional cases/helpers |
| `filesystem/fuzz_test.go` | 0 | FuzzRootRecordCanonicalRoundTrip; FuzzNodeFilenameCanonicalDecode; additional cases/helpers |
| `filesystem/prune.go` | 0 | implementation contracts and private invariants |
| `filesystem/prune_test.go` | 0 | TestStorePrunesOnlyReleasedDurableRootsAcrossReopen |
| `filesystem/retention.go` | 0 | implementation contracts and private invariants |
| `filesystem/retention_prune_internal_test.go` | 0 | TestRetentionEncodingIdentifiersAndTemporaryRecovery; TestRetentionInventoryRejectsCorruptState; additional cases/helpers |
| `filesystem/retention_test.go` | 0 | TestStorePersistsIndependentRootRetentionsAcrossReopen; TestStoreRetentionValidatesRootLimitsAndEnumerationBounds |
| `filesystem/security_storage_test.go` | 0 | TestFilesystemStorageFailuresRedactPathsAndPreserveCauses; TestFilesystemReadAndCommitDiagnosticsRemainInspectableButRedacted; additional cases/helpers |
| `filesystem/store.go` | 0 | implementation contracts and private invariants |
| `filesystem/store_internal_test.go` | 0 | TestStoreCASIterationAndClosedLifecycle; TestStoreValidatesContextsReceiversAndIteratorInputs; additional cases/helpers |
| `filesystem/store_test.go` | 0 | TestStorePersistsPublishedTrieAcrossReopen; TestPersistentConstructionMatchesOrdinaryStreamingAndRebuild |
| `fuzz_test.go` | 0 | FuzzCompactPathCanonicalRoundTrip; FuzzEnvelopeValueValidation; additional cases/helpers |
| `geth_receipt_fixture_test.go` | 0 | TestGethReceiptFixtureChecksums; TestGethReceiptRoots |
| `hardening_fuzz_test.go` | 0 | FuzzSortedBuilderMutationHistoryParity; FuzzRawAndSecureRebuildMutationParity; additional cases/helpers |
| `index.go` | 0 | implementation contracts and private invariants |
| `index_test.go` | 0 | TestRLPIndexKeyCanonicalBoundaries; TestRLPIndexKeyReturnsOwnedBytes |
| `internal/rlp/coverage_test.go` | 0 | TestValueAccessorsRejectWrongKind; TestEncodeRejectsEveryBoundAndUnknownKind; additional cases/helpers |
| `internal/rlp/mutation_contract_test.go` | 0 | TestEncodeLimitsAreInclusiveAtEveryFramingBoundary; TestListPayloadReservationBoundsEveryAppend; additional cases/helpers |
| `internal/rlp/rlp.go` | 3 | implementation contracts and private invariants |
| `internal/rlp/rlp_fuzz_test.go` | 0 | FuzzCanonicalDecodeRoundTrip |
| `internal/rlp/rlp_test.go` | 0 | TestEncodeCanonicalVectors; TestEncodeShortAndLongLengthBoundary; additional cases/helpers |
| `iteration.go` | 4 | implementation contracts and private invariants |
| `iteration_test.go` | 0 | TestRawTrieIterationOrderPrefixRangeAndLimit; TestIterationOwnsEntryBytesAndPreservesCallbackError; additional cases/helpers |
| `legacy_fixture_test.go` | 0 | TestLegacyEthereumFixtureChecksums; TestLegacyEthereumTrieRoots; additional cases/helpers |
| `memory/admission_boundaries_test.go` | 0 | TestMemoryAdmissionRequiresEachPositiveDimension; TestMemoryAdmissionAcceptsExactOneNodeCapacity; additional cases/helpers |
| `memory/coverage_test.go` | 0 | TestStoreObservesCancellationAtReadAndCommitBoundaries; TestStoreDetectsConcurrentPublication; additional cases/helpers |
| `memory/fuzz_test.go` | 0 | FuzzCommitRetentionAndPruningStateMachine |
| `memory/security_limits_test.go` | 0 | TestStoreRetentionHonorsLoweredAdmissionBound; TestStoreAdmissionCancellationDoesNotPublish; additional cases/helpers |
| `memory/store.go` | 2 | implementation contracts and private invariants |
| `memory/store_internal_test.go` | 0 | TestInternalRetentionNilAndDetachedLifecycle; TestInternalStateReaderAndUniqueRoots; additional cases/helpers |
| `memory/store_test.go` | 0 | TestStoreCommitsAndLoadsTrie; TestStoreRejectsStaleRootAtomically; additional cases/helpers |
| `multiproof.go` | 0 | implementation contracts and private invariants |
| `multiproof_internal_test.go` | 0 | TestMultiProofBuilderRejectsInvalidRootSources; TestMultiProofBuilderRejectsInvalidTraversalSources; additional cases/helpers |
| `multiproof_test.go` | 0 | TestRawMultiProofDeduplicatesSharedNodesAndBindsAllClaims; TestSecureMultiProofBindsProfileAndClaims; additional cases/helpers |
| `mutation_contract_test.go` | 0 | TestCompactAndNodePathLimitsAcceptTheirExactMaximum; TestEveryTrieLimitRejectsZeroIndependently; additional cases/helpers |
| `node.go` | 0 | implementation contracts and private invariants |
| `node_decode_test.go` | 0 | TestDecodeNodeRoundTripsCanonicalForms; TestDecodeNodeRejectsImpossibleForms; additional cases/helpers |
| `node_footprint_internal_test.go` | 0 | TestDecodedFootprintAdmissionFailsClosedOnOverflow; TestDecodedFootprintRejectsUnresolvedMaterializedGraph |
| `node_fuzz_test.go` | 0 | FuzzCanonicalNodeDecodeRoundTrip; FuzzEthereumStateAndStorageValueDecoding |
| `node_test.go` | 0 | TestCanonicalEmptyRoot; TestRootOwnsBytes; additional cases/helpers |
| `pending_boundary_characterization_internal_test.go` | 0 | TestHostedPendingEncodingExactChildByteBoundary; TestMaterializedPendingCompactionTraversalDepthBoundary |
| `pending_encoded_root_boundary_test.go` | 0 | TestHostedCanonicalRootEncodedByteLimitIsInclusive |
| `pending_limits_constructor_test.go` | 0 | TestConstructorsRequirePositivePendingLimits |
| `pending_snapshot_boundaries_internal_test.go` | 0 | TestHostedMaterializedSnapshotExactByteBoundary; TestHostedMaterializedSnapshotExactDepthBoundary; additional cases/helpers |
| `proof.go` | 0 | implementation contracts and private invariants |
| `proof_test.go` | 0 | TestRawMembershipAndAbsenceProofs; TestProofVerificationBindsClaimRootKeyValueAndProfile; additional cases/helpers |
| `rangeproof.go` | 0 | implementation contracts and private invariants |
| `rangeproof_internal_test.go` | 0 | TestRangeProofGenerationInternalFailureContracts; TestRangeProofVerificationInternalFailureContracts; additional cases/helpers |
| `rangeproof_test.go` | 0 | TestRawRangeProofEstablishesEveryLeafInExplicitInterval; TestRangeProofIncludesPrefixKeysAndUnboundedEndpoints; additional cases/helpers |
| `reachability.go` | 11 | implementation contracts and private invariants |
| `reachability_internal_test.go` | 0 | TestReachabilityLimitsRejectEachInvalidBound; TestReachabilityExactBoundsAndFailureClassification; additional cases/helpers |
| `reachability_test.go` | 0 | TestCollectReachableNodesReturnsCanonicalGraph |
| `rebuild.go` | 0 | implementation contracts and private invariants |
| `rebuild_test.go` | 0 | TestRawTrieRebuildMatchesRootAndCanMoveStores; TestSecureTrieRebuildPreservesSecureProfile; additional cases/helpers |
| `recovery.go` | 0 | implementation contracts and private invariants |
| `recovery_cancellation_test.go` | 0 | TestPendingCompactionPropagatesDecodeFailureWithoutPartialResult; TestRecoveredMutationCancellationPrecedesOverlayInheritance; additional cases/helpers |
| `recovery_internal_test.go` | 0 | TestRecoverSnapshotHonorsHashAndPostValidationCancellation; TestRecoveryPendingLayersRemainBounded; additional cases/helpers |
| `recovery_test.go` | 0 | TestRawTrieRecoversMissingNodeIntoImmutableSnapshot; TestRecoveredNodeCanParticipateInAtomicUpdate; additional cases/helpers |
| `root.go` | 0 | implementation contracts and private invariants |
| `security_cancellation_internal_test.go` | 0 | TestProofCancellationPrecedesPendingMaterialization; TestSecurityAdmissionCancellationPreservesSnapshots; additional cases/helpers |
| `security_commit_test.go` | 0 | TestStoreCommitMetadataDoesNotCopyEncodedPayloads; TestStoreCommitHandlesPreservePublicOwnershipAndReuse |
| `security_pending_test.go` | 0 | TestPendingAdmissionBoundsRawAndSecureSnapshots; TestMaterializedCacheAdmissionCountsDistinctBuffersWithSharedHashes; additional cases/helpers |
| `security_range_test.go` | 0 | TestRangeItemAdmissionCountsSharedLeafCopies; TestRangeWitnessBudgetRemainsIndependentOfItemBudget |
| `security_storage_test.go` | 0 | TestMultiProofInputAdmissionPrecedesKeyCopies; TestStorageResponsesAreBoundedBeforeOwnedWork; additional cases/helpers |
| `shared_pending_integrity_test.go` | 0 | TestSharedPendingLeafSurvivesSiblingReplacement |
| `state_profile.go` | 0 | implementation contracts and private invariants |
| `state_storage_internal_test.go` | 0 | TestDecodeStorageWordRejectsMalformedAndNonCanonicalValues; TestDecodeAccountRejectsNonceBeyondUint64; additional cases/helpers |
| `state_storage_test.go` | 0 | TestAccountValueAndStateTrieUseCanonicalEthereumEncoding; TestStateTrieDoesNotApplyEmptyAccountLifecycleRules; additional cases/helpers |
| `storage_profile.go` | 0 | implementation contracts and private invariants |
| `store.go` | 9 | implementation contracts and private invariants |
| `store_errors.go` | 0 | implementation contracts and private invariants |
| `store_test.go` | 0 | TestCommitAndLoadRawTrie; TestLoadEmptyRootDoesNotReadStore; additional cases/helpers |
| `structural_sharing_internal_test.go` | 0 | TestFinishedSnapshotsFreezeHashedChildrenForStructuralSharing; TestStructurallySharedPendingNodesRemainCanonical; additional cases/helpers |
| `trie.go` | 17 | implementation contracts and private invariants |
| `trie_internal_test.go` | 0 | TestBytesToNibblesUsesHighThenLowOrder |
| `trie_model_test.go` | 0 | TestRawTrieExhaustiveSmallOperationHistories |
| `trie_test.go` | 0 | TestRawTrieImmutableUpdateLookupAndDelete; TestRawTrieStrictPrefixesAreHistoryIndependent; additional cases/helpers |

## Investigation disposition

All 27 implementation files and every assurance source were traced for declarations,
assertions, helper ownership, tagged peer adapters, fixture inputs and work/cleanup.
The 46 initial missing declarations are documented without private-symbol filler.
Existing source prose also required the semantic corrections below.

- Roots use legacy Keccak, not FIPS SHA3; every root hashes complete node RLP,
  while child references embed complete lists below 32 bytes and hash at 32+.
  Canonical RLP/compact paths reject overlong framing, trailing input, bad flags,
  empty extensions and collapsible branches rather than normalizing commitments.
  RLP empty encoded input is rejected; zero Value is the empty string.
- Proof construction owns transport bytes with count/byte admission, not a
  verified claim. Single witnesses are root-first paths; multi witnesses use
  caller-key ordering/first shared-node use, not sorted hash order. Ranges use
  half-open endpoints and exact consecutive leaves; secure range/iteration
  endpoints are already transformed paths, unlike secure mutation/lookups.
  Range witness and item-byte ceilings are independent. Proof/root validity
  does not establish chain authorization/finality or an RPC response format.
- Raw empty keys are valid; empty values delete and absent deletion reports
  ErrAbsentKey. State/storage field helpers own encoding, not account clearing
  or execution rules. Fork envelopes bind activation/framing and receipt type,
  not signatures, transaction fields or execution results. These correspond to
  the ten retained decisions and pinned Yellow Paper, execution-specs and EIPs.
- Mutation/claim constructors without Limits copy caller-sized input before
  later validation. Loaded reads use encoded-node bounds, not mutation-value
  ceilings. Parser/hash phases and borrowed providers remain synchronous;
  cancellation cannot forcibly preempt them or enforce a heap/RSS deadline.
- Snapshot node/path/value ownership is immutable and structural sharing is
  intentional. Pending accounting subtracts replacements before additions;
  pressure/depth compaction retains shared nodes and conservatively keeps an
  unresolved store frontier. Pending and materialized ceilings apply separately,
  not to process memory across every retained snapshot. Builders charge lifetime
  encoding/hash budgets and commit state only after a successful call. Their
  depth policy admits the uncompressed key frontier, unlike compact-node
  traversal; a key fitting RawTrie can exceed the builder frontier limit.
- Load records a root and borrows a reader without eager graph validation.
  Concurrent use requires a safe reader and explicit historical retention.
  Commit prepares owned handles and remains source-store bound; Rebuild verifies
  root equality and detaches source availability. Recovery admits one canonical
  hash-addressed overlay node, not complete reachability or an automatic fetch.
- Store interfaces own context/concurrency/lifetime and publication semantics.
  Memory publishes immutable maps optimistically; it can refuse stale state
  and has process-local leases. Its iterators hold stable map snapshots and
  invoke callbacks unlocked. Callback errors propagate unchanged, possibly after
  earlier deliveries; default package diagnostics do not redact caller errors.
- Filesystem Open checks metadata/checksums and filenames, not all reachable
  nodes. One process exclusively owns the directory and trusted ancestors; no
  cross-process lock or race-free hostile-path sandbox is claimed. Existing
  ancestry permissions remain caller-owned. Filesystem iteration snapshots
  filenames, and a later prune may abort a partial delivery. Close refuses busy
  mutations/repeated closing without draining admitted reads; Root remains usable.
- Node files sync before ROOT replacement. A post-rename directory-sync error
  exposes the new complete root in process without confirming crash durability;
  callers must reconcile before retrying. Retention error/partial-handle results
  likewise require reopen/inventory. Prune rename is its logical commit point;
  later cleanup/sync failures return committed accounting, while earlier failure
  restores staged nodes. Restored lease handles alias the same durable lease.

## Assurance and immutable inputs

All 80 assurance files were inspected for entrypoints, predicates/assertions,
fixtures, private seams and resource ownership; no existing test body was changed.
Every fuzz target now identifies input limits and its actual oracle. Hostile-input
proof/recovery/value targets that discard errors establish panic safety only,
not universal rejection; model, canonical replay and surplus-node targets retain
those stronger specific assertions. Benchmark comments distinguish setup from
owned-byte lookup, mutation/construction, proof, warm local files and synced commit
work. No benchmark result or new cryptographic/conformance claim is produced.

Four interoperability-tagged sources retain Geth's separate alternate module
and pinned EthereumJS oracles; production has no tag/platform variants. No Go
source is generated. The two fixture update scripts were inspected for pinned
archive/source/checksum and license handling, not executed. Official fixture
bytes, local regression provenance, producer scripts, API baseline, dependency
pins, ten decision records/history prefixes and dated evidence remain unchanged.

## Missing versus unresolved

Structural omissions and misleading public contracts are repaired as comments.
Known limitations remain source constraints rather than newly passing policy:
cooperative cancellation, constructor copies before operation admission, lazy
root validation, post-publication filesystem errors, filename-only iteration and
historical availability dependent on retained provider nodes. The desired
MPT-DEC-009 failure/publication language is retained; provider-defined post-rename
error reconciliation is documented without redefining consensus or desired policy.
The local EIP-1186 fixture manifest names `@ethereumjs/trie`, whereas current
producer scripts/provenance use `@ethereumjs/mpt` from the pinned repository.
Historical generator identity is not reconstructed or relabeled by this source
pass; the frozen manifest is retained. No broad crypto, hostile-input, crash,
race, peer or benchmark campaign is scoped.

## Selected local verification

Go 1.27.2 formatting covers all 107 tracked Go files. Native parsing and
comparison preserve their AST, executable tokens and literals, including all
13 build/security annotations; no implicit-semicolon difference was needed.
All 386 comment groups render as HTML and all four production packages render
through `go doc -all`. The configured API mode is apidiff, so no full prose
snapshot is regenerated or immutable API baseline altered.

The native root documentation route checks structure, spelling and links; it
has no configured example operation. Its initial link pass reports zero errors.
The eight existing examples were separately executed with
`go test ./... -run '^Example' -count=1 -timeout=2m` and passed. The native
offline specification check accepts all ten unchanged decisions and histories.
No other runtime, interoperability, fuzz, benchmark or crash campaign was run.

## Objective Enforcement

The checksum-pinned public Tools v2.4.0 verifier now supplies
`make source-docs-check` as a prerequisite of the existing local and CI Make
routes. Ordinary shared CI selects strict source-documentation checking and
requires its result. Native prerequisites, operation budgets, security and
release controls are retained. Hosted exact-head evidence is recorded by the
coordination ledger; configuration alone does not establish green delivery.
