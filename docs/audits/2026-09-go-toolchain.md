---
description: "Which Go toolchain built each published Core-Geth archive, and the Go standard library advisories each v1.12.x archive carries that v1.13.0 does not."
---

# Go toolchain: v1.12.x archive to v1.13.0

This document records which Go toolchain built each published `v1.12.x` release, the Go standard
library advisories those binaries carry, what `v1.13.0` ships instead, and why an operator on the
older line should move. It is a companion to `2026-08-dependency-modernization.md`, which covers the
toolchain and module changes in the source, and to `2026-09-release-pipeline.md`, which measures the
same published archives for platform floors and provenance. `2026-03-security-audit.md` and
`2026-08-security-followup.md` cover the CVE remediation in the client's own code.

**Baseline:** the four releases `v1.12.20` through `v1.12.23`, published from `etclabscore/core-geth`,
measured by downloading every archive and reading the build information recorded in each binary
rather than the configuration that built it. The Linux x86_64 archive of each earlier stable
`v1.12.x` release was measured the same way.

**Work carried out by:** [White B0x](https://whiteb0x.com)

**On this page:** [In short](#in-short) · [Why the toolchain needed a pass of its own](#why-the-toolchain-needed-a-pass-of-its-own) · [Each release from v1.12.20 to v1.12.23 was built by two toolchains](#finding-each-release-from-v11220-to-v11223-was-built-by-two-toolchains) · [The v1.12.20 to v1.12.23 binaries carry Go standard library advisories](#finding-the-v11220-to-v11223-binaries-carry-go-standard-library-advisories) · [The v1.12.x source does not build on a supported Go](#finding-the-v112x-source-does-not-build-on-a-supported-go) · [Third-party module advisories were cleared alongside](#finding-third-party-module-advisories-were-cleared-alongside) · [What a vulnerability scanner reports against v1.13.0](#what-a-vulnerability-scanner-reports-against-v1130) · [What this means if you are upgrading](#what-this-means-if-you-are-upgrading) · [Verification](#verification) · [After v1.13.0](#after-v1130) · [Disposition of the build finding](#disposition-of-the-build-finding-17-september-2026) · [Supporting this work](#supporting-this-work)

## In short

**Upgrade to [v1.13.0 or later](https://github.com/ethereumclassic/core-geth/releases/latest).** Each release from `v1.12.20` to `v1.12.23` was built by two Go
toolchains, Go 1.21 and Go 1.22, both out of support before those releases shipped. Each of those archives
carries 55 to 61 Go standard library advisories, 43 to 47 of them with the vulnerable code present in the binary.
Every `v1.13.0` archive and image is built with go1.26.8 and carries none.

## Why the toolchain needed a pass of its own

The Go standard library is compiled into every binary, `crypto/tls` and `net/http` included. A Go
security release therefore reaches an operator only through a client release built with it, and a
binary keeps its toolchain's advisories for as long as it runs.

The build configuration names a Go version, and the binary records the one that actually built it.
`go version -m` reads that record from a downloaded file without running it. Read from the
`go-version` setting in the release workflow, the `v1.12.x` releases were built with Go 1.21. Their
binaries record two toolchains, and the difference decides which advisories each archive carries.

## Finding: each release from v1.12.20 to v1.12.23 was built by two toolchains

Every archive of `v1.12.20` to `v1.12.23` and of `v1.13.0`, read from the binaries:

| Release | Published | `linux` | `osx` | `osx-arm64` | `win64` | `arm`, `arm5`, `arm6`, `arm7`, `arm64` |
|---|---|---|---|---|---|---|
| `v1.12.20` | 2024-06-10 | go1.21.10 | go1.21.10 | not published | go1.22.1 | go1.22.1 |
| `v1.12.21` | 2026-03-18 | go1.21.13 | go1.21.13 | not published | go1.21.13 | go1.22.1 |
| `v1.12.22` | 2026-03-28 | go1.21.13 | go1.21.13 | not published | go1.21.13 | go1.22.1 |
| `v1.12.23` | 2026-08-14 | go1.21.13 | go1.21.13 | not published | go1.21.13 | go1.22.1 |
| `v1.13.0` | 2026-09-14 | go1.26.8 | go1.26.8 | go1.26.8 | go1.26.8 | go1.26.8 |

**21 of the 32 `v1.12.x` archives were built with go1.22.1, every Arm archive among them.** The
other 11 were built with Go 1.21. In these four releases the `osx` archive is an Apple Silicon build,
which [release artifacts](2026-09-release-pipeline.md#finding-the-macos-archive-contains-a-binary-most-of-its-downloaders-cannot-run)
covers; `v1.13.0` publishes the two architectures under separate names.

The split follows two build paths in the release configuration, and the pins were the same at all
four release commits. The release workflow installed Go with `go-version: '1.21'` and built the Linux
x86_64, macOS and, from `v1.12.21`, Windows archives with it. The Arm archives were built with
`build/ci.go install -dlgo`, which downloads the Go version pinned in `build/checksums.txt`: 1.22.1.
The `v1.12.20` Windows archive came from AppVeyor, which also used `-dlgo`.

Every `v1.13.0` node archive, and `geth` in the two `core-geth-docker` image tarballs attached to the
release, was built with go1.26.8.

**Both older majors are out of support.** Go's [release policy](https://go.dev/doc/devel/release):
*"Each major Go release is supported until there are two newer major releases."*

| Go major | First release | Support ended | Archives built with it, `v1.12.20` to `v1.13.0` |
|---|---|---|---:|
| Go 1.21 | 2023-08-08 | 2024-08-13, when go1.23.0 was released | 11 |
| Go 1.22 | 2024-02-06 | 2025-02-11, when go1.24.0 was released | 21 |
| Go 1.26 | 2026-02-10 | Supported until go1.28.0 is released | 9 |

**The 2026 releases shipped after both majors had left support.** The toolchains on each release
date:

| Release | Archives | Go | Major supported that day | Newest patch of that major that day | Security releases of that major the archives lacked that day |
|---|---|---|---|---|---|
| `v1.12.20` (2024-06-10) | `linux`, `osx` | go1.21.10 | yes | go1.21.11 | 1 (go1.21.11) |
| `v1.12.20` (2024-06-10) | `win64`, Arm archives | go1.22.1 | yes | go1.22.4 | 3 (go1.22.2, go1.22.3, go1.22.4) |
| `v1.12.21` (2026-03-18) | `linux`, `osx`, `win64` | go1.21.13 | no | go1.21.13 | 0 |
| `v1.12.21` (2026-03-18) | Arm archives | go1.22.1 | no | go1.22.12 | 7 |
| `v1.12.22` (2026-03-28) | `linux`, `osx`, `win64` | go1.21.13 | no | go1.21.13 | 0 |
| `v1.12.22` (2026-03-28) | Arm archives | go1.22.1 | no | go1.22.12 | 7 |
| `v1.12.23` (2026-08-14) | `linux`, `osx`, `win64` | go1.21.13 | no | go1.21.13 | 0 |
| `v1.12.23` (2026-08-14) | Arm archives | go1.22.1 | no | go1.22.12 | 7 |
| `v1.13.0` (2026-09-14) | all archives | go1.26.8 | yes | go1.26.8 | 0 |

Their Go 1.21 archives carry go1.21.13, the last Go 1.21 release, so nothing newer exists on that
branch: every fix published after August 2024 shipped only in newer Go majors. Their Arm archives
carry go1.22.1, which was missing seven later security releases on its own branch, and Go 1.22 had no
release after February 2025. `v1.12.20`, published while both majors were still supported, already
shipped behind the newest patch of each.

**The earlier releases are on unsupported majors as well.** The Linux x86_64 archive of every stable
release from `v1.12.0` to `v1.12.19` was measured too. Across those and the 32 above, the 50
`v1.12.x` archives measured were built with eight Go majors, from Go 1.13 to Go 1.22, and none of
them is supported today. Two of those releases shipped on a major that had already left support:
`v1.12.2` (go1.13.4) and `v1.12.10` (go1.18.10).

## Finding: the v1.12.20 to v1.12.23 binaries carry Go standard library advisories

The Go vulnerability database lists 62 standard library advisories that apply to at least one
archive of `v1.12.20` to `v1.12.23`, read with govulncheck against the database snapshot recorded
under [Verification](#verification). Two levels matter:

- **Module level:** the toolchain version is in the advisory's affected range, for the archive's
  operating system and architecture.
- **Symbol level:** a function the advisory lists as vulnerable is present in the binary.

Presence is not proof of reachability. Whether the client can reach the vulnerable code was not
assessed for any advisory on this page.

| Release (published) | Go | Module level | ...already published that day | Symbol level | ...already published that day |
|---|---|---:|---:|---:|---:|
| `v1.12.20` (2024-06-10) | go1.21.10, go1.22.1 | 62 | 4 | 48 | 3 |
| `v1.12.21` (2026-03-18) | go1.21.13, go1.22.1 | 62 | 36 | 48 | 29 |
| `v1.12.22` (2026-03-28) | go1.21.13, go1.22.1 | 62 | 36 | 48 | 29 |
| `v1.12.23` (2026-08-14) | go1.21.13, go1.22.1 | 62 | 62 | 48 | 48 |
| `v1.13.0` (2026-09-14) | go1.26.8 | 0 | 0 | 0 | 0 |

Counts are distinct advisories across all archives of a release. "Already published that day" counts
those whose database entry was published on or before the release date. All 62 had been published
when `v1.12.23` shipped. The `v1.12.20` count is lower only because most were published after it,
and its binaries carry all 62 regardless.

Per archive, in `v1.12.23`:

| `v1.12.23` archive | Go | Module level | Symbol level |
|---|---|---:|---:|
| `linux` | go1.21.13 | 56 | 43 |
| `osx` | go1.21.13 | 55 | 43 |
| `win64` | go1.21.13 | 56 | 44 |
| `arm`, `arm5`, `arm6`, `arm7`, `arm64` | go1.22.1 | 61 | 47 |

**From `v1.12.21` on, the go1.22.1 archives carry five advisories the Go 1.21 archives of the same
release do not.** Two were never in a Go 1.21 archive, and go1.21.13 includes the fixes for the other
three, which go1.22.1 predates:

| GO ID | Summary (verbatim) | First fixed in | Only in the go1.22.1 archives of |
|---|---|---|---|
| GO-2024-2687 | HTTP/2 CONTINUATION flood in net/http | go1.22.2 | `v1.12.20` to `v1.12.23` |
| GO-2024-2824 | Malformed DNS message can cause infinite loop in net | go1.22.3 | `v1.12.20` to `v1.12.23` |
| GO-2024-2887 | Unexpected behavior from Is methods for IPv4-mapped IPv6 addresses in net/netip | go1.21.11, go1.22.4 | `v1.12.21` to `v1.12.23` |
| GO-2024-2888 | Mishandling of corrupt central directory record in archive/zip | go1.21.11, go1.22.4 | `v1.12.21` to `v1.12.23` |
| GO-2024-2963 | Denial of service due to improper 100-continue handling in net/http | go1.21.12, go1.22.5 | `v1.12.21` to `v1.12.23` |

Two others depend on the operating system: GO-2025-3750 applies only to Windows, and GO-2026-4864
only to Linux.

The Go vulnerability database assigns no severity to these entries, and this page does not add one.

??? note "All 62 standard library advisories"
    *First fixed in* is the first Go release carrying the fix after the shipped patch. A version on
    the shipped major's own branch means a patch release of that major fixed it; two versions mean
    Go 1.21 and Go 1.22 each had one; a newer major means neither did. *Highest level* is the highest
    found in any archive, where *package* means the affected package is compiled in but none of its
    listed vulnerable functions was found.

    | # | GO ID | CVE | Package | Summary (verbatim) | Published | First fixed in | Highest level | Symbol level in |
    |---|---|---|---|---|---|---|---|---|
    | 1 | GO-2024-2687 | CVE-2023-45288 | `net/http` | HTTP/2 CONTINUATION flood in net/http | 2024-04-03 | go1.22.2 | symbol | `v1.12.20`: `win64` and Arm archives; `v1.12.21` to `v1.12.23`: Arm archives |
    | 2 | GO-2024-2824 | CVE-2024-24788 | `net` | Malformed DNS message can cause infinite loop in net | 2024-05-07 | go1.22.3 | symbol | `v1.12.20`: `win64` and Arm archives; `v1.12.21` to `v1.12.23`: Arm archives |
    | 3 | GO-2024-2887 | CVE-2024-24790 | `net/netip` | Unexpected behavior from Is methods for IPv4-mapped IPv6 addresses in net/netip | 2024-06-04 | go1.21.11, go1.22.4 | symbol | `v1.12.20`: all archives; `v1.12.21` to `v1.12.23`: Arm archives |
    | 4 | GO-2024-2888 | CVE-2024-24789 | `archive/zip` | Mishandling of corrupt central directory record in archive/zip | 2024-06-04 | go1.21.11, go1.22.4 | module | none |
    | 5 | GO-2024-2963 | CVE-2024-24791 | `net/http` | Denial of service due to improper 100-continue handling in net/http | 2024-07-02 | go1.21.12, go1.22.5 | symbol | `v1.12.20`: all archives; `v1.12.21` to `v1.12.23`: Arm archives |
    | 6 | GO-2024-3105 | CVE-2024-34155 | `go/parser` | Stack exhaustion in all Parse functions in go/parser | 2024-09-06 | go1.22.7 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 7 | GO-2024-3106 | CVE-2024-34156 | `encoding/gob` | Stack exhaustion in Decoder.Decode in encoding/gob | 2024-09-06 | go1.22.7 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 8 | GO-2024-3107 | CVE-2024-34158 | `go/build/constraint` | Stack exhaustion in Parse in go/build/constraint | 2024-09-06 | go1.22.7 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 9 | GO-2025-3373 | CVE-2024-45341 | `crypto/x509` | Usage of IPv6 zone IDs can bypass URI name constraints in crypto/x509 | 2025-01-28 | go1.22.11 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 10 | GO-2025-3420 | CVE-2024-45336 | `net/http` | Sensitive headers incorrectly sent after cross-domain redirect in net/http | 2025-01-28 | go1.22.11 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 11 | GO-2025-3503 | CVE-2025-22870 | `net/http` | HTTP Proxy bypass using IPv6 Zone IDs in golang.org/x/net | 2025-03-12 | go1.23.7 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 12 | GO-2025-3563 | CVE-2025-22871 | `net/http/internal` | Request smuggling due to acceptance of invalid chunked data in net/http | 2025-04-08 | go1.23.8 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 13 | GO-2025-3751 | CVE-2025-4673 | `net/http` | Sensitive headers not cleared on cross-origin redirect in net/http | 2025-06-11 | go1.23.10 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 14 | GO-2025-3750 | CVE-2025-0913 | `os` (Windows only), `syscall` (Windows only) | Inconsistent handling of O_CREATE\|O_EXCL on Unix and Windows in os in syscall | 2025-06-11 | go1.23.10 | symbol | `v1.12.20` to `v1.12.23`: `win64` |
    | 15 | GO-2025-3849 | CVE-2025-47907 | `database/sql` | Incorrect results returned from Rows.Scan in database/sql | 2025-08-07 | go1.23.12 | package | none |
    | 16 | GO-2025-3956 | CVE-2025-47906 | `os/exec` | Unexpected paths returned from LookPath in os/exec | 2025-09-18 | go1.23.12 | module | none |
    | 17 | GO-2025-4006 | CVE-2025-61725 | `net/mail` | Excessive CPU consumption in ParseAddress in net/mail | 2025-10-29 | go1.24.8 | module | none |
    | 18 | GO-2025-4007 | CVE-2025-58187 | `crypto/x509` | Quadratic complexity when checking name constraints in crypto/x509 | 2025-10-29 | go1.24.9 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 19 | GO-2025-4008 | CVE-2025-58189 | `crypto/tls` | ALPN negotiation error contains attacker controlled information in crypto/tls | 2025-10-29 | go1.24.8 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 20 | GO-2025-4009 | CVE-2025-61723 | `encoding/pem` | Quadratic complexity when parsing some invalid inputs in encoding/pem | 2025-10-29 | go1.24.8 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 21 | GO-2025-4010 | CVE-2025-47912 | `net/url` | Insufficient validation of bracketed IPv6 hostnames in net/url | 2025-10-29 | go1.24.8 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 22 | GO-2025-4011 | CVE-2025-58185 | `encoding/asn1` | Parsing DER payload can cause memory exhaustion in encoding/asn1 | 2025-10-29 | go1.24.8 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 23 | GO-2025-4012 | CVE-2025-58186 | `net/http` | Lack of limit when parsing cookies can cause memory exhaustion in net/http | 2025-10-29 | go1.24.8 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 24 | GO-2025-4013 | CVE-2025-58188 | `crypto/x509` | Panic when validating certificates with DSA public keys in crypto/x509 | 2025-10-29 | go1.24.8 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 25 | GO-2025-4014 | CVE-2025-58183 | `archive/tar` | Unbounded allocation when parsing GNU sparse map in archive/tar | 2025-10-29 | go1.24.8 | module | none |
    | 26 | GO-2025-4015 | CVE-2025-61724 | `net/textproto` | Excessive CPU consumption in Reader.ReadResponse in net/textproto | 2025-10-29 | go1.24.8 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 27 | GO-2025-4155 | CVE-2025-61729 | `crypto/x509` | Excessive resource consumption when printing error string for host certificate validation in crypto/x509 | 2025-12-02 | go1.24.11 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 28 | GO-2025-4175 | CVE-2025-61727 | `crypto/x509` | Improper application of excluded DNS name constraints when verifying wildcard names in crypto/x509 | 2025-12-02 | go1.24.11 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 29 | GO-2026-4340 | CVE-2025-61730 | `crypto/tls` | Handshake messages may be processed at the incorrect encryption level in crypto/tls | 2026-01-28 | go1.24.12 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 30 | GO-2026-4341 | CVE-2025-61726 | `net/url` | Memory exhaustion in query parameter parsing in net/url | 2026-01-28 | go1.24.12 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 31 | GO-2026-4342 | CVE-2025-61728 | `archive/zip` | Excessive CPU consumption when building archive index in archive/zip | 2026-01-28 | go1.24.12 | module | none |
    | 32 | GO-2026-4403 | CVE-2025-22873 | `os` | Improper access to parent directory of root in os | 2026-02-04 | go1.23.9 | package | none |
    | 33 | GO-2026-4337 | CVE-2025-68121 | `crypto/tls` | Unexpected session resumption in crypto/tls | 2026-02-05 | go1.24.13 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 34 | GO-2026-4601 | CVE-2026-25679 | `net/url` | Incorrect parsing of IPv6 host literals in net/url | 2026-03-06 | go1.25.8 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 35 | GO-2026-4602 | CVE-2026-27139 | `os` | FileInfo can escape from a Root in os | 2026-03-06 | go1.25.8 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 36 | GO-2026-4603 | CVE-2026-27142 | `html/template` | URLs in meta content attribute actions are not escaped in html/template | 2026-03-06 | go1.25.8 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 37 | GO-2026-4864 | CVE-2026-32282 | `internal/syscall/unix` (Linux only) | TOCTOU permits root escape on Linux via Root.Chmod in os in internal/syscall/unix | 2026-04-07 | go1.25.9 | package | none |
    | 38 | GO-2026-4865 | CVE-2026-32289 | `html/template` | JsBraceDepth Context Tracking Bugs (XSS) in html/template | 2026-04-07 | go1.25.9 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 39 | GO-2026-4869 | CVE-2026-32288 | `archive/tar` | Unbounded allocation for old GNU sparse in archive/tar | 2026-04-07 | go1.25.9 | module | none |
    | 40 | GO-2026-4870 | CVE-2026-32283 | `crypto/tls` | Unauthenticated TLS 1.3 KeyUpdate record can cause persistent connection retention and DoS in crypto/tls | 2026-04-07 | go1.25.9 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 41 | GO-2026-4946 | CVE-2026-32281 | `crypto/x509` | Inefficient policy validation in crypto/x509 | 2026-04-07 | go1.25.9 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 42 | GO-2026-4947 | CVE-2026-32280 | `crypto/x509` | Unexpected work during chain building in crypto/x509 | 2026-04-07 | go1.25.9 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 43 | GO-2026-4918 | CVE-2026-33814 | `net/http` | Infinite loop in HTTP/2 transport when given bad SETTINGS_MAX_FRAME_SIZE in net/http/internal/http2 in golang.org/x/net | 2026-05-07 | go1.25.10 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 44 | GO-2026-4971 | CVE-2026-39836 | `net` | Panic in Dial and LookupPort when handling NUL byte on Windows in net | 2026-05-07 | go1.25.10 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 45 | GO-2026-4976 | CVE-2026-39825 | `net/http/httputil` | ReverseProxy forwards queries with more than urlmaxqueryparams parameters in net/http/httputil | 2026-05-07 | go1.25.10 | module | none |
    | 46 | GO-2026-4977 | CVE-2026-42499 | `net/mail` | Quadratic string concatenation in consumePhrase in net/mail | 2026-05-07 | go1.25.10 | module | none |
    | 47 | GO-2026-4980 | CVE-2026-39826 | `html/template` | Escaper bypass leads to XSS in html/template | 2026-05-07 | go1.25.10 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 48 | GO-2026-4981 | CVE-2026-33811 | `net` | Crash when handling long CNAME response in net | 2026-05-07 | go1.25.10 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 49 | GO-2026-4982 | CVE-2026-39823 | `html/template` | Bypass of meta content URL escaping causes XSS in html/template | 2026-05-07 | go1.25.10 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 50 | GO-2026-4986 | CVE-2026-39820 | `net/mail` | Quadratic string concatentation in consumeComment in net/mail | 2026-05-07 | go1.25.10 | module | none |
    | 51 | GO-2026-5026 | CVE-2026-39821 | `net/http`, `net/http/internal/http2` | Invoking failure to reject ASCII-only Punycode-encoded labels in golang.org/x/net/idna | 2026-05-22 | go1.25.13 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 52 | GO-2026-5037 | CVE-2026-27145 | `crypto/x509` | Inefficient candidate hostname parsing in crypto/x509 | 2026-06-02 | go1.25.11 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 53 | GO-2026-5038 | CVE-2026-42504 | `mime` | Quadratic complexity in WordDecoder.DecodeHeader in mime | 2026-06-02 | go1.25.11 | package | none |
    | 54 | GO-2026-5039 | CVE-2026-42507 | `net/textproto` | Arbitrary inputs are included in errors without any escaping in net/textproto | 2026-06-02 | go1.25.11 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 55 | GO-2026-4970 | CVE-2026-39822 | `os` | Root escape via symlink plus trailing slash in os | 2026-07-07 | go1.25.12 | package | none |
    | 56 | GO-2026-5856 | CVE-2026-42505 | `crypto/tls` | Invoking Encrypted Client Hello privacy leak in crypto/tls | 2026-07-07 | go1.25.12 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 57 | GO-2026-5972 | CVE-2026-33818 | `encoding/asn1` | Enforce maximum recursion depth in encoding/asn1 | 2026-08-13 | go1.25.13 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 58 | GO-2026-6088 | CVE-2026-56859 | `encoding/xml` | Add recursion depth guard during decode in encoding/xml | 2026-08-13 | go1.25.13 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 59 | GO-2026-6089 | CVE-2026-56853 | `net/http` | Apply ReadHeaderTimeout when doing unencrypted HTTP/2 check in net/http | 2026-08-13 | go1.25.13 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 60 | GO-2026-6090 | CVE-2026-56862 | `crypto/tls` | Limit handshake messages we are willing to accept post-handshake in crypto/tls | 2026-08-13 | go1.25.13 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 61 | GO-2026-6091 | CVE-2026-56858 | `html/template` | Fix Javascript regexp context tracking in html/template | 2026-08-13 | go1.25.13 | symbol | `v1.12.20` to `v1.12.23`: all archives |
    | 62 | GO-2026-6218 | CVE-2026-56860 | `net/url` | Avoid quadratic complexity in resolvePath in net/url | 2026-08-13 | go1.25.13 | symbol | `v1.12.20` to `v1.12.23`: all archives |

## Finding: the v1.12.x source does not build on a supported Go

Rebuilding a `v1.12.x` release with a supported Go would take the standard library fixes if the
source built. It does not. At the `v1.12.23` release commit, built with go1.26.6 using the build tags
and C flags recorded in the published `v1.12.23` binary, the build stops at two independent
blockers, in this order.

**`blst` v0.3.11 does not compile from Go 1.24.** The release build sets the `ckzg` tag, as
`make geth` and `make all` do, and the KZG binding it enables imports `blst`. The build fails in
`blst`'s Go bindings with `cannot define new methods on non-local type SecretKey`, and the same error
for `Fp12`, `P1Affine` and `P2Affine`. From the [Go 1.24 release notes](https://go.dev/doc/go1.24):

> The compiler already disallowed defining new methods with receiver types that were cgo-generated,
> but it was possible to circumvent that restriction via an alias type. Go 1.24 now always reports an
> error if a receiver denotes a cgo-generated type, whether directly or indirectly (through an alias
> type).

Upstream go-ethereum moved to `blst` v0.3.14, whose release carries a fix for Go 1.24, in
[#31165](https://github.com/ethereum/go-ethereum/pull/31165), first released in go-ethereum v1.15.1.

**`fjl/memsize` does not link from Go 1.23.** Without the `ckzg` tag the build reaches the linker and
stops at `link: github.com/fjl/memsize: invalid reference to runtime.stopTheWorld`. memsize backs the
`/memsize/` page of the pprof debug server and refers to runtime internals through `//go:linkname`.
From the [Go 1.23 release notes](https://go.dev/doc/go1.23):

> The linker now disallows using a //go:linkname directive to refer to internal symbols in the
> standard library (including the runtime) that are not marked with //go:linkname on their
> definitions.

Upstream go-ethereum removed memsize ahead of Go 1.23, in
[#30253](https://github.com/ethereum/go-ethereum/pull/30253): "Removing because memsize will very
likely be broken by Go 1.23". That change was first released in go-ethereum v1.14.8.

**Removing memsize alone does not produce a build.** With the `ckzg` tag the release configuration
sets, the `blst` compile error comes first on any Go from 1.24. Leaving the tag out as well does
produce one, and operators have run `v1.12.23` built that way on a current Go.

**A toolchain swap alone also keeps Go 1.21's defaults.** The `go 1.21` directive in `go.mod` keeps
Go 1.21's compatibility settings when a newer Go builds the module. All 21 go1.22.1 archives record
them:

```
DefaultGODEBUG=httplaxcontentlength=1,httpmuxgo121=1,tls10server=1,tlsrsakex=1,tlsunsafeekm=1
```

**On a Go from 1.24 on, the same directive switches off two denial-of-service fixes.** Go added
`httpcookiemaxnum`, a limit on the cookies `net/http` parses, and `urlmaxqueryparams`, a limit on the
query parameters `net/url` parses, as the fixes for GO-2025-4012 and GO-2026-4341. Its compatibility
table gives a module declaring a Go older than 1.24 the value `0` for both, which, in the words of
[Go's GODEBUG history](https://go.dev/doc/godebug), accepts "an indefinite number of cookies" and
"disables the limit". go1.26.6 and go1.27.1 carry the same entries. So a `v1.12.23` rebuilt on a
current Go reports a Go version whose standard library carries those fixes, with both limits off.
Setting `GODEBUG=httpcookiemaxnum=3000,urlmaxqueryparams=10000` at run time turns them back on.

**The previous repository records both blockers.** In `etclabscore/core-geth`:

- [#680](https://github.com/etclabscore/core-geth/issues/680), an issue opened 2025-04-30, reports
  that `v1.12.20` does not compile with go1.24.0, and its build log shows the `blst` v0.3.11 errors
  above. It was closed on 2025-06-29 with the reply "Go 1.24 is not supported, please try Go 1.21".
- [#683](https://github.com/etclabscore/core-geth/pull/683), a pull request titled "Support go 1.24
  compiler" and opened 2025-06-30, changes `go.mod`, `internal/debug/flags.go`, the workflows and the
  Dockerfiles. Its description says it fixes #680 and #665, an earlier pull request that removes
  memsize. It was closed without being merged on 2026-02-06.
- [#690](https://github.com/etclabscore/core-geth/issues/690), an issue opened 2025-11-10, reports the
  AppVeyor build failing inside `blst` v0.3.11's header with `'bool' cannot be defined via 'typedef'`,
  and [#691](https://github.com/etclabscore/core-geth/pull/691), a pull request opened the same day,
  upgrades `blst` in `go.mod` to address it. That failure comes from the C compiler rather than from Go: it is the C23
  finding in
  [release artifacts](2026-09-release-pipeline.md#finding-the-v112x-line-no-longer-builds-against-a-c23-compiler).
- [#701](https://github.com/etclabscore/core-geth/pull/701), a pull request opened 2026-09-15, removes
  memsize and its debug endpoint, citing the link error above. On its own that change does not clear
  the `blst` error.

**`v1.13.0` carries `blst` v0.3.17, no memsize, and `go 1.26.0` in `go.mod`.** Its binaries record no
`DefaultGODEBUG` override.

## Finding: third-party module advisories were cleared alongside

The same scans cover the modules compiled into each binary. In `v1.12.23`, govulncheck matches 47
advisories against third-party modules at module level, 9 of them at symbol level. 44 are absent from
`v1.13.0`, 3 remain, and none is new. The main module, `github.com/ethereum/go-ethereum`, is not
counted here; [the scanner section](#what-a-vulnerability-scanner-reports-against-v1130) covers it.

| Module | `v1.12.23` | `v1.13.0` | Cleared | Cleared, symbol level in `v1.12.23` | Still matched in `v1.13.0` |
|---|---|---|---:|---:|---:|
| `github.com/consensys/gnark-crypto` | v0.12.1 | v0.21.0 | 1 | 1 | 0 |
| `github.com/golang-jwt/jwt/v4` | v4.5.0 | v4.5.2 | 2 | 2 | 0 |
| `github.com/gorilla/websocket` | v1.5.0 | v1.5.3 | 1 | 1 | 0 |
| `github.com/tidwall/gjson` | v1.6.0 | v1.19.0 | 4 | 4 | 0 |
| `golang.org/x/crypto` | v0.17.0 | v0.55.0 | 19 | 0 | 3 |
| `golang.org/x/net` | v0.18.0 | v0.58.0 | 14 | 0 | 0 |
| `golang.org/x/sys` | v0.16.0 | v0.47.0 | 1 | 0 | 0 |
| `golang.org/x/text` | v0.14.0 | v0.41.0 | 1 | 1 | 0 |
| `google.golang.org/protobuf` | v1.31.0 | v1.33.0 | 1 | 0 | 0 |
| **Total** |  |  | **44** | **9** | **3** |

The three that remain are `golang.org/x/crypto` advisories, described in the scanner section. Four of
the 44 also cover the standard library and appear in the list of 62 above.
[Dependency and toolchain modernization](2026-08-dependency-modernization.md) records the module
changes themselves.

??? note "All 44 cleared advisories"
    *Symbol level* refers to the module named in the row.

    | # | GO ID | Aliases | Module | Summary (verbatim) | First fixed version | Published | Symbol level in `v1.12.23` |
    |---|---|---|---|---|---|---|---|
    | 1 | GO-2025-4087 | GHSA-fj2x-735w-74vq | `github.com/consensys/gnark-crypto` | Unchecked memory allocation during vector deserialization in github.com/consensys/gnark-crypto | 0.18.1 | 2025-11-05 | all archives |
    | 2 | GO-2024-3250 | CVE-2024-51744, GHSA-29wx-vh33-7x7r | `github.com/golang-jwt/jwt/v4` | Improper error handling in ParseWithClaims and bad documentation may cause dangerous situations in github.com/golang-jwt/jwt | 4.5.1 | 2024-11-12 | all archives |
    | 3 | GO-2025-3553 | CVE-2025-30204, GHSA-mh63-6h87-95cp | `github.com/golang-jwt/jwt/v4` | Excessive memory allocation during header parsing in github.com/golang-jwt/jwt | 4.5.2 | 2025-03-26 | all archives |
    | 4 | GO-2026-6278 | GHSA-w67g-5rqw-f597 | `github.com/gorilla/websocket` | Gorilla WebSocket Uses Cryptographically Weak PRNG for WebSocket Mask Key in github.com/gorilla/websocket | 1.5.3 | 2026-08-25 | all archives |
    | 5 | GO-2021-0054 | CVE-2020-36067, GHSA-p64j-r5f4-pwwx | `github.com/tidwall/gjson` | Panic due to improper input validation in ForEach in github.com/tidwall/gjson | 1.6.6 | 2021-04-14 | all archives |
    | 6 | GO-2021-0059 | CVE-2020-35380, GHSA-w942-gw6m-p62c | `github.com/tidwall/gjson` | Panic due to improper input validation in Get in github.com/tidwall/gjson | 1.6.4 | 2021-04-14 | all archives |
    | 7 | GO-2021-0265 | CVE-2021-42248, CVE-2021-42836, GHSA-c9gm-7rfj-8w5h, GHSA-ppj4-34rq-v8j9 | `github.com/tidwall/gjson` | Denial of service via maliciously crafted path in github.com/tidwall/gjson | 1.9.3 | 2022-08-15 | all archives |
    | 8 | GO-2022-0957 | CVE-2020-36066, GHSA-wjm3-fq3r-5x46 | `github.com/tidwall/gjson` | Denial of service via maliciously crafted JSON in github.com/tidwall/gjson | 1.6.5 | 2022-08-25 | all archives |
    | 9 | GO-2024-3321 | CVE-2024-45337, GHSA-v778-237x-gjrc | `golang.org/x/crypto` | Misuse of connection.serverAuthenticate may cause authorization bypass in golang.org/x/crypto | 0.31.0 | 2024-12-11 | no |
    | 10 | GO-2025-3487 | CVE-2025-22869 | `golang.org/x/crypto` | Potential denial of service in golang.org/x/crypto | 0.35.0 | 2025-02-26 | no |
    | 11 | GO-2025-4116 | CVE-2025-47913 | `golang.org/x/crypto` | Potential denial of service in golang.org/x/crypto/ssh/agent | 0.43.0 | 2025-11-13 | no |
    | 12 | GO-2025-4134 | CVE-2025-58181, GHSA-j5w8-q4qc-rx2x | `golang.org/x/crypto` | Unbounded memory consumption in golang.org/x/crypto/ssh | 0.45.0 | 2025-11-19 | no |
    | 13 | GO-2025-4135 | CVE-2025-47914, GHSA-f6x5-jh6r-wrfv | `golang.org/x/crypto` | Malformed constraint may cause denial of service in golang.org/x/crypto/ssh/agent | 0.45.0 | 2025-11-19 | no |
    | 14 | GO-2026-5005 | CVE-2026-39833, GHSA-jppx-rxg9-jmrx | `golang.org/x/crypto` | Invoking key constraints not enforced in golang.org/x/crypto/ssh/agent | 0.52.0 | 2026-05-22 | no |
    | 15 | GO-2026-5006 | CVE-2026-39832, GHSA-f5wc-c3c7-36mc | `golang.org/x/crypto` | Invoking agent constraints dropped when forwarding keys in golang.org/x/crypto/ssh/agent | 0.52.0 | 2026-05-22 | no |
    | 16 | GO-2026-5013 | CVE-2026-46597 | `golang.org/x/crypto` | Invoking byte arithmetic causes underflow and panic in golang.org/x/crypto/ssh | 0.52.0 | 2026-05-22 | no |
    | 17 | GO-2026-5014 | CVE-2026-39828 | `golang.org/x/crypto` | Invoking bypass of certificate restrictions in golang.org/x/crypto/ssh | 0.52.0 | 2026-05-22 | no |
    | 18 | GO-2026-5015 | CVE-2026-39835 | `golang.org/x/crypto` | Invoking server panic during CheckHostKey/Authenticate in golang.org/x/crypto/ssh | 0.52.0 | 2026-05-22 | no |
    | 19 | GO-2026-5016 | CVE-2026-39827 | `golang.org/x/crypto` | Invoking memory leak when rejecting channels can lead to DoS in golang.org/x/crypto/ssh | 0.52.0 | 2026-05-22 | no |
    | 20 | GO-2026-5017 | CVE-2026-39830 | `golang.org/x/crypto` | Invoking client can cause server deadlock on unexpected responses in golang.org/x/crypto/ssh | 0.52.0 | 2026-05-22 | no |
    | 21 | GO-2026-5018 | CVE-2026-39829 | `golang.org/x/crypto` | Invoking pathological RSA/DSA parameters may cause DoS in golang.org/x/crypto/ssh | 0.52.0 | 2026-05-22 | no |
    | 22 | GO-2026-5019 | CVE-2026-39831 | `golang.org/x/crypto` | Invoking bypass of FIDO/U2F security keys physical interaction in golang.org/x/crypto/ssh | 0.52.0 | 2026-05-22 | no |
    | 23 | GO-2026-5020 | CVE-2026-39834 | `golang.org/x/crypto` | Invoking infinite loop on large channel writes in golang.org/x/crypto/ssh | 0.52.0 | 2026-05-22 | no |
    | 24 | GO-2026-5021 | CVE-2026-42508, GHSA-5cgq-3rg8-m6cv | `golang.org/x/crypto` | Invoking auth bypass via unenforced @revoked status in golang.org/x/crypto/ssh/knownhosts | 0.52.0 | 2026-05-22 | no |
    | 25 | GO-2026-5023 | CVE-2026-46595 | `golang.org/x/crypto` | Invoking VerifiedPublicKeyCallback permissions skip enforcement in golang.org/x/crypto/ssh | 0.52.0 | 2026-05-22 | no |
    | 26 | GO-2026-5033 | CVE-2026-46598 | `golang.org/x/crypto` | Invoking pathological inputs can lead to client panic in golang.org/x/crypto/ssh/agent | 0.52.0 | 2026-05-22 | no |
    | 27 | GO-2026-6303 | CVE-2026-56854 | `golang.org/x/crypto` | Source-address critical option not enforced for non-public-key auth callbacks in golang.org/x/crypto/ssh | 0.55.0 | 2026-08-28 | no |
    | 28 | GO-2024-2687 | CVE-2023-45288, GHSA-4v7x-pqxf-cx7m | `golang.org/x/net` | HTTP/2 CONTINUATION flood in net/http | 0.23.0 | 2024-04-03 | no |
    | 29 | GO-2024-3333 | CVE-2024-45338, GHSA-w32m-9786-jp63 | `golang.org/x/net` | Non-linear parsing of case-insensitive content in golang.org/x/net/html | 0.33.0 | 2024-12-18 | no |
    | 30 | GO-2025-3503 | CVE-2025-22870, GHSA-qxp5-gwg8-xv66 | `golang.org/x/net` | HTTP Proxy bypass using IPv6 Zone IDs in golang.org/x/net | 0.36.0 | 2025-03-12 | no |
    | 31 | GO-2025-3595 | CVE-2025-22872 | `golang.org/x/net` | Incorrect Neutralization of Input During Web Page Generation in x/net in golang.org/x/net | 0.38.0 | 2025-04-16 | no |
    | 32 | GO-2026-4440 | CVE-2025-47911 | `golang.org/x/net` | Quadratic parsing complexity in golang.org/x/net/html | 0.45.0 | 2026-02-05 | no |
    | 33 | GO-2026-4441 | CVE-2025-58190 | `golang.org/x/net` | Infinite parsing loop in golang.org/x/net | 0.45.0 | 2026-02-05 | no |
    | 34 | GO-2026-4918 | CVE-2026-33814 | `golang.org/x/net` | Infinite loop in HTTP/2 transport when given bad SETTINGS_MAX_FRAME_SIZE in net/http/internal/http2 in golang.org/x/net | 0.53.0 | 2026-05-07 | no |
    | 35 | GO-2026-5025 | CVE-2026-42506 | `golang.org/x/net` | Invoking incorrect handling of namespaced elements in foreign content in golang.org/x/net/html | 0.55.0 | 2026-05-22 | no |
    | 36 | GO-2026-5026 | CVE-2026-39821 | `golang.org/x/net` | Invoking failure to reject ASCII-only Punycode-encoded labels in golang.org/x/net/idna | 0.55.0 | 2026-05-22 | no |
    | 37 | GO-2026-5027 | CVE-2026-42502 | `golang.org/x/net` | Invoking incorrect handling of HTML elements in foreign content in golang.org/x/net/html | 0.55.0 | 2026-05-22 | no |
    | 38 | GO-2026-5028 | CVE-2026-25680 | `golang.org/x/net` | Invoking denial of service when parsing arbitrary HTML in golang.org/x/net/html | 0.55.0 | 2026-05-22 | no |
    | 39 | GO-2026-5029 | CVE-2026-25681 | `golang.org/x/net` | Invoking incorrect handling of character references in DOCTYPE nodes in golang.org/x/net/html | 0.55.0 | 2026-05-22 | no |
    | 40 | GO-2026-5030 | CVE-2026-27136 | `golang.org/x/net` | Invoking duplicate attributes can cause XSS in golang.org/x/net/html | 0.55.0 | 2026-05-22 | no |
    | 41 | GO-2026-5942 | CVE-2026-46600 | `golang.org/x/net` | Parsing an invalid SVCB or HTTPS RR can panic in golang.org/x/net/dns/dnsmessage | 0.56.0 | 2026-07-14 | no |
    | 42 | GO-2026-5024 | CVE-2026-39824 | `golang.org/x/sys` | Invoking integer overflow in NewNTUnicodeString in golang.org/x/sys/windows | 0.44.0 | 2026-05-22 | no |
    | 43 | GO-2026-5970 | CVE-2026-56852 | `golang.org/x/text` | Infinite loop on invalid input in golang.org/x/text | 0.39.0 | 2026-07-14 | all archives |
    | 44 | GO-2024-2611 | CVE-2024-24786, GHSA-8r3f-844c-mc37 | `google.golang.org/protobuf` | Infinite loop in JSON unmarshaling in google.golang.org/protobuf | 1.33.0 | 2024-03-05 | no |

## What a vulnerability scanner reports against v1.13.0

**Standard library: nothing.** No Go standard library advisory applies to any of the nine `v1.13.0`
archives, or to `geth` in either `core-geth-docker` image tarball, even at module level. On
2026-09-17 go1.26.8 was the newest Go 1.26 release, and the last Go 1.26 security release before it
was go1.26.6.

**Three `golang.org/x/crypto` advisories, at module level only.** `v1.13.0` ships
`golang.org/x/crypto` v0.55.0, which is in the affected range of three advisories:

| GO ID | CVE | Packages | Summary (verbatim) | First fixed version | Published |
|---|---|---|---|---|---|
| GO-2026-5932 | none | `golang.org/x/crypto/openpgp` and 6 subpackages | The golang.org/x/crypto/openpgp package is unmaintained, unsafe by design, and has known security issues | no fix | 2026-07-07 |
| GO-2026-6354 | CVE-2026-78662 | `golang.org/x/crypto/ssh` | Prevent DoS on deadlocked undecided channel in golang.org/x/crypto/ssh | 0.56.0 | 2026-09-02 |
| GO-2026-6355 | CVE-2026-56855 | `golang.org/x/crypto/ssh` | Prevent DoS on deadlocked established channel in golang.org/x/crypto/ssh | 0.56.0 | 2026-09-02 |

The affected packages are not compiled into any `v1.13.0` binary. No `golang.org/x/crypto/ssh` or
`golang.org/x/crypto/openpgp` function name occurs in any of them, while other `golang.org/x/crypto`
functions do.

**The two macOS archives show those three at symbol level.** govulncheck reads no symbols from the
`v1.13.0` `osx` and `osx-arm64` binaries. When it reads none, it handles the binary as stripped and
reports every listed vulnerable symbol of every module-level match, so a scan of either archive lists
symbol-level `golang.org/x/crypto/ssh` and `golang.org/x/crypto/openpgp` findings for code the binary
does not contain. Their module-level results are unaffected. The same scanner read symbols from the
`v1.12.x` macOS archives. The `v1.13.0` macOS binaries were linked with `-s`, which removes the symbols it
looks for; [After v1.13.0](#after-v1130) gives the detail and the change.

**Six go-ethereum advisories, all fixed in `v1.13.0`.** Every `v1.13.0` archive matches six
advisories filed against go-ethereum, at symbol level, and the match is structural. The Go module
path is deliberately `github.com/ethereum/go-ethereum`, which is what makes this client a drop-in
downstream. The scanner compares this client's version with go-ethereum's version ranges, and
`v1.13.0` sorts below the go-ethereum releases that fixed these, so each one matches whether or not
its fix is present. Matching on symbol names does not settle it either, because a backported fix adds
its guard inside the same function.
[Release artifacts](2026-09-release-pipeline.md#what-a-vulnerability-scanner-will-say-about-these-artifacts)
describes how to read such a scan. Each match was confirmed by reading the guard at the `v1.13.0`
tag:

| GO ID | CVE | Upstream fix govulncheck names | Issue | Guard in `v1.13.0` |
|---|---|---|---|---|
| GO-2024-2819 | CVE-2024-32972 | v1.13.15 | A header request could pull unbounded data from disk | `ReadHeaderRange` in `core/rawdb/accessors_chain.go` returns on a zero count and caps the freezer read at 2 MB |
| GO-2026-4314 | CVE-2026-22868 | v1.16.8 | KZG proof verification denial of service | `core/txpool/validation.go` returns `ErrKZGVerificationError`, and `eth/fetcher/tx_fetcher.go` disconnects the peer that sent it |
| GO-2026-4315 | CVE-2026-22862 | v1.16.8 | ECIES ciphertext length undercheck | `Decrypt` in `crypto/ecies/ecies.go` rejects a ciphertext too short to hold its public key, MAC and one cipher block, before `symDecrypt` reads it |
| GO-2026-4507 | CVE-2026-26314 | v1.16.9 | secp256k1 coordinates outside the field | `IsOnCurve` in `crypto/secp256k1/curve.go`, the return-value checks in `secp256k1_ext_scalar_mul` in `crypto/secp256k1/ext.h`, and the non-cgo check in `crypto/signature_nocgo.go` |
| GO-2026-4508 | CVE-2026-26313 | v1.17.0 | p2p message memory exhaustion | Messages are held as `rlp.RawList` until validated (`eth/protocols/eth/protocol.go`, `eth/protocols/snap/protocol.go`), and responses are checked against their pending request by `tracker.Fulfil` in both protocols' handlers |
| GO-2026-4511 | CVE-2026-26315 | v1.16.9 | ECIES public key validation in the RLPx handshake | `GenerateShared` in `crypto/ecies/ecies.go` rejects nil or off-curve public keys |

## What this means if you are upgrading

- **Rebuilding a `v1.12.x` release on a newer Go takes the standard library fixes, but not all of what
  upgrading does.** The build needs memsize removed and the KZG build tag left out. It still declares
  `go 1.21`, so two denial-of-service limits stay off unless `GODEBUG` sets them, and it keeps the
  third-party module advisories and everything else `v1.12.23` carries. `v1.13.0` declares `go 1.26.0` and
  is built with go1.26.8.
- **Arm builds of `v1.12.x` carry the most.** Every Arm archive from `v1.12.20` to `v1.12.23` carries
  61 standard library advisories at module level.
- **Go 1.26 is supported until go1.28.0 is released**, under the policy above. A Go security release
  reaches your node only through a Core-Geth release built with it, so follow the
  [`ethereumclassic/core-geth` release line](https://github.com/ethereumclassic/core-geth/releases).

## Verification

Each check below was calibrated so that it could report a negative:

- **The published files.** 59 node archives were downloaded: every archive of `v1.12.20` to
  `v1.12.23` and `v1.13.0`, and the Linux x86_64 archive of each earlier stable `v1.12.x` release.
  Each matched its published `.sha256`, and the 33 from 2026 releases also matched GitHub's
  server-side digest, which the older release assets do not carry. Both `v1.13.0`
  `core-geth-docker` image tarballs matched as well.
- **The source they came from.** Each binary's recorded `vcs.revision` equals the commit its release
  tag points to.
- **The toolchain.** Read from each binary's build information with `go version -m`. No downloaded
  binary was executed.
- **The advisories.** govulncheck v1.7.0 in binary mode, at symbol level and at module level, on all
  41 binaries of `v1.12.20` to `v1.13.0`, against one snapshot of the Go vulnerability database for
  every scan: `vuln.go.dev` as modified 2026-09-15T18:39:25Z, 4,445 entries, SHA-256
  `fe51e51645e6b1e946731cef048339972256c6417de91d4cf3985fc87b4aa89b`.
- **Calibration in both directions.** Every `v1.12.20` to `v1.12.23` binary reports 55 to 61 standard
  library advisories at module level, and a Go binary built with go1.26.1 reports 28, fixed in
  go1.26.2 to go1.26.6. A Go binary built with go1.26.6 reports none.
- **An independent re-derivation.** Each advisory's version ranges were evaluated directly against
  each binary's Go version and module versions, excluding withdrawn entries and applying the same
  operating system and architecture filter as govulncheck. The result matched govulncheck's for all
  41 binaries, including the first fixed Go version of every standard library match.
- **Whether the scanner read symbols.** govulncheck read between 49,101 and 53,753 symbols from 39 of
  the 41 binaries, and none from the two `v1.13.0` macOS binaries described above. The absence of
  `golang.org/x/crypto/ssh` and `golang.org/x/crypto/openpgp` function names was checked in all 41,
  against other `golang.org/x/crypto` names present in every one.
- **Release-date counts.** "Already published that day" uses each advisory's published date.
  Counting instead by the release date of the first Go version that fixes it gives the same figure
  for every archive of every release.
- **The failed build.** At the `v1.12.23` release commit, with go1.26.6 and modules verified against
  its `go.sum`, this command exits 1 with the `blst` compile error, and the same command without the
  `ckzg` tag exits 1 at the link step with the memsize error:

    ```
    GOTOOLCHAIN=local CGO_CFLAGS="-O2 -g -D__BLST_PORTABLE__" go build -p 4 -tags urfave_cli_no_docs,ckzg ./cmd/geth
    ```

**Not measured:**

- The `core-geth-alltools` archives and image tarballs, so the go1.26.8 result for `v1.13.0` is
  established for `geth` only.
- Container images published to a registry, for any release, including `v1.12.x` images built from
  `golang:1.22-alpine`.
- Archives of `v1.12.0` to `v1.12.19` other than Linux x86_64, and govulncheck on any of them.
- Prereleases and drafts.
- Builds on Go 1.23, 1.24, 1.25 or 1.27. The Go release each blocker arrived in comes from Go's
  release notes and source, not from building on it.
- Reachability of any advisory from the client's code paths.
- go-ethereum advisories against the `v1.12.x` binaries. Those record the main module's version as
  `(devel)`, which govulncheck never matches, so their absence is not evidence either way.
- Why the release workflow's `go-version: '1.21'` produced go1.21.10 on 2024-06-10, when go1.21.11
  had been released on 2024-06-04, and how the `v1.12.2` Linux archive came to be built with
  go1.13.4.
- Issues and pull requests in `etclabscore/core-geth` beyond those listed above, which were found by
  search and screened by title.

## After v1.13.0

Each of these is resolved on `main` for the next release. None changes what `v1.13.0` ships or how it
behaves.

- **Dependency maintenance.** `golang.org/x/crypto` moves to v0.56.0, which fixes GO-2026-6354 and
  GO-2026-6355 in `golang.org/x/crypto/ssh`, and `golang.org/x/mod` to v0.40.0, which fixes GO-2026-6179 and
  GO-2026-6180 in `golang.org/x/mod/sumdb`. `v1.13.0` compiles none of those packages into any binary, so
  they are not a concern for `v1.13.0`. GO-2026-5932, in `golang.org/x/crypto/openpgp`, has no fixed version
  and is not compiled in either.
- **A weekly scan.** `.github/workflows/govulncheck.yml` scans `main` every week, in source mode and against a
  `geth` binary built with the release Go version, and keeps one issue current with anything that needs a
  maintainer. The six go-ethereum advisories above are expected; one of them disappearing is reported too.
- **The macOS symbol table.** The release build linked the macOS binaries with `-s`, which removes Go's
  symbols from the Mach-O symbol table. govulncheck looks up the `go:func.*` symbol there and, without it,
  reports at module level, although the `__gopclntab` function table it reads next is present. `build/ci.go`
  now links them with `-w`, which strips only the debugging information. A go1.26 macOS test binary linked
  with `-w` keeps `go:func.*` and every Go symbol; one linked with `-s` keeps none.

## Disposition of the build finding, 17 September 2026

**The finding that the `v1.12.x` source does not build on a supported Go now has a proposed remedy in the
repository that publishes that line.**
[Pull request #702](https://github.com/etclabscore/core-geth/pull/702) at `etclabscore/core-geth`, opened
17 September 2026 and unmerged when this was written, makes the three changes this audit measured as necessary:
it removes `github.com/fjl/memsize`, whose `//go:linkname` use Go 1.23 restricted; it moves
`github.com/supranational/blst` from v0.3.11, which Go 1.24 rejects for a cgo alias receiver, to v0.3.16; and it
raises the `go.mod` directive from `go 1.21` to `go 1.24.0`. It also moves the `# version:golang` pin in
`build/checksums.txt` from 1.22.1 to 1.25.12 and updates `golang.org/x/crypto` from v0.17.0 to v0.48.0.

**Two measurements to repeat rather than assume, if that release ships.** Go 1.24 left support on 10 February
2026 and Go 1.25 on 19 August 2026, each when two newer major releases existed, so a binary built from either
would still carry standard library advisories, fewer than Go 1.21 carries. And the advisory count for any such
release has to be measured from its published archives, the way the counts in this document were, rather than
inferred from the `go.mod` directive: this audit found two toolchains per release where the configuration
implied one.

This paragraph records a disposition, which is what an audit finding is supposed to acquire. It is not a claim
about anyone's reasons.

## Supporting this work

The modernization this document records was carried out by [White B0x](https://whiteb0x.com) as
unfunded public-goods work for Ethereum Classic. Mining pools, centralized exchanges, issuers of
Ethereum Classic financial products, Etchash mining hardware manufacturers and large holders all
depend on this client. If your operation relies on Ethereum Classic, please help fund its
maintenance: contact <donations@ethereumclassic.com>, or donate directly to the address below, which
receives on any EVM-compatible chain:

``` { .text .copy }
0x86FE8d331A4B984B57d3e92C6F4cb9C881eC9B04
```
