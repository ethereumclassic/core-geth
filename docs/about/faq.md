---
description: "Short answers with links to the record: which repository publishes Core-Geth, which version to run, node key rotation, MESS, public RPC endpoints and how to verify a download."
---

# Core-Geth FAQ

Short answers, each linking to the record it rests on. If you run a node and want the ordered path instead,
read [Node operators: start here](../operators.md).

## Which repository publishes Core-Geth releases?

[`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth/releases), in the Ethereum Classic
community's GitHub organization, is the canonical release pipeline for Core-Geth, where the client is developed and released now. `v1.13.0`, published September 14, 2026, is the first release cut from it.
Archives published under the previous `etclabscore` namespace are not built from this source and do not carry the
fixes released here. [Project history](project-history.md) covers how maintenance moved.

## Where do I download Core-Geth?

From [coregeth.com/downloads](https://coregeth.com/downloads/), which has the archive for each platform, its checksum and
the commands. [Installation](../getting-started/installation.md) has a route for every platform, and
[How do I verify a download?](#how-do-i-verify-a-download) confirms the archive is the one this repository built.

## Is `etclabscore/core-geth` still maintained?

The ETC Cooperative owns that repository, and its board has communicated its intent to archive it. The
Cooperative has been in maintenance mode since the end of 2024 and is due to dissolve.
[The ETC Cooperative transition](../etc-cooperative-transition.md) lists where each service it maintained
continues.

## Which Core-Geth version should I run?

[The latest v1.13 release](https://coregeth.com/downloads/). Every `v1.12.20` to `v1.12.23` archive was built with an
out-of-support Go release and carries 55 to 61 Go standard library advisories, 43 to 47 of them present in the binary,
that `v1.13.0` does not ([Go toolchain audit](../audits/2026-09-go-toolchain.md)). `v1.12.23` also lacks the GraphQL depth
limit and carries the `eth_syncing` regression, and releases before `v1.12.21` carry all six client CVEs. The
[migration guide](../tutorials/v1.13.0-migration.md) covers the upgrade and how to roll back.

## Why does the upgrade ask me to rotate my P2P node key?

[CVE-2026-26315](https://github.com/ethereum/go-ethereum/commit/46bee92f9) leaks bits of the node key across repeated
handshakes. The fix stops new leaks but cannot recall old ones, so go-ethereum's advisory
([GHSA-m6j8-rg6r-7mv8](https://github.com/ethereum/go-ethereum/security/advisories/GHSA-m6j8-rg6r-7mv8)) recommends
rotating after the upgrade. A key that served `v1.12.20` or earlier should be treated as exposed. The
[migration guide](../tutorials/v1.13.0-migration.md#rotate-the-p2p-node-key) gives the steps.

## Has anything gone wrong with `v1.13.0` since it shipped?

No issue reporting a defect in `v1.13.0` has been opened in
[this repository](https://github.com/ethereumclassic/core-geth/issues) since it shipped, and operators who moved
to it have reported stable nodes. It was run on Ethereum Classic mainnet and Mordor throughout its development,
from February 2026 to the release. [The record behind the release](../release-reports/v1.13.0-record.md) sets out
what has been claimed about it and what the record shows.

## Is `ethereumclassic/core-geth` the same Core-Geth, or a fork?

It is the same client, in the Ethereum Classic community's own organization, and the canonical home for Core-Geth now.
The repository was created on December 21, 2024 from `etclabscore/core-geth` at commit `7ef3ecd7a`, and holds the full
history up to that point on the
[`archive-etclabscore-2024-12`](https://github.com/ethereumclassic/core-geth/tree/archive-etclabscore-2024-12) branch.
[Project history](project-history.md) traces the move, and
[the record](../release-reports/v1.13.0-record.md) answers the specific claims made about it.

## Which Ethereum Classic client should I run?

Core-Geth `v1.13.x` is what this project maintains and recommends today. A network where most of the hashrate
runs one client has a single point of failure whatever that client is, so running more than one implementation
across your fleet is sound practice.

## Is MESS on, and should it be?

MESS (ECBP-1100) changes which of two competing chains a node keeps during a deep reorganization, never whether a block
is valid. `--mess` turns it on and `--mess=false` turns it off; without either, the node uses the setting its release
ships with. Once the node is in sync with at least five peers, `admin.messStatus().enabled` shows whether it is on
(`admin.ecbp1100Status().enabled` before `v1.13.1`). [MESS](../operate/mess.md) sets out the choice, and the
[confirmation calculator](../guides/mess-calculator.md) shows what it means for deposit confirmations.

## What happened to `etc.rivet.link`?

It went offline in late August 2026 and no longer resolves. Until this project's own endpoints are answering, use the
lists on ChainList, for [Ethereum Classic, chain 61](https://chainlist.org/chain/61) and
[Mordor, chain 63](https://chainlist.org/chain/63).
[Public JSON-RPC](../etc-cooperative-transition.md#public-json-rpc) explains what is replacing it.

## How do I verify a download?

Check the archive against its published `.sha256`, then verify its build attestation with
`gh attestation verify <file> --repo ethereumclassic/core-geth`. For a pulled container image, use
`gh attestation verify oci://ghcr.io/ethereumclassic/core-geth:latest --repo ethereumclassic/core-geth`. Both need the
[GitHub CLI](https://cli.github.com/), signed in with `gh auth login`. `go version -m geth` prints the Go version and the
commit a binary was built from without running it.
[Verifying a download](../release-reports/v1.13.0.md#verifying-a-download) gives the commands.

## How do I check any of this for myself?

Every claim on these pages links to the record it rests on, and the records are public: commits, pull requests,
release files, advisory databases and the published archives. [Do not trust this page. Verify
it.](../release-reports/v1.13.0-record.md#do-not-trust-this-page-verify-it) carries a prompt you can paste into
any assistant that fetches URLs, which reads every document here plus the article written about the release and
reports where they agree and disagree.

## Who maintains Core-Geth, and where do I report a security issue?

The [Ethereum Classic DAO](https://ethereumclassicdao.org), a Wyoming DAO LLC, is the legal entity that maintains it,
together with Ethereum Classic's core developers and community contributors. The work is done in the
`ethereumclassic` organization. Report a security issue privately through this repository's
[private advisories](https://github.com/ethereumclassic/core-geth/security/advisories) or by email to
<security@ethereumclassic.net>, never in a public issue or on Discord. Mining pools, exchanges and service providers
should use that address as their point of contact as the ETC Cooperative dissolves. For anything else, the Ethereum
Classic Core Developers run a [Discord](https://coregeth.com/discord).
