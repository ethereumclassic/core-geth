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
    - **Every node still on a `v1.12` release should move to the `v1.13` line.** The
      [migration guide](tutorials/v1.13.0-migration.md) has the commands and checksums for each platform, and the
      steps, including the required P2P node key rotation. A `v1.12` node that cannot move yet can
      [read the peer discovery lists and bootnodes the `v1.13` line ships](#if-a-v112-node-cannot-move-yet), a small
      step that closes none of the flaws the audits found in the `v1.12` line.
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
| Peer discovery DNS trees | The trees published from [`etclabscore/discv4-dns-lists`](https://github.com/etclabscore/discv4-dns-lists) | The trees under `ethereumclassic.net`, `ethclassic.net` and `ethereumclassic.network`, published from [`ethereumclassic/discv4-dns-lists`](https://github.com/ethereumclassic/discv4-dns-lists) | Published since September 10, 2026 and built into `v1.13.0`. The previous trees still resolve, and the board has communicated that they will stop being updated as the Cooperative dissolves |
| Bootnodes | The bootnodes built into `v1.12.x` | Three Classic and two Mordor bootnodes run by the Core-Geth maintainers, built into `v1.13.0` | Running |
| Public JSON-RPC endpoints | `etc.rivet.link`, which [ethereumclassic.org](https://ethereumclassic.org/knowledge/metamask/) describes as provided by Rivet under contract with the ETC Cooperative, and `rpc.mordor.etccooperative.org` for Mordor, on the Cooperative's own domain | The lists on ChainList today; `rpc.ethereumclassic.net` and `rpc-mordor.ethereumclassic.net` as they are stood up, from [`ethereumclassic/public-rpc`](https://github.com/ethereumclassic/public-rpc) and [`ethereumclassic/nodes`](https://github.com/ethereumclassic/nodes) | `etc.rivet.link` went offline in late August 2026. See [Public JSON-RPC](#public-json-rpc) |
| Security contact for stakeholders | The ETC Cooperative | <security@ethereumclassic.net> | Active |

[Bootnodes and peer discovery](guides/bootnodes-and-discovery.md) explains how the discovery trees are built from a
crawl of the network, and how to point a node at any other list.

### If a v1.12 node cannot move yet

Every node still on a `v1.12` release should move to the `v1.13` line. One that cannot move yet can take a small step:
read the peer discovery lists and bootnodes the `v1.13` line ships.
[`ethereumclassic/discv4-dns-lists`](https://github.com/ethereumclassic/discv4-dns-lists) records that the lists a
`v1.12` release reads by default stop being updated when the ETC Cooperative dissolves. The node stays exposed: the step
closes none of the flaws [the audits](releases.md#v1130) found in the `v1.12` line. It is for `v1.12.17` and later
releases: a node on an older one follows an obsolete chain and
[resyncs instead](tutorials/v1.13.0-migration.md#if-the-node-runs-a-release-before-v11217).

After a [clean stop](getting-started/run-classic-node.md#stop-it-safely), add both flags to the command or unit you
already start the node with:

```sh
geth --classic --discovery.dns "enrtree://APDLRZ2T7ERXPWXX4D5USB32NIFYHXMFVZQ3DZALK6JJJ5L4VSYIQ@all.classic.ethereumclassic.net,enrtree://APDLRZ2T7ERXPWXX4D5USB32NIFYHXMFVZQ3DZALK6JJJ5L4VSYIQ@all.classic.ethclassic.net,enrtree://APDLRZ2T7ERXPWXX4D5USB32NIFYHXMFVZQ3DZALK6JJJ5L4VSYIQ@all.classic.ethereumclassic.network" --bootnodes "enode://5c1f4f1c8036879c3186f0632ce0fdf3514e67679f0b5e1814335f5093a21c6c2fda20266ca2b54aac3305085fa46d9711493da1b564489ecec7b37d1aebb5f6@178.104.220.119:30303,enode://82100619a51fed73bce973a5f7ead3a9468535bdd4c87488037f99bf53b741e4701c5909b7c92897c863c23e494e27c83470a2fcbcc51f8b49655ed377dbd962@135.148.46.64:30303,enode://d77bc6696fc2a1419871c8819baef44089cd1847b1b895266382892fd04571cc01a97c45d4df6cbe4c37004195a3c578a76e706d12548916ca5c4341768b1972@62.72.47.101:30303"
```

On Mordor, Ethereum Classic's test network, use this line instead:

```sh
geth --mordor --discovery.dns "enrtree://APDLRZ2T7ERXPWXX4D5USB32NIFYHXMFVZQ3DZALK6JJJ5L4VSYIQ@all.mordor.ethereumclassic.net,enrtree://APDLRZ2T7ERXPWXX4D5USB32NIFYHXMFVZQ3DZALK6JJJ5L4VSYIQ@all.mordor.ethclassic.net,enrtree://APDLRZ2T7ERXPWXX4D5USB32NIFYHXMFVZQ3DZALK6JJJ5L4VSYIQ@all.mordor.ethereumclassic.network" --bootnodes "enode://d6f2e535a2d183f930bd780913d1d4169d3479f22c2941117d10184d65adbb614866db83049438c17657d1b24a53f3560280f563305426e93868a9a98a02c5fc@2.29.31.123:30303,enode://b0b8da5708a22a0bc92e993c6fb93801d1bc5956efaf3ef7cafcd3af2167f634a2147abb6913a719d61ed03e180aa2045ffea7bd2e367db589165e83058df487@40.160.140.180:30303"
```

`--discovery.dns` takes the place of the lists built into the release, and `--bootnodes` the place of its bootnodes,
with the ones the Core-Geth maintainers run. If the command already has either flag, replace that flag's values with
these. The node needs no working bootnode to use the lists: it dials the nodes in them directly. Keep the node's other
flags: a node on `--datadir` that runs this line as it stands starts an empty chain in the default directory. On
Windows, run `.\geth.exe` in place of `geth`. A container cannot take a new flag: stop it with `docker stop -t 300` and
replace it, on the same data, with a container from the same `v1.12` image that has both flags after the image name.
[Step 5 of the Docker guide](tutorials/migration/docker.md#5-start-and-verify) shows the commands; use the old image,
which `docker inspect <container> --format '{{.Config.Image}}'` prints, in place of `$IMAGE`. Under systemd, run
`sudo systemctl daemon-reload` after you edit the unit.

A node started from a config file that names discovery lists, as one `geth dumpconfig` wrote does, needs the file
changed as well: `--discovery.dns` sets only `EthDiscoveryURLs`, so the node would keep the file's `SnapDiscoveryURLs`,
the `v1.12` lists. Replace the values of the file's `EthDiscoveryURLs` and `SnapDiscoveryURLs` lines, under `[Eth]`,
with the three lists, and delete its `BootstrapNodes` and `BootstrapNodesV5` lines, under `[Node.P2P]`: `--bootnodes`
takes their place. Then start the node with `--bootnodes` and without `--discovery.dns`.
[Using the trees](https://github.com/ethereumclassic/discv4-dns-lists/blob/main/docs/using-the-trees.md#core-geth-v112-line)
shows the file.

When the node moves to the `v1.13` line, take both flags out, and delete the two list lines this step put in a config
file. That line reads these lists and its own bootnodes by default: `--bootnodes` would hold it to the bootnodes it
names, and `--discovery.dns` would leave it no lists for snap peers.
[Bootnodes and peer discovery](guides/bootnodes-and-discovery.md) covers both flags.

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
