# Dependency and toolchain modernization — v1.12.x archive to v1.13.x

This document records the dependency and toolchain changes that take Core-Geth
from its December 2024 archive state to the `v1.13.x` line. It is a companion to
`2026-03-security-audit.md` and `2026-08-security-followup.md`, which cover the
source-level CVE remediation; this one covers what changed underneath the code.

**Baseline:** commit `7ef3ecd7a`, 2024-12-16, the last commit of the archived
`v1.12.x` development line.

## Why this work was necessary

The archived line was pinned to Go 1.21 across every build surface. Under Go's
support policy — *"each major Go release is supported until there are two newer
major releases"* — Go 1.21 left support in 2024, so every artifact built from that
line shipped a standard library receiving no security fixes.

That failure mode is silent. An end-of-life toolchain still compiles, still passes
tests, and emits no deprecation warning; nothing in an ordinary build reports it.
The same applies to a dependency that is merely old rather than known-vulnerable.
Neither condition announces itself, which is why this pass was scheduled rather
than triggered by an incident.

## Toolchain

| Surface | Before | After |
|---|---|---|
| `go.mod` directive | `go 1.21` | `go 1.26` |
| CI workflows | `1.21` | `1.26` |
| Docker builder | `golang:1.22-alpine` | `golang:1.26-alpine` |
| `build/checksums.txt` (`version:golang`) | `1.22.1` | `1.26.6` |

There is no `toolchain` directive; the `go` directive is the whole statement.

Three Go versions were in force simultaneously on the archived line — 1.21 in CI,
1.22.1 for the `-dlgo` download path, and 1.22 in the Docker builder. They now
agree.

**`build/checksums.txt` carries a second Go pin, `version:ppa-builder`, which was
deliberately left alone.** It is read only by the Debian source-package path,
which no active workflow invokes. It is documented in `AGENTS.md` as knowingly
insufficient: bootstrapping a modern Go from source requires a compiler two majors
back, and closing that gap needs the recursive builder the file's own comment
anticipates rather than a version bump.

**The linter pin was not moved.** `version:golangci` remains at 1.55.2. Raising it
is a separate change with its own diff, since a newer golangci-lint reports
findings the pinned one does not, and mixing that into a dependency pass would
obscure both.

## Dependency delta

Measured against the baseline commit: **26 modules changed, 6 removed, none
added.** The manifest went from 176 modules (83 direct) to 170 (82 direct).

Nothing was added, which is the number worth noting: this was a currency and
security pass, not a change in what the client depends on.

### Direct dependencies changed

| Module | Before | After |
|---|---|---|
| `golang.org/x/crypto` | v0.17.0 | v0.54.0 |
| `golang.org/x/sys` | v0.16.0 | v0.47.0 |
| `golang.org/x/tools` | v0.15.0 | v0.48.0 |
| `golang.org/x/text` | v0.14.0 | v0.40.0 |
| `golang.org/x/sync` | v0.5.0 | v0.22.0 |
| `golang.org/x/time` | v0.3.0 | v0.15.0 |
| `github.com/consensys/gnark-crypto` | v0.12.1 | v0.21.0 |
| `github.com/crate-crypto/go-kzg-4844` | v0.7.0 | v1.1.0 |
| `github.com/supranational/blst` | v0.3.11 | v0.3.17 |
| `github.com/holiman/uint256` | v1.2.4 | v1.3.2 |
| `github.com/btcsuite/btcd/btcec/v2` | v2.2.0 | v2.5.0 |
| `github.com/graph-gophers/graphql-go` | v1.3.0 | v1.10.2 |
| `github.com/gorilla/websocket` | v1.5.0 | v1.5.3 |
| `github.com/golang-jwt/jwt/v4` | v4.5.0 | v4.5.2 |
| `github.com/golang/protobuf` | v1.5.3 | v1.5.4 |
| `github.com/stretchr/testify` | v1.8.4 | v1.11.1 |
| `github.com/tidwall/gjson` | v1.6.0 | v1.18.0 |

Four of these carry consensus weight — `gnark-crypto`, `go-kzg-4844`, `blst` and
`uint256`. Each was read before being taken rather than accepted on version number
alone, and the elliptic-curve and KZG changes were checked against the consensus
test suites named under Verification below.

`golang-jwt/jwt` is an authentication boundary rather than a test dependency: it
backs Engine API JWT verification. Its shortest module path runs through test
tooling, which understates it.

