---
description: "The dependency changes in Core-Geth v1.13.1: Go 1.26.10, golang.org/x/net v0.62.0 and other Go modules, security updates to the documentation build, the CI tooling and the Alpine 3.24.3 container base, with the security advisories each change fixes."
---

# Dependency bump: v1.13.0 to v1.13.1

This page lists the dependency changes between `v1.13.0` and `v1.13.1`: the Go toolchain, the Go
modules, the Python packages that build this documentation site, and the container base images. For
each change it records which security advisories it fixes and how the result was checked. An advisory
is a published record of a vulnerability that names the affected versions and the version that fixes
it. IDs that start with `GO-` come from the Go vulnerability database, and IDs that start with `GHSA-`
come from the GitHub Advisory Database.

Two earlier pages cover what came before.
[Dependency and toolchain modernization](2026-08-dependency-modernization.md) covers the move from the
December 2024 code base to `v1.13.0`, and [Go toolchain](2026-09-go-toolchain.md) counts the Go
advisories that affected each release's files when `v1.13.0` was released. Advisories published after
that date are covered here.

**Baseline:** `v1.13.0`, released September 14, 2026 and built with go1.26.8. It was tested for about
six months on Ethereum Classic mainnet and Mordor before release, and no critical issue has been found
in it since.

**Agentic auditing methods and tooling:** the strongest available Claude models at each stage of the work, October 2026

