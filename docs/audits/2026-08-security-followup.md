# Core-Geth Security Follow-Up — August 2026

**Follows:** [`2026-03-security-audit.md`](./2026-03-security-audit.md)
**Subject:** `etclabscore/core-geth` v1.12.23 ("Argos"), released 14 August 2026
**Prepared by:** The core-geth Authors
**Method:** direct measurement against `etclabscore/core-geth` at tag `v1.12.23`, and
against the advisory records at OSV and the GitHub Advisory Database

---

## What operators need to do

**Track releases at [`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth)
and upgrade to v1.13.0.** Every finding in this document and in the
[March 2026 audit](./2026-03-security-audit.md) is resolved there.

v1.12.23 is real security work, and it does not close the two findings that matter most to a
node operator:

- **The Go toolchain is still 1.21**, which reached end of life in August 2024. Every node built
  from any v1.12.x release, including this one, carries whatever runtime vulnerabilities have
  accumulated in twenty-four months, independent of the client's own code.
- **The `eth_syncing` regression introduced in v1.12.22 is untouched.**
  [#697](https://github.com/etclabscore/core-geth/issues/697) reports `highestBlock` incorrectly
  against the network head, which misleads exchanges, explorers and monitoring that read it to
  decide whether a node is synced.
- **One response cap in this release can disconnect peers that are answering correctly.** The
  release corrects that class of error for account ranges and leaves it in place for storage
  ranges; see the finding below.
- **Its release notes carry no advisory identifiers**, and the notes for the two releases before
  it attach CVE-2026-26315 to work the records assign to CVE-2026-26314.

**The organization that owns the repository has published that it does not expect to continue
this work.** Its
[2024 retrospective](https://etccooperative.org/etc-cooperative-retrospective-2024.pdf) states:

> "Without any way to obtain funding, the ETC Coop will be in maintenance mode, with spending
> minimized, until the funding runs out. At that time, it will be up to other stakeholders to
> take on any required maintenance of the ETC client, unless a new plan materializes."

**v1.13.0 is that continuation**, developed in the open at `ethereumclassic/core-geth` with
review before merge. It carries the Go 1.26 toolchain, every CVE fix, the same p2p hardening
series measured here, and the storage-range correction this release does not make.

---

## Summary

A third release has been cut at `etclabscore/core-geth` since the March 2026 audit.
v1.12.23 is substantive security work: 32 commits carrying a delayed-decoding hardening
series for the p2p protocols, plus seven go-ethereum backports.

Two findings from the March audit are unchanged by it. The client remains on a Go toolchain
that reached end-of-life in August 2024, now 24 months. The `eth_syncing` regression
introduced in v1.12.22 is not addressed.

This document records what v1.12.23 contains, what it leaves open, and where v1.13.0 stands.

---

## What v1.12.23 contains

Measured across `v1.12.22..v1.12.23`: 32 commits, 71 files, +3,036 / -796 lines.

**Delayed-decoding series.** The release defers decoding of p2p message contents until
after the message has been validated, so a peer cannot force allocation proportional to a
claimed item count before anything has checked it:

| Commit | Change |
|---|---|
| `b25f25cf1` | `rlp: add RawList for working with un-decoded lists` |
| `887da91c6` | delay block body decoding until content validation |
| `717942ef3` | delay receipts decoding until content validation |
| `de544c0c3` | delay transaction packet decoding |
| `4a2b0effc` | delay snap message decoding |
| `4cdebb034` | delay propagated block decoding until content validation |
| `4a8cd33ea` | discard message before size check |
| `0f8924182` | validate responses against pending requests in `p2p/tracker` |
| `ff7d98f1a` | fix crash in `p2p/tracker` clean when the tracker is stopped |
| `85ac4d685` | recalibrate snap account-range response item cap |

**Upstream backports.** Seven changes from go-ethereum: `#34787` (stop serving on
unavailable responses), `#34745` (drop peers sending invalid bodies or receipts), `#34870`
(batch index in deliver reconstruct), `#34976` (snap error message), `#33803` (embedded
node size validation), `#32210` (announcement drop logic), `#30918` (prevent hanging
dispatch).

**The delayed-decoding refactor is adaptable to core-geth.** go-ethereum's `rlp.RawList`
approach does not cherry-pick cleanly onto this codebase — 13 merge conflicts — and
v1.12.23 shows it can nonetheless be carried across with manual work.

---

## What v1.12.23 does not change

**The Go toolchain is unchanged and remains end-of-life.** Measured at tag `v1.12.23`:

| Location | Value |
|---|---|
| `go.mod` | `go 1.21` |
| `.github/workflows/test-linux.yml` | `go-version: '1.21'` |
| `Dockerfile` | `FROM golang:1.22-alpine` |
| `build/checksums.txt` | `version:golang 1.22.1` |

Go 1.21 reached end-of-life in August 2024 and receives no security patches. Every node
built from this release inherits whatever unpatched runtime vulnerabilities the toolchain
carries, independent of the client's own code. This was the March audit's finding and it is
unchanged across all three 2026 releases.

