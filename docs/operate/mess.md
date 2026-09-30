---
title: MESS
description: "MESS (ECBP-1100) in Core-Geth: what it does during a deep chain reorganization, the two defaults ECBP-1100 and ECBP-1110 describe, and how to turn it on or off."
---

MESS, Modified Exponential Subjective Scoring (ECBP-1100), makes a node resist deep chain
reorganizations. When a competing chain would replace blocks the node has already accepted, MESS
asks that chain for more total difficulty the older the point where the two chains split. A short
reorganization passes as it always did. A deep one, the kind a majority-hashrate attack produces,
needs far more work than the chain it replaces. The [MESS confirmation calculator](../guides/mess-calculator.md)
shows how much.

MESS is not a consensus rule. It never changes whether a block is valid, only which of two valid
chains the node prefers. Nodes that must agree on the chain, such as an exchange's deposit and
withdrawal nodes, should all run with the same MESS setting: during a deep reorganization, a node
with MESS on and a node with MESS off can prefer different chains.

## The two defaults

Two Best Practice documents describe when MESS applies, and a release bundles one of their two
windows for each network:

| Network | MESS applies from block | ECBP-1100's deactivation block | ECBP-1110's deactivation block |
|---|---|---|---|
| Ethereum Classic | 11,380,000 | none | 19,250,000 |
| Mordor | 2,380,000 | none | 10,400,000 |

- **ECBP-1100's default:** MESS applies from the activation block and never stops.
- **ECBP-1110's default:** MESS applies from the activation block and stops at the deactivation block.

