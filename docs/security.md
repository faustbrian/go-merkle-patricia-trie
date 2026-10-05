# Security

## Trust boundary

Keys, values, compact paths, encoded nodes, proof nodes, storage responses, and
fixture updates are untrusted. Stores may omit, corrupt, substitute, replay, or
partially persist nodes. Contexts, explicit byte/count/depth/hash/read/write
limits, and checked arithmetic bound work before conversion or allocation.

## Threats

The implementation must defend against profile or double-hash confusion; RLP
and compact-path malleability; the 31/32-byte reference boundary; proof
substitution, truncation, reordering, duplication, and surplus nodes; cycles
and path amplification; corrupt storage; stale or partial root publication;
unsafe pruning across shared roots; preimage disclosure; integer and allocation
overflow; CPU, memory, storage-read, and disk amplification; slice aliasing and
concurrent mutation; secret leakage in errors; and compromised dependencies or
fixtures.

Errors expose bounded metadata and typed causes, never complete keys, values,
nodes, proofs, preimages, or credentials.

## Proof limitation

A valid proof establishes a key/value, absence, or exact range-completeness
claim under the supplied root. It does not establish that the root is
canonical, finalized, recent, or authorized.

The complete family threat model, bound ownership, mitigations, and review
triggers are recorded in [the threat model](threat-model.md). Backend causes
remain explicitly inspectable but their default storage-error strings are
redacted; callers must not turn cause inspection into automatic sensitive logs.

## Reviewed mutation dispositions

The [exact equivalent-mutant inventory](../.verification/mutation/equivalent-inventory.json)
retains five source-proven limitations of the pinned Gremlins verifier, not a
runtime or release pass. Independent reviewer
`postgres_safe_config_v2_final_review` reviewed the unchanged production at
`4ffaad7b09b42fd8f01f3ef3850c1177addc3678` on 2026-10-04. The dispositions
expire on 2027-01-02, or immediately upon a named source, verifier, dependency
invariant, or caller-domain change. Maintainers must reassess them before use
after expiry or such a change. Native `LIVED` results remain visible; only
the five exact coordinates and contract domains are selected. All unreviewed
survivors and incomplete statuses remain failures.