### Indirect dependencies changed

`bits-and-blooms/bitset` v1.10.0 to v1.24.6 · `cespare/xxhash/v2` v2.2.0 to v2.3.0 ·
`decred/dcrd/dcrec/secp256k1/v4` v4.0.1 to v4.4.0 · `rogpeppe/go-internal` v1.9.0 to
v1.12.0 · `tidwall/match` v1.0.1 to v1.1.1 · `tidwall/pretty` v1.0.0 to v1.2.0 ·
`golang.org/x/mod` v0.14.0 to v0.38.0 · `golang.org/x/net` v0.18.0 to v0.57.0 ·
`google.golang.org/protobuf` v1.31.0 to v1.33.0

### Modules removed

`consensys/bavard` · `mmcloughlin/addchain` · `rsc.io/tmplfunc` ·
`opentracing/opentracing-go` · `google/go-cmp` · `fjl/memsize`

All six are unreferenced in source. The first three were `gnark-crypto`'s
code-generation dependencies, which its current release no longer requires; the
rest fell out of `go mod tidy` during the toolchain upgrade. None was removed by
hand.

## What deliberately did not move

Recording these matters as much as recording the changes, because an unexplained
old pin invites a later contributor to "fix" it.

**The AWS SDK v2 stack (twelve modules) is held at its archived versions**, which
match go-ethereum's own pins exactly. These are reachable only from `cmd/devp2p`,
the DNS node-list publishing tool, and not from the node itself. Advancing them
independently would diverge from upstream for a peripheral tool with no exposure in
the client.

**`c-kzg-4844` remains at v0.4.0.** No advisory affects it. A major-version bump on
a cgo binding is the highest-risk change available in the consensus tier, and
upstream go-ethereum has since moved to a different module path
(`c-kzg-4844/v2`), so catching up is an import-path migration rather than a version
change. That belongs in its own release.

**The Verkle modules — `go-ipa` and `go-verkle` — remain at their archived
versions.** Verkle trees are not live on any Ethereum network, including Ethereum
Classic, and both modules are pre-1.0 upstream.

**`crypto/secp256k1/libsecp256k1` is unchanged and is tracked as outstanding
work.** It is vendored C, not a Go module, so it is invisible to `go list`,
`go.sum` and `govulncheck` by construction — and on a normal cgo build it, rather
than the Go elliptic-curve libraries, is the active signature-verification path.
Its last content synchronization predates this fork's divergence by years.
Re-vendoring C that binds through cgo requires build and correctness verification
beyond the scope of a dependency pass, so it is deliberately not bundled here.

## Verification

Every change was gated on the consensus suites, not on a successful build:

```bash
go test ./params/... ./core/... ./consensus/... ./crypto/... -count=1
go test ./tests/  -run TestETCDifficultyBombVectors -count=1
go test ./core/   -run TestECBP1100Vectors          -count=1
make test
make test-coregeth
```

The ETC difficulty vectors bind each fork label to its own rule set, so a change in
elliptic-curve or integer-arithmetic behavior surfaces as a specific label failing
rather than as a single opaque total. Both networks were additionally synchronized
from genesis on the resulting build.

`govulncheck` is run in both source and binary mode. A set of advisories is
reported permanently and has been adjudicated; `AGENTS.md` carries the current set
and the reasoning. Two properties of that check are worth stating because they
invert the usual reading:

- **An exit status of 0 is a change, not a pass.** Findings are expected here.
- **An advisory disappearing is as much a finding as one appearing.** It means
  either the adjudication has gone stale or the artifact scanned is not the one
  that was adjudicated.

Most reported advisories are structural. This client retains go-ethereum's module
path, so its pseudo-version always sorts below any upstream fixed tag, and
`govulncheck` matches on symbol names while a backported fix adds guards inside the
same function. Each was confirmed by reading the guard in this tree rather than by
comparing version strings.

## Outstanding

- The vendored `libsecp256k1` resynchronization described above.
- The `version:ppa-builder` bootstrap chain, which needs a recursive builder rather
  than a version bump.
- The `version:golangci` linter pin, deferred so its findings do not land mixed
  into a dependency diff.

Dependency automation is configured and deliberately disabled; `.github/dependabot.yml`
and `AGENTS.md` carry that decision and its rationale. Go's own tooling —
`go list -m -u all` for retractions and `govulncheck` for advisories — is what is
run against this repository, by hand.
