---
description: "The short path for Ethereum Classic node operators, pools and exchanges: what to do now, what to read in what order, and how to verify all of it yourself."
---

# Node operators: start here

**If you run an Ethereum Classic node, a mining pool, an exchange or an RPC endpoint, this page is the shortest
path through everything else on this site.** It takes about a minute to read, and it links the rest in the order
that answers the questions operators actually arrive with.

## Do these four things

1. **Upgrade to [v1.13.0 or later](https://github.com/ethereumclassic/core-geth/releases/latest).** Every published
   `v1.12.20` to `v1.12.23` archive, the newest included, was built on an unsupported Go release and carries 55
   to 61 Go standard library advisories that `v1.13.0` does not. The releases before `v1.12.21` also carry all
   six client CVEs, two of which were exploited against Ethereum Classic bootnodes in March 2026. On
   17 September 2026, 71 percent of the Core-Geth nodes on the network were running a release older than
   `v1.12.23`, and 30 percent were running one older than `v1.12.21`. The
   [migration guide](tutorials/v1.13.0-migration.md) has a guide per platform. Your chain data carries over, and
   the upgrade costs about twenty minutes of downtime.

    **What the network was running on 17 September 2026**, from [etcnodes.org](https://etcnodes.org),
    524 Core-Geth nodes of 550:

    | Running | Nodes | Share | Carries |
    | --- | ---: | ---: | --- |
    | `v1.12.20` and older | 158 | 30.2% | all six client CVEs, unpatched. Two were exploited against ETC bootnodes in March 2026 |
    | `v1.12.21` | 50 | 9.5% | the ECIES crash and key oracle closed. Still missing two curve checks, the RLP work and the GraphQL limit. Go 1.21 |
    | `v1.12.22` | 165 | 31.5% | the rest of the CVE backports, but CVE-2026-26313 only partly mitigated. Introduces the `eth_syncing` regression. Go 1.21 |
    | `v1.12.23` | 138 | 26.3% | the delayed decoding series, and nothing else here is fixed: 55 to 61 Go advisories in the binary, no GraphQL depth limit, the `eth_syncing` regression untouched, and a response cap that can disconnect honest peers |
    | `v1.12.24` | 4 | 0.8% | a development build of the previous repository's `master`, not a release |
    | `v1.13.0` | 9 | 1.7% | **Recommended client:** six CVEs resolved and the GraphQL limit fixed, Go 1.26.8, zero Go advisories |

    Each row is measured in the [Go toolchain audit](audits/2026-09-go-toolchain.md) and the
    [March 2026 security audit](audits/2026-03-security-audit.md). The figures move:
    `https://api.etcnodes.org/peers` is the instrument.
2. **Rotate the P2P node key.** This one is required rather than precautionary: one of the fixed issues leaks
   bits of that key across repeated handshakes, and the fix cannot recall what already leaked.
   [How, and which peer lists to update](tutorials/v1.13.0-migration.md#rotate-the-p2p-node-key).
3. **Track releases at [`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth/releases).** A
   node tracking the previous repository will not see `v1.13.0`.
4. **Decide your MESS setting, and give every node you run the same one.** The bundled default follows
   ECIP-1110: MESS is inactive from the Spiral block, and `--mess` turns it on. MESS decides which of two
   competing chains a node prefers during a deep reorganization and never whether a block is valid, so in
   normal operation a node with it on and a node with it off follow the same chain. An exchange, a custodian
   or a pool that credits deposits is the one to think hardest about it: a deep reorganization is the attack
   aimed at them. A miner carries a cost either way, and that trade is worth reading before deciding rather
   than after.
   [Which setting fits which operator](operate/mess.md#which-setting-fits-which-operator) has it per node type.
   Whichever you choose, nodes in one fleet that disagree can prefer different chains, which is worse than
   either setting.

## Then read, in this order

**If you have ten minutes.** [Questions and answers](about/faq.md) is the short version of this whole site: which
repository, which version, why the key rotation, MESS, public RPC and how to verify a download. The
[v1.13.0 release report](release-reports/v1.13.0.md) is what changed, what to expect and the known issues.

**If you are deciding whether the upgrade is worth the downtime.** The
[Go toolchain audit](audits/2026-09-go-toolchain.md) lists the advisories present in the binaries you are running
right now, measured from the published archives rather than from any build configuration. The
[March 2026 security audit](audits/2026-03-security-audit.md) covers the six CVEs and the GraphQL denial of
service, per release, with the disclosure timeline. The
[August 2026 follow-up](audits/2026-08-security-followup.md) covers what `v1.12.23` fixed and what it left open.

**If you have heard conflicting claims about this release.**
[v1.13.0: the record behind the release](release-reports/v1.13.0-record.md) answers each published claim from
commit metadata, pull requests, release files and advisory databases, and concedes the ones that are accurate.
Do not take it on faith:
[check it yourself](release-reports/v1.13.0-record.md#do-not-trust-this-page-verify-it) with a prompt that reads
every document here, and the article written about the release, and reports where they agree and disagree.

**If you depend on other Ethereum Classic services.**
[The ETC Cooperative transition](etc-cooperative-transition.md) lists where the client, the peer discovery lists,
the bootnodes and the public JSON-RPC endpoints continue as the ETC Cooperative winds down, and which public
endpoint to point an application at today.
[Project history](about/project-history.md) traces who has maintained this client since 2018 and how it reached
the community organization.

**If you want the artifacts checked rather than described.**
[Release artifacts](audits/2026-09-release-pipeline.md) measures what the published archives actually contain,
and [dependency and toolchain modernization](audits/2026-08-dependency-modernization.md) records what moved
underneath the code between the December 2024 archive point and this release.

## Then set your node up properly

| You run | Start with | Then |
| --- | --- | --- |
| A mining pool | [Mining pool node](guides/mining-pool-node.md) | [Production operations](guides/production-operations.md) |
| An exchange or custody service | [Production operations](guides/production-operations.md) | [MESS confirmation calculator](guides/mess-calculator.md) |
| A public RPC endpoint | [Public RPC endpoint](guides/public-rpc-endpoint.md) | [Security and network exposure](operate/security.md) |
| An explorer or indexer | [Archive node](guides/archive-node.md) | [Sync modes and data retention](operate/sync-modes.md) |
| A node for yourself | [Running a node](getting-started/run-a-node.md) | [Maintenance, backup and upgrades](operate/maintenance.md) |

[Choose your role](guides/index.md) has the rest.

## Questions, and reporting something

Report a security issue privately, never in a public issue: through this repository's
[private advisories](https://github.com/ethereumclassic/core-geth/security/advisories), or by email to
<security@ethereumclassic.com>. Mining pools, exchanges and service providers should use that address as their
point of contact as the ETC Cooperative dissolves. A person answers it: one of the core developers who maintain
this client.

Anything else, including a defect in `v1.13.0`, belongs in
[the repository's issues](https://github.com/ethereumclassic/core-geth/issues), where it is public and checkable.

To talk to the people who maintain this client, the Ethereum Classic Core Developers run a
[Discord](https://ethereumclassic.com/discord). It is the place for questions about running a node, an upgrade, or anything on these pages.
**It is not the place to report a vulnerability**: a public channel discloses it to everyone at once, which is
what the private routes above exist to avoid.