**Work carried out by:** [White B0x](https://whiteb0x.com)

**Work funded by:** [Ethereum Classic DAO](https://ethereumclassicdao.org)

**Last edited:** November 1, 2026

## In short

**Upgrade to [the latest v1.13 release](https://github.com/ethereumclassic/core-geth/releases/latest).**

Each release in the `v1.13` line is built with the latest patch release of its Go version as of the
day it is cut. `v1.13.0` was built with go1.26.8. As of November 1, 2026, `v1.13.1` is built with
go1.26.10, and every module in its `go.mod` that has a published security fix moves to a fixed
version.

Since `v1.13.0` was released, the Go team has published 14 security advisories for the Go standard
library that affect go1.26.8:

- 13 were published on October 8, 2026 and are fixed in go1.26.9. Most concern HTTP and HTTP/2
  handling in `net/http`.
- 1 was published on October 21, 2026 and is fixed in go1.26.10. With the fix, `crypto/x509` rejects
  intermediate certificates whose private keys are known to be compromised.

go1.26.10 includes the fixes from go1.26.9, so `v1.13.1` fixes all 14. It also includes:

- `golang.org/x/net` v0.62.0. This module has its own copy of Go's HTTP/2 code, and five of the
  October 8 issues need fixing there too. This version also fixes two advisories in the module's
  `html` package, which no `geth` binary includes. Other `golang.org/x` modules are updated with it.
- `golang.org/x/image` v0.45.0 and `github.com/klauspost/compress` v1.18.7, which fix six advisories in
  packages that no `geth` binary includes.
- Security updates to three Python packages that build this documentation site. They are not part of
  `geth`.
- Container images built on Alpine 3.24.3, which includes security fixes to musl (the C library),
  OpenSSL and the system's trusted certificates. Both base images are pinned by digest.
- Newer versions of the linter, the vulnerability scanner and five GitHub Actions used in CI and
  releases, one of them with a security fix.

**Result.** A scan of a `geth` binary built from `v1.13.1` reports no Go standard library advisory, as
of that date. The same scan of the published `v1.13.0` Linux binary reports 13 of the 14; the
other one applies only to Windows. When `v1.13.0` was released, the scan reported none, because none of
the 14 had been published yet.

Updating the Go toolchain and the dependencies with each release is how the `v1.13` line stays
current. The `v1.12.x` line was not maintained this way, and the
[March 2026 security audit](2026-03-security-audit.md) and [Go toolchain](2026-09-go-toolchain.md)
record the result.

## The Go toolchain

| Where the Go version is set | `v1.13.0` | `v1.13.1` |
|---|---|---|
| `go.mod` | `go 1.26` | `go 1.26.0` (the same version, written in full) |
| Release workflow (`go-version`) | `1.26.8` | `1.26.10` |
| Docker builder image | `golang:1.26.8-alpine` | `golang:1.26.10-alpine3.24` |
| `build/checksums.txt` (`version:golang`) | `1.26.8` | `1.26.10` |

The line in `go.mod` sets the Go language version. It stays at Go 1.26; only the way it is written
changes. The other three settings decide which Go release builds the published files, so they always
move together; [Go toolchain](2026-09-go-toolchain.md) explains why. `build/checksums.txt` also holds the SHA-256 checksum
of each go1.26.10 download, which the build checks before it uses the download.

Go 1.26 receives security fixes until go1.28.0 is released.

The go1.26.9 release also fixes two issues in the `go` command itself, GO-2026-6601 and GO-2026-6602.
They affect building software with Go, not the code inside a `geth` binary.

**The linter.** `build/checksums.txt` also pins golangci-lint, which moves from 2.13.1 to 2.15.0. The
linter checks the source code in CI and is not part of any binary.

## Go modules

**`golang.org/x/net` moves to v0.62.0.** The `golang.org/x/net` module has its own copy of Go's HTTP/2
code. Five of the October 8 advisories (GO-2026-6603, GO-2026-6610, GO-2026-6611, GO-2026-6612 and
GO-2026-6617) affect both copies, so fixing them needs both go1.26.9 or later and
`golang.org/x/net` v0.60.0 or later. Version v0.62.0 also fixes GO-2026-6860 and GO-2026-6866, published
on October 27, 2026, in the module's `html` package, which no package in this repository imports.

**Six related modules move with it.** `golang.org/x/net` v0.62.0 requires newer versions of
`golang.org/x/crypto`, `golang.org/x/sys` and `golang.org/x/text`, and `golang.org/x/text` v0.43.0
requires newer versions of `golang.org/x/mod`, `golang.org/x/sync` and `golang.org/x/tools`.
`golang.org/x/crypto` and `golang.org/x/mod` had already been updated once on `main` after `v1.13.0`.
Every change is an upgrade, and none of the new versions is deprecated or retracted.

| Module | `v1.13.0` | `v1.13.1` |
|---|---|---|
| `golang.org/x/net` | v0.58.0 | v0.62.0 |
| `golang.org/x/crypto` | v0.55.0 | v0.58.0 |
| `golang.org/x/text` | v0.41.0 | v0.43.0 |
| `golang.org/x/sys` | v0.47.0 | v0.49.0 |
| `golang.org/x/mod` | v0.39.0 | v0.41.0 |
| `golang.org/x/sync` | v0.22.0 | v0.24.0 |
| `golang.org/x/tools` | v0.49.0 | v0.51.0 |

**Three modules change where they come from.** `v1.13.0` used `replace` directives in `go.mod` to fetch
`github.com/etclabscore/go-openrpc-reflect`, `github.com/etclabscore/go-jsonschema-walk` and
`github.com/yusufpapurcu/wmi` from another location. In `v1.13.1`, `wmi` comes from its own repository,
and the two `etclabscore` modules come from `ethereumclassic/archive-reference-material`, at the same
commits.

**Two more modules move to versions that fix advisories.** Neither is compiled into any `geth` binary.

| Module | `v1.13.0` | `v1.13.1` | Why it is in the module list |
|---|---|---|---|
| `golang.org/x/image` | v0.41.0 | v0.45.0 | Used only by the MESS plotting tools in `core`, which build with `-tags plot` |
| `github.com/klauspost/compress` | v1.16.7 | v1.18.7 | Used only by builds without cgo, through the pebble database; `geth` is always built with cgo |

The `golang.org/x/image` advisories were published between April and August 2026, before `v1.13.0`
was released; the `github.com/klauspost/compress` advisory was published on October 20, 2026 in the
GitHub Advisory Database. They are listed in the table below.

### Which change fixes which advisory

| Advisory | Package | Fixed by |
|---|---|---|
| GO-2026-6599, GO-2026-6600 | `html/template` | go1.26.9 or later |
| GO-2026-6604 (Windows only) | `os` | go1.26.9 or later |
| GO-2026-6605, GO-2026-6608, GO-2026-6609, GO-2026-6613 | `net/http`, `net/textproto`, `mime/multipart` | go1.26.9 or later |
| GO-2026-6607 | `crypto/tls` | go1.26.9 or later |
| GO-2026-6603, GO-2026-6610, GO-2026-6611, GO-2026-6612, GO-2026-6617 | `net/http` and `golang.org/x/net` | go1.26.9 or later, together with `golang.org/x/net` v0.60.0 or later |
| GO-2026-6809 | `crypto/x509` | go1.26.10 |
| GO-2026-6860, GO-2026-6866 | `golang.org/x/net/html` | `golang.org/x/net` v0.62.0 |
| GO-2026-6354, GO-2026-6355 | `golang.org/x/crypto/ssh` | `golang.org/x/crypto` v0.56.0 or later |
| GO-2026-6179, GO-2026-6180 | `golang.org/x/mod/sumdb` | `golang.org/x/mod` v0.40.0 or later |
| GO-2026-4961, GO-2026-5061, GO-2026-5062, GO-2026-5066, GO-2026-6222 | `golang.org/x/image` (`webp`, `tiff`, `vp8l`) | `golang.org/x/image` v0.45.0 |
| GO-2026-5841 (CVE-2026-63209) | `github.com/klauspost/compress/s2` | `github.com/klauspost/compress` v1.18.7 |

The first six rows are the 14 standard library advisories described above. The last five rows are twelve
advisories in packages that no Core-Geth binary includes, so they never affected a running node. They
are listed because a scanner that checks module versions, rather than the code inside a binary,
reports them.

## The documentation build

`requirements-mkdocs.txt` lists the Python packages that build this documentation site. None of them is
part of `geth`. Three of them receive security updates:

| Package | `v1.13.0` | `v1.13.1` | Advisory |
|---|---|---|---|
| `certifi` | 2026.7.22 | 2026.10.23 | GHSA-6q48-6c7w-9gw2, critical. The certificate bundle included a root certificate, `emSign ECC Root CA - G3`, under which eight intermediate certificates had their private keys compromised. 2026.10.19 removes that root, and 2026.10.23 removes two more emSign roots, `emSign ECC Root CA - C3` and `emSign Root CA - C1`. |
| `urllib3` | 2.7.0 | 2.8.0 | GHSA-8988-9cw3-xx77, high. The TLS settings for an HTTPS proxy could be overridden by the destination server's settings. The same release fixes two more issues. |
| `GitPython` | 3.1.60 | 3.2.1 | GHSA-59cr-6r3x-644w, moderate. A crafted submodule path could create directories outside the repository. Fixed in 3.1.62. GHSA-83vg-56qc-22m7, high. GitPython checked a submodule's `.git` path only after cloning into it, which allowed remote code execution. Fixed in 3.2.1. |

The other packages move to their current versions by the procedure written at the top of
`requirements-mkdocs.txt`. `python-dateutil` moves to 4.0.0, which adds `tzdata` as a new package and no
longer needs `six`. It, `pymdown-extensions` 12.3 and `watchdog` 7.0.0 are new major versions. None of the
new versions has been withdrawn (yanked) by its publisher.

## Container images

The Docker images are built in two stages. The first compiles `geth` in the `golang:1.26.10-alpine3.24`
image and links it statically, so the binary includes that image's copy of musl, the C library. The
second copies `geth` into an `alpine:3.24.3` image and installs the `ca-certificates` package, which
holds the trusted certificates.

Both base images are named by version and by digest, so every build uses the same base images until
those lines change. In `v1.13.0` they followed the current Alpine release, through the
`golang:1.26.8-alpine` and `alpine:latest` tags. Alpine 3.24.3, released October 21, 2026, includes
security fixes to musl (CVE-2026-107891), to OpenSSL, and to `ca-certificates` for the same eMudhra
certificate compromise. Both digests cover `linux/amd64` and `linux/arm64`. The packages each stage
installs with `apk` still come from the Alpine 3.24 repository when the image is built.

## CI and release tooling

These run in GitHub Actions and are not part of any binary or image.

| Tool | `v1.13.0` | `v1.13.1` | Note |
|---|---|---|---|
| `docker/build-push-action` | v7.3.0 | v7.4.0 | Includes a fix for workflow command injection through metadata logs |
| `docker/setup-buildx-action` | v4.3.0 | v4.4.1 | |
| `docker/login-action` | v4.6.0 | v4.7.0 | |
| `actions/upload-artifact` | v7.0.1 | v7.0.2 | Retries rate-limited uploads |
| `actions/setup-python` | v7.0.0 | v7.1.0 | |
| Buildx, which builds the images | v0.37.0, from the runner image | v0.38.1 | Named in the workflow |
| BuildKit, which Buildx runs | `buildx-stable-1`, a floating tag | v0.34.1 | Named by version and digest. Built with go1.26.10 |
| `stringer`, which `make devtools` installs | v0.49.0 | v0.51.0 | Follows `golang.org/x/tools` |
| govulncheck, in the weekly scan | not present | v1.9.0 | The weekly scan was added after `v1.13.0` was released, with v1.8.0 |
| gcc-9, in the EVMC job | 9.5.0 | removed | The job uses the runner's compiler, as the test jobs do |
| solc, in the `go generate` job | 0.6.0 | removed | Nothing in the job ran it |

Each action is pinned to a full commit SHA, with the version in a comment beside it.

The EVMC job runs the state tests through two external virtual machines, hera and evmone. hera is
checked out by commit, and evmone's download is checked against the SHA-256 its release publishes.

## What a vulnerability scanner reports for v1.13.1

**Go standard library: nothing.** A scan of a `geth` binary built from `v1.13.1` with go1.26.10 finds
no Go standard library advisory in the database recorded under [Verification](#verification).

**Six go-ethereum advisories that Core-Geth already fixes.** Core-Geth keeps go-ethereum's module
path, `github.com/ethereum/go-ethereum`, for compatibility with code written against go-ethereum. A
scanner therefore compares Core-Geth's version number with go-ethereum's. Core-Geth's 1.13 version
numbers are lower than the go-ethereum releases that fixed six advisories, so the scanner reports
them: GO-2024-2819, GO-2026-4314, GO-2026-4315, GO-2026-4507, GO-2026-4508 and GO-2026-4511. The fixes
for all six are in Core-Geth's code, and
[Go toolchain](2026-09-go-toolchain.md#what-a-vulnerability-scanner-reports-against-v1130) shows where
each one is. `v1.13.0` reports the same six.

**One advisory with no fix, GO-2026-5932.** It marks the `golang.org/x/crypto/openpgp` package as
unmaintained. There is no fixed version, and no Core-Geth binary includes that package.

## How the v1.13 line stays current

- **A supported Go version.** Go 1.26 still receives security fixes, and each Core-Geth release is built
  with the latest Go 1.26 patch release. The `v1.12.x` line was built with Go versions whose support ended in 2024 and
  2025, so its binaries will never receive these fixes. [Go toolchain](2026-09-go-toolchain.md) measures
  the difference.
- **A weekly scan.** The workflow `.github/workflows/govulncheck.yml` scans `main` every week, both the
  source code and a `geth` binary built with the release Go version, and keeps a single issue up to date
  with anything a maintainer needs to see.
- **Stable, checked updates.** Every version this release moves to is a stable release that is not
  deprecated, retracted or yanked, and each update was built and scanned before it was merged.
- **Measured results.** The figures on this page come from scanning built binaries with govulncheck,
  the same tool and method the other audits use, so anyone can reproduce them.

**The `v1.12.x` line shows what happens without this.** In
[`etclabscore/core-geth`](https://github.com/etclabscore/core-geth) it stayed on Go versions that
reached end of life in 2024 and 2025, and its dependencies were not updated after December 2024, for 21
months without security maintenance. [Go toolchain](2026-09-go-toolchain.md) counts the Go standard
library advisories its binaries are affected by, a number that keeps growing. The
[March 2026 security audit](2026-03-security-audit.md) records the flaw that was exploited against the
Core-Geth `v1.12` bootnodes during that time. Backported fixes closed individual CVEs, but they left the
old Go version and the old dependencies in place. The `v1.13` line updates the client's own code, its Go
version and its dependencies together.

## Verification

Two terms used here: a scan **at module level** reports an advisory when a binary includes an affected
version of Go or of a module. A scan **at symbol level** reports it only when the vulnerable function
itself is in the binary.

- **Go version.** `go version -m` reads the build information recorded in a `geth` binary without
  running it. A binary built from the release branch with go1.26.10 shows go1.26.10, `golang.org/x/net`
  v0.62.0 and the other module versions above.
- **Checksums.** All 14 go1.26.10 checksums in `build/checksums.txt` match the values go.dev publishes, and
  `build/ci.go install -dlgo`, which builds the Arm release files, downloads go1.26.10 and checks it
  against them.
- **Build and lint.** `go build ./...`, `go vet ./...` and `make lint` pass, as do the two special builds
  that use `golang.org/x/image` and `github.com/klauspost/compress`.

**Advisories.** govulncheck v1.9.0 scanned three `geth` binaries on November 1, 2026, at module level and
at symbol level, against the Go vulnerability database (`vuln.go.dev`) as last modified
2026-10-27T21:16:18Z, with 4,712 entries:

| Binary | Built with | Standard library advisories, module level | Standard library advisories, symbol level |
|---|---|---:|---:|
| Release branch for `v1.13.1` | go1.26.10 | 0 | 0 |
| The same source | go1.26.9 | 1 (GO-2026-6809) | 1 |
| Published `v1.13.0` Linux binary | go1.26.8 | 13 | 10 |

The go1.26.9 and `v1.13.0` rows show that the scan does find standard library advisories when they are
present. The `v1.13.0` file's checksum and build attestation were verified before it was scanned. All
three binaries also report the six go-ethereum advisories described above, and GO-2026-5932 at module
level only. A scan of the release branch's source code finds no vulnerable function that the code
calls.

**Not checked:**

- The published `v1.13.1` files and container images, which do not exist yet. The figures above come
  from a build of the release branch. The published files will be checked the way
  [Go toolchain](2026-09-go-toolchain.md#verification) checks a release, once `v1.13.1` ships.
- Whether the client's code can actually reach any of the vulnerable functions.

## Supporting this work

The client is maintained by Ethereum Classic's core developers and open community contributions, in
the [`ethereumclassic`](https://github.com/ethereumclassic) GitHub organization. The Ethereum Classic
DAO is the legal entity standing behind that work, and the floor under it when nobody else is paying.
Mining pools, centralized exchanges, issuers of Ethereum Classic financial products, Etchash mining
hardware manufacturers and large holders all depend on this client. If your operation relies on
Ethereum Classic, please consider sponsoring its maintenance. [Support this work](../support.md) says
who maintains and funds it and how to sponsor it.
