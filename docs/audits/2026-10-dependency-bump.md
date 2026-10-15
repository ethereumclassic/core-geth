---
description: "The v1.13.1 dependency bump: Go 1.26.9, golang.org/x/net v0.60.0 and its cascade, and the documentation-build security updates, measured against the current Go vulnerability database, and the standing practice that keeps the line current."
---

# Dependency bump: v1.13.0 to v1.13.1

This document records the dependency changes between `v1.13.0` and `v1.13.1`, and the standing practice
that produces them. It is a companion to
[dependency and toolchain modernization](2026-08-dependency-modernization.md), which covers the larger
move from the December 2024 archive state to `v1.13.0`, and to
[Go toolchain](2026-09-go-toolchain.md), which measures the Go standard library advisories each archive
carries. Where that modernization was a one-time overhaul, this is the routine currency work a
supported line takes on each release.

**Baseline:** `v1.13.0`, built with go1.26.8, released September 14, 2026, with about six months of
testing on Ethereum Classic mainnet and Mordor behind it and a month in production since, no critical
issue found.

**Agentic auditing methods and tooling:** the strongest available Claude models at each stage of the work, October 2026

**Work carried out by:** [White B0x](https://whiteb0x.com)

**Work funded by:** [Ethereum Classic DAO](https://ethereumclassicdao.org)

**Last edited:** October 12, 2026

## In short

**Upgrade to [the latest v1.13 release](https://github.com/ethereumclassic/core-geth/releases/latest).**
`v1.13.1` moves the build from go1.26.8 to **go1.26.9**, the Go security release of October 8, 2026, which
carries the October 2026 Go standard library batch. It moves `golang.org/x/net` to **v0.60.0**, which
carries the same HTTP/2 fixes on the module side, and the dependency graph follows. It takes two
security updates in the documentation build. A go1.26.9 build of the line carries **no Go standard
library advisory**, at module or symbol level, against the current database; the go1.26.6 artifact
measured for comparison carries 12 at module level. This is the supported-Go path doing its routine
job: a Go security release reaches the node on the next build. Bumping dependencies on each release,
rather than letting the client decay and ship vulnerable ones, is the maintenance commitment behind
the `v1.13` line. The earlier `v1.12.x` line did the opposite, and the
[March 2026 security audit](2026-03-security-audit.md) and [Go toolchain](2026-09-go-toolchain.md)
measure what that cost.

## The toolchain

| Surface | `v1.13.0` | `v1.13.1` |
|---|---|---|
| `go.mod` directive | `go 1.26.0` | `go 1.26.0` (unchanged) |
| Release workflow `go-version` | `1.26.8` | `1.26.9` |
| Docker builder | `golang:1.26.8-alpine` | `golang:1.26.9-alpine` |
| `build/checksums.txt` (`version:golang`) | `1.26.8` | `1.26.9` |

The `go.mod` directive is the language version and does not move; the three release pins, which decide
what actually builds the published archives, move together, as
[Go toolchain](2026-09-go-toolchain.md) requires. `build/checksums.txt` carries the SHA-256 of each
go1.26.9 toolchain archive that the `-dlgo` download path verifies before use.

**go1.26.9 carries the October 2026 Go batch.** [Go toolchain](2026-09-go-toolchain.md#finding-the-v11220-to-v11223-binaries-carry-go-standard-library-advisories)
records that batch of 13 advisories, published October 8, 2026 and fixed in go1.26.9. The published `v1.13.0`
archives were built with go1.26.8 and carried none when the release shipped; go1.26.8 predates the fix,
so a scan today matches 12 of the batch against it at module level. `v1.13.1` is built with go1.26.9 and
matches none. Go 1.26 is supported until go1.28.0 is released.

## The module graph

**`golang.org/x/net` moves to v0.60.0.** Five of the October `net/http` advisories match both the Go
standard library and `golang.org/x/net`, because the `x/net` HTTP/2 code is a second copy of the same
source. The toolchain bump clears the standard library side; the `x/net` bump clears the module side.
Both are needed to clear the match entirely.

**Its direct requirements pull the rest of the `golang.org/x` graph up by one patch.** `x/net` v0.60.0
requires newer `x/crypto`, `x/sys` and `x/text`, and resolving the graph to a minimal set moves `x/mod`
and `x/sync` as well. Each is an upgrade; none is deprecated or retracted.

| Module | `v1.13.0` | `v1.13.1` |
|---|---|---|
| `golang.org/x/net` | v0.58.0 | v0.60.0 |
| `golang.org/x/crypto` | v0.56.0 | v0.57.0 |
| `golang.org/x/text` | v0.41.0 | v0.42.0 |
| `golang.org/x/sys` | v0.47.0 | v0.48.0 |
| `golang.org/x/mod` | v0.40.0 | v0.41.0 |
| `golang.org/x/sync` | v0.22.0 | v0.23.0 |

What each change clears, confirmed against the built binary:

| Advisory | Package | Cleared by |
|---|---|---|
| GO-2026-6599, GO-2026-6600 | `html/template` | go1.26.9 |
| GO-2026-6604 (Windows only) | `os` | go1.26.9 |
| GO-2026-6605, GO-2026-6608, GO-2026-6609, GO-2026-6613 | `net/http`, `net/textproto`, `mime/multipart` | go1.26.9 |
| GO-2026-6607 | `crypto/tls` | go1.26.9 |
| GO-2026-6603, GO-2026-6610, GO-2026-6611, GO-2026-6612, GO-2026-6617 | `net/http` and `golang.org/x/net` | go1.26.9 and `x/net` v0.60.0 together |
| GO-2026-6354, GO-2026-6355 | `golang.org/x/crypto/ssh` | already current at v0.56.0; not compiled into any binary |
| GO-2026-6179, GO-2026-6180 | `golang.org/x/mod/sumdb` | already current at v0.40.0; not compiled into any binary |

`golang.org/x/crypto/ssh`, `golang.org/x/crypto/openpgp` and `golang.org/x/mod/sumdb` are not compiled
into any `v1.13` binary, so the advisories against them do not reach a node; they are tracked here
because a scanner reports them at module level.

## The documentation build

`requirements-mkdocs.txt`, which builds the documentation site and no part of the client, takes two
security updates. Neither is in any `geth` binary.

| Package | Before | After | Advisory |
|---|---|---|---|
| `urllib3` | 2.7.0 | 2.8.0 | GHSA-8988-9cw3-xx77 (HTTPS proxy TLS confusion), high, with two further fixes in the same release |
| `GitPython` | 3.1.60 | 3.1.62 | GHSA-59cr-6r3x-644w (path traversal through a submodule path), moderate |

Both arrived as automated dependency pull requests and were confirmed against the advisory records and
the release notes before being taken.

## What a vulnerability scanner reports against v1.13.1

**Standard library: none.** A go1.26.9 build of the `v1.13.1` line reports no Go standard library
advisory, at module level or symbol level, against the database recorded under
[Verification](#verification).

**Six go-ethereum advisories, all structural, as in `v1.13.0`.** The Go module path is deliberately
`github.com/ethereum/go-ethereum`, so a scanner compares this client's version against go-ethereum's
version ranges and matches six advisories whose fixes are present in the code: GO-2024-2819,
GO-2026-4314, GO-2026-4315, GO-2026-4507, GO-2026-4508 and GO-2026-4511.
[Go toolchain](2026-09-go-toolchain.md#what-a-vulnerability-scanner-reports-against-v1130) explains the
match and confirms each guard. They are the same six `v1.13.0` matches; this bump neither adds nor
removes any.

**One `golang.org/x/crypto` advisory at module level, with no fix.** GO-2026-5932 marks
`golang.org/x/crypto/openpgp` unmaintained. It has no fixed version, and the package is not compiled
into any binary, so no symbol-level finding reaches a node.

## The standing practice

This bump is routine, not an incident response, and the practice that produces it is the point of the
`v1.13` line.

- **A supported Go toolchain.** Go 1.26 receives security releases, and a client release built with the
  newest one takes them. The `v1.12.x` line sat on Go majors that left support in 2024 and 2025 and
  never will; [Go toolchain](2026-09-go-toolchain.md) measures the difference, which only grows.
- **A weekly scan.** `.github/workflows/govulncheck.yml` scans `main` every week, in source mode and
  against a `geth` binary built with the release Go version, and keeps one issue current with anything
  a maintainer needs to see. It is what surfaced the October batch against the shipped go1.26.8 binary,
  which this bump answers.
- **Version discipline on every change.** A dependency is taken only at a version that is a stable
  release, old enough to have been exercised, and not deprecated or retracted, and the exact fix
  version named in an advisory is preferred to a speculative jump to latest. Each bump is built and
  scanned before it lands.
- **Measured, not assumed.** The figures here were read from the built binary with govulncheck, the
  same tool and method the other audits use, so the result can be reproduced rather than taken on
  trust.

**The alternative is decay, and its cost is on the record.** The `v1.12.x` line in
[`etclabscore/core-geth`](https://github.com/etclabscore/core-geth) was left on Go toolchains that
reached end of life in 2024 and 2025, and on a module graph unchanged since December 2024, through 21
months without security maintenance. [Go toolchain](2026-09-go-toolchain.md) counts the Go standard
library advisories its binaries carry, a total that only grows, and the
[March 2026 security audit](2026-03-security-audit.md) records the flaw that was exploited against
the Core-Geth `v1.12` bootnodes while that line sat unmaintained. Backporting a fix into it closed
individual CVEs and never the toolchain or the dependency graph underneath them. The `v1.13` line is
the one where the client's own CVEs, the Go standard library and the module graph are current at once,
so its binaries are the complete resolution.

## Verification

Each check below was calibrated so that it could report a negative.

- **The toolchain.** A `geth` binary built with go1.26.9 from the release branch records `go1.26.9`,
  read with `go version -m` without running it, alongside `golang.org/x/net` v0.60.0 and the other
  module versions in the table above.
- **The download path.** The build used `build/ci.go install -dlgo`, which downloads the Go toolchain
  named in `build/checksums.txt` and verifies it against the SHA-256 written there before building, so
  the go1.26.9 checksums were exercised rather than only copied. The go1.26.9 source tarball was hashed
  independently and matched the go.dev value.
- **The advisories.** govulncheck v1.8.0 in binary mode, at module level and at symbol level, against
  `vuln.go.dev` as modified 2026-10-08T22:31:09Z, 4,613 entries, scanned October 11, 2026. The go1.26.9 build
  reports 0 Go standard library advisories at both levels, 6 go-ethereum advisories at symbol level, and
  GO-2026-5932 at module level only.
- **Calibration.** The go1.26.6 artifact measured for comparison reports 12 Go standard library
  advisories at module level against the same database, 9 of them with the code in the binary; the
  go1.26.9 build reports none.
- **Lint.** `make lint` reports no issue, as expected of a change that touches manifests and pins and no
  Go source.

**Not measured:**

- The published `v1.13.1` archives, which did not exist when this was written; the figures are from a
  build of the release branch, and the published archives are to be confirmed the way
  [Go toolchain](2026-09-go-toolchain.md#verification) confirms a release once `v1.13.1` ships.
- Reachability of any advisory from the client's code paths.

## Supporting this work

The client is maintained by Ethereum Classic's core developers and open community contributions, in
the [`ethereumclassic`](https://github.com/ethereumclassic) GitHub organization. The Ethereum Classic
DAO is the legal entity standing behind that work, and the floor under it when nobody else is paying.
Mining pools, centralized exchanges, issuers of Ethereum Classic financial products, Etchash mining
hardware manufacturers and large holders all depend on this client. If your operation relies on
Ethereum Classic, please consider sponsoring its maintenance. [Support this work](../support.md) says
who maintains and funds it and how to sponsor it.
