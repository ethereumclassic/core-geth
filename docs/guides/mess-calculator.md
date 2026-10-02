---
title: MESS confirmation calculator
---

On a node running [MESS](../operate/mess.md), a reorganization needs more total difficulty the older the
block where the two chains split. This calculator shows how much more. Exchanges and other integrators can
use it to see what a confirmation policy asks of an attacker.

## Calculator

<form id="mess-calculator" onsubmit="return false;">
<p><label>Time from the split to your node's head:
<input name="hours" type="number" min="0" step="0.25" value="1" style="width:6em"> hours</label></p>
<p>A competing chain needs at least <strong><output name="ratio">2.66</output> times</strong> the difficulty
your node's chain added since the split.</p>
</form>

## Required ratio by time

| Time from the split to the head | Required ratio |
|---|---|
| 1 minute | 1.00 |
| 10 minutes | 1.05 |
| 30 minutes | 1.44 |
| 1 hour | 2.66 |
| 2 hours | 6.97 |
| 3 hours | 12.85 |
| 4 hours | 19.26 |
| 5 hours | 25.12 |
| 6 hours | 29.38 |
| 6 hours 58 minutes 52 seconds, or more | 31.00 |

## How it is computed

MESS compares the two chains from the block where they split, their common ancestor. With `x` the number
of seconds from the common ancestor's timestamp to the head's timestamp, capped at 25,132:

```text
ratio = (128 + (3*x^2 - 2*x^3 // 25132) * 3840 // 25132^2) / 128
```

`//` is integer division. The node accepts the reorganization only if the competing chain added at least
`ratio` times the difficulty its own chain added since the split. The ratio starts at 1 and reaches its cap
of 31 at 25,132 seconds, just under 7 hours. The calculator follows `ecbp1100PolynomialV` in
`core/blockchain_af.go`.

## Reading it as hashrate

When both chains are mined over the same stretch of time, the difficulty each adds is roughly proportional
to its hashrate. So the ratio is roughly the multiple of the rest of the network's hashrate an attacker needs
to replace blocks that old on nodes running MESS. Without MESS, any chain with more total difficulty wins.

## Limits

- It applies only where MESS is in force. A node switches MESS off when it has fewer than five peers or its
  head goes stale: [When the node switches it off by itself](../operate/mess.md#when-the-node-switches-it-off-by-itself).
  `admin.messStatus()` shows whether it is on.
- The time comes from block timestamps, which miners set, not from your clock.
- The hashrate reading is an approximation. Difficulty adjusts block by block, so a chain's difficulty
  tracks its hashrate only over many blocks.
- Nodes without MESS follow the chain with the most total difficulty.
