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
$ geth --classic attach --exec 'admin.messStatus().enabled' <datadir>/geth.ipc
true
```

`true` means MESS is on. `false` means it is off, or paused while the node is out of sync or has
fewer than five peers.

## Turn MESS on or off

```shell
$ geth --classic --mess          # on; default turns on.
$ geth --classic --mess=true     # on
$ geth --classic --mess=false    # off
```

Without `--mess`, the node uses the setting its release ships with. Give every node you run the
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

## Three ways to set MESS

| Way | Use it to | Takes effect | Lasts |
|---|---|---|---|
| A [flag](#flags), such as `geth --classic --mess` | set MESS for one run of the node | when the node starts | until the node stops |
| A [config file line](#config-file), written with `geth --classic --mess dumpconfig node.toml` | keep a setting across restarts | each time the node starts with `--config node.toml` | until you remove the line |
| An [admin call](#admin-calls), such as `admin.mess(true)` | change a running node without a restart | at once | until the node stops |

All three set the same things by the same rules: the block MESS turns on from, the block it turns
off from, and whether the node may switch MESS off while it is out of sync. A flag on the command
line wins over the config file. An admin call wins until the node stops; at the next start, the
flags and the config file decide again.

### Flags

| Flag | What it does |
|---|---|
| `--mess` or `--mess=true` | Turns MESS on |
| `--mess=false` | Turns MESS off |
| `--mess.activate=<block>` | Turns MESS on from that block |
| `--mess.deactivate=<block>` | Turns MESS off from that block |
| `--mess.nodisable` or `--mess.nodisable=true` | Turns MESS on, and keeps it on while the node is out of sync or short of peers |
| `--mess.nodisable=false` | Lets the node switch MESS off again while it is out of sync or short of peers, over a config file's `ECBP1100NoDisable = true` |

- `--mess.activate` and `--mess.nodisable` turn MESS on with no end, unless you give
  `--mess.deactivate`. With both blocks given, the later one wins.
- `--mess.activate` wins over `--mess=false`. Flags that contradict each other log a warning at
  startup.
- The older spellings also work: `--ecbp1100`, `--override.ecbp1100.deactivate` and
  `--ecbp1100.nodisable`.
- Each flag can also be set with an environment variable, such as `GETH_MESS=true` or
  `GETH_MESS_ACTIVATE=15000000`. The older names, such as `GETH_ECBP1100`, also work.

### Config file

`geth dumpconfig` writes a config file from the flags you give it, and starts no node. A node
started with `--config` reads the file each time it starts:

```shell
$ geth --classic --mess dumpconfig node.toml    # writes node.toml, starts nothing
$ geth --classic --config node.toml             # starts the node with MESS on
```

| Flag given to `geth dumpconfig` | Line it writes under `[Eth]` |
|---|---|
| `--mess` or `--mess=true` | `OverrideECBP1100Deactivate = 18446744073709551614` |
| `--mess=false` | `OverrideECBP1100 = 18446744073709551614` |
| `--mess.activate=<block>` | `OverrideECBP1100 = <block>` |
| `--mess.deactivate=<block>` | `OverrideECBP1100Deactivate = <block>` |
| `--mess.nodisable` or `--mess.nodisable=true` | `ECBP1100NoDisable = true` |
| `--mess.nodisable=false` | No line, and it leaves out a file's `ECBP1100NoDisable` line |

- A line means what its flag means, and you can also write one by hand.
- The file keeps the setting until you remove the line. For example, a file dumped with `--mess`
  keeps MESS on until you remove its `OverrideECBP1100Deactivate` line.
- A flag on the command line wins over the file.
- `geth dumpconfig` never reads a running node, so it does not show a change made with an admin
  call.

### Admin calls

```shell
$ geth --classic attach --exec 'admin.mess(true).enabled' <datadir>/geth.ipc
true
```

| Call | What it does | Like the flag |
|---|---|---|
| `admin.mess(true)` | Turns MESS on | `--mess` |
| `admin.mess(false)` | Turns MESS off | `--mess=false` |
| `admin.messActivate(<block>)` | Turns MESS on from that block | `--mess.activate=<block>` |
| `admin.messDeactivate(<block>)` | Turns MESS off from that block | `--mess.deactivate=<block>` |
| `admin.messNoDisable(true)` | Turns MESS on, and keeps it on while the node is out of sync or short of peers | `--mess.nodisable` |
| `admin.messNoDisable(false)` | Lets the node switch MESS off again while it is out of sync or short of peers | `--mess.nodisable=false` |
| `admin.messStatus()` | Shows the current setting, and changes nothing | none |

- A block is `"latest"`, which means the current block, or a block number in hex, such as
  `"0x125bb50"` for block 19,250,000. A decimal number is refused.
- The calls follow the flags' rules. `admin.messActivate` and `admin.messNoDisable(true)` turn MESS
  on with no end, unless a deactivation block was given, with `--mess.deactivate`, a config file
  line or `admin.messDeactivate`. With both blocks given, the later one wins, and when the
  deactivation block keeps MESS off, the node logs a warning naming it.
- `admin.mess(true)` and `admin.messNoDisable(true)` move an activation block that keeps MESS off
  back to the block MESS first applied from on your network, as `--mess` does.
- Called while the node is out of sync or short of peers, `admin.messNoDisable(true)` turns MESS on
  at the next sync, and keeps it on from then.
- The ECBP-1100 names also work: `admin.ecbp1100`, `admin.ecbp1100Deactivate`,
  `admin.ecbp1100NoDisable` and `admin.ecbp1100Status`.
- The calls use the `admin` namespace, which the node serves over IPC.
- None of them is saved. To keep a change, use the flag or its config file line.

The [admin module](../JSON-RPC-API/modules/admin.md) page lists each call with its parameters and
examples.

## When the node switches it off by itself

MESS pauses while the node has fewer than five peers, or while its newest block is more than 390
seconds old. It resumes once the node is back in sync with enough peers. The log shows
`Disabled artificial finality features` with the reason, and
[Troubleshooting](troubleshooting.md#what-do-disabled-artificial-finality-features-and-reorg-disallowed-mean)
explains it. `--mess.nodisable` and `admin.messNoDisable(true)` keep MESS on through both.

## The two defaults

A release ships one of two defaults, from two Best Practice documents:

| Default | Ethereum Classic | Mordor |
|---|---|---|
| ECBP-1100's | on from block 11,380,000 | on from block 2,380,000 |
| ECBP-1110's | on from block 11,380,000, off from block 19,250,000 | on from block 2,380,000, off from block 10,400,000 |

Both chains are past the ECBP-1110 blocks, so under ECBP-1110's default MESS is off, and under
ECBP-1100's it is on. `--mess` and `--mess=false` override either one.

## After an upgrade

When a node starts with MESS on where the client that last used its database had it off, or off
where it had it on, it logs one warning. An upgrade to a release with the other default does this:

```text
WARN [...] ECBP1100 (MESS) is configured differently from the client that last used this database was=on now=off ... keep_previous=--mess
```

`was` and `now` show the setting before and after. To keep the setting the node had, add the flag
that `keep_previous` names to the command that starts the node. The warning appears once, because
the node stores the new setting with its chain data. It does not appear when you set MESS with a
flag or a config file line.
