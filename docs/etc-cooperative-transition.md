---
description: "Where the Ethereum Classic services the ETC Cooperative maintained continue as it winds down: the client, peer discovery, bootnodes and public JSON-RPC."
---

# The ETC Cooperative transition

**The ETC Cooperative is winding down, and the services it maintained are moving to the
[`ethereumclassic`](https://github.com/ethereumclassic) GitHub organization.** The Cooperative has been in
maintenance mode since the end of 2024. Its
[2024 retrospective](https://etccooperative.org/etc-cooperative-retrospective-2024.pdf) describes
"cooperative-maintained services, including RPC endpoints", and states that the Cooperative will stay in
maintenance mode until its funding runs out: *"At that time, it will be up to other stakeholders to take on any
required maintenance of the ETC client, unless a new plan materializes."* The ETC Cooperative board has
communicated that the Cooperative will dissolve.

Long-time Ethereum Classic core developers are standing up replacements in the `ethereumclassic` organization, so
that operators have somewhere to move before the Cooperative closes. One service has already stopped.

!!! warning "What to do now"
    - **Get Core-Geth releases from [`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth/releases).**
      The ETC Cooperative board has communicated its intent to archive `etclabscore/core-geth`.
    - **For a public endpoint today, use the lists on ChainList**, for
      [Ethereum Classic, chain 61](https://chainlist.org/chain/61) and
      [Mordor, chain 63](https://chainlist.org/chain/63). `https://etc.rivet.link` no longer resolves. This
      project's own endpoints are being stood up to replace it: see [Public JSON-RPC](#public-json-rpc).
    - **Use <security@ethereumclassic.net> as your security contact.** A person answers it: one of the core
      developers who maintain Core-Geth.
    - **Check the organization before you follow a notice.** Before you move a node, a pool or an application to a
      new source of software or data, confirm which organization publishes it. The services on this page continue
      in the `ethereumclassic` organization.

## Where each service continues

| Service | Before | Continues at | Status |
| --- | --- | --- | --- |
| Core-Geth client | [`etclabscore/core-geth`](https://github.com/etclabscore/core-geth), owned by the ETC Cooperative | [`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth) | `v1.13.0` released on September 14, 2026 |
| Peer discovery DNS trees | The trees published from [`etclabscore/discv4-dns-lists`](https://github.com/etclabscore/discv4-dns-lists) | The trees under `ethereumclassic.net`, `ethclassic.net` and `ethereumclassic.network`, published from [`ethereumclassic/discv4-dns-lists`](https://github.com/ethereumclassic/discv4-dns-lists) | Published since September 10, 2026 and built into `v1.13.0`. The previous trees still resolve, and the board has communicated that they are to be archived as the Cooperative dissolves |
| Bootnodes | The bootnodes built into `v1.12.x` | Three Classic and two Mordor bootnodes run by the Core-Geth maintainers, built into `v1.13.0` | Running |
| Public JSON-RPC endpoints | `etc.rivet.link`, which [ethereumclassic.org](https://ethereumclassic.org/knowledge/metamask/) describes as provided by Rivet under contract with the ETC Cooperative, and `rpc.mordor.etccooperative.org` for Mordor, on the Cooperative's own domain | The lists on ChainList today; `rpc.ethereumclassic.net` and `rpc-mordor.ethereumclassic.net` as they are stood up, from [`ethereumclassic/public-rpc`](https://github.com/ethereumclassic/public-rpc) and [`ethereumclassic/nodes`](https://github.com/ethereumclassic/nodes) | `etc.rivet.link` went offline in late August 2026. See [Public JSON-RPC](#public-json-rpc) |
| Security contact for stakeholders | The ETC Cooperative | <security@ethereumclassic.net> | Active |

[Bootnodes and peer discovery](guides/bootnodes-and-discovery.md) explains how the discovery trees are built from a
crawl of the network, and how to point a node at any other list.

### Public JSON-RPC

**Use ChainList today.** It lists the public endpoints Ethereum Classic operators publish, for
[Ethereum Classic, chain 61](https://chainlist.org/chain/61) and
[Mordor, chain 63](https://chainlist.org/chain/63), with the chain parameters each one needs. `etc.rivet.link`
went offline in late August 2026 and no longer resolves, so any configuration still naming it needs changing now.
`rpc.mordor.etccooperative.org` still answers, from the dissolving ETC Cooperative's own domain.

**This project is standing up replacements**, `rpc.ethereumclassic.net` for Ethereum Classic and
`rpc-mordor.ethereumclassic.net` for Mordor, built from
[`ethereumclassic/public-rpc`](https://github.com/ethereumclassic/public-rpc) and
[`ethereumclassic/nodes`](https://github.com/ethereumclassic/nodes). Until they answer, ChainList is the
recommendation. This page says so when that changes.

**They are run from the organization's repositories on purpose.** An endpoint a wallet, an explorer or a payment
processor depends on should not disappear because one company's contract ended, which is what happened to
`etc.rivet.link`. Keeping the configuration and the admission list in public repositories in the
[`ethereumclassic`](https://github.com/ethereumclassic) organization means more than one maintainer can rebuild
or move the service, and anyone can read how it is run. The `ethereumclassic.net` domain is designated for the
network's public goods as the ETC Cooperative dissolves, and `v1.13.0` already publishes its
[discovery trees](guides/bootnodes-and-discovery.md) under it.

**A public endpoint is someone else's node.** It sees every request you send it, and it can rate limit you,
fall behind or stop. Anything you depend on should run against a node you control.
[Public RPC endpoint](guides/public-rpc-endpoint.md) covers serving JSON-RPC from your own node.

## Timeline

| When | What |
| --- | --- |
| 2024-07-30 | The Cooperative's published roadmap lists moving Core-Geth to the community GitHub repository |
| End of 2024 | The ETC Cooperative enters maintenance mode and sets out moving its public goods, Core-Geth among them, to the `ethereumclassic` organization |
| 2024-12-21 | `ethereumclassic/core-geth` is created from `etclabscore/core-geth` at [`7ef3ecd7a`](https://github.com/ethereumclassic/core-geth/commit/7ef3ecd7a716354589c2f27ff8f4b74d7a5edf2e) |
| 2026-03-18 | `v1.12.21` is published from the previous repository |
| 2026-08-14 | [`v1.12.23`](https://github.com/etclabscore/core-geth/releases/tag/v1.12.23) is published from the previous repository, its last legacy release before v1.13.0 |
| 2026-08-28 | `ethereumclassic/discv4-dns-lists` is created |
| Late August 2026 | `etc.rivet.link` goes offline |
| 2026-08-31 | `ethereumclassic/public-rpc` and `ethereumclassic/nodes` are created |
| 2026-09-10 | The first crawl of the network is published from `ethereumclassic/discv4-dns-lists` |
| 2026-09-14 | [`v1.13.0`](https://github.com/ethereumclassic/core-geth/releases/tag/v1.13.0) is released from `ethereumclassic/core-geth` with the organization's bootnodes and discovery trees, establishing the canonical release pipeline for Core-Geth, where the client is now developed and released |

## Why the `ethereumclassic` organization

**The `ethereumclassic` organization is where Ethereum Classic keeps its community resources:** the
[ECIPs](https://github.com/ethereumclassic/ECIPs), the
[website](https://github.com/ethereumclassic/ethereumclassic.github.io) and, since December 2024, Core-Geth. The
Cooperative's 2024 retrospective reproduces a conference slide by its Senior Editor, Donald McIntyre, proposing
"Move Core Geth to the Ethereum Classic community repository" (page 24 of the PDF), and the Cooperative's
published roadmap says the same: [Ethereum Classic Pathways](https://ethereumclassic.org/blog/2024-07-30-ethereum-classic-pathways-by-etc-cooperative-istora-and-donald-mcIntyre),
July 30, 2024, lists *"Move Core Geth to the Ethereum Classic community GitHub repository: This would restore
sovereignty of the community on the node client software."* Moving the services there keeps their code,
configuration and history in public repositories under the community's name, with more than one maintainer, after
the Cooperative closes.

**Core-Geth's releases and documentation are published from this repository.** The
[releases page](https://github.com/ethereumclassic/core-geth/releases/latest) carries the current release, and
[Installation](getting-started/installation.md) has the verified download and build steps.

The [project history](https://github.com/ethereumclassic/core-geth#project-history) traces Core-Geth's code and its
maintainers, and [v1.13.0: the record behind the release](release-reports/v1.13.0-record.md) sets out the record
for this release.
