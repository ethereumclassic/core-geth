# core/forkid/testdata

`etc_fork_identifiers.json` — EIP-2124 fork-identifier vectors for Ethereum Classic
mainnet, consumed by `forkid_etc_vectors_test.go`.

**Upstream:** `fukuii-tests`,
`networks/ethereumclassic/mainnet/forkid/mainnet_fork_identifiers.json`. Vendored so the
test is self-contained; re-sync upstream rather than editing this copy.

**Non-circular.** The file's `_info.oracle` states the values are computed independently
from EIP-2124 — a CRC-32 accumulated over the genesis hash and each passed fork block as a
big-endian 64-bit value — and only then cross-checked against a client. That ordering is
what makes running them here a conformance test rather than a change-detector.

Fork identifiers gate peering: two nodes that disagree will not connect. A change here is
therefore a network-visible change, which is why it is worth asserting against an external
oracle rather than against this client's own output.
