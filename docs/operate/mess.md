---
title: MESS
description: "MESS (ECBP-1100) in Core-Geth: what it does during a deep chain reorganization, which releases apply it, and how to set it for your node."
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

## What changed, and when

| When | What |
| --- | --- |
| 11 October 2020, block 11,380,000 | MESS activates on Ethereum Classic and becomes the default |
| `v1.12.17`, December 2023 | a deactivation is scheduled at block 19,250,000 on Ethereum Classic |
| `v1.12.18` | the same is scheduled at block 10,400,000 on Mordor, which `v1.12.17` did not carry |
| 5 February 2024, block 19,250,000 | the chain reaches Spiral and the default changes: nodes on `v1.12.17` or later stop applying MESS |
| `v1.13.0`, 14 September 2026 | the deactivation is removed and the activation kept, so MESS applies again |

`v1.12.16` and earlier carry no deactivation. Between the middle rows MESS applied for three years
and four months; from Spiral onward the published default is that it does not, and `--mess` is how an
operator chooses otherwise.

**The table is the answer to "does my node apply MESS".** Find the release you run and read across.
A release is added here when it ships, so a version absent from the table is one this page predates.

## Why MESS exists

**MESS exists because of what happened the last time Ethereum Classic's client base thinned.** The record, with
sources:

- **31 July 2020: a 51% attack that ran for 12 hours.** The
  [MESS testing report](https://medium.com/etc-core/mess-testing-results-report-4cba96ed92fa), published by ETC
  Labs and written by OpenRelay, carries Bitquery's estimate that the attack cost about $192,000 to run and double
  spent about $5.6 million of ETC. The same report calculates that with MESS in place the attacker would have
  needed **31 times more hashing power** to hold that reorganization, which would have cost more than it took.
- **25 September 2020:** an entire core developers' call was
  [given over to 51% attack solutions](https://ethereumclassic.org/blog/2020-09-25-core-devs-call-51-attack-solutions).
  Seven proposals competed, from merged mining to checkpointing to VeriBlock. MESS is the one that shipped.
- **28 September and 10 October 2020:** MESS activated on Mordor at block 238,000 and on Ethereum Classic at
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

**What preceded the 2020 attacks was a thinning client base, not a change in the chain's rules.** That is worth
knowing when deciding your own setting: the hashrate distribution and client mix ECBP-1110's reasoning depends on
are not fixed, and the case for MESS is strongest in the troughs.

**Which window your client ships is a release decision, and the timeline above has each one.**
Whatever it is, `--mess` and `--mess=false` override it, and `admin.ecbp1100Status()` reports what your
node is actually doing, which is the only answer that accounts for the flags you pass.
Exchanges, custodians and payment processors are what a deep reorganization is aimed at,
the July 2020 attacker having taken $5.6 million from them rather than from the protocol, so those operators are
the ones for whom the flag is worth the thought. [Which setting fits which operator](#which-setting-fits-which-operator)
has the trade per node type.

**This is a client default, not a rule.** ECBP-1100 and ECBP-1110 are both Best Practice documents, so each
client chooses its own, and a client shipping a different default does not fork: it disagrees about which of two
competing chains to prefer during a deep reorganization. **Give every node in a fleet the same setting**, because
that disagreement is only visible once it already matters. The status of both documents is open in
[ECIPs #580](https://github.com/ethereumclassic/ECIPs/pull/580).

## Where it applies

| Network | MESS applies from block | Deactivation block |
|---|---|---|
| Ethereum Classic | 11,380,000 | 19,250,000 |
| Mordor | 2,380,000 | 10,400,000 |

**Both chains are long past their deactivation blocks**, so on any release that bundles one, MESS is
not applying there. Whether yours does depends on the release and on the flags you pass;
[MESS on this node](../getting-started/run-classic-node.md#mess-on-this-node) shows how
`admin.ecbp1100Status()` answers it for your node.

## When the node switches it off by itself

MESS is not unconditional. The node switches it off when either of these holds, and logs
`Disabled artificial finality features` with the reason:

- **`low peers`:** the node has fewer than five peers.
- **`stale safety interval`:** the node's newest block is more than 390 seconds old. The node checks
  on the same interval, so this happens roughly 6.5 to 13 minutes after the head stops advancing.

It switches MESS back on once the node is in sync with enough peers and no peer advertises a chain
with more total difficulty.

That is the limit of the defense. An eclipse attack, a network partition or a network-wide stall
can switch MESS off, and a heavier chain advertised by a peer keeps it off through the
reorganization. A node in that state follows the chain with the most total difficulty, as it would
without MESS. [Troubleshooting](troubleshooting.md#what-do-disabled-artificial-finality-features-and-reorg-disallowed-mean)
explains the log lines.

`--mess.nodisable` keeps MESS on through both conditions, so it also applies while the node is out
of sync or short of peers. Where the node would have switched it off, the log shows
`Preventing disable artificial finality`.

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

`--mess=false` moves the activation block out of reach, which is the method ECIP-1110 itself
documents for disabling MESS. On a release that bundles a deactivation the chain has already passed,
MESS is not applying anyway and the flag changes nothing.

## Turn MESS on

Pass `--mess`.

The bundled configuration activates MESS and then stops applying it at the block ECIP-1110 names,
so `--mess` moves both ends of that window: it clears an off switch on the activation side and
pushes the deactivation out of reach. Without the second half the flag would report success and
leave MESS off past the deactivation block. `--mess` also overrides an off switch written into a
config file, and `--mess.nodisable=false` overrides `ECBP1100NoDisable = true`.

An explicit `--mess.activate` or `--mess.deactivate` is applied after `--mess` and wins over it,
so a specific block number is always honored.

## Flags and config file keys

| Flag | Older spelling, still accepted | Config file key, under `[Eth]` | Effect |
|---|---|---|---|
| `--mess=false` | none | `OverrideECBP1100 = 18446744073709551614` | Turns MESS off |
| `--mess.activate=<block>` | `--ecbp1100` | `OverrideECBP1100` | Sets the activation block, and wins over `--mess=false` |
| `--mess.deactivate=<block>` | `--override.ecbp1100.deactivate` | `OverrideECBP1100Deactivate` | Sets a deactivation block |
| `--mess.nodisable` | `--ecbp1100.nodisable` | `ECBP1100NoDisable = true` | Keeps MESS on, bypassing both automatic switch-offs |

A negative block number, or one of 2^64 or more, is refused when the flags are parsed.

## Change it on a running node

`admin_ecbp1100` sets the activation block in the running node. It is in the `admin` namespace,
which the node serves over IPC. At the next start, the flags and the config file decide again. Its
one parameter is a block number:

- a hex quantity such as `0x1`, up to `0x7fffffffffffffff`. A decimal number is refused, and so is
  the `--mess=false` value;
- `earliest` for block 0, or `latest` and `pending` for the current head. `finalized` and `safe` are
  refused.

Check the result with `admin.ecbp1100Status()`. The [admin module](../JSON-RPC-API/modules/admin.md)
lists both methods.
