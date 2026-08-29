# params/mutations/testdata

`etc_era_emission_schedule.json` — ECIP-1017 monetary-policy vectors, consumed by
`rewards_etc_vectors_test.go`.

**Upstream:** `fukuii-tests`,
`networks/ethereumclassic/mainnet/blocks/era_emission_schedule.json`. Vendored so the test
is self-contained; re-sync upstream rather than editing this copy.

**Non-circular.** Computed independently from ECIP-1017 and only then cross-checked against
a client.

These assert Ethereum Classic's issuance schedule: the 5,000,000-block era length, the 4/5
reduction at each boundary, the 1/32 includer bonus per ommer, and the distance-dependent
ommer reward that applies in era 0 only. The era boundary is the subtle part — block
5,000,000 is still era 0 and 5,000,001 is era 1 — and getting it wrong shifts every
subsequent reward by one era while remaining entirely plausible.