**Both chains are past their deactivation blocks, so the two defaults no longer agree.** A node on
ECBP-1110's default does not apply MESS, and a node on ECBP-1100's default does. `--mess` turns MESS
on and `--mess=false` turns it off, whichever default the node would otherwise run.
[Check your node](#check-your-node) shows which one it runs.

## Check your node

`admin.ecbp1100Status()` reports the window a node has and whether MESS applies at its head, and
changes nothing. It is in the `admin` namespace, which the node serves over IPC:

```shell
$ geth --classic attach --exec 'admin.ecbp1100Status()' <datadir>/geth.ipc
```

- **`activatedAtBlock` and `defaultDisabledAtBlock`** are the two ends of the window.
  `defaultDisabledAtBlock` is `null` on ECBP-1100's window, and a block number on ECBP-1110's or on
  one set with `--mess.deactivate`. `--mess` sets it to `0xfffffffffffffffe`, a block out of reach,
  and `--mess=false` sets `activatedAtBlock` to the same value.
- **`enabled`** says whether MESS applies right now: `nodeSwitch` is on, and the node's `head` is at
  or past `activatedAtBlock` and short of `defaultDisabledAtBlock`, if that is a block number.
- **`nodeSwitch`** is the node's own switch. The node turns it on once it is in sync with enough
  peers, whatever the window, and off as
  [When the node switches it off by itself](#when-the-node-switches-it-off-by-itself) describes.

A synced Ethereum Classic node on ECBP-1110's default reports `enabled` as `false` and
`defaultDisabledAtBlock` as `0x125bb50`, block 19,250,000. On ECBP-1100's default it reports
`enabled` as `true` and `defaultDisabledAtBlock` as `null`.

## What changed, and when

| When | What |
| --- | --- |
| 11 October 2020, block 11,380,000 | MESS activates on Ethereum Classic with no deactivation block, as ECBP-1100 describes |
| `v1.12.17`, December 2023 | a deactivation is scheduled at block 19,250,000, as ECBP-1110 describes |
| 5 February 2024, block 19,250,000 | the chain reaches Spiral, and nodes on `v1.12.17` or later stop applying MESS |
| `v1.13.0`, 14 September 2026 | the deactivation is removed and the activation kept, so MESS applies again |

`v1.12.16` and earlier carry no deactivation. The two defaults agree until the deactivation block:
on Ethereum Classic, a node on either one applied MESS from 11 October 2020 to 5 February 2024, three
years and four months.

## Why the two defaults exist

**ECBP-1100's default dates from the 2020 attacks, and ECBP-1110's from January 2024.** The record,
with sources:

- **31 July 2020: a 51% attack that ran for 12 hours.** The
  [MESS testing report](https://medium.com/etc-core/mess-testing-results-report-4cba96ed92fa), published by ETC
  Labs and written by OpenRelay, carries Bitquery's estimate that the attack cost about $192,000 to run and double
  spent about $5.6 million of ETC. The same report calculates that with MESS in place the attacker would have
  needed **31 times more hashing power** to hold that reorganization, which would have cost more than it took.
- **25 September 2020:** an entire core developers' call was
  [given over to 51% attack solutions](https://ethereumclassic.org/blog/2020-09-25-core-devs-call-51-attack-solutions).
  Seven proposals competed, from merged mining to checkpointing to VeriBlock. MESS is the one that shipped.
- **28 September and 11 October 2020:** MESS activated on Mordor at block 2,380,000 and on Ethereum Classic at
  **block 11,380,000**, shipped in Core-Geth v1.11.15
  ([release announcement](https://medium.com/etc-core/ethereum-classic-stakeholders-critical-security-release-to-prevent-51-attacks-aa83596a0903),
  [client upgrade](https://ethereumclassic.org/blog/2020-10-10-mess-client-upgrade)). This client still activates
  it at that block, in
  [`params/config_classic.go`](https://github.com/ethereumclassic/core-geth/blob/main/params/config_classic.go).
- **It existed in exactly one client.** That announcement states it plainly: *"The MESS Security Feature is ONLY
  AVAILABLE in the Core-Geth client."* Hyperledger Besu did not carry it, and by the
  [Thanos upgrade in November 2020](https://medium.com/etc-core/ethereum-classic-prepares-for-the-thanos-hard-fork-upgrade-2e1e52633cc)
  Parity Ethereum, OpenEthereum, Multi-Geth and Geth Classic were deprecated and would no longer follow the chain.
  That is checkable rather than remembered: the archived OpenEthereum chainspec for Ethereum Classic carries the
  Phoenix transitions at block 10,500,839 and contains no ECIP-1099 epoch change at all, in
  `openethereum/openethereum/ethcore/res/ethereum/classic.json` of the
  [client archive](https://github.com/fukuii-project/archive-reference-material). Those clients kept up through
  Phoenix in June 2020 and stopped at Thanos.
- **January 2024:** [ECBP-1110](https://ecips.ethereumclassic.org/ECIPs/ecip-1110) recommended that clients ship
  MESS off by default. Its reasoning rests on Ethereum Classic holding roughly 85% of apparent compatible hashrate
  at the time, which is a claim about market conditions rather than about the code.

**Neither is a rule.** ECBP-1100 and ECBP-1110 are both Best Practice documents, so each client
chooses its own default and each operator can override it. The status of both is open in
[ECIPs #580](https://github.com/ethereumclassic/ECIPs/pull/580).

## When the node switches it off by itself

MESS is not unconditional. Inside its window, it applies only while the node's own switch is on, and
the node turns that switch off when either of these holds. Where the window is open at the head, it
logs `Disabled artificial finality features` with the reason:

- **`low peers`:** the node has fewer than five peers.
- **`stale safety interval`:** the node's newest block is more than 390 seconds old. The node checks
  on the same interval, so this happens roughly 6.5 to 13 minutes after the head stops advancing.

It turns the switch back on once the node is in sync with enough peers and no peer advertises a chain
with more total difficulty.

That is the limit of the defense. An eclipse attack, a network partition or a network-wide stall
can switch MESS off, and a heavier chain advertised by a peer keeps it off through the
reorganization. A node in that state follows the chain with the most total difficulty, as it would
without MESS. [Troubleshooting](troubleshooting.md#what-do-disabled-artificial-finality-features-and-reorg-disallowed-mean)
explains the log lines.

`--mess.nodisable` keeps the switch on through both conditions once it is on, so inside the window
MESS also applies while the node is out of sync or short of peers. It does not move the window: past
the deactivation block, MESS still does not apply. Where the node would have switched it off, the log
shows `Preventing disable artificial finality`. `--mess.nodisable=false` undoes
`ECBP1100NoDisable = true` in a config file.

## Which setting fits which operator

**The choice only matters during a deep reorganization.** In normal operation a node with MESS on and
a node with MESS off follow the same chain and produce the same blocks. What the setting decides is
which chain the node keeps when a competing chain arrives that would replace blocks it already
accepted, which is what a majority-hashrate attack produces.

| If you run | The setting that fits | Why |
|---|---|---|
| An exchange, a custodian, a payment processor | On, and the same on every node you operate | A deep reorganization is the attack aimed at you: the July 2020 attacker took about $5.6 million from exchanges rather than from the protocol. Deposit and withdrawal nodes that disagree about MESS can prefer different chains, which is worse than either setting |
| A mining pool that also credits deposits | On, and the same on every node | You carry the exchange's exposure as well as the miner's |
| A mining pool or solo miner | Either, deliberately | This is the one case with a cost on both sides. See below |
| A public RPC endpoint, or your own node | On | You serve, or follow, the chain a deep reorganization would have replaced |

**The miner's trade, stated plainly.** With MESS on, your node refuses a deep reorganization that
other nodes may accept. If that happens, you keep mining on the chain you had while others move, and
the blocks you find in that window are orphaned for you and not for them. With MESS off, your node
follows whichever chain carries the most total difficulty, including a chain an attacker paid to
produce. Neither is free, and the choice is about which failure you would rather have. What it is
not is a question about block validity: no block becomes invalid either way.

**Whichever you choose, the mix of settings across the network is a market condition and not a
property of the code.** ECBP-1110's reasoning rests on Ethereum Classic holding roughly 85% of
apparent compatible hashrate when it was written. That figure is checkable today, and it is what the
recommendation depends on.

## Turn MESS off

```sh
geth --classic --mess=false
```

`--mess=false` moves the activation block out of reach, so MESS does not apply under either default.
`geth dumpconfig` writes it into a config file as:

```toml
[Eth]
OverrideECBP1100 = 18446744073709551614
```

## Turn MESS on

```sh
geth --classic --mess
```

`--mess` moves both ends of the window: it clears the off switch that `--mess=false` writes into a
config file, and it moves the deactivation block out of reach, so MESS applies from the activation
block under either default. An activation block set any other way is kept, and an explicit
`--mess.activate` or `--mess.deactivate` wins over `--mess`.

`geth dumpconfig` run with `--mess` writes that deactivation into the config file, under `[Eth]`:

```toml
[Eth]
OverrideECBP1100Deactivate = 18446744073709551614
```

A node started from that file keeps MESS on, whether or not `--mess` is passed again, until the line
is removed.

## Flags and config file keys

| Flag | Older spelling, still accepted | Config file key, under `[Eth]` | Effect |
|---|---|---|---|
| `--mess` | none | `OverrideECBP1100Deactivate = 18446744073709551614` | Turns MESS on |
| `--mess=false` | none | `OverrideECBP1100 = 18446744073709551614` | Turns MESS off |
| `--mess.activate=<block>` | `--ecbp1100` | `OverrideECBP1100` | Sets the activation block, and wins over `--mess` and `--mess=false` |
| `--mess.deactivate=<block>` | `--override.ecbp1100.deactivate` | `OverrideECBP1100Deactivate` | Sets the deactivation block, and wins over `--mess` |
| `--mess.nodisable` | `--ecbp1100.nodisable` | `ECBP1100NoDisable = true` | Keeps the node's switch on once it is on, bypassing both automatic switch-offs |

A negative block number, or one of 2^64 or more, is refused when the flags are parsed.

## Change it on a running node

`admin_ecbp1100` sets the activation block in the running node. It is in the `admin` namespace,
which the node serves over IPC. At the next start, the flags and the config file decide again. Its
one parameter is a block number:

- a hex quantity such as `0x1`, up to `0x7fffffffffffffff`. A decimal number is refused, and so is
  the `--mess=false` value;
- `earliest` for block 0, or `latest` and `pending` for the current head. `finalized` and `safe` are
  refused.

It moves only the activation block, so it cannot turn MESS on once the head is past the deactivation
block: on a node running ECBP-1110's default, it returns `false` and MESS stays off. To apply MESS
there, restart the node with `--mess`.

Check the result with `admin.ecbp1100Status()`. The [admin module](../JSON-RPC-API/modules/admin.md)
lists both methods.
