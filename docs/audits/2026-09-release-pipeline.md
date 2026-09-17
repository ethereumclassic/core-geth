---
description: "What the published Core-Geth archives actually contain: platform floors, architectures and provenance, measured from the files themselves."
---

# Release artifacts: v1.12.x archive to v1.13.0

This document records what the published `v1.12.x` releases actually contain, what
`v1.13.0` ships instead, and why an operator on the older line should move. It is a
companion to `2026-03-security-audit.md` and `2026-08-security-followup.md`, which
cover source-level CVE remediation, and to `2026-08-dependency-modernization.md`,
which covers the toolchain and module graph. This one covers the artifacts those
produce: the files an operator downloads and runs.

**Baseline:** the four published releases `v1.12.20` through `v1.12.23`, measured by
downloading and inspecting the archives rather than by reading the configuration that
built them. `v1.12.20` (June 2024) is the last release before this repository was
created from the archived development line; the other three were published in 2026.

**Work carried out by:** [White B0x](https://whiteb0x.com)

**On this page:** [In short](#in-short) · [Why the artifacts needed a pass of their own](#why-the-artifacts-needed-a-pass-of-their-own) · [The platform floor was lost across the v1.12.x releases](#finding-the-platform-floor-was-lost-across-the-v112x-releases) · [The macOS archive contains a binary most of its downloaders cannot run](#finding-the-macos-archive-contains-a-binary-most-of-its-downloaders-cannot-run) · [The v1.12.x line no longer builds against a C23 compiler](#finding-the-v112x-line-no-longer-builds-against-a-c23-compiler) · [The v1.12.x releases were published from outside this organization](#finding-the-v112x-releases-were-published-from-outside-this-organization) · [No container image was ever published from this repository](#finding-no-container-image-was-ever-published-from-this-repository) · [The release files did not identify themselves](#finding-the-release-files-did-not-identify-themselves) · [What a vulnerability scanner will say about these artifacts](#what-a-vulnerability-scanner-will-say-about-these-artifacts) · [How v1.13.0 sets its platform floor](#how-v1130-sets-its-platform-floor) · [What this means if you are upgrading](#what-this-means-if-you-are-upgrading) · [Verification](#verification) · [Outstanding](#outstanding) · [Supporting this work](#supporting-this-work)

## In short

**Upgrade to [v1.13.0 or later](https://github.com/ethereumclassic/core-geth/releases/latest).** Measured by downloading them, the published `v1.12.20` to
`v1.12.23` archives raised the glibc version a Linux host needs without saying so, shipped an Apple Silicon
binary under a name that does not say so, and carried nothing that identifies what built them. Every `v1.13.0`
archive and image ships with a checksum and a build attestation, and names its platform.

## Why the artifacts needed a pass of their own

A release is not what the build configuration says it builds. It is what the archive
contains, and the two can disagree for a long time without anything reporting it.

The defects below are invisible in the source, in the build log, and in the release
notes. The build succeeds, the tests pass, the archive is well-formed and correctly
named, and the checksum matches. They are visible only by opening the published file
and reading the binary inside, which is what this pass did, and which is why they had
survived across releases.

## Finding: the platform floor was lost across the v1.12.x releases

A binary linked against glibc runs on that version or newer, never older. The
requirement is recorded in the artifact, so it can be read directly:

| Release | Published | Highest glibc symbol required |
|---|---|---|
| `core-geth-linux-v1.12.20.zip` | 2024-06 | **GLIBC_2.17** |
| `core-geth-linux-v1.12.21.zip` | 2026-03 | **GLIBC_2.34** |
| `core-geth-linux-v1.12.22.zip` | 2026-03 | **GLIBC_2.34** |
| `core-geth-linux-v1.12.23.zip` | 2026-08 | **GLIBC_2.34** |

**Every Arm archive moved the same way, in the same release.** These are built by a
separate cross-compilation path, so they are an independent instance of the defect
rather than the same artifact measured twice:

| Archive | v1.12.20 | v1.12.23 |
|---|---|---|
| `arm5` | 2.28 | **2.34** |
| `arm7` | 2.28 | **2.34** |
| `arm64` | **2.17** | **2.34** |

The floor did not drift. It moved once, between `v1.12.20` and `v1.12.21`, and the
mechanism is visible in the symbols themselves. Every symbol that requires 2.34 in the
later binaries belongs to one set:

```
__libc_start_main   dlopen  dlsym  dlclose  dlerror
pthread_create  pthread_join  pthread_detach  pthread_key_create  ...
```

glibc 2.34 merged `libpthread` and `libdl` into `libc`, so any binary linked against
glibc 2.34 or newer acquires that whole family at that version. The later releases were
therefore built on a host whose glibc was 2.34 or newer, while `v1.12.20` was not. No
change in the source causes this, and nothing in a build log reports it.

Systems that can run `v1.12.20` but cannot run any later `v1.12.x` release:

| Distribution | glibc |
|---|---|
| Ubuntu 20.04 LTS | 2.31 |
| Debian 11 | 2.31 |
| RHEL / Rocky / AlmaLinux 8 | 2.28 |
| Amazon Linux 2 | 2.26 |

Anything older than 2.34 is affected; those are the widely deployed cases. An operator
on one of them who followed the security guidance to upgrade could not run the result.

**`v1.13.0` requires GLIBC_2.17 on x86_64**, restoring every system above and matching
the last release before this repository was created. Each Arm target is restored to the
floor `v1.12.20` shipped for it (2.17 on arm64, 2.28 on the 32-bit targets), and every
one of them is checked at build time.

**The difference that matters is not the number but how it is held.** `v1.12.20`'s 2.17
was a byproduct of whichever machine built it, which is exactly why it was lost without
anyone deciding to lose it. `v1.13.0` states the floor as a build input and **fails the
release** if the artifact does not meet it.

## Finding: the macOS archive contains a binary most of its downloaders cannot run

| Release | `core-geth-osx-<version>.zip` contains |
|---|---|
| v1.12.20 | Mach-O 64-bit **arm64** |
| v1.12.21 | Mach-O 64-bit **arm64** |
| v1.12.22 | Mach-O 64-bit **arm64** |
| v1.12.23 | Mach-O 64-bit **arm64** |

Unlike the glibc floor, this one predates this repository. It begins at `v1.12.20`:
`v1.12.19` (February 2024) contains an x86_64 binary, and every release since contains an
Arm64 one under the same name.

The archive name records no architecture, and earlier releases under that name were
x86_64. The build selected a floating image label, which changed meaning when the
provider moved it to Apple Silicon. An Intel Mac user following the same download path
as always receives a binary their machine cannot execute, and the failure appears at
execution rather than at download.

**`v1.13.0` publishes both architectures and names them.** The `osx` name keeps the
architecture earlier releases shipped, so an existing download path does not silently
change meaning a second time, and Apple Silicon is published beside it under its own
name.

## Finding: the v1.12.x line no longer builds against a C23 compiler

This affects anyone who builds rather than downloads, including a distributor producing
their own packages.

`v1.12.x` carries `blst` v0.3.11 (confirmed at both `v1.12.20` and `v1.12.23`), whose
public header defines `bool` with a `typedef`, guarded only against C++ and pre-C99
compilers. C23 makes `bool` a language keyword, so that definition became illegal. The
build fails inside a vendored dependency's header, at a line no change in this repository
touched.

**Compilers have been moving their default to C23 independently of each other**, so this
is not confined to one platform. This project met it first in its Alpine-based Docker
build, recorded in `2026-03-security-audit.md`; it is also reported against the compiler
bundled with current Windows images. Any toolchain that defaults to C23 reaches the same
line.

**`v1.13.0` carries `blst` v0.3.17**, whose header guards the definition on the standard
version and uses the keyword under C23. The upgrade began in the CVE work and was carried
forward by the dependency pass in `2026-08-dependency-modernization.md`; it is recorded
here because it is a concrete way the older line has stopped building rather than merely
aged.

## Finding: the v1.12.x releases were published from outside this organization

Every `v1.12.x` release was published from `etclabscore/core-geth`, the repository this
one was created from, which sits outside this organization. `v1.12.23` was cut there in
August 2026. If you are running a `v1.12.x` binary today, that is where it came from, and
three things follow:

- **It was produced by a pipeline this project does not run and had not reviewed.** That
  is how the platform floor above moved across three releases with nobody observing it.
- **Its checksum does not close the gap.** Each archive is published with a `.sha256`
  beside it, which establishes that the download was not corrupted in transit. It says
  nothing about what produced the file, because whoever publishes the artifact publishes
  the checksum.
- **Nothing binds it to the source it claims to come from.** That line publishes no
  signature and no build attestation, so the link between a release archive and the commit
  it was built from rests entirely on the publisher's assertion.

**Of the three findings above, the floor regression is the only one drawn on the handoff
boundary, and there it is exact.** `7ef3ecd7a` (2024-12-16) is the last commit on the
predecessor's master before this repository was created on 2024-12-21; the tag
`archive/etclabscore-2024-12` marks it. `v1.12.20`, published June 2024 and the last
release before that boundary, carries a 2.17 Linux floor and the Arm floors `v1.13.0`
restores. All three releases carrying the regression were published in 2026, well after
maintenance had moved here. **Nothing that was handed over introduced it.**

**The other two findings are older, and are not drawn on that boundary.** The macOS archive
has contained an arm64 binary under an architecture-free name across all four releases
measured, and `blst` v0.3.11 is carried by the whole `v1.12.x` line, confirmed at both
ends of it. Both predate the handoff, and neither follows from it.

**`v1.13.0` is published from this repository, and every archive carries a build
attestation.** The client, the build and the artifacts are back in one place the project
controls, and each file can be traced to the run that produced it rather than resting on
anyone's word.

The attestation is minted against a short-lived certificate issued to the workflow run
itself, so no signing key is stored anywhere and none can be stolen. It records which
workflow, at which commit, produced a given file: the question a checksum cannot answer.
You can check it yourself against a downloaded archive:

```bash
gh attestation verify core-geth-linux-v1.13.0.zip --repo ethereumclassic/core-geth
```

What the guarantee still rests on is control of the release tag, recorded below.

## Finding: no container image was ever published from this repository

`2026-08-dependency-modernization.md` records the cause: `build/ci.go` carried a
complete image-publishing implementation that nothing ever called. The images under
the previous namespace came from a registry-side integration configured outside the
source tree, so the capability was absent rather than broken, and nothing in the
repository looked wrong.

**`v1.13.0` publishes images from this repository**, for `linux/amd64` and
`linux/arm64`, each built on a native runner rather than under emulation and merged
into a multi-architecture manifest. Two variants are published: the client alone,
and `alltools-` carrying the full set of executables.

**`:latest` is reserved for full releases.** It is what a bare `docker pull`
resolves to, so a release candidate must never take it. An operator who omits a
tag is asking for the current stable client, not the newest thing that exists. The
tag is applied only when the version carries no prerelease suffix:

```
v1.13.0        -> also tagged :latest
v1.13.0-rc1    -> published under its own name only
```

That rule is explicit rather than inferred. The mechanism it replaces
(`docker/metadata-action`'s `latest=auto`) reads as though it withholds the moving
tag from a prerelease and does not: `auto` keys off the tag-ref rule, which has no
concept of one. Measured during a pipeline rehearsal, where a build named
`pipeline-test` was published as `:latest`. A defect of this shape is invisible in
a green pipeline and visible only in what an untagged pull returns.

## Finding: the release files did not identify themselves

Measured on the published `v1.12.20` archives:

- **No license text.** Each archive held only its binary. The binaries are distributed under the
  GPL-3.0 in `COPYING`, which asks that recipients get a copy of the license with the program.
- **No version information on Windows.** `geth.exe` carried an application manifest and no version
  resource, so its properties named no product, company or version.

`v1.13.0` changes both, and adds what the images need:

- Every archive carries `COPYING` beside its binaries.
- Every Windows executable carries version information naming Core-Geth, its version, the Ethereum
  Classic DAO LLC and the copyright holders, and a manifest that identifies the executable and runs it
  at the caller's privilege level. The release fails if Windows reads no version information back from
  an executable.
- The container images name their source, version, revision, vendor and documentation in their labels,
  and the multi-arch index carries its description, source, vendor and documentation as annotations,
  which is where GitHub reads a multi-arch image's description from.

## What a vulnerability scanner will say about these artifacts

**A scanner run against a `v1.13.0` binary reports several go-ethereum advisories
that are already fixed in it.** This is expected, and it is worth understanding
before it is read as a finding.

The Go module path is deliberately `github.com/ethereum/go-ethereum`. That is what
makes this a drop-in downstream, and changing it would break every consumer. A
binary scanner reads the module version out of the build and resolves it to a
go-ethereum pseudo-version, which sorts *below* every upstream tag that fixed
anything. So each upstream advisory fixed above that point matches, regardless of
whether the fix is present in this source.

**Tagging sharpens the match rather than resolving it.** A clean `v1.13.0` sorts
below the upstream tags that carry those fixes, so the report does not improve with
the release.

Two consequences worth holding:

- **An exit status of 0 from such a scan is a change, not a pass.** Findings are
  expected here. An advisory *disappearing* is as much a signal as one appearing.
  It means either the adjudication has gone stale or the artifact scanned is not
  the one that was adjudicated.
- **Each match is confirmed by reading the guard in this tree**, not by comparing
  version strings. [Go toolchain](2026-09-go-toolchain.md#what-a-vulnerability-scanner-reports-against-v1130)
  records the adjudicated set and the reasoning.

**The client's own `version-check` command is a separate matter and used to be
wrong in the other direction.** It queries go-ethereum's feed, whose patterns begin
`Geth/` without a leading anchor, and this client identifies as `Core-Geth/`,
which contains `Geth`. An operator running the built-in checker on a fresh release
was told, falsely, that it carried a High-severity advisory. Fixed by requiring the
match to begin at the start of the version string; a genuine go-ethereum version
still matches, and this client no longer matches itself.

## How v1.13.0 sets its platform floor

Pinning to a specific build image is an improvement and not a solution: the floor is
still inherited, and images are retired on their own schedule, which is how a floor moves
with nobody deciding to move it.

`v1.13.0` states the floor instead of inheriting it. Linux release archives are built
with a cross-compilation toolchain that targets a chosen glibc directly, pinned to an
exact release and verified against a checksum recorded in `build/checksums.txt` alongside
the Go pin.

**The release is gated on the measurement, not on the intent.** The build fails, rather
than publishing, if the artifact's floor rises above target or if the C23 symbol
redirections that raise it reappear. The regression above shipped three times because
nothing checked; a fix with no check is a defect waiting to recur.

**macOS is set the same way.** The macOS archives target macOS 12, the oldest version Go 1.26 runs on,
through the deployment target rather than whatever the build machine's macOS would stamp, and the
release fails if any binary declares a different minimum.

### Why the compiler is the lever

The obvious alternative (build without cgo and link statically) is not available in
either line. Two packages require it unconditionally:

- `github.com/ethereum/evmc/v7/bindings/go/evmc` excludes all of its Go files when cgo
  is disabled.
- `consensus/lyra2` defines its hashing entry point only in a cgo file that compiles the
  algorithm's C sources, with no alternative implementation.

So a C toolchain is always involved, and the only question is which one and what it
targets. That is what makes controlling the compiler the remedy rather than one option
among several.

## What this means if you are upgrading

- **If you run a distribution with glibc older than 2.34** (Ubuntu 20.04, Debian 11,
  RHEL 8, Amazon Linux 2 among them), no `v1.12.x` release after `v1.12.20` runs on your
  system. `v1.13.0` does.
- **If you run an Intel Mac**, none of the four `osx` archives measured here is for your
  machine. `v1.13.0` publishes one that is, under the same name.
- **If you build from source on Windows**, the `v1.12.x` line no longer builds on a
  current image. `v1.13.0` does.
- **Everyone else**: the `v1.13.0` Linux binary has a lower glibc requirement than any
  `v1.12.x` release after `v1.12.20`, so an upgrade cannot reduce the set of systems it
  runs on.

Nothing about how the client is invoked changes, and the executable is still named
`geth`. The migration guide covers the upgrade itself.

## Verification

Each check below was calibrated so that it could report a negative:

- Platform claims read from the **published release archives**, downloaded and unpacked,
  rather than inferred from build configuration.
- glibc requirements read from each binary's dynamic symbol table with `objdump -T`; the
  2.34 attribution confirmed by listing the symbols carrying that version rather than by
  reading the maximum alone.
- Architecture read from the Mach-O header.
- The `v1.12.20` Windows resources read from the published `geth.exe` with Go's `debug/pe`, and each
  archive's contents listed from the published zip.
- The floor gates tested in both directions, against real artifacts rather than synthetic
  values: each accepts the corresponding `v1.13.0` build and the `v1.12.20` baseline, and
  each rejects the `v1.12.23` artifact for its own target.
- The `v1.13.0` build verified through the same entry point the release uses, with the
  cryptographic dependency affected by the C23 change actually linked, rather than
  through a reduced build that would not have exercised it.

**One claim is not a measurement of ours, and is marked rather than blended in.** The
`blst` headers were read directly at both versions, and the version each line carries was
confirmed at the `v1.12.20` and `v1.12.23` release tags. That part is measured, as is the
Alpine Docker failure this project met itself. That current Windows images also default to
C23 is taken from a reported build failure against `blst.h:27` and was not reproduced on a
Windows image here. The remedy does not depend on which toolchain you meet it with:
`v1.13.0` carries a header that is correct under either default.

## Outstanding

- **The container images are attested, not signed.** Provenance is recorded for the published
  image index and for each image tarball. Signing an image is a separate mechanism, and it has not
  been applied here.
- **The Windows and macOS binaries are not code-signed.** Windows names their publisher as unknown,
  and macOS blocks a downloaded binary until it is allowed.
  [Installation](../getting-started/installation.md) says how to check an archive and proceed.
- **Container package visibility is a registry setting, not a repository one.** A registry
  creates a new package private by default, so the first published image is unreachable
  until someone makes it public, and making it public is the moment every tag it carries
  starts serving real traffic. Check what `:latest` points at before flipping it, not
  after.
- **The container runtime base image is a floating tag**, carried forward from the dependency
  pass, where the same objection is recorded. The builder names an exact Go patch.

## Supporting this work

The release work this document records was carried out by [White B0x](https://whiteb0x.com) as
unfunded public-goods work for Ethereum Classic. Mining pools, centralized exchanges, issuers of
Ethereum Classic financial products, Etchash mining hardware manufacturers and large holders all
depend on this client. If your operation relies on Ethereum Classic, please help fund its
maintenance: contact <donations@ethereumclassic.com>, or donate directly to the address below, which
receives on any EVM-compatible chain:

``` { .text .copy }
0x86FE8d331A4B984B57d3e92C6F4cb9C881eC9B04
```
