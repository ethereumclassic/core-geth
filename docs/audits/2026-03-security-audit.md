---
description: "Six CVEs and a GraphQL denial of service in Core-Geth v1.12.x, with the per-release breakdown and the March 2026 attack on the Core-Geth `v1.12` bootnodes."
---

# Core-Geth v1.12 line Security Audit: March 2026

- **Client:** Core-Geth (Ethereum Classic)
- **Audit Date:** March 2026
- **Repository audited:** [github.com/etclabscore/core-geth](https://github.com/etclabscore/core-geth)
- **Last substantive release there:** v1.12.20 (June 10, 2024)
- **Emergency patches there:** [v1.12.21](https://github.com/etclabscore/core-geth/releases/tag/v1.12.21) (March 18, 2026) · [v1.12.22](https://github.com/etclabscore/core-geth/releases/tag/v1.12.22) (March 28, 2026): partial CVE backports that left the KZG denial of service and the GraphQL depth limit open; still built on the EOL Go 1.21 and Go 1.22 toolchains
- **Patched repository:** [github.com/ethereumclassic/core-geth](https://github.com/ethereumclassic/core-geth)
- **Patched release:** v1.13.0
- **Patched release date:** September 14, 2026
- **Patched by:** The core-geth Authors and the Ethereum Classic Core Developers team, `@ethereumclassic/core-developers`
- **Agentic auditing methods and tooling:** the strongest available Claude models at each stage of the work, February through September 2026
- **Reported over 2025-2026 by:** the Ethereum Classic community and security researchers (see [Acknowledgments](#acknowledgments))
- **Remediation carried out by:** [White B0x](https://whiteb0x.com)
- **Work funded by:** [Ethereum Classic DAO](https://ethereumclassicdao.org)
- **As of:** the `v1.13.0` release, September 14, 2026

**Findings**

The client's own code, detailed in this audit:

- [CVE-2025-24883](#cve-2025-24883-off-curve-public-key-in-unmarshalpubkey): off-curve public key accepted by `UnmarshalPubkey`
- [CVE-2026-22862](#cve-2026-22862-ecies-elliptic-curve-integrated-encryption-scheme-decrypt-ciphertext-length-undercheck): ECIES decrypt ciphertext-length undercheck, exploited against the Core-Geth `v1.12` bootnodes
- [CVE-2026-26315](#cve-2026-26315-ecies-generateshared-accepts-unvalidated-public-key): ECIES `GenerateShared` accepts an unvalidated public key
- [CVE-2026-26314](#cve-2026-26314-secp256k1-isoncurve-field-boundary-bypass): secp256k1 `IsOnCurve` field-boundary bypass
- [CVE-2026-22868](#cve-2026-22868-kzg-kate-zaverucha-goldberg-blob-proof-verification-dos): KZG blob proof-verification denial of service
- [CVE-2026-26313](#cve-2026-26313-p2p-rlp-recursive-length-prefix-item-count-memory-exhaustion): P2P RLP item-count memory exhaustion
- [GraphQL query depth](#graphql-query-depth-dos): unbounded query nesting denial of service

The rest of the `v1.12.x` line's exposure, measured in the companion audits:

- Go standard library advisories carried by the end-of-life Go 1.21 and Go 1.22 toolchains, and the unmaintained `golang.org/x/crypto/openpgp` dependency: [Go toolchain](./2026-09-go-toolchain.md)
- Release artifacts: glibc and macOS platform floors that left older systems and many macOS users unable to run the published binaries, a macOS archive most of its downloaders cannot run, and release files that did not identify themselves or carry a build attestation: [release artifacts](./2026-09-release-pipeline.md)
- The `eth_syncing` regression and the snap storage-range response cap that can disconnect honest peers: [August follow-up](./2026-08-security-followup.md)
- The third-party module graph, unmoved since December 2024: [dependency and toolchain modernization](./2026-08-dependency-modernization.md)

---

## What operators need to do

**Track releases at [`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth)
and upgrade to [v1.13.0 or later](https://github.com/ethereumclassic/core-geth/releases/latest).** Every finding in this document, and every finding in the
[August 2026 follow-up](./2026-08-security-followup.md), is resolved there.

**Findings apply to every release in the v1.12.x line, including the most recent.**

Each entry gives what the release fixed, what it left open, and what it introduced, so no release
can be read as having closed everything.

- **v1.12.20** (June 10, 2024): the baseline. **Open:** all six client CVEs (CVE-2025-24883,
  CVE-2026-22862, CVE-2026-26315, CVE-2026-26314, CVE-2026-26313, CVE-2026-22868) and the GraphQL
  depth-limit denial of service, none patched; CVE-2026-22862 was later exploited against the Core-Geth
  `v1.12` bootnodes. **Toolchain:** Go 1.21 for the Linux and macOS archives and Go 1.22 for the
  Windows and Arm archives, which reached end of life two and eight months later; its binaries carry
  62 Go standard library advisories across the release's archives.
- **v1.12.21 "Aegis"** (March 18, 2026): cut during a live attack on the Core-Geth `v1.12` bootnodes,
  five hours after the issue was opened. **Fixed:** CVE-2026-22862, the exploited ECIES crash, and
  CVE-2026-26315, the ECIES key oracle. **Left open:** CVE-2025-24883, CVE-2026-26314, CVE-2026-26313
  and CVE-2026-22868, and the GraphQL depth-limit DoS. **Introduced:** nothing. **Toolchain:**
  unchanged, still end-of-life Go 1.21 and Go 1.22, carrying 62 Go standard library advisories across
  its archives.
- **v1.12.22 "Hermes"** (March 28, 2026): **Fixed:** CVE-2025-24883 and CVE-2026-26314, and a partial
  mitigation of CVE-2026-26313 in which the out-of-memory crash is prevented but a CPU-amplification
  denial of service remains. **Left open:** CVE-2026-22868 (declared not applicable to Ethereum
  Classic, though the `v1.12.x` line still ships the Ethereum mainnet, Sepolia and Holesky
  configurations and selects them with `--mainnet`, `--sepolia` and `--holesky`, where EIP-4844 and
  the blob path are active) and the GraphQL depth-limit DoS. **Introduced:** the `eth_syncing` regression that reports
  `highestBlock` incorrectly, reported as
  [#697](https://github.com/etclabscore/core-geth/issues/697). **Toolchain:** unchanged, still
  end-of-life Go 1.21 and Go 1.22, carrying 62 Go standard library advisories across its archives.
- **v1.12.23 "Argos"** (August 14, 2026): **Fixed:** the delayed-decoding p2p hardening series
  completes the CVE-2026-26313 remediation, bounding a response by its pending request. **Left open:**
  CVE-2026-22868, the GraphQL depth-limit DoS, and the `eth_syncing` regression carried from
  v1.12.22. **Introduced:** a snap storage-range response cap that can disconnect peers answering
  correctly, measured in the [follow-up](./2026-08-security-followup.md). **Toolchain:** unchanged,
  still Go 1.21 and, for the Arm archives, Go 1.22, now twenty-four and eighteen months past end of
  life and carrying 62 Go standard library advisories across its archives.
- **v1.13.0** (September 14, 2026, released from `ethereumclassic/core-geth`): **Resolution.** All six
  client CVEs and the GraphQL depth-limit DoS are patched with test coverage; the build moves to the
  supported Go 1.26 toolchain, at its newest release when `v1.13.0` was cut, which
  [Go toolchain](./2026-09-go-toolchain.md) measures at zero Go standard library advisories; the p2p
  hardening series is carried with its original authorship; and
  the v1.12.23 storage-range cap is gone. Nothing in this document is left open. **Operators on any
  `v1.12.x` release should move their tracking to `ethereumclassic/core-geth` and run the latest v1.13
  release.** v1.13.0 ships with MESS active (ECBP-1100, a network best practice that affects which
  chain a node prefers in a deep reorganization, never block validity). MESS is an operator flag
  whose default each release sets for itself, so keep it consistent across your own nodes.
  [MESS](../operate/mess.md) has the detail.
- **Advisory identifiers in the v1.12.22 release notes are not the ones the advisory records
  assign.** CVE-2026-26315 is attached to work that is CVE-2026-26314. An operator matching
  advisories against those notes will not find 26314 addressed anywhere and may conclude a fix
  is missing when it shipped. The v1.12.21 and v1.12.23 notes carry no advisory identifiers at
  all.

**Two properties of that release line are worth stating plainly, because they bear on how much
assurance a release carries rather than on any individual finding.** The emergency releases
were each authored, reviewed and merged by one person, in about seventy minutes and under two
minutes respectively; the pull requests record this. And the organization that owns the repository has
published that it does not expect to continue this work. Its
[2024 retrospective](https://etccooperative.org/etc-cooperative-retrospective-2024.pdf) states:

> "Without any way to obtain funding, the ETC Coop will be in maintenance mode, with spending
> minimized, until the funding runs out. At that time, it will be up to other stakeholders to
> take on any required maintenance of the ETC client, unless a new plan materializes."

**The `ethereumclassic/core-geth` v1.13 release line is that continuation.** Its first release,
v1.13.0, carries every fix in this document, the Go 1.26 toolchain, the p2p hardening series, and
ETC-specific network tooling, and it is developed in the open at `ethereumclassic/core-geth`. It
has about six months of testing on Ethereum Classic mainnet and Mordor behind it, with no critical
issue found.
Operators tracking the v1.12.x line should move their tracking there.

---

## Executive Summary

During cross-client testing, the Ethereum Classic Core Developers identified that `etclabscore/core-geth` (the primary Ethereum Classic execution client) had received no security maintenance since June 2024, a 21-month gap. Six CVEs and one GraphQL depth-limit DoS were found unpatched, spanning cryptographic key validation, P2P protocol memory exhaustion, and query processing. CVE-2026-22868 is not reached on Ethereum Classic in normal operation, but the `v1.12.x` line also ships and selects the Ethereum Foundation networks (`--mainnet` selected Ethereum mainnet there), where EIP-4844 and the KZG path are active, so it is a live issue in a supported configuration of the client and is patched in v1.13.0. The Go toolchains underpinning the client had also reached end-of-life (EOL), Go 1.21 in August 2024 and Go 1.22, which built the Arm archives, in February 2025, exposing all deployed nodes to unpatched Go standard library vulnerabilities for at least 19 months.

Disclosures to `etclabscore/core-geth` in early 2025 received no response. A public security disclosure by a Ledger security researcher in February 2026 likewise received no response until an active attack on the Core-Geth `v1.12` bootnodes in March 2026 forced emergency patches (v1.12.21 and v1.12.22 at `etclabscore/core-geth`). Those patches backported some of the CVE fixes but did not close all of them (the KZG denial of service and the GraphQL depth limit stayed open) and left the client on the EOL Go 1.21 and Go 1.22 toolchains with no ETC-specific modernization. The full remediation (Go 1.26 upgrade, ETC network tooling, and all CVE fixes) is released as v1.13.0 under the `ethereumclassic` GitHub organization.

---

## Background

Cross-client interoperability testing surfaced the Core-Geth maintenance gap in early 2026. The Ethereum Classic Core Developers' review of the go-ethereum security advisory database and Go vulnerability database against the v1.12.20 codebase found all six CVEs unpatched. The security remediation and Go toolchain upgrade are released as v1.13.0 together with Go 1.26 and ETC-specific network tooling.

**Where this work is done, and why.** The `etclabscore` organization is controlled by ETC
Cooperative, which entered maintenance mode at the end of 2024 and said so in its own published
reporting. Its
[2024 retrospective](https://etccooperative.org/etc-cooperative-retrospective-2024.pdf), covering
the year ended December 31, 2024, gives the cause and the position. Its recurring income had
ended: the arrangement under which Grayscale LLC contributed a third of the fees from its
Ethereum Classic Trust ran to the end of its two-year term, with **March 2022 the last month
the Cooperative earned fees from it**, and no replacement funding followed. Its statement of
the resulting position (that maintenance of the ETC client would fall to other stakeholders)
is quoted at the top of this document. Its
[Q1 2025 report](https://etccooperative.org/posts/2025-06-24-q1-report-en) confirms that
position was in effect: *"we are now in maintenance mode and spending has decreased
significantly."*

Moving the client to a community repository was also proposed publicly from inside the
organization before that happened. The same retrospective reproduces, on page 24 of the PDF (printed page 18), a slide from
a conference talk by the Cooperative's Senior Editor, headed "ETC Pathway Proposed by Donald
McIntyre", whose second item reads **"Move Core Geth to the Ethereum Classic community
repository"**. It is presented there as one person's proposal rather than as adopted policy,
but it is carried in the Cooperative's own annual report.

[`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth) was created on
**December 21, 2024** as a fork of `etclabscore/core-geth`, under the same GitHub organization
that hosts the ECIP repository and other ETC-native tooling. It is the repository this
remediation is developed in and the one v1.13.0 is released from, and it is where contributions
and release artifacts for this series belong. `etclabscore/core-geth` remains readable for
history and for the v1.12.x series.

This audit takes no position on the Cooperative's funding situation. The reporting is cited
because it is the organization's own statement of who was expected to maintain the client, which
is the question a reader deciding where to contribute needs answered.

### Disclosure Timeline

| Date | Event |
|------|-------|
| 2024-06-10 | `etclabscore/core-geth` v1.12.20 released: last substantive release from ETC Cooperative staff |
| August 2024 | Go 1.21 reaches end-of-life; the toolchain of Core-Geth's Linux x86_64 and macOS archives enters EOL status |
| December 2024 | The ETC Cooperative announces maintenance mode for Core-Geth. Its [2024 retrospective](https://etccooperative.org/etc-cooperative-retrospective-2024.pdf) and published roadmap ([Ethereum Classic Pathways](https://ethereumclassic.org/blog/2024-07-30-ethereum-classic-pathways-by-etc-cooperative-istora-and-donald-mcIntyre), July 30, 2024) set out moving Core-Geth and the network's other public goods to the Ethereum Classic community GitHub organization for the community to maintain |
| 2024-12-21 | [`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth) is created from `etclabscore/core-geth` at commit [`7ef3ecd7a`](https://github.com/ethereumclassic/core-geth/commit/7ef3ecd7a716354589c2f27ff8f4b74d7a5edf2e), the community repository for that handoff |
| 2025-01-23 | Last commit to `etclabscore/core-geth`: a GitHub Actions dependency bump (not a code change) |
| 2025-01-30 | [CVE-2025-24883](https://github.com/ethereum/go-ethereum/commit/159fb1a1d) published (go-ethereum GHSA-q26p-9cq4-7fc2) |
| February 2025 | Go 1.22 reaches end-of-life; the toolchain of Core-Geth's Windows and Arm archives enters EOL status |
| Throughout 2025 | Private security disclosures sent to `etclabscore/core-geth`; no response received |
| 2025-06-30 | Community member @tornadocontrib opens [PR #683](https://github.com/etclabscore/core-geth/pull/683) to `etclabscore/core-geth`: Go 1.24 upgrade and [CVE-2025-24883](https://github.com/ethereum/go-ethereum/commit/159fb1a1d) fix included; closed without merge when contributor deleted their fork in February 2026 |
| 2026-01-13 | [CVE-2026-22862](https://github.com/ethereum/go-ethereum/commit/638741b08) and [CVE-2026-22868](https://github.com/ethereum/go-ethereum/commit/fdfd1235a) published |
| 2026-02-04 | Ledger security researcher [@niooss-ledger](https://github.com/niooss-ledger) opens [issue #692](https://github.com/etclabscore/core-geth/issues/692) publicly documenting [CVE-2025-24883](https://github.com/ethereum/go-ethereum/commit/159fb1a1d), [CVE-2026-22862](https://github.com/ethereum/go-ethereum/commit/638741b08), [CVE-2026-22868](https://github.com/ethereum/go-ethereum/commit/fdfd1235a): no response |
| 2026-02-17 | go-ethereum v1.17.0 released, carrying the fixes for [CVE-2026-26313](https://github.com/ethereum/go-ethereum/pull/33835), [CVE-2026-26314](https://github.com/ethereum/go-ethereum/commit/895a8597c) and [CVE-2026-26315](https://github.com/ethereum/go-ethereum/commit/46bee92f9) |
| 2026-02-18 | GHSA-689v-6xwf-5jf3, GHSA-2gjw-fg97-vg3r and GHSA-m6j8-rg6r-7mv8 published; the CVE records follow on February 19 |
| 2026-02-26 | [CVE-2025-24883](https://github.com/ethereum/go-ethereum/commit/159fb1a1d) patch authored by [White B0x](https://whiteb0x.com) |
| Early March 2026 | Go toolchain upgrade and remaining CVE patches authored by White B0x |
| February-March 2026 | White B0x reports the security work privately to the `etclabscore/core-geth` maintainers |
| 2026-03-18 | Active attack on the Core-Geth `v1.12` bootnodes: ECIES handshake crash-loop (`crypto/ecies.symDecrypt` panic) exploited in production; [@diega](https://github.com/diega) opens [PR #694](https://github.com/etclabscore/core-geth/pull/694) and self-merges it about 70 minutes later with no peer review, releasing [v1.12.21](https://github.com/etclabscore/core-geth/releases/tag/v1.12.21) ("Aegis") approximately 5 hours after the issue was first reported. This was the first code response there in 21 months, forced by the live attack rather than prior disclosures. A cryptographic security patch authored, reviewed, and merged by one person, in a repository with no other active reviewers, is itself a supply-chain risk: the process has no second reviewer able to catch a defective or malicious change shipped under cover of an emergency. |
| 2026-03-19 | @niooss-ledger [documents remaining unpatched CVEs](https://github.com/etclabscore/core-geth/pull/694#issuecomment-4089185353) after v1.12.21: [CVE-2025-24883](https://github.com/ethereum/go-ethereum/commit/159fb1a1d), [CVE-2026-26313](https://github.com/ethereum/go-ethereum/pull/33835), and the secp256k1 coordinate check still unaddressed. The comment gives the coordinate check [CVE-2026-26315](https://github.com/ethereum/go-ethereum/commit/46bee92f9); the advisory records assign it [CVE-2026-26314](https://github.com/ethereum/go-ethereum/commit/895a8597c) |
| March 20-21, 2026 | White B0x submits the security work to `ethereumclassic/core-geth` as individually scoped pull requests ([#10](https://github.com/ethereumclassic/core-geth/pull/10)-[#36](https://github.com/ethereumclassic/core-geth/pull/36)), one per CVE with test coverage and linked CVE references, and cross-references the set in [issue #692](https://github.com/etclabscore/core-geth/issues/692) |
| 2026-03-28 | [v1.12.22 "Hermes"](https://github.com/etclabscore/core-geth/releases/tag/v1.12.22) released at `etclabscore/core-geth` ([PR #696](https://github.com/etclabscore/core-geth/pull/696)): backports CVE-2025-24883 and CVE-2026-26314 and partly mitigates CVE-2026-26313, but leaves the KZG denial of service (CVE-2026-22868) and the GraphQL depth limit open and introduces the `eth_syncing` regression; EOL Go 1.21 and Go 1.22 toolchains unchanged, no ETC-specific modernization |
| 2026-08-14 | [v1.12.23 "Argos"](https://github.com/etclabscore/core-geth/releases/tag/v1.12.23) released at `etclabscore/core-geth`: closes CVE-2026-26313 with the delayed-decoding p2p hardening series, but leaves the KZG denial of service (CVE-2026-22868), the GraphQL depth limit, the `eth_syncing` regression and the EOL Go 1.21 and Go 1.22 toolchains unaddressed, and adds a snap storage-range cap that can disconnect honest peers. The [August 2026 follow-up](./2026-08-security-followup.md#status-by-release) measures it |
| 2026-09-14 | [v1.13.0](https://github.com/ethereumclassic/core-geth/releases/tag/v1.13.0) is released from `ethereumclassic/core-geth`, beginning the v1.13 line: all six client CVEs and the GraphQL denial of service patched, on the Go 1.26 toolchain with ETC network tooling and DNS discovery. It completes the migration of Core-Geth to the Ethereum Classic community organization and establishes the canonical release pipeline there, where the client is now developed and released; the [ETC Cooperative transition](../etc-cooperative-transition.md) covers the remaining service handoffs |

**Note on v1.12.21 and v1.12.22:** Those emergency patches backported some of the CVE fixes and are a safer option than v1.12.20 for operators who have not yet migrated, but they did not close every finding: the KZG denial of service (CVE-2026-22868) and the GraphQL depth limit stayed open, the `eth_syncing` regression was introduced, and they remain on the EOL Go 1.21 and Go 1.22 toolchains without the ETC network tooling or DNS discovery updates in v1.13.0. Operators running v1.12.x should upgrade to v1.13.0, released from [ethereumclassic/core-geth](https://github.com/ethereumclassic/core-geth).

### Prior Maintainers

Core-Geth is a fork of [multi-geth](https://github.com/multi-geth/multi-geth), originally created and maintained by **Wei Tang** ([@sorpaas](https://github.com/sorpaas)) from 2018 to 2020. Multi-geth was the first multi-network go-ethereum fork with first-class ETC support, and its architecture is the direct ancestor of Core-Geth's chain configuration system.

The Core-Geth fork was then developed by ETC Labs until they left the ETC ecosystem in 2021. ETC Cooperative-paid staff maintained the client through the Spiral hard fork up until announcing maintenance mode for the client in December 2024:

- **Isaac Ardis** ([@meowsbits](https://github.com/meowsbits)): primary architect and long-term maintainer (2020 to 2024)
- **Chris Ziogas** ([@ziogaschr](https://github.com/ziogaschr)): contributor and maintainer (2022 to 2024)
- **Diego López León** ([@diega](https://github.com/diega)): release manager; cut the v1.12.20 release (2022 to 2024)

---

## Vulnerability Summary

| CVE | Severity | Component | etclabscore/core-geth | v1.13.0 (ethereumclassic) |
|-----|----------|-----------|------------------------|--------------------------|
| [CVE-2025-24883](https://github.com/ethereum/go-ethereum/commit/159fb1a1d) | High | crypto: secp256k1 key deserialization | Backported in v1.12.22 | Patched |
| [CVE-2026-22862](https://github.com/ethereum/go-ethereum/commit/638741b08) | High | crypto/ecies: ECIES decrypt length check | Backported in v1.12.21 | Patched |
| [CVE-2026-26315](https://github.com/ethereum/go-ethereum/commit/46bee92f9) | Medium | crypto/ecies + secp256k1: ECIES GenerateShared / IsOnCurve | Backported in v1.12.21 | Patched |
| [CVE-2026-26314](https://github.com/ethereum/go-ethereum/commit/895a8597c) | High | crypto/secp256k1: coordinate field boundary bypass | Backported in v1.12.22 | Patched |
| [CVE-2026-22868](https://github.com/ethereum/go-ethereum/commit/fdfd1235a) | Medium | txpool / P2P: KZG DoS (blob/KZG proof verification) | Declared N/A to ETC¹ | Patched (peer disconnect on invalid proof) |
| [CVE-2026-26313](https://github.com/ethereum/go-ethereum/pull/33835) | High | P2P: RLP item count memory exhaustion | Mitigated in v1.12.22²  | Patched |
| [GraphQL depth](https://github.com/graph-gophers/graphql-go/security/advisories/GHSA-mh3m-8c74-74xh) | Medium | RPC: unbounded query nesting DoS | Not addressed | Patched |

¹ [@diega stated](https://github.com/etclabscore/core-geth/issues/692) that CVE-2026-22868 is "not applicable to ETC" because ETC does not support EIP-4844 blob/KZG transactions. That holds for the Ethereum Classic network, but the `v1.12.x` line also ships the Ethereum mainnet, Sepolia and Holesky configurations and selects them with `--mainnet`, `--sepolia` and `--holesky`, where EIP-4844 and the KZG path are active, so the finding applies to a supported configuration of the client, not only in theory. The v1.13.0 patch implements the go-ethereum fix (peer disconnect on invalid KZG proof via the `ErrKZGVerificationError` sentinel) and refuses those three networks, with `--mainnet` selecting Ethereum Classic.

² Both mitigations prevent the out-of-memory crash by refusing to allocate per declared item. Neither avoids scanning the payload: `rlp.CountValues` walks the tag bytes of the whole list with no early exit, and it runs inside `RawList`'s own decode as well as in v1.12.22's pre-decode count. A peer can therefore still spend CPU proportional to message size on a message that will be rejected. The distinction between the two releases is where the bound comes from: a static item ceiling in v1.12.22, the pending request in v1.13.0, not whether the payload is read.

**Read the fourth column by release, not as a single verdict.** "Backported in v1.12.22" means the fix is present in v1.12.22 and later, not that v1.12.22 closed every finding above. By v1.12.22 four of the six CVEs were fixed and CVE-2026-26313 was only partly mitigated, while the KZG denial of service and the GraphQL depth limit were still open, the `eth_syncing` regression had been introduced, and the Go toolchains remained out of support. The [August 2026 follow-up](./2026-08-security-followup.md#status-by-release) gives the full per-release matrix across the client CVEs, the Go toolchain and the regressions. No v1.12.x release, v1.12.23 included, closes all of them; the v1.13 line does.

---

## Vulnerability Details

### CVE-2025-24883: Off-Curve Public Key in UnmarshalPubkey

- **Severity:** High
- **CVSS v3.1:** `AV:N/AC:H/PR:N/UI:N/S:U/C:H/I:H/A:H`: **7.4 High**
- **GHSA:** GHSA-q26p-9cq4-7fc2
- **Component:** `crypto/crypto.go`: `UnmarshalPubkey()`
- **Affected:** etclabscore/core-geth ≤ v1.12.21
- **Patched:** ethereumclassic/core-geth v1.13.0
- **Commit:** `681c915f0`
- **Upstream reference:** go-ethereum PR #31100 / commit `159fb1a1d`

**Description:**
`UnmarshalPubkey()` decoded a 65-byte uncompressed public key into `(x, y)` field elements but did not verify that the resulting point lies on the secp256k1 curve. A malicious peer could supply an off-curve point that passes unmarshaling without error, then produces invalid or undefined results in all subsequent ECDSA or ECDH operations that consume the deserialized key.

**Impact:**
Any code path that deserializes an untrusted public key (including P2P node identity validation and RLPx (Recursive Length Prefix transport) handshake processing) could be supplied an off-curve point. Downstream cryptographic operations on the invalid key can produce incorrect signatures, incorrect shared secrets, or panics depending on the caller.

**Fix:**
Added an `IsOnCurve(x, y)` check immediately after coordinate extraction in `UnmarshalPubkey()`. Off-curve points now return `errInvalidPubkey` before any downstream use.

```go
if !S256().IsOnCurve(x, y) {
    return nil, errInvalidPubkey
}
```

---

### CVE-2026-22862: ECIES (Elliptic Curve Integrated Encryption Scheme) Decrypt Ciphertext Length Undercheck

- **Severity:** High
- **CVSS v3.1:** `AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H`: **7.5 High** (confirmed live exploitation, March 18, 2026)
- **GHSA:** GHSA-mr7q-c9w9-wh4h
- **Component:** `crypto/ecies/ecies.go`: `Decrypt()`
- **Affected:** etclabscore/core-geth ≤ v1.12.20
- **Patched:** ethereumclassic/core-geth v1.13.0
- **Commit:** `c46834dd8`
- **Upstream reference:** go-ethereum commit `638741b08`

**Description:**
The ECIES `Decrypt()` function validated ciphertext minimum length using `rLen + hLen + 1`, where the `+ 1` accounts for only one byte beyond the point and HMAC fields. The correct minimum is `rLen + hLen + params.BlockSize` (AES block size = 16 bytes for the default ECIES parameters). The off-by-fifteen gap allows a ciphertext between 2 and 15 bytes shorter than a valid AES block to pass the length guard, after which array indexing proceeds into out-of-bounds memory.

**Impact:**
A malicious peer sending a crafted RLPx `auth` message with an undersized ECIES payload can trigger an out-of-bounds read during handshake processing, crashing the node (remote DoS). The RLPx handshake accepts unauthenticated ECIES ciphertexts from any connecting peer, so no prior authentication is required.

**Fix:**
Changed the minimum length check from `rLen + hLen + 1` to `rLen + hLen + params.BlockSize`.

```go
// Before:
if len(c) < (rLen + hLen + 1) {
// After:
if len(c) < (rLen + hLen + params.BlockSize) {
```

**Observed Attack (March 18, 2026):**
This vulnerability was actively exploited against the Core-Geth `v1.12` bootnodes (ams3, sfo3) on March 18, 2026. Malicious P2P traffic sent crafted `auth` messages with undersized ECIES payloads, triggering the out-of-bounds slice allocation and crashing each node on the next inbound handshake attempt. Because the crash occurred in `listenLoop`, the node process exited and restarted under the service manager, only to crash again on the next malicious connection, producing a crash-loop.

The following stack trace, reported by node operator [@shrikus](https://github.com/shrikus) in [issue #692](https://github.com/etclabscore/core-geth/issues/692), confirms the call path:

```
panic: runtime error: makeslice: len out of range

goroutine 42797 [running]:
github.com/ethereum/go-ethereum/crypto/ecies.symDecrypt(...)
        crypto/ecies/ecies.go:224
github.com/ethereum/go-ethereum/crypto/ecies.(*PrivateKey).Decrypt(...)
        crypto/ecies/ecies.go:322
github.com/ethereum/go-ethereum/p2p/rlpx.(*handshakeState).readMsg(...)
        p2p/rlpx/rlpx.go:612
github.com/ethereum/go-ethereum/p2p/rlpx.(*handshakeState).runRecipient(...)
        p2p/rlpx/rlpx.go:415
github.com/ethereum/go-ethereum/p2p/rlpx.(*Conn).Handshake(...)
        p2p/rlpx/rlpx.go:308
github.com/ethereum/go-ethereum/p2p.(*rlpxTransport).doEncHandshake(...)
        p2p/transport.go:132
github.com/ethereum/go-ethereum/p2p.(*Server).setupConn(...)
        p2p/server.go:987
github.com/ethereum/go-ethereum/p2p.(*Server).listenLoop.func2()
        p2p/server.go:921
```

Per the description of [PR #694](https://github.com/etclabscore/core-geth/pull/694), the v1.12.21 release pull request: bootnode **sfo3** had accumulated **805+ restart cycles** on v1.12.20 at the time of the patch, confirming the attack was ongoing and automated. The patched binary was deployed to **ams3** first and was running stable there while sfo3, still on v1.12.20, kept crash-looping.

---

### CVE-2026-26315: ECIES GenerateShared Accepts Unvalidated Public Key

- **Severity:** Medium
- **CVSS v3.1:** `AV:N/AC:H/PR:N/UI:N/S:U/C:H/I:N/A:N`: **5.9 Medium** (key-oracle attack requires repeated unauthenticated handshake attempts). The go-ethereum advisory [GHSA-m6j8-rg6r-7mv8](https://github.com/ethereum/go-ethereum/security/advisories/GHSA-m6j8-rg6r-7mv8) rates it Moderate.
- **GHSA:** GHSA-m6j8-rg6r-7mv8
- **Component:** `crypto/ecies/ecies.go`: `GenerateShared()`
- **Affected:** etclabscore/core-geth ≤ v1.12.20
- **Patched:** ethereumclassic/core-geth v1.13.0
- **Commit:** `c19892395`
- **Upstream reference:** go-ethereum commit `46bee92f9`

**Description:**
The RLPx handshake uses ECIES decryption on unauthenticated input from the network. `GenerateShared()` (called during ECDH (Elliptic Curve Diffie-Hellman) shared-secret derivation) accepted a `*PublicKey` without verifying it lies on the curve. An ephemeral public key with `X == nil`, `Y == nil`, or coordinates not satisfying the secp256k1 curve equation would proceed into ECDH multiplication and fail only at MAC verification. The MAC failure reveals to the attacker whether the faulty key survived ECDH, which can be used as an oracle to leak bits of the node's static P2P private key across multiple handshake attempts.

**Impact:**
A remote attacker making repeated unauthenticated RLPx handshake attempts with crafted ephemeral keys can potentially recover bits of the target node's P2P private key through timing or MAC-oracle analysis. No authentication or prior connection is needed.

**Fix:**
Added an explicit nil and `IsOnCurve` guard at the start of `GenerateShared()`:

```go
if pub.X == nil || pub.Y == nil || !pub.Curve.IsOnCurve(pub.X, pub.Y) {
    return nil, ErrInvalidPublicKey
}
```

---

### CVE-2026-26314: secp256k1 IsOnCurve Field Boundary Bypass

- **Severity:** High
- **CVSS v3.1:** `AV:N/AC:H/PR:N/UI:N/S:U/C:H/I:H/A:H`: **8.1 High**
- **GHSA:** GHSA-2gjw-fg97-vg3r
- **Component:** `crypto/secp256k1/curve.go`: `IsOnCurve()`; `crypto/secp256k1/ext.h`: `secp256k1_ext_scalar_mul()`; `crypto/signature_nocgo.go`: `btCurve.IsOnCurve()`
- **Affected:** etclabscore/core-geth ≤ v1.12.21
- **Patched:** ethereumclassic/core-geth v1.13.0
- **Commit:** `c19892395` (bundled with CVE-2026-26315)
- **Upstream reference:** go-ethereum commit `895a8597c`

**Description:**
`IsOnCurve()` verified the curve equation `y² ≡ x³ + b (mod P)` but did not first verify that the coordinates `x` and `y` are within the field, i.e., strictly less than the curve prime `P`. Due to modular arithmetic, coordinates equal to or greater than `P` may still satisfy the curve equation when reduced, but they represent invalid (non-canonical) points. Additionally, the C-level `secp256k1_ext_scalar_mul` function did not check the return value of `secp256k1_fe_set_b32`, which returns 0 when a coordinate is out-of-field. A crafted out-of-field coordinate could therefore bypass `IsOnCurve` and proceed into scalar multiplication, producing undefined or attacker-influenced results.

**Impact:**
An attacker supplying points with out-of-field coordinates that still satisfy the naive curve check can pass a gate that is supposed to reject invalid public keys, potentially leading to node crash or, in a consensus context, divergent state computation.

**Identifier mismatch.** `etclabscore/core-geth` labels this issue CVE-2026-26315. The
advisory records assign that identifier to the ECIES handshake issue documented in the previous
section, and assign **CVE-2026-26314** to this one: OSV lists `895a8597c` ("crypto/secp256k1: fix
coordinate check") as the fix for [GHSA-2gjw-fg97-vg3r](https://osv.dev/vulnerability/GHSA-2gjw-fg97-vg3r),
aliased to CVE-2026-26314 and GO-2026-4507, first released in go-ethereum v1.16.9.
CVE-2026-26313 and CVE-2026-26314 share the title "Go Ethereum affected by DoS via malicious p2p
message", so this issue is distinguished by its CWE class and fix commit rather than by its title,
which is what makes its identifier straightforward to transpose.

The label appears in commit `46bba8dfc`, a doc comment and the Go test function name
`TestIsOnCurveRejectsCoordinatesAboveP_CVE_2026_26315` in `crypto/secp256k1/curve_cve_test.go`,
the description of [PR #696](https://github.com/etclabscore/core-geth/pull/696), the published
release notes for v1.12.22, and a
[reply on issue #692](https://github.com/etclabscore/core-geth/issues/692#issuecomment-4146498131).
No commit, source file, pull request description or release note there uses CVE-2026-26314; that
reply does, giving it to the ECIES fix in v1.12.21. Operators matching advisories against release
notes should treat the two identifiers as covering both issues.

**Fix:**
Added field-boundary checks before the curve equation test in both the Go and C implementations:

```go
// Go
if x.Cmp(BitCurve.P) >= 0 || y.Cmp(BitCurve.P) >= 0 {
    return false
}
```

```c
// C
if (!secp256k1_fe_set_b32(&feX, point) ||
    !secp256k1_fe_set_b32(&feY, point+32)) {
    return 0;
}
```

---

### CVE-2026-22868: KZG (Kate-Zaverucha-Goldberg) Blob Proof Verification DoS

- **Severity:** Medium
- **CVSS v3.1:** `AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:L`: **5.3 Medium** (not reached on Ethereum Classic in normal operation; the score reflects the EIP-4844 networks the client also ships and selects)
- **Component:** `core/txpool/validation.go`: `validateBlobSidecar()`; `eth/fetcher/tx_fetcher.go`: `Enqueue()`
- **Affected:** etclabscore/core-geth ≤ v1.12.23 (code present; inactive on Ethereum Classic, active on the Ethereum Foundation configurations the client also ships)
- **Patched:** ethereumclassic/core-geth v1.13.0
- **Commit:** `9985c33fb`
- **Upstream reference:** go-ethereum commit `fdfd1235a` (v1.16.8)

**ETC Applicability:** [@diega stated that this is not applicable to ETC](https://github.com/etclabscore/core-geth/issues/692). ETC does not support EIP-4844 blob transactions, so the KZG code path is not reached in normal operation on the ETC network. The v1.12.22 release at `etclabscore/core-geth` does not address it. The v1.13.0 patch follows the go-ethereum approach: introduces `ErrKZGVerificationError` as a sentinel error and disconnects any peer that delivers a transaction with an invalid KZG proof, preventing repeated DoS attempts from the same peer.

**Reach across the client's networks.** The finding is not applicable only in the narrow sense that the KZG path is not reached on Ethereum Classic, Mordor or MintMe: none of those configurations activates EIP-4844, and the transaction pool rejects a blob transaction before that upgrade. It is reached on any configuration that activates EIP-4844. `etclabscore/core-geth` v1.12.23 still ships the Ethereum mainnet, Sepolia and Holesky configurations, selects them with `--mainnet`, `--sepolia` and `--holesky`, and does not carry the fix. v1.13.0 carries the fix and refuses those three networks, where `--mainnet` selects Ethereum Classic.

**Description:**
KZG blob proof verification is computationally expensive. A malicious peer could repeatedly broadcast blob transactions with invalid KZG proofs, causing the node to perform the full expensive cryptographic verification on each delivery attempt before rejecting the transaction. Because the error was not distinguished from other validation failures, the peer was not disconnected and could continue flooding the node indefinitely.

**Impact:**
On ETC: no active exposure in normal operation (blob transactions are rejected at an earlier validation stage). The code path is nonetheless present in the inherited go-ethereum codebase and could be reached by a crafted peer interaction.

**Remediation:**
Ported from go-ethereum commit `fdfd1235a` (v1.16.8). Adds `ErrKZGVerificationError` as a sentinel error in `core/txpool/errors.go` and wraps KZG validation failures in `validateBlobSidecar()` with it. In `eth/fetcher/tx_fetcher.go`, any delivery flagged with this violation immediately triggers `f.dropPeer(delivery.origin)`, preventing the offending peer from repeatedly triggering expensive cryptographic verification.

```go
case errors.Is(err, txpool.ErrKZGVerificationError):
    violation = err
// ...
if delivery.violation != nil {
    f.dropPeer(delivery.origin)
}
```

---

### CVE-2026-26313: P2P RLP (Recursive Length Prefix) Item Count Memory Exhaustion

- **Severity:** High
- **CVSS v3.1:** `AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H`: **7.5 High**
- **GHSA:** GHSA-689v-6xwf-5jf3
- **Component:** `eth/protocols/eth/`, `eth/protocols/snap/`, `p2p/tracker/`
- **Affected:** etclabscore/core-geth ≤ v1.12.22 (partial mitigation in v1.12.22; fully fixed in v1.12.23)
- **Patched:** ethereumclassic/core-geth v1.13.0
- **Commits:** the delayed-decoding series carried from v1.12.23: `68469cb19`, `863a59733`, `b77579985`,
  `92e92dc72`, `7a4988919`, `14b116bc1`, `665b53252`, `105dc459d` and `62489811f`. `7a4988919` is the one that
  bounds a response by its pending request.
- **Upstream reference:** go-ethereum PR #33835

**Description:**
The P2P message handler validated message size against a 10 MiB cap (`maxMessageSize`) before decoding, but did not validate the number of items declared in the RLP list header. A malicious peer could craft a valid RLP list header claiming millions of tiny items within the 10 MiB budget. When `msg.Decode` ran, it would allocate a pointer or struct object for each declared item before any further validation, causing out-of-memory crashes proportional to the declared item count rather than the actual payload size.

The attack requires only a valid peer connection at the devp2p (Ethereum device peer-to-peer wire protocol) handshake level. No authenticated session is needed to send crafted message payloads.

**Impact:**
Remote out-of-memory (OOM) crash of any reachable Core-Geth node. An attacker with a single peer connection could crash the node by sending one crafted message on the eth or snap protocol sub-channel.

**Fix:**
Response contents are decoded lazily. `rlp.RawList` holds a list without decoding its elements, so a message is validated before any per-item object is allocated, and a claimed item count no longer drives allocation. On the snap protocol a response is additionally matched against its pending request before its contents are read. This is go-ethereum's own remediation, PR #33835.

An earlier form of this patch counted items before decoding instead, because the upstream refactor could not be applied to this tree at the time. Both are no longer needed together, and the counting layer was removed once the refactor landed; go-ethereum carries no equivalent.

Response messages are bounded by the request they answer rather than by a fixed ceiling, which is both tighter than a static limit and incapable of refusing an honest peer. Transaction broadcasts keep an explicit limit of 5,000, and block announcements, which are unsolicited and not lazily decoded, keep one of 2,048.

- **Status in v1.12.22:** Partially mitigated: OOM crash prevented, CPU amplification DoS remains open.
- **Status in v1.13.0:** Patched. Response messages are bounded by their pending request, transaction broadcasts and block announcements by explicit ceilings.

---

### GraphQL Query Depth DoS

- **Severity:** Medium
- **CVSS v3.1:** `AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H`: **7.5 High** (assumes GraphQL endpoint exposed; off by default)
- **Related advisory:** GHSA-mh3m-8c74-74xh / CVE-2022-21708: a stack-overflow DoS in `graphql-go` reachable only when `MaxDepth` is enabled. Affected `<v1.3.0`; patched in v1.3.0
- **Component:** `graphql/service.go`; `go.mod`: `graphql-go v1.3.0 → v1.9.0`
- **Affected:** etclabscore/core-geth ≤ v1.12.23
- **Patched:** ethereumclassic/core-geth v1.13.0
- **Commit:** `6a046ee910`

**Disclosure:** Identified by White B0x during cross-client security review, combining review of the published `graphql-go` advisory with direct inspection of `graphql/service.go`, which confirmed that no depth limit was configured at the application layer. No separate responsible disclosure to `etclabscore` was required as the dependency advisory was already public; the finding was included in the patch set submitted to `ethereumclassic/core-geth`.

**Description:**
The GraphQL endpoint (`--graphql` flag) had no query complexity or depth limit. Deeply nested queries (for example, a query recursively nesting block references) would cause the server to perform unbounded recursive schema traversal, exhausting CPU and memory on the serving node.

**Impact:**
Any peer or client with access to the GraphQL endpoint can crash or heavily degrade a running node with a single deeply nested query. The GraphQL endpoint is off by default but is commonly enabled on infrastructure nodes.

**Fix:**
Added a `MaxDepth(20)` limit to the GraphQL schema parser. `graphql-go` is raised v1.3.0 → v1.9.0 for currency: GHSA-mh3m-8c74-74xh describes a stack-overflow DoS reachable only once `MaxDepth` is enabled, which is what this change does, and its fix is already present in v1.3.0, so the upgrade adds margin rather than closing an exposure:

```go
const maxQueryDepth = 20

s, err := graphql.ParseSchema(schema, &q, graphql.MaxDepth(maxQueryDepth))
```

---

## Go Toolchain End-of-Life

**Issue:** Core-Geth v1.12.20 was built and shipped on two Go toolchains: Go 1.21 for its Linux x86_64 and macOS archives, which reached end-of-life in August 2024, and Go 1.22 for its Windows and Arm archives, which reached end-of-life in February 2025. From those dates onward, vulnerabilities in the Go standard library (including `net/http`, `crypto/tls`, and `crypto/x509`) received no patches from the Go team for those toolchain versions, and all binaries compiled against them remained exposed. The Go 1.22 archives were already three security releases behind when v1.12.20 shipped.

**Impact:** Node operators running those binaries were exposed to Go standard library security issues for a minimum of 19 months by the time of this audit (August 2024 through March 2026), and v1.12.21 and v1.12.22 kept both toolchains. The Go vulnerability database lists multiple advisories published in this period whose affected versions include the shipped toolchains, three of them in the TLS stack. [Go toolchain](./2026-09-go-toolchain.md) measures which toolchain built each archive and lists the advisories each one carries.

**Fix:** Core-Geth v1.13.0 builds on Go 1.26 (current stable as of March 2026). The upgrade proceeded in two steps: 1.21 → 1.24 (removing the incompatible `fjl/memsize` dependency and fixing `go vet` format string errors introduced by 1.24's stricter checks), then 1.24 → 1.26 (updating all `golang.org/x/` dependencies for compatibility). The `blst` cryptography dependency was simultaneously upgraded from v0.3.11 to fix a C23 `typedef bool` incompatibility with the Alpine-based Docker build environment. That upgrade landed at v0.3.16 and a later dependency pass carried it to **v0.3.17**, which is what v1.13.0 ships; the header fix is present in both.

**Commit:** `b7d164018`

---

## Release Timeline

*For the full disclosure and maintenance history, see the [Disclosure Timeline](#disclosure-timeline) above.*

| Date | Release | Notes |
|------|---------|-------|
| 2024-06-10 | [v1.12.20](https://github.com/etclabscore/core-geth/releases/tag/v1.12.20) | Final substantive release at `etclabscore/core-geth` |
| August 2024 | - | Go 1.21 reaches end-of-life; no further toolchain security patches from the Go team |
| February 2025 | - | Go 1.22, which built the Windows and Arm archives of v1.12.20, reaches end-of-life |
| 2026-03-18 | [v1.12.21 "Aegis"](https://github.com/etclabscore/core-geth/releases/tag/v1.12.21) | Emergency patch for active ECIES crash-loop attack ([PR #694](https://github.com/etclabscore/core-geth/pull/694)); still built on EOL Go 1.21 and Go 1.22 |
| 2026-03-28 | [v1.12.22 "Hermes"](https://github.com/etclabscore/core-geth/releases/tag/v1.12.22) | Backports CVE-2025-24883 and CVE-2026-26314 and partly mitigates CVE-2026-26313; the KZG denial of service (CVE-2026-22868) and the GraphQL depth limit left open ([PR #696](https://github.com/etclabscore/core-geth/pull/696)); EOL Go 1.21 and Go 1.22 toolchains unchanged; `eth_syncing` regression introduced (issue [#697](https://github.com/etclabscore/core-geth/issues/697), open) |
| 2026-08-14 | [v1.12.23 "Argos"](https://github.com/etclabscore/core-geth/releases/tag/v1.12.23) | 32 commits: the delayed-decoding p2p hardening series (`rlp.RawList`, deferred block/body/receipt/transaction/snap decoding, request-response validation in `p2p/tracker`) plus seven upstream go-ethereum backports. EOL Go 1.21 and Go 1.22 toolchains unchanged. No CVE identifiers in the release notes |
| 2026-09-14 | v1.13.0 | Released from `ethereumclassic/core-geth`, the canonical release pipeline for Core-Geth: all six CVEs and the GraphQL DoS patched, Go 1.26 toolchain |

---

## Risk Assessment

| Risk Area | Severity | Description | Mitigation |
|-----------|----------|-------------|------------|
| Six CVEs + GraphQL DoS | Critical | Six CVEs and one GraphQL depth-limit DoS unaddressed for 21 months in the primary ETC client ([CVE-2026-22868](https://github.com/ethereum/go-ethereum/commit/fdfd1235a) is not reached on Ethereum Classic but is active on the Ethereum Foundation networks the client also ships and selects) | All patched in v1.13.0 |
| Cryptographic oracle ([CVE-2026-26315](https://github.com/ethereum/go-ethereum/commit/46bee92f9)) | Medium | Repeated handshake attempts could leak P2P node key bits | Patched; key rotation required when upgrading from any v1.12.x release |
| Remote crash via ECIES ([CVE-2026-22862](https://github.com/ethereum/go-ethereum/commit/638741b08)) | High | Any peer can crash a node during RLPx handshake with a malformed ECIES payload | Patched in v1.13.0 |
| Remote OOM via RLP ([CVE-2026-26313](https://github.com/ethereum/go-ethereum/pull/33835)) | High | Any peer can OOM-crash a node with a single crafted P2P message | Patched in v1.13.0 |
| CPU amplification DoS: [CVE-2026-26313](https://github.com/ethereum/go-ethereum/pull/33835) residual (v1.12.22 only) | Medium | v1.12.22 mitigation scans full RLP payload before rejecting oversized messages; ~2,500× more work per attack message than v1.13.0; malicious peers can exhaust CPU without causing OOM | Unmitigated in v1.12.22; patched in v1.13.0 |
| Go Runtime EOL | High | By March 2026, 19 months on unsupported Go 1.21, and Go 1.22 unsupported since February 2025; standard library CVEs accumulated unpatched | Upgraded to Go 1.26 in v1.13.0 |
| Unreviewed emergency patches (supply chain risk) | High | v1.12.21 and v1.12.22 were each authored, reviewed, and merged by the same individual with no independent peer review: v1.12.21 in about 70 minutes, v1.12.22 in under 2 minutes. With no second reviewer, a defective or malicious change shipped under cover of an emergency has nothing to catch it. The risk is structural and applies to any future emergency patch cut this way. | v1.13.0 is developed in the open under the `ethereumclassic` org, and its changes are measured in these audits |

---

## Scope

- **Audit target:** `etclabscore/core-geth` at tag `v1.12.20` (commit `c2fb44129`)
- **Audit date range:** February - March 2026
**In scope:**
- All Go source packages inherited from go-ethereum with known CVE exposure
- Go toolchain version and dependency security posture
- P2P protocol handler input validation (devp2p, RLPx, eth, snap sub-protocols)
- RPC endpoint security (GraphQL, JSON-RPC)

**Out of scope:**
- Consensus layer correctness and ETC protocol compliance
- Other client codebases
- EVM execution correctness
- Dependencies not listed in the go-ethereum security advisory database
- Infrastructure (bootnode operators, DNS, CDN)

---

## Methodology

The audit was initiated during cross-client interoperability testing. The `etclabscore/core-geth` codebase at tag `v1.12.20` was cross-referenced against the go-ethereum security advisory database (GitHub Security Advisories for `github.com/ethereum/go-ethereum`) and the Go vulnerability database (`vuln.go.dev`). Each published advisory was assessed for applicability by tracing shared code ancestry between Core-Geth and go-ethereum at the affected package level.

**Exploitability assessment:** Where a CVE had confirmed exploitation records (CVE-2026-22862: active ECIES crash-loop on the Core-Geth `v1.12` bootnodes, March 18, 2026), exploitability was treated as confirmed. For the remaining CVEs, static code analysis of the affected call paths was used to establish reachability from unauthenticated network input and assess preconditions.

**Patch development:** Applicable fixes were cherry-picked from upstream go-ethereum commits where possible. In two cases (CVE-2026-26313, CVE-2026-26314), structural divergence from go-ethereum produced merge conflicts incompatible with clean cherry-pick; fixes were manually ported by adapting the upstream approach to Core-Geth's code structure.

**Tooling:** `govulncheck` (Go vulnerability scanner), manual code review, go-ethereum advisory DB, `vuln.go.dev`.

**Testnet validation:** All patches were validated against the Mordor testnet. Validation covered: successful node sync resumption after patch application, P2P handshake stability under normal peer traffic, no regression in JSON-RPC method responses (`eth_syncing`, `eth_blockNumber`, `net_peerCount`), and block processing correctness against known Mordor block hashes.

---

## Upgrade Path

**If you are running any v1.12.x release, upgrade to v1.13.0 immediately.** v1.13.0 (the full Go 1.26 upgrade, ETC-specific modernization and every fix in this audit) is released from [`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth).

---

## Recommendations

- **Node operators (any v1.12.x release):** Upgrade to v1.13.0 from [github.com/ethereumclassic/core-geth](https://github.com/ethereumclassic/core-geth) immediately; it patches every CVE in this audit. Nodes on v1.12.20 or earlier are exposed to remote crash and potential key-oracle attacks from any peer.
- **Infrastructure providers and exchanges:** Treat the upgrade to v1.13.0 as a security-critical update, not a routine version bump.
- **Every node upgrading from v1.12.x:** Rotate the P2P node key (`--nodekey`). This is required, not precautionary: CVE-2026-26315 is an oracle against the node key, so a key used by an unpatched node should be treated as exposed. Any node reachable from the public internet over the 21-month gap was potentially targeted. The [migration guide](../tutorials/v1.13.0-migration.md#rotate-the-p2p-node-key) gives the procedure. None of the `etclabscore/core-geth` emergency releases note this: `v1.12.21` patches the oracle without recording that the key the node used under `v1.12.20` was exposed and should be replaced, and `v1.12.22` and `v1.12.23` are silent on it, so an operator who upgraded within the `v1.12.x` line and did not rotate is still running the exposed key. The go-ethereum advisory for the fix does recommend rotation.
- **Multi-client operation:** Running more than one independent implementation is what limits the blast radius of any single client going unmaintained.
- **GraphQL endpoint:** If `--graphql` is enabled on public-facing nodes, disable it until the node runs v1.13.0, which adds the query depth limit.

---

## Acknowledgments

The issues this document records were reported, and in several cases fixed, in public before `v1.13.0`.
Community members and security researchers opened them at `etclabscore/core-geth`, where each was closed or
left unaddressed. This work builds on theirs:

- **@niooss-ledger**, a Ledger security researcher, opened
  [issue #692](https://github.com/etclabscore/core-geth/issues/692) in February 2026, documenting
  CVE-2025-24883, CVE-2026-22862 and CVE-2026-22868.
- **@shrikus** reported the ECIES crash-loop stack trace in that issue, and opened
  [pull request #701](https://github.com/etclabscore/core-geth/pull/701) to remove the memsize debug
  endpoint, one of the two blockers to building on a supported Go.
- **@tornadocontrib** opened [pull request #683](https://github.com/etclabscore/core-geth/pull/683) in June
  2025, "Support go 1.24 compiler," which also carried a fix for CVE-2025-24883; it was closed without merge.
- **@cloorc** opened [issue #690](https://github.com/etclabscore/core-geth/issues/690) and
  [pull request #691](https://github.com/etclabscore/core-geth/pull/691) in November 2025, upgrading `blst`
  to fix the failing build.
- **@blackmennewstyle** reported in [issue #680](https://github.com/etclabscore/core-geth/issues/680) in
  April 2025 that `v1.12.20` does not compile on Go 1.24.
- **@solopool** reported the `eth_syncing` regression in
  [issue #697](https://github.com/etclabscore/core-geth/issues/697) in April 2026; it remains open.

White B0x carried the remediation through to the released `v1.13.0`, funded by the Ethereum Classic DAO.

---

## Supporting this work

The client is maintained by Ethereum Classic's core developers and open community contributions, in
the [`ethereumclassic`](https://github.com/ethereumclassic) GitHub organization. The Ethereum Classic
DAO is the legal entity standing behind that work, and the floor under it when nobody else is paying.
Mining pools, centralized exchanges, issuers of Ethereum Classic financial products, Etchash mining
hardware manufacturers and large holders all depend on this client. If your operation relies on
Ethereum Classic, please consider sponsoring its maintenance. [Support this work](../support.md) says
who maintains and funds it and how to sponsor it.

---

## References

- Patched client: https://github.com/ethereumclassic/core-geth
- etclabscore/core-geth (the repository audited): https://github.com/etclabscore/core-geth
- ETC Cooperative 2024 retrospective (maintenance mode; client maintenance passing to other stakeholders): https://etccooperative.org/etc-cooperative-retrospective-2024.pdf
- ETC Cooperative Q1 2025 report (maintenance mode in effect): https://etccooperative.org/posts/2025-06-24-q1-report-en
- Issue #692 (public CVE disclosure, Ledger security researcher): https://github.com/etclabscore/core-geth/issues/692
- PR #683 (community Go 1.24 + CVE-2025-24883 fix, closed without merge): https://github.com/etclabscore/core-geth/pull/683
- PR #694 (v1.12.21 emergency patch discussion and CVE analysis): https://github.com/etclabscore/core-geth/pull/694
- PR #696 (v1.12.22 CVE backports): https://github.com/etclabscore/core-geth/pull/696
- Issue #697 (eth_syncing regression introduced by v1.12.22): https://github.com/etclabscore/core-geth/issues/697
- v1.12.21 ("Aegis") release: https://github.com/etclabscore/core-geth/releases/tag/v1.12.21
- v1.12.22 ("Hermes") release: https://github.com/etclabscore/core-geth/releases/tag/v1.12.22
- Go vulnerability database: https://vuln.go.dev
- go-ethereum security advisories: https://github.com/ethereum/go-ethereum/security/advisories
- CVE-2025-24883 (GHSA-q26p-9cq4-7fc2): https://github.com/ethereum/go-ethereum/security/advisories/GHSA-q26p-9cq4-7fc2
- CVE-2026-22862 (GHSA-mr7q-c9w9-wh4h): https://github.com/ethereum/go-ethereum/security/advisories/GHSA-mr7q-c9w9-wh4h
- CVE-2026-26315 (GHSA-m6j8-rg6r-7mv8): https://github.com/ethereum/go-ethereum/security/advisories/GHSA-m6j8-rg6r-7mv8
- CVE-2026-26314 (GHSA-2gjw-fg97-vg3r): https://github.com/ethereum/go-ethereum/security/advisories/GHSA-2gjw-fg97-vg3r
- CVE-2026-26313 (GHSA-689v-6xwf-5jf3): https://github.com/ethereum/go-ethereum/security/advisories/GHSA-689v-6xwf-5jf3
- graphql-go GHSA-mh3m-8c74-74xh: https://github.com/graph-gophers/graphql-go/security/advisories/GHSA-mh3m-8c74-74xh
