---
hide:
  - navigation # Hide navigation
title: Core-Geth
description: "Core-Geth is an Ethereum Classic execution client: install and run an ETC or Mordor node, upgrade from v1.12.x, and read the security audits behind v1.13.0."
---

# Core-Geth: an Ethereum Classic execution client

> A [go-ethereum](https://github.com/ethereum/go-ethereum) downstream that keeps chain
> configuration data-driven rather than hard-coded, so one binary serves Ethereum
> Classic, Mordor and private chains.

**Upstream go-ethereum has removed support for Ethereum Classic.** ETC consensus rules
are maintained here rather than inherited, which is what this client is for.

## Start here

| If you are | Go to |
| --- | --- |
| **Running a node, a pool, an exchange or an RPC endpoint** | [Node operators: start here](operators.md): what to do now, and what to read in what order |
| **Pages for your role** | [Choose your role](guides/index.md): miners, pools, exchanges, RPC providers and more |
| **Upgrading from v1.12.x** | [Migration guide](tutorials/v1.13.0-migration.md): read this before you upgrade |
| Using a service the ETC Cooperative maintained, such as `etc.rivet.link` | [ETC Cooperative transition](etc-cooperative-transition.md): where each one continues |
| Asking a short question | [Questions and answers](about/faq.md): which repository, which version, node keys, MESS, RPC and verification |
| Wondering who maintains this client | [Project history](about/project-history.md): the lineage of the code and of the organizations behind it |
| Running a node for the first time | [Installation](getting-started/installation.md), then [Running a node](getting-started/run-a-node.md): a guide for Linux, macOS, Windows and Docker, and the flags for each kind of node |
| Trying the Mordor test network | [Choose your network](getting-started/run-a-node.md#choose-your-network) |
| Building it yourself | [Build from source](developers/build-from-source.md) |
| Reviewing what was audited | [The five audit reports](#what-was-audited) |

!!! danger "The v1.12.x line is insecure: upgrade to the latest v1.13 release"
    The v1.12.x line carries six CVEs and a GraphQL denial of service, stands on
    out-of-support Go toolchains, ships binaries with release-pipeline defects, and its
    emergency patches introduced regressions. One of the CVEs was exploited against the Core-Geth
    `v1.12` bootnodes in March 2026. No v1.12.x release, v1.12.23 included, closes all
    of it; the v1.13 line does.

    | Identifier | Severity | What it allows |
    | --- | --- | --- |
    | [CVE-2026-22862](https://github.com/ethereum/go-ethereum/commit/638741b08) | High | Any peer crashes the node during the RLPx handshake. **Exploited against the Core-Geth `v1.12` bootnodes, March 18, 2026.** |
    | [CVE-2026-26315](https://github.com/ethereum/go-ethereum/commit/46bee92f9) | Medium | Repeated handshakes leak bits of the node's own P2P key |
    | [CVE-2026-26314](https://github.com/ethereum/go-ethereum/commit/895a8597c) | High | secp256k1 coordinate field boundary bypass |
    | [CVE-2025-24883](https://github.com/ethereum/go-ethereum/commit/159fb1a1d) | High | Off-curve public key accepted by `UnmarshalPubkey` |
    | [CVE-2026-26313](https://github.com/ethereum/go-ethereum/pull/33835) | High | One crafted p2p message exhausts the node's memory |
    | [CVE-2026-22868](https://github.com/ethereum/go-ethereum/commit/fdfd1235a) | Medium | KZG proof verification denial of service |
    | [GraphQL query depth](https://github.com/graph-gophers/graphql-go/security/advisories/GHSA-mh3m-8c74-74xh) | Medium | Unbounded GraphQL query nesting; no CVE identifier assigned |

    Beyond the CVEs, the same line carries:

    - **Out-of-support Go toolchains.** Every archive was built on Go 1.21 or Go 1.22, whose
      support ended in August 2024 and February 2025, and carries 67 to 73 Go standard library
      advisories those toolchains will never receive fixes for.
      [Go toolchain audit](audits/2026-09-go-toolchain.md).
    - **Release-pipeline and binary defects.** glibc and macOS platform floors left older
      systems and many macOS users unable to run the published binaries, a macOS archive most
      of its downloaders cannot run, and release files that did not identify themselves or carry
      a build attestation. [Release artifacts audit](audits/2026-09-release-pipeline.md).
    - **Regressions from the emergency patches.** The emergency v1.12.21, v1.12.22 and v1.12.23
      releases backported some CVE fixes but left the KZG denial of service and the GraphQL depth
      limit open, introduced the `eth_syncing` regression (v1.12.22) and a snap storage-range cap
      that can disconnect honest peers (v1.12.23).
      [March 2026 audit](audits/2026-03-security-audit.md) · [August 2026 follow-up](audits/2026-08-security-followup.md).

    Upgrade to [the latest v1.13 release](https://github.com/ethereumclassic/core-geth/releases/latest); the
    [migration guide](tutorials/v1.13.0-migration.md) covers the steps, the required node-key rotation, and how
    to roll back.

## Supported networks

| Network | Chain ID | Consensus | Flag |
| --- | --- | --- | --- |
| Ethereum Classic | 61 | Proof of Work (Etchash) | `--classic`, `--mainnet` or no flag |
| Mordor testnet | 63 | Proof of Work (Etchash) | `--mordor` |
| MintMe.com Coin | 24734 | Proof of Work | `--mintme` |
| Private chains | configurable | PoW / PoA | genesis configuration |

Ethereum Classic runs every hard fork from Frontier through Spiral. The
[README's consensus table](https://github.com/ethereumclassic/core-geth#etc-consensus-history)
lists each upgrade with its activation block and the EIPs it included; `params/` in the
source is the authority, and a fork schedule recalled from memory is a guess.

**With no network flag, `geth` runs Ethereum Classic mainnet**, and `--mainnet` is the
same as `--classic`. A data directory initialized with a private network's genesis keeps
running that network.

### Networks this client does not support

**Ethereum mainnet, Sepolia and Holesky are not supported.** This client implements
Ethereum upgrades only through Cancun, so it cannot follow any of them. `--ethereum`,
`--sepolia` and `--holesky` are deprecated and refuse to start, rather than start a node
that would fall off its chain.

**MintMe is carried deliberately.** `--mintme` works, and this release includes the
MintMe hardfork enabling PUSH0 and MCOPY. It is scheduled for deprecation in a later
release; it is here so that community has a modernized client to build from rather than
a fork of an abandoned one.

## Where releases come from

Releases are published from
[`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth/releases).
Archives published under the previous `etclabscore` namespace are not built from this
source and do not carry the fixes released here.

## What was audited

Five reports, each measuring a different layer. Together they are the evidence for the
paragraph above. No single one of them answers "what differs".

| Report | What it measures |
| --- | --- |
| [March 2026 security audit](audits/2026-03-security-audit.md) | The six CVEs and the GraphQL denial of service, with the per-release breakdown |
| [August 2026 security follow-up](audits/2026-08-security-followup.md) | `v1.12.23` measured at the tag against the advisory records: what it fixed and what it left open |
| [Dependency and toolchain modernization](audits/2026-08-dependency-modernization.md) | What changed underneath the code between the December 2024 archive and this release: the Go toolchain, 83 modules, the linter |
| [Release artifacts](audits/2026-09-release-pipeline.md) | What the published archives actually contain: platform floors, architectures, provenance |
| [Go toolchain](audits/2026-09-go-toolchain.md) | Which Go toolchain built each published archive, and the Go standard library advisories each one carries |

The first three describe the source. **The last two describe the files you download**,
which can disagree with the source: a release is not what the build configuration says it builds.

## Security reporting

Report a security issue privately, never in a public issue: through this repository's
[private advisories](https://github.com/ethereumclassic/core-geth/security/advisories), or by
email to <security@ethereumclassic.net>. To follow security updates and releases, track the
[`ethereumclassic/core-geth` release line](https://github.com/ethereumclassic/core-geth/releases).
[`SECURITY.md`](https://github.com/ethereumclassic/core-geth/blob/main/SECURITY.md) has the
disclosure policy and the details of each channel.

With the ETC Cooperative's dissolution, Ethereum Classic stakeholders such as mining pools,
exchanges and service providers should use <security@ethereumclassic.net> as their point of
contact. A person answers it: one of the core developers who maintain this repository and have
been with the network since its inception.

## Support this work

Maintaining Core-Geth since it moved to the community repository is public-goods work, sustained by
sponsorship and backed by the Ethereum Classic DAO. Mining pools, centralized exchanges, issuers of
Ethereum Classic financial products, Etchash
mining hardware manufacturers and large holders all depend on an Ethereum Classic client that
follows the network correctly and receives security fixes promptly. If your operation relies on
Ethereum Classic, please help fund that work. The [support page](support.md) carries the routes:
sponsorship on GitHub, an invoiced maintenance agreement for organizations, or a direct transfer to
an address provided on request.
