# Contributing

## Before Editing

1. Read [`AGENTS.md`](AGENTS.md) and the affected package's docs.
2. Run the narrow standalone baseline for the affected package.
3. Identify owned dependencies and reverse dependants in `modules.json`.
4. Preserve unrelated work and generated/corpus provenance.

The repository keeps package-specific verification inputs under source-owned
paths. Shared repository checks, tool versions, evidence handling, and CI are
provided by the pinned `go-library-tools` release rather than copied scripts.

## Changes

Keep commits focused and conventional. Update every affected changelog with
the behavior and migration impact. Public API changes require compatibility
evidence and documentation. Specification behavior requires a decision record,
fixture coverage, and interoperability evidence.

New direct dependencies and dependency updates must follow the
[dependency governance policy](AGENTS.md#dependencies-and-supply-chain). Package-local
update bots are forbidden; the root policy owns every module and action update.

Specification-backed changes must follow the
[specification governance contract](AGENTS.md#design), update
the affected stable entries in the
[specification decision register](docs/specification-decisions.md), and complete
the Specification Decisions section of the pull request template. An unresolved
interpretation or stale source pin is release-blocking; peer behavior cannot
silently select policy.

Selected mutation gates must finish with zero surviving viable mutants.

Do not add package-local workflows, permanent replacements, machine-specific
paths, bypass flags, broad mutation exclusions, or aggregate quality metrics
that hide a failing package.

## Verification

Use Go 1.27.2 for development and ordinary CI; the public minimum remains
Go 1.27.0. A standalone checkout supports ordinary tests and a focused active
fuzz smoke without sibling repositories or goal ledgers:

```bash
GOWORK=off go test ./...
GOWORK=off go test . -run '^$' -fuzz '^FuzzProofVerificationRejectsHostileInput$' -fuzztime=2s
```

Select the native fuzz target affected by your change. Ordinary tests exercise
deterministic fuzz seeds and corpus; base-free local checks do not start an
active fuzz campaign. Install the public `golib` v2.1.1 release declared in
`.golib.yaml` and verify its checksums for the shared routes:

```bash
make check
make ci
```

Both Make routes are bounded local checks, including one filesystem race run
that asserts writer serialization, nonblocking readers and close refusal during
commit. The module declares no services or test tags; filesystem tests use
temporary local directories. Ordinary CI uses the published verifier without
building it from source; changed-source CI can select a bounded active fuzz
smoke. Optional mutation analysis can select one package:

```bash
golib mutation --module . --package ./filesystem
```

Use `make check-full` for explicit full assurance, including online specification
validation, the unchanged 10000x fuzz targets, full race/coverage/mutation,
API, conformance, peer interoperability and benchmarks. Encoding, commitment
and proof changes need the pinned Ethereum fixtures and peer comparisons;
storage changes need directly affected retention, recovery and publication
assertions. The explicit interoperability target is
`make -f verification/package.mk interoperability`; it requires Node.js/npm and
the pinned EthereumJS packages and uses the owned `go.interop.mod` Geth harness
with the `interoperability` tag, not a sibling checkout. See
[benchmark reproduction](docs/benchmarks.md) for selected performance work.
Full campaigns remain available for selected scheduled or release qualification;
ordinary scheduled CI remains bounded. Release rehearsal and root qualification
are explicitly selected, not ordinary CI results. Report unavailable or failing
checks without claiming release readiness. CI/configuration-only changes do not
require a library release.

## Selected Weekly Assurance

[`assurance.yml`](.github/workflows/assurance.yml) runs on main each Saturday
at 05:17 UTC. It selects one of the 21 declared root-module fuzz targets using
the shared UTC-week selector, not all targets each week. Explicit main-only
dispatch accepts one exact declared selector, for example
`^FuzzCanonicalNodeDecodeRoundTrip$`.

The selected risks include canonical node/RLP/path/record encodings, bounded
hostile proof and recovery inputs, trie mutation models, iteration cancellation
and retention/pruning state. Existing round-trip targets assert canonical
bytes; model targets compare stored or iterated values with their reference.
Hostile proof targets exercise panic safety but discard verification results:
their pass does not prove that every invalid proof is rejected. These are not
blanket cryptographic, root/hash, Ethereum conformance or deployment claims.
Retain affected pinned fixtures, peer interoperability, filesystem concurrency,
public API and supported/experimental contract checks.

The immutable shared runner bounds active fuzzing to 30 seconds, two workers
and two CPUs, with a five-minute combined native session, 2 GiB memory and
1 MiB output limits inside a ten-minute job. It installs the declared public
CLI without a source build, uses private disposable caches and keeps failures,
timeouts and exhausted budgets visible. Ordinary deterministic tests still run
seeds/corpus; ordinary CI and its Required contract are unchanged. Full/manual
operations retain all 21 targets at 10000x/2m each, conformance, interoperability
and benchmark budgets in `.golib.yaml`. The selected route needs no live service
or test tag; it does not substitute for the Node.js/npm and tagged Geth harness
used by explicit interoperability qualification.

Treat findings as owner-attributed diagnostics: promptly repair viable defects
with focused regressions; retain reviewed narrow equivalent/invalid dispositions
rather than weakening tests. Review cadence against useful findings and actual
cost; do not interpret a selected pass as full or release qualification.

## Adding A Module

Follow [repository structure policy](AGENTS.md#repository-structure). New modules
require an explicit purpose, ownership boundary, dependency review, package
catalog entry, full quality gates, documentation, changelog, license, security
policy, compatibility plan, and release dry-run.
