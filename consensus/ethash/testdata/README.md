# consensus/ethash/testdata

`etchash_epoch_schedule.json` — ECIP-1099 epoch-schedule and seed-continuity vectors,
consumed by `etchash_epoch_vectors_test.go`.

**Upstream:** `fukuii-tests`,
`networks/ethereumclassic/mainnet/pow/etchash_epoch_schedule.json`. Vendored so the test is
self-contained; re-sync upstream rather than editing this copy.

**Non-circular.** Derived from ECIP-1099 and the ethash specification, not read back from a
client.

The seed-continuity half is the valuable part. It pins the identity
`seed(postForkEpoch e, 60000) == seed(legacyEpoch 2e, 30000)` and names the natural wrong
implementation — dividing the epoch start block by the epoch length *in force* rather than
by 30,000. That mistake does not raise an error: it silently reuses a seed from roughly six
million blocks earlier and generates a real dataset for the wrong epoch, so a miner produces
work nobody accepts while every internal check passes.
