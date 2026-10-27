---
title: Releases
description: "Every Core-Geth release published from this repository, with its release notes, upgrade notes, audits and reports, and which release lines are maintained."
---

# Releases

Core-Geth releases are published at
[`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth/releases). **Run the
[latest release](https://github.com/ethereumclassic/core-geth/releases/latest).**
[Installation](getting-started/installation.md) downloads and verifies it, and
[Upgrading](upgrading.md) moves a running node to it.

## Release lines

| Line | Status |
| --- | --- |
| `v1.13.x` | Maintained. Published from `ethereumclassic/core-geth` |
| `v1.12.x` and earlier | Not maintained here. Published from the previous repository, `etclabscore/core-geth`, and affected by the vulnerabilities the [v1.13.0 audits](#v1130) document |

## Every release

One row per release, newest first. Each row links what was published for that release.

| Release | Date | Release notes | Upgrading | Audits and reports |
| --- | --- | --- | --- | --- |
| [`v1.13.0`](https://github.com/ethereumclassic/core-geth/releases/tag/v1.13.0) | 2026-09-14 | [Release report](release-reports/v1.13.0.md) | [Leaving the v1.12.x line](upgrading.md#leaving-the-v112x-line) | [Listed below](#v1130) |

## v1.13.0

The first release published from this repository.

| Document | What it covers |
| --- | --- |
| [Release report](release-reports/v1.13.0.md) | What the release fixes and changes, the files it publishes, and how to verify a download |
| [Migrating to v1.13.0](tutorials/v1.13.0-migration.md) | The upgrade from `v1.12.x`, with a guide for Linux, macOS, Windows and Docker |
| [March 2026 security audit](audits/2026-03-security-audit.md) | The six CVEs and the GraphQL denial of service, with the per-release breakdown |
| [August 2026 security follow-up](audits/2026-08-security-followup.md) | `v1.12.23` measured at the tag against the advisory records: what it fixed and what it left open |
| [Dependency and toolchain modernization](audits/2026-08-dependency-modernization.md) | What changed underneath the code between the December 2024 archive and `v1.13.0`: the Go toolchain, 83 modules, the linter |
| [Release artifacts](audits/2026-09-release-pipeline.md) | What the published archives actually contain: platform floors, architectures, provenance |
| [Go toolchain](audits/2026-09-go-toolchain.md) | Which Go toolchain built each published archive, and the Go standard library advisories each one carries |
| [The record behind the release](release-reports/v1.13.0-record.md) | Questions raised about the release after it shipped, answered from commit metadata, pull requests, release files and advisory records |

The two security audits and the modernization report describe the source. **The release artifacts
and Go toolchain reports describe the files you download**, which can disagree with the source: a
release is not what the build configuration says it builds.

## In preparation

`v1.13.1` is the next release in the line. Like every release in it, it is built with the newest
patch release of its Go version as of the day it is cut. When it ships, its release notes and reports
are added above.

## Earlier

| Document | What it covers |
| --- | --- |
| [2021 chain split postmortem](postmortems/2021-08-22-split-postmortem.md) | go-ethereum's account of the minority chain split of August 2021, kept from upstream |
