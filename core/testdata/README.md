# core/testdata

## `ecbp1100_mess_vectors.json`

Conformance vectors for ECBP-1100 (MESS), consumed by
`blockchain_af_vectors_test.go`.

**Upstream:** `fukuii-tests`, at
`networks/ethereumclassic/mainnet/chainselection/mess_artificial_finality.json`.
Vendored rather than referenced so the test is self-contained in a clone and in CI.
Re-sync from upstream rather than editing this copy.

**The vectors are not derived from this client.** The file's own `_info.oracle` states
they are computed from ECIP-1100's pseudocode and that no value is read back from any
client, and four implementations were read while producing them. That is what makes
running them here a conformance test rather than a change-detector.

**Validated against the specification before adoption**, 2026-08-28, by transcribing
ECIP-1100's `get_curve_function_numerator` and re-deriving every value: 11/11
`curveVectors`, 11/11 `decisionVectors` verdicts, and 6/6 `subchainVectors` including
their `derivedIntermediates`.

**`windowVectors` is deliberately NOT consumed.** It encodes `deactivationBlock`
19,250,000, and this client ships ECBP-1100 with no deactivation. That layer asserts the
behavior this release deliberately changed; the curve, decision and derivation layers are
unaffected, because what changed is the window and not the arithmetic.
