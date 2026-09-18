---
description: "Who built Core-Geth and who maintains it: the lineage from go-ethereum through multi-geth, and the organizations that funded and ran the client from 2018 to today."
---

# Project history

**The code and the maintaining team have separate lineages, and neither is a straight line.** This page traces
both, because knowing where a consensus client comes from is part of deciding whether to run it.

## The code

go-ethereum began in December 2013. [multi-geth](https://github.com/multi-geth/multi-geth) forked it in March
2018 to support several networks from a single client, holding chain rules as configuration rather than as code
branches. Core-Geth forked from multi-geth in 2020 and kept that model, which is why adding a network here is a
chain configuration rather than a new client.

The Go module path is still `github.com/ethereum/go-ethereum`, and every internal import uses it. That is
deliberate downstream compatibility, not a leftover.

## The maintainers

ETC Labs Core formed in December 2018, many of its developers having previously been part of
[ETCDEV](https://web.archive.org/web/20190330063218/https://www.etcdevteam.com/), which supported the Classic
Geth client. As Classic Geth was retired the team supported multi-geth, and then Core-Geth from 2020, publishing
as [ETC Labs](https://web.archive.org/web/20200425081322/https://etclabs.org/) and
[ETC Core](https://web.archive.org/web/20200426174445/https://etccore.io/).

ETC Labs left the Ethereum Classic ecosystem in 2021. From January 2022 the work was funded by the
[ETC Cooperative](https://web.archive.org/web/20250205192722/https://etccooperative.org/), as it
[announced in December 2021](https://web.archive.org/web/20211222232926/https://etccooperative.org/posts/2021-12-22-coop-now-funding-core-geth).
The Cooperative entered maintenance mode at the end of 2024, as its
[2024 retrospective](https://etccooperative.org/etc-cooperative-retrospective-2024.pdf) states and its
[Q1 2025 report](https://web.archive.org/web/20250811083834/https://etccooperative.org/posts/2025-06-24-q1-report-en)
confirms: *"we are now in maintenance mode and spending has decreased significantly."* It stopped developing
Core-Geth in 2025 and wound down its team.

## The organizations

| Organization | What it did for Core-Geth | When |
| --- | --- | --- |
| ETCDEV | supported Classic Geth, the client whose developers went on to form ETC Labs Core | to 2018 |
| ETC Labs and ETC Core | developed multi-geth and then Core-Geth | 2018 to 2021 |
| ETC Cooperative | funded and maintained Core-Geth at `etclabscore/core-geth`, and maintained other Ethereum Classic services | January 2022 to 2024, then maintenance mode |
| [Ethereum Classic DAO](https://ethereumclassicdao.org) | a Wyoming DAO LLC launched in May 2025 that succeeds the ETC Cooperative | May 2025 onward |
| [`ethereumclassic`](https://github.com/ethereumclassic) GitHub organization | holds Ethereum Classic's public goods: the ECIPs, the website, the discovery lists and Core-Geth | 2024 onward |
| [White B0x](https://whiteb0x.com) | took up development, security work, disclosures and the modernization releases, for the Ethereum Classic DAO | February 2026 onward |

[The ETC Cooperative transition](../etc-cooperative-transition.md) lists which services are moving where as the
Cooperative winds down.

## This repository

Maintenance moved to [`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth), the Ethereum
Classic community repository, created on 21 December 2024 from the preceding repository at commit
[`7ef3ecd7a`](https://github.com/ethereumclassic/core-geth/commit/7ef3ecd7a716354589c2f27ff8f4b74d7a5edf2e)
(16 December 2024). That commit is preserved on the
[`archive-etclabscore-2024-12`](https://github.com/ethereumclassic/core-geth/tree/archive-etclabscore-2024-12)
branch and tagged
[`archive/etclabscore-2024-12`](https://github.com/ethereumclassic/core-geth/releases/tag/archive%2Fetclabscore-2024-12),
so the full history before the move is intact and checkable.

The Cooperative's own 2024 retrospective reproduces a conference slide by its Senior Editor, Donald McIntyre,
proposing "Move Core Geth to the Ethereum Classic community repository" (page 24 of the PDF). Keeping a
reference client in the project's own organization is ordinary practice, as with Bitcoin Core at
[`bitcoin/bitcoin`](https://github.com/bitcoin/bitcoin) and go-ethereum at
[`ethereum/go-ethereum`](https://github.com/ethereum/go-ethereum).

The same proposal is on the Cooperative's own published roadmap.
[Ethereum Classic Pathways](https://ethereumclassic.org/blog/2024-07-30-ethereum-classic-pathways-by-etc-cooperative-istora-and-donald-mcIntyre), published on
30 July 2024 and attributed to the ETC Cooperative, Istora and Donald McIntyre, lists: *"Move Core Geth to the
Ethereum Classic community GitHub repository: This would restore sovereignty of the community on the node client
software."*

**The community website's client entry has not followed.**
[ethereumclassic.org/development/clients](https://ethereumclassic.org/development/clients) still sends readers to the
previous project: its documentation site, its installation guide, its releases and its support channel. An
operator who follows the website lands there rather than here.
[#1678](https://github.com/ethereumclassic/ethereumclassic.github.io/pull/1678), opened on 20 March 2026 by a member of the
[Cooperative's board](https://etccooperative.org/people), would point that entry at this repository and remove
the previous organization's link from the site footer. It has not been merged. The same pull request carries a draft
[migration announcement](https://github.com/ethereumclassic/ethereumclassic.github.io/blob/f6b004e9d3c7325f8d4b7133e47722959be4bf3b/content/blog/2026-08-01-coregeth-repository-migration/index.md) for the website, written by its author on 2 August 2026 and readable on the
site's own [deploy preview](https://deploy-preview-1678--ethereumclassic.netlify.app/blog/2026-08-01-coregeth-repository-migration), which places the client's
move in the consolidation the Cooperative announced in July 2024 and counts the community Discord's
[November 2025 move](https://ethereumclassic.org/blog/2025-11-05-discord-migration/) as part of it. That
announcement has not been published either. Three days before it,
the website's footer link to this repository was reverted to the previous
organization's GitHub page in [#1676](https://github.com/ethereumclassic/ethereumclassic.github.io/pull/1676), a revert three reviewers approved, one of them recording that
"the etclabs repo is not maintained" and that "a longer term solution needs to be found". `v1.12.21` was
published from the previous repository the following day.

## The maintenance gap

After maintenance moved, the client went unfunded and unmaintained. The
[March 2026 security audit](../audits/2026-03-security-audit.md) counts 21 months without security maintenance,
from the `v1.12.20` release in June 2024 to March 2026. The previous repository received no commit between
January 2025 and March 2026, and security disclosures sent to it privately in 2025 went unanswered.

## 2026: the security work and the modernization

In February 2026 White B0x began modernizing the client here for the Ethereum Classic DAO: patching the
outstanding CVEs, moving to a supported Go toolchain, and rebuilding the release pipeline. In February and March
2026 it reported that security work privately to the ETC Cooperative, which owns the previous repository. A live
attack on Ethereum Classic bootnodes in March 2026 then confirmed the exposure.

The work reported privately was cut into the emergency `v1.12.21` and `v1.12.22` releases in the previous
repository, which gave `v1.12.x` operators immediate relief while `v1.13` was tested on Ethereum Classic mainnet
and Mordor from February to September 2026. The same fixes were submitted here as 27 individually scoped pull
requests, [#10](https://github.com/ethereumclassic/core-geth/pull/10) through
[#36](https://github.com/ethereumclassic/core-geth/pull/36), on 20 and 21 March 2026.

`v1.13.0` shipped on 14 September 2026, the first release published from this repository. What it fixes, and how
each claim about it checks out against the record, are in the
[release report](../release-reports/v1.13.0.md) and
[v1.13.0: the record behind the release](../release-reports/v1.13.0-record.md).

## Timeline

**This is where the code and its links have lived.** The disclosure record, which vulnerability was reported
when and where each fix landed, is the
[March 2026 security audit's own timeline](../audits/2026-03-security-audit.md#disclosure-timeline) and is not
repeated here.

| When | What |
| --- | --- |
| 10 June 2024 | [`v1.12.20`](https://github.com/etclabscore/core-geth/releases/tag/v1.12.20) is released from `etclabscore/core-geth`, the last release there before a 21-month gap |
| 30 July 2024 | The ETC Cooperative's published roadmap lists moving Core-Geth to the community GitHub repository |
| 16 December 2024 | [`7ef3ecd7a`](https://github.com/ethereumclassic/core-geth/commit/7ef3ecd7a716354589c2f27ff8f4b74d7a5edf2e), the last commit on the previous repository before this one is created |
| 21 December 2024 | `ethereumclassic/core-geth` is created from that commit |
| 23 January 2025 | The previous repository's last commit before 2026: a CI dependency bump, not a code change |
| 10 November 2025 | The website's footer GitHub link is changed from the previous organization to this repository, inside [#1648](https://github.com/ethereumclassic/ethereumclassic.github.io/pull/1648), a pull request about the community's Discord |
| 17 March 2026 | That footer change is reverted ([#1676](https://github.com/ethereumclassic/ethereumclassic.github.io/pull/1676)). Three reviewers approve it, one recording that "the etclabs repo is not maintained" and that "a longer term solution needs to be found" |
| 18 March 2026 | [`v1.12.21`](https://github.com/etclabscore/core-geth/releases/tag/v1.12.21) is published from the previous repository. Its release pull request records no review |
| 20 March 2026 | [#1678](https://github.com/ethereumclassic/ethereumclassic.github.io/pull/1678) proposes pointing the website's client entry at this repository and removing the previous organization's footer link. It has not been merged |
| 2 August 2026 | A draft [migration announcement](https://github.com/ethereumclassic/ethereumclassic.github.io/blob/f6b004e9d3c7325f8d4b7133e47722959be4bf3b/content/blog/2026-08-01-coregeth-repository-migration/index.md) for the website is added to that pull request. It has not been published |
| 20 to 21 March 2026 | The security work is submitted here as pull requests [#10](https://github.com/ethereumclassic/core-geth/pull/10) to [#36](https://github.com/ethereumclassic/core-geth/pull/36), one per vulnerability, with tests and linked advisories |
| 28 March 2026 | [`v1.12.22`](https://github.com/etclabscore/core-geth/releases/tag/v1.12.22) is published from the previous repository |
| 14 August 2026 | [`v1.12.23`](https://github.com/etclabscore/core-geth/releases/tag/v1.12.23) is published from the previous repository |
| 14 September 2026 | [`v1.13.0`](https://github.com/ethereumclassic/core-geth/releases/tag/v1.13.0) is released from this repository |

## What comes next

**`v1.13` is the last Core-Geth release line**, maintained through the transition. Core-Geth is scheduled to
sunset gradually in favor of two efforts. Ethereum Classic network extensions, overlays or plugins for Ethereum
clients are in development and none is offered here as an option yet. [Fukuii](https://fukuii.org), a client
native to Ethereum Classic, removes upstream third-party client dependencies from Ethereum Classic's core
software; moving to it starts with its first release, and the migration guide's
[Fukuii section](../tutorials/v1.13.0-migration.md#migrating-to-fukuii) says when. Fukuii is inspired by the work
of earlier Ethereum Classic development teams: ETCDEV's
[Classic Geth](https://github.com/ethereumproject/go-ethereum) and Orbita vision, and
[IOHK](https://iohk.io/)'s [Mantis](https://web.archive.org/web/20211026113958/https://mantisclient.io/).

Maintaining Core-Geth since it moved to the community repository has been unfunded public-goods work.
[Support this work](../index.md#support-this-work) says how to fund it.