**The `eth_syncing` regression is not addressed.** Issue
[#697](https://github.com/etclabscore/core-geth/issues/697), introduced in v1.12.22, reports
`highestBlock` as incorrect relative to the network head. No commit in
`v1.12.22..v1.12.23` touches `internal/ethapi/`.

**The cryptographic fixes from v1.12.22 are present and unchanged.** Verified at tag: the
ECIES public-key validation in `GenerateShared`, and the field-boundary check in
`crypto/secp256k1/curve.go`, are both in place and byte-identical to this project's own
implementation.

---

## Finding: one response cap in v1.12.23 can still reject honest peers

v1.12.23 recalibrates the AccountRange item cap and documents why: the previous value
rejected honest responses. `snapResponseLimits` now reads

```go
AccountRangeMsg:  2*maxRequestSize/common.HashLength + 1,   // 32769
StorageRangesMsg: maxCodeLookups * 10,                      // 10240
ByteCodesMsg:     maxCodeLookups,                           //  1024
TrieNodesMsg:     maxTrieNodeLookups,                       //  1024
```

with a comment deriving the AccountRange figure from the request tracker's byte bound.

**StorageRanges is bounded the same way and did not get the same treatment.** Reading
`ServiceGetStorageRangesQuery` at that tag, the server appends slots until `size >= req.Bytes`,
with a `hardLimit` of `req.Bytes * (1 + stateLookupSlack)` above it. Neither is an item count.
A storage slot costs a hash plus its value, so at least 33 bytes on the wire, and a 512 KiB
budget therefore admits on the order of **15,900 slots** before the byte bound stops the
server — against a cap of **10,240**.

The two caps that remain item-derived are correct: `ByteCodesMsg` and `TrieNodesMsg` are
genuinely capped server-side at `maxCodeLookups` and `maxTrieNodeLookups`, so those numbers
are the honest maximum.

**Consequence.** A peer answering a large storage-range request with small slots can exceed
the cap without doing anything unusual, and is disconnected. This is the failure the
AccountRange recalibration was made to fix, in the message type beside it.

**This client carries no snap item ceiling at all.** A response is matched against the pending
request that produced it, in bytes for the two byte-bounded messages and in item count for the
other two, so the bound is whatever was actually asked for. That cannot refuse an honest reply,
which is the failure a fixed cap produces.

---

## CVE identifier reconciliation

The advisory records assign these identifiers, verified against OSV and each fix commit
confirmed in a full-history go-ethereum clone:

| Identifier | Issue | Fix commit | First release |
|---|---|---|---|
| CVE-2026-26313 · GHSA-689v-6xwf-5jf3 · GO-2026-4508 | p2p message memory exhaustion (CWE-770) | go-ethereum `#33835` | v1.17.0 |
| CVE-2026-26314 · GHSA-2gjw-fg97-vg3r · GO-2026-4507 | secp256k1 coordinate check (CWE-20) | `895a8597c` | v1.16.9 |
| CVE-2026-26315 · GHSA-m6j8-rg6r-7mv8 | ECIES public key validation in the RLPx handshake | `46bee92f9` | v1.16.9 |

**`etclabscore/core-geth` labels the secp256k1 issue CVE-2026-26315.** That identifier belongs
to the ECIES handshake issue. The correct identifier for the coordinate check is
CVE-2026-26314, which does not appear anywhere in that repository. The label appears in four
places: commit `46bba8dfc`, a doc comment and the Go test function name
`TestIsOnCurveRejectsCoordinatesAboveP_CVE_2026_26315` in `crypto/secp256k1/curve_cve_test.go`,
and the published release notes for v1.12.21 and v1.12.22. v1.12.23's notes carry no CVE
identifiers.

**Why this is easy to get wrong, and why it matters.** CVE-2026-26313 and CVE-2026-26314 are
published under the same title, "Go Ethereum affected by DoS via malicious p2p message". They
are distinguished by CWE class and fix commit, not by their titles. An operator matching
advisory identifiers against release notes to decide whether they are exposed will not find
CVE-2026-26314 addressed at `etclabscore/core-geth`, and may conclude it is outstanding when the
underlying fix shipped in v1.12.22. Both issues are in fact addressed there.

---

## Status of v1.13.0

v1.13.0 is being prepared at
[`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth), the repository
created on 21 December 2024 under the GitHub organization that also hosts the ECIP
repository. It is where contributions and release artifacts for this series belong. ETC
Cooperative's [2024 retrospective](https://etccooperative.org/etc-cooperative-retrospective-2024.pdf)
records that with the organization in maintenance mode, "it will be up to other stakeholders to
take on any required maintenance of the ETC client". Relative to
v1.12.23 it carries the Go 1.26 toolchain upgrade and the dependency graph that requires, the
six CVE patches with their tests, ETC network tooling, and the p2p hardening series. It is
the current release series. The network is migrating to [Fukuii](https://fukuii.org) as its
ETC-native execution client.

**Recommendation for operators is unchanged in direction and updated in target.** Nodes on
v1.12.20 or earlier should upgrade immediately; v1.12.23 is the latest release there and
is a safer position than v1.12.20 on both the cryptographic and p2p surfaces. It does not
resolve the toolchain exposure. Track `ethereumclassic/core-geth` for v1.13.0.

---

## References

- v1.12.23 release: https://github.com/etclabscore/core-geth/releases/tag/v1.12.23
- Delayed decoding PR: https://github.com/etclabscore/core-geth/pull/699
- `eth_syncing` regression: https://github.com/etclabscore/core-geth/issues/697
- CVE-2026-26313: https://osv.dev/vulnerability/GHSA-689v-6xwf-5jf3
- CVE-2026-26314: https://osv.dev/vulnerability/GHSA-2gjw-fg97-vg3r
- CVE-2026-26315: https://github.com/advisories/GHSA-m6j8-rg6r-7mv8
- go-ethereum v1.16.9 release: https://github.com/ethereum/go-ethereum/releases/tag/v1.16.9
- go-ethereum PR #33835 (delayed decoding): https://github.com/ethereum/go-ethereum/pull/33835
- Go release policy and end-of-life: https://go.dev/doc/devel/release
