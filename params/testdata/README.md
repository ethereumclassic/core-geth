# params/testdata

`etc_required_block_headers.json` — required-block-hash vectors, consumed by
`config_etc_required_hashes_test.go`.

**Upstream:** `fukuii-tests`,
`networks/ethereumclassic/mainnet/blocks/required_block_headers.json`. Vendored so the test
is self-contained; re-sync upstream rather than editing this copy.

**Non-circular.** Its `_info.oracle` records that every hash was read from an archive node
rather than from a client's configuration.

The vectors matter because the hash required at block 1,920,000 is what makes this chain
Ethereum Classic. The rejected alternatives include the hash of the forked chain's block at
the same height, so the test fails if the configuration ever names the wrong side of the DAO
fork — a change that would otherwise show up only as a node syncing the wrong chain.
