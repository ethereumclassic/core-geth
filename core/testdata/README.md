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

**The window layer is no longer in this file.** `fukuii-tests` split activation, deactivation
and `windowVectors` into `mess_window.json`, so what remains here is the conformance half
only: the curve, the decision, and the derivation of the decision's inputs from two competing
segments.

That split resolved a real framing problem rather than a formatting one. ECBP-1100's window is
*client configuration* — a client shipping MESS permanently on does not fork from one shipping
it off — while the arithmetic is not. One file was pinning both, so running it against any
client mixed a conformance claim with a policy claim. This repository ships MESS enabled and
therefore diverges from the window layer deliberately; it never diverged on the arithmetic.

**Do not re-add the window layer here.** If this client's MESS policy needs asserting, that is
`params/config_etc_test.go`'s job against this repository's own configuration, not a
conformance fixture's.
