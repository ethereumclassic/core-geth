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

**The linter pin moved separately, and last.** `version:golangci` went from
1.55.2 to 2.13.1 in its own commit, after the dependency work, for the reason the
deferral anticipated: a newer golangci-lint reports findings the pinned one does
not, and mixing that into a dependency diff would obscure both. See "Linter" below
for what it found.

The old pin was not merely stale. golangci-lint 1.55.2 cannot read Go 1.26
standard-library export data, so it failed on every package with a `typecheck`
error and reported no real findings at all. A linter that cannot parse the tree
reports clean by failing, which is the same shape as the end-of-life toolchain
problem this document opens with.

## Dependency delta

Measured against the baseline commit with `go mod edit -json`: **83 modules
changed, 17 removed, 16 added.** The manifest went from 176 modules (83 direct,
93 indirect) to 175 (82 direct, 93 indirect).

**No direct dependency was added. All 16 additions are indirect**, and that is the
number worth noting: this was a currency and security pass, not a change in what
the client depends on.

Seven of the 16 are not new code at all but upstream import-path renames, each
pairing with an entry in the removed list:

| Removed | Successor |
|---|---|
| `StackExchange/wmi` | `yusufpapurcu/wmi` |
| `cockroachdb/sentry-go` | `getsentry/sentry-go` |
| `deepmap/oapi-codegen` | `oapi-codegen/runtime` |
| `go-fonts/liberation` | `codeberg.org/go-fonts/liberation` |
| `go-latex/latex` | `codeberg.org/go-latex/latex` |
| `go-pdf/fpdf` | `codeberg.org/go-pdf/fpdf` |
| `rivo/uniseg` | `clipperhouse/displaywidth`, `clipperhouse/uax29/v2` |

The remaining nine arrived transitively behind modules that were updated:
`apapsch/go-jsonmerge/v2`, `cockroachdb/fifo`, `emicklei/dot`, `olekukonko/cat`,
`olekukonko/errors`, `olekukonko/ll`, `go.mongodb.org/mongo-driver` and
`go.yaml.in/yaml/v3`.

### Direct dependencies changed

| Module | Before | After |
|---|---|---|
| `golang.org/x/crypto` | v0.17.0 | v0.55.0 |
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

Ten modules left with no successor: `consensys/bavard` · `mmcloughlin/addchain` ·
`rsc.io/tmplfunc` · `opentracing/opentracing-go` · `google/go-cmp` ·
`fjl/memsize` · `campoy/embedmd` · `hashicorp/go-cleanhttp` ·
`hashicorp/go-retryablehttp` · `pmezard/go-difflib`

The other seven were renames rather than removals; see the table above.

All ten are unreferenced in source. `bavard`, `addchain` and `tmplfunc` were
`gnark-crypto`'s code-generation dependencies, which its current release no longer
requires. **`fjl/memsize` is the one removed for cause rather than by tidy:** it
reaches into runtime internals through `//go:linkname`, which Go 1.26 tightened,
and the link failed until it was dropped. The rest fell out of `go mod tidy`
during the toolchain upgrade. None was removed by hand.

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

## Linter

Raising `version:golangci` from 1.55.2 to 2.13.1 required migrating
`.golangci.yml` to the v2 schema, and the migration changed what runs in two ways
that were not intended and had to be undone.

**golangci-lint v2 folded `gosimple`, `stylecheck` and a new `quickfix` family
into `staticcheck`.** This tree enabled the first two under v1 and never enabled
`stylecheck`; `quickfix` did not exist. The migration therefore switched on 34
style findings nobody had opted into, most of them in code inherited from
upstream. `ST*` and `QF*` are disabled again, so the check set matches what was
in force before.

**`exportloopref` was replaced by `copyloopvar`,** which has a much wider surface:
it reports every now-unnecessary `x := x` loop-variable copy, obsolete since Go
1.22 changed loop scoping. Those 35 were taken rather than suppressed — the copies
are dead weight in a tree targeting Go 1.26 — along with the comments explaining
them, which described a copy that no longer existed.

**Several existing exclusions had silently stopped matching.** staticcheck v2
reports fully-qualified names, so patterns written against `event.TypeMux is
deprecated` no longer matched `github.com/ethereum/go-ethereum/event.TypeMux is
deprecated`. They were rewritten to match loosely. A broken exclusion is worse
than a missing one: it looks like a considered decision while doing nothing.

Fixed rather than excluded:

- **`crypto/bn256/cloudflare/gfp_decl.go`** — `// go:noescape` carried a space, so
  the compiler directive was inert and `gfpNeg` silently lost it. Its three
  siblings in the same file were correct.
- **`core/vm/evmc.go`** — an assignment to `status` that every path below returns
  past.
- **`reflect.PtrTo` to `reflect.PointerTo`** across `rlp/`, deprecated since Go
  1.22. These were masked by the `govet` `inline` analyzer, which reports where a
  call could not be inlined — a compiler diagnostic, not a defect, and one that
  fired only because the linter binary is built against a newer Go than this
  module targets. Disabling it surfaced the real findings underneath.
- An `S1009` redundant nil check, an `S1011` append loop, and assorted stray
  blank lines.

**What is excluded, and why it is a deferral rather than a dismissal.** Go 1.26
deprecated the raw `big.Int` fields on `crypto/ecdsa` keys and the low-level
`crypto/elliptic` operations, directing callers to `crypto/ecdh`. That package
covers the NIST curves only, and this client is secp256k1 throughout, so for most
of the flagged sites there is no replacement to migrate to. Two of them are the
CVE fixes recorded in the companion documents: the `IsOnCurve` calls in
`crypto/crypto.go` and `crypto/ecies/ecies.go` answer CVE-2025-24883 and
CVE-2026-26315, where the on-curve check is the fix. The sites in
`p2p/discover/v4wire` and `p2p/enode/urlv4` encode the devp2p discovery wire
format, which is fixed by the protocol rather than chosen here.

`runtime.GOROOT` and `go/parser.ParseDir` are likewise excluded in
`internal/build/` and `cmd/geth/misccmd.go`. Both are correct there: that code
runs from the source tree during a build, which is the case their replacements do
not cover.

`make lint` reports zero issues on the resulting tree.

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
- The `crypto/ecdh` migration for the deprecated `crypto/ecdsa` and
  `crypto/elliptic` call sites, once a secp256k1 path exists or upstream
  go-ethereum moves first. Scoped per package rather than tree-wide: the keystore
  and scwallet sites are the most tractable, the wire-format ones the least.
- Whether to adopt the `stylecheck` and `quickfix` families that the golangci-lint
  v2 migration switched on and this pass switched back off. They are reasonable
  checks; taking them means a large mechanical diff against code inherited from
  upstream, so it is a trade to weigh rather than an oversight.

Dependency automation is configured and deliberately disabled; `.github/dependabot.yml`
and `AGENTS.md` carry that decision and its rationale. Go's own tooling —
`go list -m -u all` for retractions and `govulncheck` for advisories — is what is
run against this repository, by hand.
