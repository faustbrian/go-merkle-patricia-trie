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

## Adding A Module

Follow [repository structure policy](AGENTS.md#repository-structure). New modules
require an explicit purpose, ownership boundary, dependency review, package
catalog entry, full quality gates, documentation, changelog, license, security
policy, compatibility plan, and release dry-run.
