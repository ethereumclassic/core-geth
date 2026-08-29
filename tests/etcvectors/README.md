# tests/etcvectors

ETC conformance vectors vendored from `fukuii-tests`, consumed by tests in package `tests`.

**Why not `tests/testdata-etc/`:** that path is a git submodule, so files cannot be added to
it from this repository. This directory is plain and versioned here.

## `etc_bomb_pause_and_removal.json`

**Upstream:** `networks/ethereumclassic/mainnet/difficulty/bomb_pause_and_removal.json`.

**Non-circular.** Its `_info.oracle` records three independent implementations read as source
rather than run, and a fourth check was added during the exchange described below.

It covers the difficulty bomb across every ETC ruleset: absent below the first period,
introduced at 200,000, growing, paused at 3,000,000 by ECIP-1010, and removed at 5,900,000 by
ECIP-1041.

### The file is LABEL-KEYED — bind a configuration per label

Each `ETC_*` section states expectations under **that upgrade's rule set**, with
`currentBlockNumber` as a free parameter. That is the published corpus's own convention:
`dfFrontier` and `dfByzantium` both span 100,000 to 4,900,000.

**Running these against a single `params.ClassicChainConfig` is wrong and scores 65 of 192.**
That resolves the rules from `CurrentBlockNumber` and discards the label, so every
Atlantis-and-later case at a low height is evaluated under pre-EIP-100 arithmetic. Binding a
configuration per label scores **192 of 192**.

**This repository reported the fixture as defective on exactly that mistake** and was wrong.
The disagreement is one adjustment unit because EIP-100 changes the interval divisor from 10
to 9, so a nine-second interval gives `max(1 - 9//9, -99) = 0` where Frontier and EIP-2 both
give `+1`. Both values are in the fixture, under the labels where each is correct.

Two halves of the mapping, and they are not symmetric:

- **The label's own adjustment rule is IN FORCE** — its transition is `0`, not its mainnet
  height.
- **Rules the label merely carries keep their REAL parameters** — ECIP-1010's pause window and
  ECIP-1041's removal height are parameters *of those rules*, not label selectors.
