# Security Policy

## Reporting

Do not open a public issue for a suspected vulnerability. Report it privately
through [GitHub Security Advisories for
`faustbrian/go-merkle-patricia-trie`](https://github.com/faustbrian/go-merkle-patricia-trie/security/advisories/new)
before public disclosure. Do not include credentials, private fixtures, or
affected deployment information in a public report.

Include the affected module and version, impact, reproduction, preconditions,
and any suggested mitigation. Reports are acknowledged as soon as practical;
timelines depend on severity and verification.

## Supported Versions

The latest stable `v1` release receives security fixes. Older releases and the
`main` branch are not supported; upgrade before reporting. At and after v1,
incompatible exported API or documented behavior changes require a new major
version under the [compatibility policy](COMPATIBILITY.md) and
[deprecation policy](DEPRECATION.md).

## Security Gates

Releases require isolated tests, race and hostile-input checks, exact coverage
and mutation results, `govulncheck`, secret scanning, license verification,
SBOM generation, provenance validation, and clean-consumer resolution. A
missing scanner or unavailable service is a failed gate, not a warning.

Security fixes MUST include a regression test that does not publish weaponized
details or real secrets. Credentials MUST be redacted from logs and evidence.

## Repository Assurance

The repository [safety and concurrency policy](AGENTS.md#safety-and-concurrency)
and [supply-chain policy](AGENTS.md#dependencies-and-supply-chain) define shared
trust boundaries and release requirements. Package-specific security guidance
refines those rules for its owned boundary.
