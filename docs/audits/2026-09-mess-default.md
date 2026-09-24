---
description: "Why the bundled MESS default returns to ECIP-1110's, what each v1.12.x release actually shipped, and the two cases where an upgrade changes a node's chain-selection setting."
---

# MESS: the bundled default returns to ECIP-1110's

This document records why `v1.13.1` restores the ECBP-1100 (MESS) deactivation blocks that
`v1.13.0` removed, what every `v1.12.x` release actually shipped, and the narrow cases where an
upgrade changes a node's setting rather than leaving it alone. It is the companion to
[MESS](../operate/mess.md), which is the operator reference for the mechanism and its flags, and to
the [v1.13.0 release report](../release-reports/v1.13.0.md), which records the default that release
shipped.

**None of this is consensus.** ECBP-1100 and ECBP-1110 are Ethereum Classic Best Practice
documents, not validity rules. Every difference below changes which of two competing chains a node
prefers during a deep reorganization, never whether a block is valid. A client shipping one default
does not fork from a client shipping the other, and the fork ID is unchanged.

**Baseline:** every `v1.12.x` release from `v1.12.17` to `v1.12.23` was built from source and run,
and the chain configuration each one loads was read from its own startup output rather than from its
source. The same was done for this release. Both networks were measured for every binary.

**Work carried out by:** [White B0x](https://whiteb0x.com)

## What changed, and why

[ECIP-1110](https://ecips.ethereumclassic.org/ECIPs/ecip-1110) sets the bundled default to inactive
from block 19,250,000 on Ethereum Classic, the Spiral block, and from block 10,400,000 on Mordor.
Releases from `v1.12.17` carried those blocks, and both chains passed them long ago: read from a
local copy of each chain, Mordor reached its block on 12 January 2024 and Ethereum Classic reached
Spiral on 5 February 2024.

`v1.13.0` removed the deactivation blocks and kept the activation, so a node on that release applies
MESS past Spiral. `v1.13.1` restores them. The mechanism stays in the client and `--mess` turns it
on, which is what ECIP-1110 itself contemplates: the document sets a default, not a prohibition.

## What each release ships

Measured from each binary's own chain configuration output, both networks:

| Release | Ethereum Classic | Mordor |
|---|---|---|
| `v1.12.17` | 11,380,000 to 19,250,000 | 2,380,000 onward, **no deactivation** |
| `v1.12.18` through `v1.12.23` | 11,380,000 to 19,250,000 | 2,380,000 to 10,400,000 |
| `v1.13.0` | 11,380,000 onward, no deactivation | 2,380,000 onward, no deactivation |
| `v1.13.1` | 11,380,000 to 19,250,000 | 2,380,000 to 10,400,000 |

**`v1.13.1` and `v1.12.18` through `v1.12.23` ship identical windows on both networks.** A node
moving from any of those releases to this one sees no MESS change.

**`v1.12.17` is the one release that differs, and only on Mordor.** It carries the Ethereum Classic
deactivation but not Mordor's, which arrived in `v1.12.18`. A Mordor node on `v1.12.17` applies MESS
today and stops after upgrading.

## Two upgrade paths, and they are not alike

| | Migrating from `v1.12.x` | Updating from `v1.13.0` |
|---|---|---|
| Also involves | a different repository and a new release source to track, plus the required node key rotation | nothing beyond the release itself |
| MESS before | inactive, except the two cases below | active past Spiral |
| MESS after | inactive | inactive |
| Changes? | no, for almost every node | **yes, for every node** |

The ordering is worth stating plainly, because it inverts what an operator would expect: the path
with the large operational change carries no MESS change, and the path that is only an update is the
one whose chain-selection setting moves.

## Where a block number stopped being read as written

Releases on the `v1.12.x` line read fork block numbers as signed values. A number at or above 2^63,
which is 9,223,372,036,854,775,808, came out negative, and both ends of the MESS window inverted.

| What an operator set | What it means | `v1.12.x` did | This release does |
|---|---|---|---|
| Activation at or above 2^63 | never activate, so MESS off | MESS **on** | MESS off |
| Deactivation at or above 2^63 | never deactivate, so MESS on | MESS **off** | MESS on |
| Any value below 2^63 | as written | as written | as written |

**The activation row is the one that matters**, because pushing the activation out of reach is the
method ECIP-1110 documents for disabling MESS. An operator who followed the document, using a value
at or above 2^63, had MESS running the whole time. This release reads the value as written, so such
a node starts doing what its configuration always said.

Ordinary block numbers were never affected. This reading changed in `v1.13.0`, as part of reading
chainspec block numbers as unsigned, and is recorded here because the `v1.12.x` line is where nodes
are upgrading from.

## A node could not report its own setting

**No `v1.12.x` release has `admin_ecbp1100Status`.** There was no way to ask a running node whether
MESS was in force. That is why a misconfiguration of the kind above could persist unnoticed: the
setting was not readable, so nothing contradicted the operator's belief about it.

`v1.13.0` added the status call. This release adds the rest of the surface, so that every flag has a
runtime counterpart:

| Call | The flag it mirrors |
|---|---|
| `admin.mess(true\|false)` | `--mess` |
| `admin.ecbp1100(<block>)` | `--mess.activate` |
| `admin.ecbp1100Deactivate(<block>)` | `--mess.deactivate` |
| `admin.ecbp1100Status()` | none, it reports and changes nothing |

**`admin.mess` is new and it is the one to use.** MESS is a window with two ends, and on the bundled
default it is the deactivation end that holds it shut, so setting the activation alone does not turn
MESS on. `admin_ecbp1100` sets one end by design, mirroring `--mess.activate`; before this release
there was no call that reached the other, which meant no way to turn MESS on at runtime once a chain
was past its deactivation block.

## What an upgrading node reports

A node whose MESS setting changes across an upgrade now says so once, when it starts:

```
WARN ECBP1100 (MESS) is configured differently from the client that last used this database
     was=on now=off head=17,036,428
     stored.activation=2,380,000 stored.deactivation=none
     current.activation=2,380,000 current.deactivation=10,400,000
     keep_previous=--mess
```

It compares what the database records against what is in force, so it stays quiet for the nodes that
see no change, and it does not appear when an operator has set the value themselves. **It appears
once.** Opening the database records the new configuration, so later starts have nothing to compare
against. An operator whose upgrade is unattended should not rely on seeing it.

## Related

- [MESS](../operate/mess.md), the operator reference: what it does, the flags, and the trade per node type
- [Migrating from v1.12.x](../tutorials/v1.13.0-migration.md), which covers the upgrade in full
- [March 2026 security audit](2026-03-security-audit.md) and
  [August 2026 security follow-up](2026-08-security-followup.md), for the CVE work in this client
- [Go toolchain](2026-09-go-toolchain.md), which compares the same releases for build provenance
