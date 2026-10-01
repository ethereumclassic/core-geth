---
title: MESS
description: "How to check whether MESS (ECBP-1100) is on in Core-Geth, turn it on or off, and choose the setting for your node."
---

MESS (Modified Exponential Subjective Scoring, ECBP-1100) makes a node resist deep chain
reorganizations. It never changes whether a block is valid. It only decides which of two valid
chains the node keeps when a competing chain would replace blocks it already accepted. The
[MESS confirmation calculator](../guides/mess-calculator.md) shows how much protection it gives.

## Check whether MESS is on

```shell
$ geth --classic attach --exec 'admin.ecbp1100Status().enabled' <datadir>/geth.ipc
true
```

`true` means MESS is on. `false` means it is off, or paused while the node is out of sync or has
fewer than five peers.

## Turn MESS on or off

```shell
$ geth --classic --mess          # on
$ geth --classic --mess=false    # off
```

Without either flag, the node uses the setting its release ships with. Give every node you run the
same setting: during a deep reorganization, nodes with different settings can follow different
chains.

## Which setting fits which operator

The setting only matters during a deep reorganization, when a competing chain would replace blocks
your node already accepted. With MESS on, your node refuses that chain. With MESS off, it follows
whichever chain has the most total difficulty, including one an attacker produced.

| If you run | What the setting decides for you |
|---|---|
| An exchange, a custodian or a payment processor | Whether your node follows a deep reorganization that reverses deposits it already accepted |
| A mining pool that credits deposits | The same as for an exchange, and the miner's cost below |
| A mining pool or solo miner | The miner's cost below |
| A public RPC endpoint, or your own node | Which chain you serve or follow during a deep reorganization |

For a miner, MESS on has a cost: if other nodes accept a deep reorganization, the blocks your node
mines on the chain it kept are orphaned.

## Flags

| Flag | What it does |
|---|---|
| `--mess` | Turns MESS on |
| `--mess=false` | Turns MESS off |
| `--mess.activate=<block>` | Turns MESS on from that block |
| `--mess.deactivate=<block>` | Turns MESS off from that block |
| `--mess.nodisable` | Turns MESS on, and keeps it on while the node is out of sync or short of peers |

- If you give both `--mess.activate` and `--mess.deactivate`, the later block wins.
- `--mess.activate` wins over `--mess=false`. Flags that contradict each other log a warning at
  startup.
- The older spellings also work: `--ecbp1100`, `--override.ecbp1100.deactivate` and
  `--ecbp1100.nodisable`.
- Each flag can also be set with an environment variable, such as `GETH_MESS=true` or
  `GETH_MESS_ACTIVATE=15000000`. The older names, such as `GETH_ECBP1100`, also work.

## Config file

Each flag has a line under `[Eth]` in a config file that does the same thing, and
`geth dumpconfig` run with your flags writes them for you:

| Flag | Config file line |
|---|---|
| `--mess` | `OverrideECBP1100Deactivate = 18446744073709551614` |
| `--mess=false` | `OverrideECBP1100 = 18446744073709551614` |
| `--mess.activate=<block>` | `OverrideECBP1100 = <block>` |
| `--mess.deactivate=<block>` | `OverrideECBP1100Deactivate = <block>` |
| `--mess.nodisable` | `ECBP1100NoDisable = true` |

A node started from the file keeps the setting until you remove the line. For example, a config
file dumped with `--mess` keeps MESS on until you remove its `OverrideECBP1100Deactivate` line.
A flag given on the command line wins over the file.

## Turn MESS on without a restart

```shell
$ geth --classic attach --exec 'admin.ecbp1100("latest")' <datadir>/geth.ipc
true
```

This turns MESS on from the current block. It uses the `admin` namespace, which the node serves over
IPC. At the next start, your flags and config file decide again.

## When the node switches it off by itself

MESS pauses while the node has fewer than five peers, or while its newest block is more than 390
seconds old. It resumes once the node is back in sync with enough peers. The log shows
`Disabled artificial finality features` with the reason, and
[Troubleshooting](troubleshooting.md#what-do-disabled-artificial-finality-features-and-reorg-disallowed-mean)
explains it. `--mess.nodisable` keeps MESS on through both.

## The two defaults

A release ships one of two defaults, from two Best Practice documents:

| Default | Ethereum Classic | Mordor |
|---|---|---|
| ECBP-1100's | on from block 11,380,000 | on from block 2,380,000 |
| ECBP-1110's | on from block 11,380,000, off from block 19,250,000 | on from block 2,380,000, off from block 10,400,000 |

Both chains are past the ECBP-1110 blocks, so under ECBP-1110's default MESS is off, and under
ECBP-1100's it is on. `--mess` and `--mess=false` override either one.
