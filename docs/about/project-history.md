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

[Sorpaas](https://github.com/sorpaas) led multi-geth's development from its start in 2018 until 2020.

ETC Labs Core formed in December 2018, many of its developers having previously been part of
[ETCDEV](https://web.archive.org/web/20190330063218/https://www.etcdevteam.com/), which supported the Classic
Geth client. It led development through the fork and rebrand of multi-geth to Core-Geth in 2020, publishing
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
| [Sorpaas](https://github.com/sorpaas) | led multi-geth's development from its start | 2018 to 2020 |
| ETC Labs and ETC Core | led development through the fork and rebrand of multi-geth to Core-Geth | 2020 to 2021 |
| ETC Cooperative | funded and maintained Core-Geth at `etclabscore/core-geth`, and maintained other Ethereum Classic services | January 2022 to 2024, then maintenance mode |
| [Ethereum Classic DAO](https://ethereumclassicdao.org) | a Wyoming DAO LLC launched in May 2025; the 2026 modernization was carried out for it | 2026 |
| [`ethereumclassic`](https://github.com/ethereumclassic) GitHub organization | holds Ethereum Classic's public goods: the ECIPs, the website, the discovery lists and Core-Geth, maintained there by Ethereum Classic core developers and community contributors | 2024 onward |
| [White B0x](https://whiteb0x.com) | carried out the 2026 security work, disclosures and the `v1.13` modernization | 2026 |

[The ETC Cooperative transition](../etc-cooperative-transition.md) lists which services are moving where as the
Cooperative winds down.

## This repository

Maintenance moved to [`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth), the Ethereum
Classic community repository, created on December 21, 2024 from the preceding repository at commit
[`7ef3ecd7a`](https://github.com/ethereumclassic/core-geth/commit/7ef3ecd7a716354589c2f27ff8f4b74d7a5edf2e)
(December 16, 2024). That commit is preserved on the
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
July 30, 2024 and attributed to the ETC Cooperative, Istora and Donald McIntyre, lists: *"Move Core Geth to the
Ethereum Classic community GitHub repository: This would restore sovereignty of the community on the node client
software."*

## The maintenance gap

After maintenance moved, the client went unmaintained. The
[March 2026 security audit](../audits/2026-03-security-audit.md) counts 21 months without security maintenance,
from the `v1.12.20` release in June 2024 to March 2026. The previous repository received no commit between
January 2025 and March 2026, and security disclosures sent to it privately in 2025 went unanswered.

## 2026: the security work and the modernization

In February 2026 White B0x began modernizing the client here for the Ethereum Classic DAO: patching the
outstanding CVEs, moving to a supported Go toolchain, and rebuilding the release pipeline. In February and March
2026 it reported that security work privately to the `etclabscore/core-geth` maintainers. A live
attack on the Core-Geth `v1.12` bootnodes in March 2026 then confirmed the exposure.

The work reported privately was cut into the emergency `v1.12.21` and `v1.12.22` releases in the previous
repository, which gave `v1.12.x` operators immediate relief while `v1.13` was tested on Ethereum Classic mainnet
and Mordor from February to September 2026. The same fixes were submitted here as 27 individually scoped pull
requests, [#10](https://github.com/ethereumclassic/core-geth/pull/10) through
[#36](https://github.com/ethereumclassic/core-geth/pull/36), on March 20 and 21, 2026.

`v1.13.0` shipped on September 14, 2026, the first release published from this repository and the start of the canonical v1.13 release line, where the client is now developed and released. What it fixes, and how
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
| 2024-06-10 | [`v1.12.20`](https://github.com/etclabscore/core-geth/releases/tag/v1.12.20) is released from `etclabscore/core-geth`, the last release there before a 21-month gap |
| 2024-07-30 | The ETC Cooperative's published roadmap lists moving Core-Geth to the community GitHub repository |
| December 2024 | The ETC Cooperative announces maintenance mode for Core-Geth and sets out moving it, with the network's other public goods, to the `ethereumclassic` community organization |
| 2024-12-16 | [`7ef3ecd7a`](https://github.com/ethereumclassic/core-geth/commit/7ef3ecd7a716354589c2f27ff8f4b74d7a5edf2e), the last commit on the previous repository before this one is created |
| 2024-12-21 | `ethereumclassic/core-geth` is created from that commit |
| 2025-01-23 | The previous repository's last commit before 2026: a CI dependency bump, not a code change |
| 2026-03-18 | [`v1.12.21`](https://github.com/etclabscore/core-geth/releases/tag/v1.12.21) is published from the previous repository |
| 2026-03-20 to 2026-03-21 | The security work is submitted here as pull requests [#10](https://github.com/ethereumclassic/core-geth/pull/10) to [#36](https://github.com/ethereumclassic/core-geth/pull/36), one per vulnerability, with tests and linked advisories |
| 2026-03-28 | [`v1.12.22`](https://github.com/etclabscore/core-geth/releases/tag/v1.12.22) is published from the previous repository |
| 2026-08-14 | [`v1.12.23`](https://github.com/etclabscore/core-geth/releases/tag/v1.12.23) is published from the previous repository |
| 2026-09-14 | [`v1.13.0`](https://github.com/ethereumclassic/core-geth/releases/tag/v1.13.0) is released from this repository, beginning the v1.13 line and establishing the canonical release pipeline for Core-Geth in the `ethereumclassic` organization, where the client is now developed and released |

## Supporting the work

Maintaining Core-Geth since it moved to the community repository is public-goods work, sustained by
sponsorship and backed by the Ethereum Classic DAO. [Support this work](../index.md#support-this-work)
says how to fund it.
