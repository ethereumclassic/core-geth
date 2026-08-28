# Core-Geth Security Audit — March 2026

**Client:** Core-Geth (Ethereum Classic)  
**Audit Date:** March 2026  
**Repository audited:** [github.com/etclabscore/core-geth](https://github.com/etclabscore/core-geth)  
**Last substantive release there:** v1.12.20 (10 June 2024)  
**Emergency patches there:** [v1.12.21](https://github.com/etclabscore/core-geth/releases/tag/v1.12.21) (18 March 2026) · [v1.12.22](https://github.com/etclabscore/core-geth/releases/tag/v1.12.22) (28 March 2026) — CVE-only backports; Go 1.21 EOL toolchain unchanged  
**Patched repository:** [github.com/ethereumclassic/core-geth](https://github.com/ethereumclassic/core-geth)  
**Patched release:** v1.13.0  
**Patched by:** The core-geth Authors  
**Auditors:** Ethereum Classic Core Developers

---

## What operators need to do

**Track releases at [`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth)
and upgrade to v1.13.0.** Every finding in this document, and every finding in the
[August 2026 follow-up](./2026-08-security-followup.md), is resolved there.

**Findings apply to every release in the v1.12.x line, including the most recent.**

- **v1.12.20** (10 June 2024) — six CVEs and a GraphQL depth-limit denial of service, all
  unpatched. Built on a Go toolchain that reached end of life two months later.
- **v1.12.21** (18 March 2026) — cut during a live attack on ETC bootnodes, five hours after
  the issue was opened. Backported two CVEs. Toolchain unchanged.
- **v1.12.22** (28 March 2026) — backported the remaining CVEs and introduced a regression in
  `eth_syncing` that reports `highestBlock` incorrectly, reported as
  [#697](https://github.com/etclabscore/core-geth/issues/697). Toolchain unchanged.
- **v1.12.23** (14 August 2026) — a substantive p2p hardening series. The `eth_syncing`
  regression is not addressed, and the toolchain is still Go 1.21, now twenty-four months past
  end of life. One response cap in it can disconnect peers that are answering correctly; see
  the follow-up.
- **Advisory identifiers in the v1.12.21 and v1.12.22 release notes are not the ones the
  advisory records assign.** CVE-2026-26315 is attached to work that is CVE-2026-26314. An
  operator matching advisories against those notes will not find 26314 addressed anywhere and
  may conclude a fix is missing when it shipped. v1.12.23's notes carry no advisory
  identifiers at all.

**Two properties of that release line are worth stating plainly, because they bear on how much
assurance a release carries rather than on any individual finding.** The emergency releases
were each authored, reviewed and merged by one person, in sixty seconds and under two minutes
respectively; the pull requests record this. And the organization that owns the repository has
published that it does not expect to continue this work. Its
[2024 retrospective](https://etccooperative.org/etc-cooperative-retrospective-2024.pdf) states:

> "Without any way to obtain funding, the ETC Coop will be in maintenance mode, with spending
> minimized, until the funding runs out. At that time, it will be up to other stakeholders to
> take on any required maintenance of the ETC client, unless a new plan materializes."

**v1.13.0 is that continuation.** It carries every fix in this document, the Go 1.26 toolchain,
the p2p hardening series, and ETC-specific network tooling, and it is developed in the open at
`ethereumclassic/core-geth` with review before merge. Operators tracking the v1.12.x line
should move their tracking there.

---

## Executive Summary

During cross-client testing, the Ethereum Classic Core Developers identified that `etclabscore/core-geth` — the primary Ethereum Classic execution client — had received no security maintenance since June 2024, a 21-month gap. Six CVEs and one GraphQL depth-limit DoS were found unpatched, spanning cryptographic key validation, P2P protocol memory exhaustion, and query processing. CVE-2026-22868 is not applicable to ETC in normal operation but is patched in v1.13.0 regardless. The Go toolchain underpinning the client had also reached end-of-life (EOL) in August 2024, exposing all deployed nodes to unpatched runtime vulnerabilities for 19 months.

Disclosures to `etclabscore/core-geth` in early 2025 received no response. A public security disclosure by a Ledger security researcher in February 2026 likewise received no response until an active attack on ETC bootnodes in March 2026 forced emergency patches (v1.12.21 and v1.12.22 at `etclabscore/core-geth`). Those patches addressed the CVE backports but left the client on the Go 1.21 EOL toolchain with no ETC-specific modernization. The full remediation — Go 1.26 upgrade, ETC network tooling, and all CVE fixes — is released as v1.13.0 under the `ethereumclassic` GitHub organization.

---

## Background

Cross-client interoperability testing surfaced the Core-Geth maintenance gap in early 2026. The Ethereum Classic Core Developers' review of the go-ethereum security advisory database and Go vulnerability database against the v1.12.20 codebase found all six CVEs unpatched. The security remediation and Go toolchain upgrade are released as v1.13.0 together with Go 1.26 and ETC-specific network tooling.

**Where this work is done, and why.** The `etclabscore` organization is controlled by ETC
Cooperative, which entered maintenance mode at the end of 2024 and said so in its own published
reporting. Its
[2024 retrospective](https://etccooperative.org/etc-cooperative-retrospective-2024.pdf), covering
the year ended 31 December 2024, gives the cause and the position. Its recurring income had
ended: the arrangement under which Grayscale LLC contributed a third of the fees from its
Ethereum Classic Trust ran to the end of its two-year term, with **March 2022 the last month
the Cooperative earned fees from it**, and no replacement funding followed. Its statement of
the resulting position — that maintenance of the ETC client would fall to other stakeholders —
is quoted at the top of this document. Its
[Q1 2025 report](https://etccooperative.org/posts/2025-06-24-q1-report-en) confirms that
position was in effect: *"we are now in maintenance mode and spending has decreased
significantly."*

Moving the client to a community repository was also proposed publicly from inside the
organization before that happened. The same retrospective reproduces, on page 24, a slide from
a conference talk by the Cooperative's Senior Editor, headed "ETC Pathway Proposed by Donald
McIntyre", whose second item reads **"Move Core Geth to the Ethereum Classic community
repository"**. It is presented there as one person's proposal rather than as adopted policy,
but it is carried in the Cooperative's own annual report.

[`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth) was created on
**21 December 2024** as a fork of `etclabscore/core-geth`, under the same GitHub organization
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
| June 10, 2024 | `etclabscore/core-geth` v1.12.20 released — last substantive release from ETC Cooperative staff |
| August 2024 | Go 1.21 reaches end-of-life; core-geth toolchain enters EOL status |
| January 23, 2025 | Last commit to `etclabscore/core-geth` — a GitHub Actions dependency bump (not a code change) |
| January 30, 2025 | CVE-2025-24883 published (go-ethereum GHSA-q26p-9cq4-7fc2) |
| Throughout 2025 | Private security disclosures sent to `etclabscore/core-geth`; no response received |
| June 30, 2025 | Community member [@tornadocontrib](https://github.com/tornadocontrib) opens [PR #683](https://github.com/etclabscore/core-geth/pull/683) to `etclabscore/core-geth` — Go 1.24 upgrade and CVE-2025-24883 fix included; closed without merge when contributor deleted their fork in February 2026 |
| January 13, 2026 | CVE-2026-22862 and CVE-2026-22868 published |
| February 4, 2026 | Ledger security researcher [@niooss-ledger](https://github.com/niooss-ledger) opens [issue #692](https://github.com/etclabscore/core-geth/issues/692) publicly documenting CVE-2025-24883, CVE-2026-22862, CVE-2026-22868 — no response |
| February 17, 2026 | go-ethereum v1.17.0 released, carrying the fixes for CVE-2026-26313, CVE-2026-26314 and CVE-2026-26315 |
| February 18, 2026 | GHSA-689v-6xwf-5jf3, GHSA-2gjw-fg97-vg3r and GHSA-m6j8-rg6r-7mv8 published; the CVE records follow on February 19 |
| February 26, 2026 | CVE-2025-24883 patch authored |
| Early March 2026 | Go toolchain upgrade and remaining CVE patches authored |
| March 18, 2026 | Active attack on ETC bootnodes — ECIES handshake crash-loop (`crypto/ecies.symDecrypt` panic) exploited in production; [@diega](https://github.com/diega) opens [PR #694](https://github.com/etclabscore/core-geth/pull/694) and self-merges it 60 seconds later with no peer review, releasing [v1.12.21](https://github.com/etclabscore/core-geth/releases/tag/v1.12.21) ("Aegis") approximately 5 hours after the issue was first reported. This was the first code response there in 21 months, forced by the live attack rather than prior disclosures. A cryptographic security patch authored, reviewed, and merged by one person, in a repository with no other active reviewers, is itself a supply-chain risk: the process has no second reviewer able to catch a defective or malicious change shipped under cover of an emergency. |
| March 18, 2026 | @niooss-ledger [documents remaining unpatched CVEs](https://github.com/etclabscore/core-geth/pull/694#issuecomment-4089185353) after v1.12.21: CVE-2025-24883, CVE-2026-26313, and CVE-2026-26315 still unaddressed |
| March 20–21, 2026 | The security work is submitted to `ethereumclassic/core-geth` as individually scoped pull requests ([#10](https://github.com/ethereumclassic/core-geth/pull/10)–[#36](https://github.com/ethereumclassic/core-geth/pull/36)), one per CVE with test coverage and linked CVE references, and cross-references the set in [issue #692](https://github.com/etclabscore/core-geth/issues/692) |
| March 28, 2026 | [v1.12.22 "Hermes"](https://github.com/etclabscore/core-geth/releases/tag/v1.12.22) released at `etclabscore/core-geth` ([PR #696](https://github.com/etclabscore/core-geth/pull/696)) — remaining CVE backports; Go 1.21 EOL toolchain unchanged, no ETC-specific modernization |
| May 2026 | No further activity at `etclabscore/core-geth`; `ethereumclassic/core-geth` continues toward v1.13.0 |

**Note on v1.12.21 and v1.12.22:** Those emergency patches address the CVE backports and are a safer option than v1.12.20 for operators who have not yet migrated. However, they remain on the Go 1.21 EOL toolchain and do not include the ETC network tooling or DNS discovery updates in v1.13.0. Operators running v1.12.x should upgrade to the latest release there at minimum and track [ethereumclassic/core-geth](https://github.com/ethereumclassic/core-geth) for v1.13.0. Plan migration to [Fukuii](https://fukuii.org).

### Prior Maintainers

Core-Geth is a fork of [multi-geth](https://github.com/multi-geth/multi-geth), originally created and maintained by **Wei Tang** ([@sorpaas](https://github.com/sorpaas)). Multi-geth was the first multi-network go-ethereum fork with first-class ETC support, and its architecture is the direct ancestor of core-geth's chain configuration system.

The core-geth fork was then developed by ETC Labs until they left the ETC ecosystem in 2022. ETC Cooperative-paid staff maintained the client through the Spiral hard fork up until announcing maintenance mode for the client in December 2024:

- **Isaac Ardis** ([@meowsbits](https://github.com/meowsbits)) — primary architect and long-term maintainer
- **Diego López León** ([@diega](https://github.com/diega)) — release manager; cut the v1.12.20 release
- **Chris Ziogas** ([@ziogaschr](https://github.com/ziogaschr)) — contributor and maintainer

---

## Vulnerability Summary

| CVE | Severity | Component | etclabscore/core-geth | v1.13.0 (ethereumclassic) |
|-----|----------|-----------|------------------------|--------------------------|
| CVE-2025-24883 | High | crypto — secp256k1 key deserialization | Backported in v1.12.22 | Patched |
| CVE-2026-22862 | High | crypto/ecies — ECIES decrypt length check | Backported in v1.12.21 | Patched |
| CVE-2026-26315 | High | crypto/ecies + secp256k1 — ECIES GenerateShared / IsOnCurve | Backported in v1.12.22 | Patched |
| CVE-2026-26314 | High | crypto/secp256k1 — coordinate field boundary bypass | Backported in v1.12.21 | Patched |
| CVE-2026-22868 | Medium | txpool / P2P — KZG DoS (blob/KZG proof verification) | Declared N/A to ETC¹ | Patched (peer disconnect on invalid proof) |
| CVE-2026-26313 | High | P2P — RLP item count memory exhaustion | Mitigated in v1.12.22²  | Patched |
| — (GraphQL depth) | Medium | RPC — unbounded query nesting DoS | Not addressed | Patched |

¹ [@diega stated](https://github.com/etclabscore/core-geth/issues/692) that CVE-2026-22868 is "not applicable to ETC" because ETC does not support EIP-4844 blob/KZG transactions. The v1.13.0 patch nonetheless implements the go-ethereum fix: peer disconnect on invalid KZG proof via `ErrKZGVerificationError` sentinel.

² Both mitigations prevent the out-of-memory crash by refusing to allocate per declared item. Neither avoids scanning the payload: `rlp.CountValues` walks the tag bytes of the whole list with no early exit, and it runs inside `RawList`'s own decode as well as in v1.12.22's pre-decode count. A peer can therefore still spend CPU proportional to message size on a message that will be rejected. The distinction between the two releases is where the bound comes from — a static item ceiling in v1.12.22, the pending request in v1.13.0 — not whether the payload is read.

---

## Vulnerability Details

### CVE-2025-24883 — Off-Curve Public Key in UnmarshalPubkey

**Severity:** High
**CVSS v3.1:** `AV:N/AC:H/PR:N/UI:N/S:U/C:H/I:H/A:H` — **7.4 High**
**GHSA:** GHSA-q26p-9cq4-7fc2
**Component:** `crypto/crypto.go` — `UnmarshalPubkey()`
**Affected:** etclabscore/core-geth ≤ v1.12.20
**Patched:** ethereumclassic/core-geth v1.13.0
**Commit:** `8e40b7e41`
**Upstream reference:** go-ethereum PR #31100 / commit `159fb1a1d`

**Description:**
`UnmarshalPubkey()` decoded a 65-byte uncompressed public key into `(x, y)` field elements but did not verify that the resulting point lies on the secp256k1 curve. A malicious peer could supply an off-curve point that passes unmarshaling without error, then produces invalid or undefined results in all subsequent ECDSA or ECDH operations that consume the deserialized key.

**Impact:**
Any code path that deserializes an untrusted public key — including P2P node identity validation and RLPx (Recursive Length Prefix transport) handshake processing — could be supplied an off-curve point. Downstream cryptographic operations on the invalid key can produce incorrect signatures, incorrect shared secrets, or panics depending on the caller.

**Fix:**
Added an `IsOnCurve(x, y)` check immediately after coordinate extraction in `UnmarshalPubkey()`. Off-curve points now return `errInvalidPubkey` before any downstream use.

```go
if !S256().IsOnCurve(x, y) {
    return nil, errInvalidPubkey
}
```

---

### CVE-2026-22862 — ECIES (Elliptic Curve Integrated Encryption Scheme) Decrypt Ciphertext Length Undercheck

**Severity:** High
**CVSS v3.1:** `AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H` — **7.5 High** (confirmed live exploitation, March 18, 2026)
**GHSA:** GHSA-mr7q-c9w9-wh4h
**Component:** `crypto/ecies/ecies.go` — `Decrypt()`
**Affected:** etclabscore/core-geth ≤ v1.12.20
**Patched:** ethereumclassic/core-geth v1.13.0
**Commit:** `dc73f2e4f`
**Upstream reference:** go-ethereum commit `638741b08`

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

**Observed Attack — 18 March 2026:**
This vulnerability was actively exploited against the ETC mainnet classic bootnodes (ams3, sfo3) on 18 March 2026. Malicious P2P traffic sent crafted `auth` messages with undersized ECIES payloads, triggering the out-of-bounds slice allocation and crashing each node on the next inbound handshake attempt. Because the crash occurred in `listenLoop`, the node process exited and restarted under the service manager, only to crash again on the next malicious connection — producing a crash-loop.

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

Per the v1.12.21 release notes ([PR #694](https://github.com/etclabscore/core-geth/pull/694)): bootnode **sfo3** had accumulated **805+ restart cycles** on v1.12.20 at the time of the patch, confirming the attack was ongoing and automated. The patched binary was deployed to **ams3** first and confirmed stable before sfo3 was upgraded.

---

### CVE-2026-26315 — ECIES GenerateShared Accepts Unvalidated Public Key

**Severity:** High
**CVSS v3.1:** `AV:N/AC:H/PR:N/UI:N/S:U/C:H/I:N/A:N` — **5.9 Medium** (key-oracle attack requires repeated unauthenticated handshake attempts)
**GHSA:** GHSA-m6j8-rg6r-7mv8
**Component:** `crypto/ecies/ecies.go` — `GenerateShared()`
**Affected:** etclabscore/core-geth ≤ v1.12.20
**Patched:** ethereumclassic/core-geth v1.13.0
**Commit:** `2d3528803`
**Upstream reference:** go-ethereum commit `46bee92f9`

**Description:**
The RLPx handshake uses ECIES decryption on unauthenticated input from the network. `GenerateShared()` — called during ECDH (Elliptic Curve Diffie-Hellman) shared-secret derivation — accepted a `*PublicKey` without verifying it lies on the curve. An ephemeral public key with `X == nil`, `Y == nil`, or coordinates not satisfying the secp256k1 curve equation would proceed into ECDH multiplication and fail only at MAC verification. The MAC failure reveals to the attacker whether the faulty key survived ECDH, which can be used as an oracle to leak bits of the node's static P2P private key across multiple handshake attempts.

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

### CVE-2026-26314 — secp256k1 IsOnCurve Field Boundary Bypass

**Severity:** High
**CVSS v3.1:** `AV:N/AC:H/PR:N/UI:N/S:U/C:H/I:H/A:H` — **8.1 High**
**GHSA:** GHSA-2gjw-fg97-vg3r
**Component:** `crypto/secp256k1/curve.go` — `IsOnCurve()`; `crypto/secp256k1/ext.h` — `secp256k1_ext_scalar_mul()`; `crypto/signature_nocgo.go` — `btCurve.IsOnCurve()`
**Affected:** etclabscore/core-geth ≤ v1.12.20
**Patched:** ethereumclassic/core-geth v1.13.0
**Commit:** `2d3528803` (bundled with CVE-2026-26315)
**Upstream reference:** go-ethereum commit `895a8597c`

**Description:**
`IsOnCurve()` verified the curve equation `y² ≡ x³ + b (mod P)` but did not first verify that the coordinates `x` and `y` are within the field, i.e., strictly less than the curve prime `P`. Due to modular arithmetic, coordinates equal to or greater than `P` may still satisfy the curve equation when reduced, but they represent invalid (non-canonical) points. Additionally, the C-level `secp256k1_ext_scalar_mul` function did not check the return value of `secp256k1_fe_set_b32`, which returns 0 when a coordinate is out-of-field. A crafted out-of-field coordinate could therefore bypass `IsOnCurve` and proceed into scalar multiplication, producing undefined or attacker-influenced results.

**Impact:**
An attacker supplying points with out-of-field coordinates that still satisfy the naive curve check can pass a gate that is supposed to reject invalid public keys, potentially leading to node crash or, in a consensus context, divergent state computation.

**Identifier mismatch.** `etclabscore/core-geth` labels this issue CVE-2026-26315. The
advisory records assign that identifier to the ECIES handshake issue documented in the previous
section, and assign **CVE-2026-26314** to this one: OSV lists `895a8597c` — "crypto/secp256k1: fix
coordinate check" — as the fix for [GHSA-2gjw-fg97-vg3r](https://osv.dev/vulnerability/GHSA-2gjw-fg97-vg3r),
aliased to CVE-2026-26314 and GO-2026-4507, first released in go-ethereum v1.16.9. The two
advisories share the title "Go Ethereum affected by DoS via malicious p2p message" and are
distinguished by their CWE class and fix commit rather than by their titles, which is what makes
them straightforward to transpose.

The label appears in four places there — commit `46bba8dfc`, a doc comment and the Go test
function name `TestIsOnCurveRejectsCoordinatesAboveP_CVE_2026_26315` in
`crypto/secp256k1/curve_cve_test.go`, and the published release notes for v1.12.21 and v1.12.22.
CVE-2026-26314 appears nowhere in that repository. Operators matching advisories against release
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

### CVE-2026-22868 — KZG (Kate-Zaverucha-Goldberg) Blob Proof Verification DoS

**Severity:** Medium
**CVSS v3.1:** `AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:L` — **5.3 Medium** (inactive on ETC in normal operation; rated against theoretical worst-case exposure)
**Component:** `core/txpool/validation.go` — `validateBlobSidecar()`; `eth/fetcher/tx_fetcher.go` — `Enqueue()`
**Affected:** etclabscore/core-geth ≤ v1.12.20 (code present but inactive on ETC)
**Patched:** ethereumclassic/core-geth v1.13.0
**Commit:** `1419c5310`
**Upstream reference:** go-ethereum commit `fdfd1235a` (v1.16.8)

**ETC Applicability:** [@diega stated that this is not applicable to ETC](https://github.com/etclabscore/core-geth/issues/692) — ETC does not support EIP-4844 blob transactions, so the KZG code path is not reached in normal operation on the ETC network. The v1.12.22 release at `etclabscore/core-geth` does not address it. The v1.13.0 patch follows the go-ethereum approach: introduces `ErrKZGVerificationError` as a sentinel error and disconnects any peer that delivers a transaction with an invalid KZG proof, preventing repeated DoS attempts from the same peer.

**Description:**
KZG blob proof verification is computationally expensive. A malicious peer could repeatedly broadcast blob transactions with invalid KZG proofs, causing the node to perform the full expensive cryptographic verification on each delivery attempt before rejecting the transaction. Because the error was not distinguished from other validation failures, the peer was not disconnected and could continue flooding the node indefinitely.

**Impact:**
On ETC: no active exposure in normal operation (blob transactions are rejected at an earlier validation stage). The code path is nonetheless present in the inherited go-ethereum codebase and could be reached by a crafted peer interaction.

**Remediation:**
Ported from go-ethereum commit `fdfd1235a` (v1.16.8). Adds `ErrKZGVerificationError` as a sentinel error in `core/txpool/errors.go` and wraps KZG validation failures in `validateBlobSidecar()` with it. In `eth/fetcher/tx_fetcher.go`, any delivery flagged with this violation immediately triggers `f.dropPeer(delivery.origin)` — preventing the offending peer from repeatedly triggering expensive cryptographic verification.

```go
case errors.Is(err, txpool.ErrKZGVerificationError):
    violation = err
// ...
if delivery.violation != nil {
    f.dropPeer(delivery.origin)
}
```

---

### CVE-2026-26313 — P2P RLP (Recursive Length Prefix) Item Count Memory Exhaustion

**Severity:** High
**CVSS v3.1:** `AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H` — **7.5 High**
**GHSA:** GHSA-689v-6xwf-5jf3
**Component:** `eth/protocols/eth/`, `eth/protocols/snap/`, `p2p/tracker/`
**Affected:** etclabscore/core-geth ≤ v1.12.20
**Patched:** ethereumclassic/core-geth v1.13.0
**Commit:** `5d0cb8b34`
**Upstream reference:** go-ethereum PR #33835

**Description:**
The P2P message handler validated message size against a 10 MiB cap (`maxMessageSize`) before decoding, but did not validate the number of items declared in the RLP list header. A malicious peer could craft a valid RLP list header claiming millions of tiny items within the 10 MiB budget. When `msg.Decode` ran, it would allocate a pointer or struct object for each declared item before any further validation, causing out-of-memory crashes proportional to the declared item count rather than the actual payload size.

The attack requires only a valid peer connection at the devp2p (Ethereum device peer-to-peer wire protocol) handshake level — no authenticated session is needed to send crafted message payloads.

**Impact:**
Remote out-of-memory (OOM) crash of any reachable Core-Geth node. An attacker with a single peer connection could crash the node by sending one crafted message on the eth or snap protocol sub-channel.

**Fix:**
Response contents are decoded lazily. `rlp.RawList` holds a list without decoding its elements, so a message is validated before any per-item object is allocated, and a claimed item count no longer drives allocation. On the snap protocol a response is additionally matched against its pending request before its contents are read. This is go-ethereum's own remediation, PR #33835.

An earlier form of this patch counted items before decoding instead, because the upstream refactor could not be applied to this tree at the time. Both are no longer needed together, and the counting layer was removed once the refactor landed; go-ethereum carries no equivalent.

Response messages are bounded by the request they answer rather than by a fixed ceiling, which is both tighter than a static limit and incapable of refusing an honest peer. Transaction broadcasts keep an explicit limit of 5,000, and block announcements, which are unsolicited and not lazily decoded, keep one of 2,048.


**Status in v1.12.22:** Partially mitigated — OOM crash prevented, CPU amplification DoS remains open.  
**Status in v1.13.0:** Patched. Response messages are bounded by their pending request, transaction broadcasts and block announcements by explicit ceilings.

---

### GraphQL Query Depth DoS

**Severity:** Medium
**CVSS v3.1:** `AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H` — **7.5 High** (assumes GraphQL endpoint exposed; off by default)
**Related advisory:** GHSA-mh3m-8c74-74xh / CVE-2022-21708 — a stack-overflow DoS in `graphql-go` reachable only when `MaxDepth` is enabled. Affected `<v1.3.0`; patched in v1.3.0
**Component:** `graphql/service.go`; `go.mod` — `graphql-go v1.3.0 → v1.9.0`
**Affected:** etclabscore/core-geth ≤ v1.12.20
**Patched:** ethereumclassic/core-geth v1.13.0
**Commit:** `6c2d383fa`

**Disclosure:** Identified during cross-client security review, combining review of the published `graphql-go` advisory with direct inspection of `graphql/service.go`, which confirmed that no depth limit was configured at the application layer. No separate responsible disclosure to `etclabscore` was required as the dependency advisory was already public; the finding was included in the patch set submitted to `ethereumclassic/core-geth`.

**Description:**
The GraphQL endpoint (`--graphql` flag) had no query complexity or depth limit. Deeply nested queries — for example, a query recursively nesting block references — would cause the server to perform unbounded recursive schema traversal, exhausting CPU and memory on the serving node.

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

**Issue:** Core-Geth v1.12.20 was built and shipped on Go 1.21, which reached end-of-life in August 2024. From that date onward, vulnerabilities in the Go standard library — including `net/http`, `crypto/tls`, and `encoding/json` — received no patches from the Go team for this toolchain version, and all binaries compiled against it remained exposed.

**Impact:** Node operators running that binary were exposed to Go runtime security issues for a minimum of 19 months (August 2024 through March 2026). The Go vulnerability database lists multiple advisories against Go 1.21 in this period, including issues in the TLS stack and HTTP/2 server.

**Fix:** Core-Geth v1.13.0 builds on Go 1.26 (current stable as of March 2026). The upgrade proceeded in two steps — 1.21 → 1.24 (removing the incompatible `fjl/memsize` dependency and fixing `go vet` format string errors introduced by 1.24's stricter checks), then 1.24 → 1.26 (updating all `golang.org/x/` dependencies for compatibility). The `blst` cryptography dependency was simultaneously upgraded from v0.3.11 to v0.3.16 to fix a C23 `typedef bool` incompatibility with the Alpine-based Docker build environment.

**Commit:** `8385cf8e8`

---

## Release Timeline

*For the full disclosure and maintenance history, see the [Disclosure Timeline](#disclosure-timeline) above.*

| Date | Release | Notes |
|------|---------|-------|
| June 10, 2024 | [v1.12.20](https://github.com/etclabscore/core-geth/releases/tag/v1.12.20) | Final substantive release at `etclabscore/core-geth` |
| August 2024 | — | Go 1.21 reaches end-of-life; no further toolchain security patches from the Go team |
| March 18, 2026 | [v1.12.21 "Aegis"](https://github.com/etclabscore/core-geth/releases/tag/v1.12.21) | Emergency patch for active ECIES crash-loop attack ([PR #694](https://github.com/etclabscore/core-geth/pull/694)); Go 1.21 EOL toolchain unchanged |
| March 28, 2026 | [v1.12.22 "Hermes"](https://github.com/etclabscore/core-geth/releases/tag/v1.12.22) | Remaining CVE backports ([PR #696](https://github.com/etclabscore/core-geth/pull/696)); Go 1.21 EOL toolchain unchanged; `eth_syncing` regression introduced (issue [#697](https://github.com/etclabscore/core-geth/issues/697), open) |
| August 14, 2026 | [v1.12.23 "Argos"](https://github.com/etclabscore/core-geth/releases/tag/v1.12.23) | 32 commits: the delayed-decoding p2p hardening series (`rlp.RawList`, deferred block/body/receipt/transaction/snap decoding, request-response validation in `p2p/tracker`) plus seven upstream go-ethereum backports. Go 1.21 EOL toolchain unchanged. No CVE identifiers in the release notes |
| — | v1.13.0 | To be released at `ethereumclassic/core-geth` — all six CVEs and the GraphQL DoS patched, Go 1.26 toolchain |

---

## Risk Assessment

| Risk Area | Severity | Description | Mitigation |
|-----------|----------|-------------|------------|
| Six CVEs + GraphQL DoS | Critical | Six CVEs and one GraphQL depth-limit DoS unaddressed for 21 months in the primary ETC client (CVE-2026-22868 inactive on ETC but patched regardless) | All patched in v1.13.0 |
| Cryptographic oracle (CVE-2026-26315) | High | Repeated handshake attempts could leak P2P node key bits | Patched; key rotation recommended for long-running nodes |
| Remote crash via ECIES (CVE-2026-22862) | High | Any peer can crash a node during RLPx handshake with a malformed ECIES payload | Patched in v1.13.0 |
| Remote OOM via RLP (CVE-2026-26313) | High | Any peer can OOM-crash a node with a single crafted P2P message | Patched in v1.13.0 |
| CPU amplification DoS — CVE-2026-26313 residual (v1.12.22 only) | Medium | v1.12.22 mitigation scans full RLP payload before rejecting oversized messages; ~2,500× more work per attack message than v1.13.0; malicious peers can exhaust CPU without causing OOM | Unmitigated in v1.12.22; patched in v1.13.0 |
| Go Runtime EOL | High | 19 months on unsupported Go toolchain; runtime CVEs accumulated unpatched | Upgraded to Go 1.26 in v1.13.0 |
| Unreviewed emergency patches (supply chain risk) | High | v1.12.21 and v1.12.22 were each authored, reviewed, and merged by the same individual with no independent peer review — v1.12.21 in 60 seconds, v1.12.22 in under 2 minutes. With no second reviewer, a defective or malicious change shipped under cover of an emergency has nothing to catch it. The risk is structural and applies to any future emergency patch cut this way. | v1.13.0 is developed under the `ethereumclassic` org with multi-contributor review; Fukuii migration eliminates the dependency entirely |

---

## Scope

**Audit target:** `etclabscore/core-geth` at tag `v1.12.20` (commit `c2fb44129`)  
**Audit date range:** February – March 2026  
**In scope:**
- All Go source packages inherited from go-ethereum with known CVE exposure
- Go toolchain version and dependency security posture
- P2P protocol handler input validation (devp2p, RLPx, eth, snap sub-protocols)
- RPC endpoint security (GraphQL, JSON-RPC)

**Out of scope:**
- Consensus layer correctness and ETC protocol compliance
- Fukuii client codebase
- EVM execution correctness
- Dependencies not listed in the go-ethereum security advisory database
- Infrastructure (bootnode operators, DNS, CDN)

---

## Methodology

The audit was initiated during cross-client interoperability testing. The `etclabscore/core-geth` codebase at tag `v1.12.20` was cross-referenced against the go-ethereum security advisory database (GitHub Security Advisories for `github.com/ethereum/go-ethereum`) and the Go vulnerability database (`vuln.go.dev`). Each published advisory was assessed for applicability by tracing shared code ancestry between core-geth and go-ethereum at the affected package level.

**Exploitability assessment:** Where a CVE had confirmed exploitation records (CVE-2026-22862: active ECIES crash-loop on ETC bootnodes, March 18, 2026), exploitability was treated as confirmed. For the remaining CVEs, static code analysis of the affected call paths was used to establish reachability from unauthenticated network input and assess preconditions.

**Patch development:** Applicable fixes were cherry-picked from upstream go-ethereum commits where possible. In two cases (CVE-2026-26313, CVE-2026-26314), structural divergence from go-ethereum produced merge conflicts incompatible with clean cherry-pick; fixes were manually ported by adapting the upstream approach to core-geth's code structure.

**Tooling:** `govulncheck` (Go vulnerability scanner), manual code review, go-ethereum advisory DB, `vuln.go.dev`.

**Testnet validation:** All patches were validated against the Mordor testnet. Validation covered: successful node sync resumption after patch application, P2P handshake stability under normal peer traffic, no regression in JSON-RPC method responses (`eth_syncing`, `eth_blockNumber`, `net_peerCount`), and block processing correctness against known Mordor block hashes.

---

## Network Migration Path

The ETC network is migrating to [Fukuii](https://fukuii.org) ([github.com/fukuii-project/fukuii-cli](https://github.com/fukuii-project/fukuii-cli)) as its ETC-native execution client. Core-Geth is maintained to give that transition a stable path; operators should plan their migration.

**If you are running v1.12.20 or earlier, upgrade to v1.12.22 immediately** to address the most critical CVEs. v1.13.0 — the full Go 1.26 upgrade, ETC-specific modernization and every fix in this audit — is released from [`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth).

---

## Recommendations

- **Node operators (v1.12.20 or earlier):** Upgrade to v1.12.22 immediately to patch the most critical CVEs. Then track [github.com/ethereumclassic/core-geth](https://github.com/ethereumclassic/core-geth) for v1.13.0. Nodes on v1.12.20 or earlier are exposed to remote crash and potential key-oracle attacks from any peer.
- **Infrastructure providers and exchanges:** Treat the upgrade to v1.13.0 as a security-critical update, not a routine version bump. Begin planning migration to Fukuii.
- **Long-running nodes:** Consider rotating the P2P node key (`--nodekey`) as a precaution against CVE-2026-26315 oracle exposure. Any node reachable from the public internet over the 21-month gap was potentially targeted.
- **Multi-client operation:** Run at minimum two independent clients (e.g., Core-Geth + Fukuii, or Core-Geth + Besu) for redundancy. Multi-client operation is what limits the blast radius of a single client going unmaintained.
- **GraphQL endpoint:** If `--graphql` is enabled on public-facing nodes, disable it until v1.13.0 is released or verify v1.12.22 is in place before re-opening the port.

---


## References

- Patched client: https://github.com/ethereumclassic/core-geth
- etclabscore/core-geth (the repository audited): https://github.com/etclabscore/core-geth
- ETC Cooperative 2024 retrospective (maintenance mode; client maintenance passing to other stakeholders): https://etccooperative.org/etc-cooperative-retrospective-2024.pdf
- ETC Cooperative Q1 2025 report (maintenance mode in effect): https://etccooperative.org/posts/2025-06-24-q1-report-en
- Issue #692 — public CVE disclosure (Ledger security researcher): https://github.com/etclabscore/core-geth/issues/692
- PR #683 — community Go 1.24 + CVE-2025-24883 fix, closed without merge: https://github.com/etclabscore/core-geth/pull/683
- PR #694 — v1.12.21 emergency patch discussion and CVE analysis: https://github.com/etclabscore/core-geth/pull/694
- PR #696 — v1.12.22 CVE backports: https://github.com/etclabscore/core-geth/pull/696
- Issue #697 — eth_syncing regression introduced by v1.12.22: https://github.com/etclabscore/core-geth/issues/697
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
