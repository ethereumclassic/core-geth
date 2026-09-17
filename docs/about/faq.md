---
description: "Short answers with links to the record: which repository publishes Core-Geth, which version to run, node key rotation, MESS, public RPC endpoints and how to verify a download."
---

# Questions and answers

Short answers, each linking to the record it rests on. If you run a node and want the ordered path instead,
read [Node operators: start here](../operators.md).

## Which repository publishes Core-Geth releases?

[`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth/releases), in the Ethereum Classic
community's GitHub organization. `v1.13.0`, published 14 September 2026, is the first release cut from it.
Archives published under the previous `etclabscore` namespace are not built from this source and do not carry the
fixes released here. [Project history](project-history.md) covers how maintenance moved.

## Is `etclabscore/core-geth` still maintained?

The ETC Cooperative owns that repository, and its board has communicated its intent to archive it. The
Cooperative has been in maintenance mode since the end of 2024 and is scheduled to dissolve by the end of 2026.
[The ETC Cooperative transition](../etc-cooperative-transition.md) lists where each service it maintained
continues.

## Which version should I run?

`v1.13.0`. Every release in the `v1.12.x` line carries at least one unpatched CVE, and each `v1.12.20` to
`v1.12.23` archive also carries 55 to 61 Go standard library advisories that `v1.13.0` does not. The
[migration guide](../tutorials/v1.13.0-migration.md) covers the upgrade and how to roll back, and the
[Go toolchain audit](../audits/2026-09-go-toolchain.md) lists every advisory.

## Why does the upgrade ask me to rotate my P2P node key?

Because CVE-2026-26315 is an oracle against the node key itself: an invalid-curve ephemeral key in the RLPx
handshake reached ECDH before failing, which leaks bits of the key across repeated handshakes. The fix stops new
leaks and cannot recall bits that already leaked, so go-ethereum's own advisory
([GHSA-m6j8-rg6r-7mv8](https://github.com/ethereum/go-ethereum/security/advisories/GHSA-m6j8-rg6r-7mv8))
recommends rotating after the upgrade. A key that served a `v1.12.20` or earlier node should be treated as
exposed. The [migration guide](../tutorials/v1.13.0-migration.md#rotate-the-p2p-node-key) gives the procedure,
including which peer lists to update.

## Has anything gone wrong with `v1.13.0` since it shipped?

No issue reporting a defect in `v1.13.0` has been opened in
[this repository](https://github.com/ethereumclassic/core-geth/issues) since it shipped, and operators who moved
to it have reported stable nodes. It was run on Ethereum Classic mainnet and Mordor throughout its development,
from February 2026 to the release. [The record behind the release](../release-reports/v1.13.0-record.md) sets out
what has been claimed about it and what the record shows.

## Is `ethereumclassic/core-geth` a fork of the "real" client?

It is the same client, in the Ethereum Classic community's own organization. The repository was created on
21 December 2024 from `etclabscore/core-geth` at commit `7ef3ecd7a`, holds the full history up to that point on
the [`archive-etclabscore-2024-12`](https://github.com/ethereumclassic/core-geth/tree/archive-etclabscore-2024-12)
branch, and carries the `v1.12.23` p2p hardening with its original authorship.
[Project history](project-history.md) traces the move, and
[the record](../release-reports/v1.13.0-record.md) answers the specific claims made about it.

## Which Ethereum Classic client should I run?

Core-Geth `v1.13.x` is what this project maintains and recommends today, and
[Fukuii](https://fukuii.org) is its intended successor once it releases. A network where most of the hashrate runs
one client has a single point of failure whatever that client is, so running more than one implementation across
your fleet is sound practice.

## Is MESS on by default, and should I turn it off?

`v1.13.0` ships MESS (ECBP-1100) on by default, and `--mess=false` turns it off. MESS changes which of two
competing chains a node prefers during a deep reorganization and never whether a block is valid. [MESS](../operate/mess.md)
explains the choice and its trade-offs, and the [confirmation calculator](../guides/mess-calculator.md) shows what
it means for deposit confirmations.

## What happened to `etc.rivet.link`?

It went offline in late August 2026 and no longer resolves. ethereumclassic.org describes it as provided by Rivet
under contract with the ETC Cooperative. Until this project's own endpoints are answering, use the lists on
ChainList, for [Ethereum Classic, chain 61](https://chainlist.org/chain/61) and
[Mordor, chain 63](https://chainlist.org/chain/63).
[Public JSON-RPC](../etc-cooperative-transition.md#public-json-rpc) explains what is being stood up to replace it,
and why an endpoint that products depend on is run from the organization's repositories rather than one
company's.

## How do I verify a download?

Check the archive against its published `.sha256`, then verify its build attestation with
`gh attestation verify <file> --repo ethereumclassic/core-geth`. `go version -m geth` prints the Go version and
the commit a binary was built from without running it. [Verifying a download](../release-reports/v1.13.0.md#verifying-a-download)
gives the commands, and [Release artifacts](../audits/2026-09-release-pipeline.md) records what each published
archive actually contains.

## How do I check any of this for myself?

Every claim on these pages links to the record it rests on, and the records are public: commits, pull requests,
release files, advisory databases and the published archives. [Do not trust this page. Verify
it.](../release-reports/v1.13.0-record.md#do-not-trust-this-page-verify-it) carries a prompt you can paste into
any assistant that fetches URLs, which reads every document here plus the article written about the release and
reports where they agree and disagree.

## Who maintains Core-Geth, and where do I report a security issue?

Long-time Ethereum Classic core developers maintain it in the `ethereumclassic` organization, and the work is
done for the [Ethereum Classic DAO](https://ethereumclassicdao.org). Report a security issue privately through
this repository's [private advisories](https://github.com/ethereumclassic/core-geth/security/advisories) or by
email to <security@ethereumclassic.com>, never in a public issue. Mining pools, exchanges and service providers
should use that address as their point of contact as the ETC Cooperative dissolves. For anything that is not a
vulnerability, the Ethereum Classic Core Developers' [Discord](https://ethereumclassic.com/discord) reaches the same people, and a vulnerability
does not belong there because a public channel discloses it to everyone at once.
