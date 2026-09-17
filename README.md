## Core-Geth: Ethereum Classic Execution Client

> A [go-ethereum](https://github.com/ethereum/go-ethereum) fork providing the production Ethereum Classic (ETC) execution client.

---

## ⚠️ Node operators: upgrade to v1.13.0, and rotate your node key

**Every release in the v1.12.x line carries at least one unpatched CVE, and the earliest
carry all six.** Every archive measured was built on a Go version that is no longer supported,
most recently [Go 1.21 and Go 1.22](docs/audits/2026-09-go-toolchain.md), whose support ended in
August 2024 and February 2025. Upgrade, then perform the one cleanup step below.

**1. Upgrade, and change where you track releases.** Releases are cut from
[`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth). A node tracking
the previous repository will not see v1.13.0.

**2. Rotate the P2P node key. This is required, not precautionary.** CVE-2026-26315 is an
oracle: an invalid-curve ephemeral key in the RLPx handshake previously reached ECDH and
failed only at MAC verification, which leaks bits of the node key across repeated
handshakes. Any key used by an unpatched node should be treated as exposed.

```bash
# with the node stopped
mv <datadir>/geth/nodekey <datadir>/geth/nodekey.old-rotated-$(date +%F)
```

The client generates a fresh key on the next start. **Your enode ID changes**, so update
every static-peer, trusted-peer or bootnode list that names this node, on your own
machines and with anyone peering with you. Capture the old enode ID before rotating if you
need it to find those references; it cannot be re-derived afterwards.

**3. Re-check your RPC exposure while the node is down.** `--http.addr` should be loopback
unless a trusted proxy sits in front of it, and `admin`, `debug` and `personal` do not
belong in `--http.api` on any reachable interface. `--http.corsdomain="*"` is not a safe
default.

**A resync is not required.** No finding in either audit corrupts chain data, and none is
known to have been exploited against this network. Verify your head matches another source
before concluding otherwise:

```bash
geth attach --exec 'eth.blockNumber' <datadir>/geth.ipc
```

Compare against a block explorer or another node you operate. Resync only if it diverges, and
if it does, that is a finding worth reporting.

**Nothing here requires touching your keystore.** These are network-layer and handshake
issues; account keys are not implicated by any of them. Rotate account keys only if your
`admin` or `personal` RPC was reachable from an untrusted network, which is its own
exposure rather than one of these CVEs.

Full detail: the [March 2026 audit](docs/audits/2026-03-security-audit.md), the
[August 2026 follow-up](docs/audits/2026-08-security-followup.md), the
[Go toolchain audit](docs/audits/2026-09-go-toolchain.md), and the
[v1.13.0 migration guide](docs/tutorials/v1.13.0-migration.md).

**The release artifacts were audited separately**, and that one matters before you
download rather than after: the published `v1.12.x` Linux and Arm binaries require a
newer glibc than the systems many operators run, and the macOS archive has contained an
Apple Silicon binary under an architecture-free name since June 2024. Each `v1.12.20`
to `v1.12.23` archive also carries 55 to 61 Go standard library advisories that `v1.13.0` does not. See
[release artifacts](docs/audits/2026-09-release-pipeline.md),
[Go toolchain](docs/audits/2026-09-go-toolchain.md) and
[dependency and toolchain modernization](docs/audits/2026-08-dependency-modernization.md).

---

Core-Geth is a production execution client for the Ethereum Classic network. It implements every ETC hard fork from Frontier through Spiral.

**Note:** Upstream go-ethereum has removed support for Ethereum Classic, so ETC consensus rules are maintained here rather than inherited.

## Supported Networks

| Network | Chain ID | Consensus | Flag |
|---------|----------|-----------|------|
| Ethereum Classic (ETC) | 61 | Proof of Work (ETChash) | `--classic`, `--mainnet` or no flag |
| Mordor Testnet | 63 | Proof of Work (ETChash) | `--mordor` |
| MintMe.com Coin | 24734 | Proof of Work | `--mintme` |
| Private chains | configurable | PoW / PoA | genesis config |

**With no network flag, Core-Geth runs Ethereum Classic mainnet**, and `--mainnet` is the
same as `--classic`. A data directory initialized with a private network's genesis keeps
running that network.

**Ethereum and its test networks are not supported.** This client implements Ethereum
upgrades only through Cancun, so it cannot follow Ethereum mainnet, Sepolia or Holesky.
`--ethereum`, `--sepolia` and `--holesky` are deprecated and refuse to start.

### ETC consensus history

Activation blocks, governing specifications and included EIPs, aligned with
[ECIP-1066](https://ecips.ethereumclassic.org/ECIPs/ecip-1066), which is the network description this table follows. The Mordor
column is client detail and is not part of that specification.

| Upgrade | ETC Block | Mordor Block | Date | Specs | Included EIPs |
|---------|----------:|-------------:|------|-------|---------------|
| Spiral <br><sub>Shanghai</sub> | 19,250,000 | 9,957,000 | 2024-02-04 | [ECIP-1109](https://ecips.ethereumclassic.org/ECIPs/ecip-1109) | Shanghai: [EIP-3651](https://eips.ethereum.org/EIPS/eip-3651), [EIP-3855](https://eips.ethereum.org/EIPS/eip-3855), [EIP-3860](https://eips.ethereum.org/EIPS/eip-3860), [EIP-6049](https://eips.ethereum.org/EIPS/eip-6049) |
| MESS Default: Off <br><sub>ECBP-1110</sub> | 19,250,000 | 10,400,000 | 2024-02-04 | [ECIP-1110](https://ecips.ethereumclassic.org/ECIPs/ecip-1110) | None |
| Mystique <br><sub>London</sub> | 14,525,000 | 5,520,000 | 2022-02-12 | [ECIP-1104](https://ecips.ethereumclassic.org/ECIPs/ecip-1104) | London: [EIP-3529](https://eips.ethereum.org/EIPS/eip-3529), [EIP-3541](https://eips.ethereum.org/EIPS/eip-3541) |
| Magneto <br><sub>Berlin</sub> | 13,189,133 | 3,985,893 | 2021-07-23 | [ECIP-1103](https://ecips.ethereumclassic.org/ECIPs/ecip-1103) | Berlin: [EIP-2565](https://eips.ethereum.org/EIPS/eip-2565), [EIP-2718](https://eips.ethereum.org/EIPS/eip-2718), [EIP-2929](https://eips.ethereum.org/EIPS/eip-2929), [EIP-2930](https://eips.ethereum.org/EIPS/eip-2930) |
| Thanos <br><sub>ECIP-1099</sub> | 11,700,000 | 2,520,000 | 2020-11-28 | [ECIP-1099](https://ecips.ethereumclassic.org/ECIPs/ecip-1099) | None |
| MESS Default: On <br><sub>ECBP-1100</sub> | 11,380,000 | 2,380,000 | 2020-10-09 | [ECIP-1100](https://ecips.ethereumclassic.org/ECIPs/ecip-1100) | None |
| Phoenix <br><sub>Istanbul</sub> | 10,500,839 | 999,983 | 2020-06-01 | [ECIP-1088](https://ecips.ethereumclassic.org/ECIPs/ecip-1088) | Istanbul: [EIP-152](https://eips.ethereum.org/EIPS/eip-152), [EIP-1108](https://eips.ethereum.org/EIPS/eip-1108), [EIP-1344](https://eips.ethereum.org/EIPS/eip-1344), [EIP-1884](https://eips.ethereum.org/EIPS/eip-1884), [EIP-2028](https://eips.ethereum.org/EIPS/eip-2028), [EIP-2200](https://eips.ethereum.org/EIPS/eip-2200) |
| Agharta <br><sub>Constantinople+Petersburg</sub> | 9,573,000 | 301,243 | 2020-01-11 | [ECIP-1056](https://ecips.ethereumclassic.org/ECIPs/ecip-1056) | Constantinople+Petersburg: [EIP-145](https://eips.ethereum.org/EIPS/eip-145), [EIP-1014](https://eips.ethereum.org/EIPS/eip-1014), [EIP-1052](https://eips.ethereum.org/EIPS/eip-1052) |
| Atlantis <br><sub>Byzantium</sub> | 8,772,000 | 0 | 2019-09-12 | [ECIP-1054](https://ecips.ethereumclassic.org/ECIPs/ecip-1054) | Spurious Dragon: [EIP-161](https://eips.ethereum.org/EIPS/eip-161), [EIP-170](https://eips.ethereum.org/EIPS/eip-170) <br> Byzantium: [EIP-100](https://eips.ethereum.org/EIPS/eip-100), [EIP-140](https://eips.ethereum.org/EIPS/eip-140), [EIP-196](https://eips.ethereum.org/EIPS/eip-196), [EIP-197](https://eips.ethereum.org/EIPS/eip-197), [EIP-198](https://eips.ethereum.org/EIPS/eip-198), [EIP-211](https://eips.ethereum.org/EIPS/eip-211), [EIP-214](https://eips.ethereum.org/EIPS/eip-214), [EIP-658](https://eips.ethereum.org/EIPS/eip-658) |
| Defuse Difficulty Bomb <br><sub>ECIP-1041</sub> | 5,900,000 | 0 | 2018-05-29 | [ECIP-1041](https://ecips.ethereumclassic.org/ECIPs/ecip-1041) | None |
| Gotham <br><sub>ECIP-1017</sub> | 5,000,000 | 0 | 2017-12-11 | [ECIP-1017](https://ecips.ethereumclassic.org/ECIPs/ecip-1017) <br> [ECIP-1039](https://ecips.ethereumclassic.org/ECIPs/ecip-1039) | None |
| Die Hard <br><sub>Spurious Dragon</sub> | 3,000,000 | 0 | 2017-01-13 | [ECIP-1010](https://ecips.ethereumclassic.org/ECIPs/ecip-1010) | Spurious Dragon: [EIP-155](https://eips.ethereum.org/EIPS/eip-155), [EIP-160](https://eips.ethereum.org/EIPS/eip-160) |
| Gas Reprice <br><sub>Tangerine Whistle</sub> | 2,500,000 | 0 | 2016-10-24 | [ECIP-1015](https://ecips.ethereumclassic.org/ECIPs/ecip-1015) | Tangerine Whistle: [EIP-150](https://eips.ethereum.org/EIPS/eip-150) |
| ~~DAO Fork~~ | 1,920,000 | — | 2016-07-20 | [HFM-779](https://eips.ethereum.org/EIPS/eip-779) | **Rejected on ETC** |
| Homestead | 1,150,000 | 0 | 2016-03-14 | [HFM-606](https://eips.ethereum.org/EIPS/eip-606) | Homestead: [EIP-2](https://eips.ethereum.org/EIPS/eip-2), [EIP-7](https://eips.ethereum.org/EIPS/eip-7), [EIP-8](https://eips.ethereum.org/EIPS/eip-8) |
| Frontier Thawing | 200,000 | 0 | 2015-09-07 | Genesis | None |
| Frontier | 1 | 0 | 2015-07-30 | Genesis | None |

**MESS** (Modified Exponential Subjective Scoring, ECBP-1100) is a chain-selection
defense against deep reorganizations. It is a client-side policy rather than a state
transition, which is why it appears as two rows: the block from which ECBP-1100 turns it
on by default, and the block from which ECBP-1110 recommends clients ship it off by
default.

**Core-Geth v1.13.x ships MESS on by default, and does not follow the ECBP-1110 row.**
That is the client maintainers' decision to make: both documents are Best Practices
rather than consensus rules, so each client chooses its own default. Centralized exchanges
requested MESS as a tool for periods of network hashrate instability, and shipping it on
lets the default lean toward security during a period when maintenance of the client is moving
between organizations, which is the condition that preceded the 2020 reorganizations MESS was
built for. Operators can always disable it with `--mess=false`.
[Why it is on by default](docs/operate/mess.md#why-it-is-on-by-default) gives the history, with
sources and the activation block. The [MESS confirmation calculator](docs/guides/mess-calculator.md) shows how much hashrate a
reorganization of a given age needs.

**Through Spiral (the head configuration this client implements), Ethereum Classic has
not adopted EIP-1559**, so transactions are the legacy and EIP-2930 access-list types.
That is the state of the fork schedule in `params/`, not a permanent property of the
network: adopting it is a protocol decision for Ethereum Classic to make, and a client
question after that.

### Wire protocol

**This client speaks `eth/68`, and that is the version it will serve until it is
retired.** The v1.13 series exists to close the security gap and carry Ethereum Classic
on a supported toolchain while the client is sunset, so it takes no new protocol
version.

`eth/69` (EIP-7642) removes Total Difficulty from the handshake, which this client's
proof-of-work chain selection reads, so adopting it here would mean reworking that path
in a client already scheduled for retirement. **That is a scoping decision about this
client, not a limitation of Ethereum Classic.** Later versions are left to the clients that
succeed it: [Fukuii](https://fukuii.org) is the preferred successor, and the migration guide's
[Fukuii section](docs/tutorials/v1.13.0-migration.md#migrating-to-fukuii) says when to move to
it.

## Build

```bash
make geth
```

The binary is written to `./build/bin/geth`.

## Test

```bash
# ETC-specific unit tests
go test ./params/... -run TestETC -v
go test ./core/... -run "TestGasLimit|TestECIP1017|TestETCForkCompliance" -v
go test ./consensus/ethash/... -run "TestDifficultyETC|TestDifficultyECIP" -v

# Live-network tests, build-tagged; requires a running Mordor or ETC node
go test -tags live ./tests/live_etc/ -v
```

The consensus fixtures live in three submodules: `tests/testdata`,
`tests/testdata-etc` and `tests/evm-benchmarks`. Run
`git submodule update --init --recursive` before any suite that reads them, or
the failures look like consensus errors rather than missing files.

## Run a node

```bash
./build/bin/geth --classic --datadir <path>     # Ethereum Classic mainnet
./build/bin/geth --mordor  --datadir <path>     # Mordor testnet
```

`--http.addr` defaults to loopback. Do not widen it, enable `--http.corsdomain`, or add
`admin` or `debug` to `--http.api` on a node reachable from an untrusted network.

## Mining

Core-Geth supports Ethash/ETChash proof-of-work mining:

```bash
./build/bin/geth --classic --mine --miner.etherbase <address>
```

For testing with fake PoW, which skips DAG generation:

```bash
./build/bin/geth --classic --mine --miner.etherbase <address> --fakepow
```

## Documentation

- Core-Geth documentation is published from this repository at
  [docs.coregeth.com](https://docs.coregeth.com/).
  + Getting Started: [Installation](https://docs.coregeth.com/getting-started/installation/) and [CLI](https://docs.coregeth.com/getting-started/run-cli/)
  + [JSON-RPC API](https://docs.coregeth.com/JSON-RPC-API/)
  + [Developers](https://docs.coregeth.com/developers/build-from-source/)
  + [Tutorials](https://docs.coregeth.com/tutorials/private-network/)
  + The source is `docs/` in this repository. It is built with `mkdocs` and
    deployed by [`.github/workflows/docs-deploy.yml`](.github/workflows/docs-deploy.yml).
- Further [ethereum/go-ethereum](https://github.com/ethereum/go-ethereum) documentation about can be found [here](https://geth.ethereum.org/docs/).
- Documentation about documentation lives [here](./docs/developers/documentation.md).

## Security

This release backports six CVE fixes from upstream go-ethereum and adds a GraphQL query
depth limit. Each is a separate commit; run `git log --grep CVE-` for the full detail,
including the upstream commit every fix derives from.

| Identifier | Component | Issue |
|-----|-----------|-------|
| CVE-2026-26313 | `eth/protocols` | Unbounded RLP message decoding |
| CVE-2026-26314 | `crypto/secp256k1` | `IsOnCurve` accepted coordinates outside the field |
| CVE-2026-26315 | `crypto/ecies` | Invalid-curve attack on the RLPx handshake |
| CVE-2026-22862 | `crypto/ecies` | Minimum ciphertext length check |
| CVE-2026-22868 | `eth` | Invalid KZG blob proofs |
| CVE-2025-24883 | `crypto` | `UnmarshalPubkey` accepted off-curve points |
| — | `graphql` | Unbounded query depth |

**Rotate P2P node keys after upgrading.** CVE-2026-26315 is an oracle against the node key
itself, so a key used by an unpatched node should be treated as exposed. Stop the node,
remove `<datadir>/geth/nodekey`, and restart; the client generates a new one. The node's
enode ID changes, so update any static-peer or trusted-peer list that names it.

**Track the [`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth/releases)
release line.** To subscribe to security updates and releases, or to report a security issue
privately, use GitHub or email <security@ethereumclassic.com>. [`SECURITY.md`](SECURITY.md) has
the disclosure policy and the details of each channel.

With the ETC Cooperative's dissolution, Ethereum Classic stakeholders such as mining pools,
exchanges and service providers should use <security@ethereumclassic.com> as their point of
contact. A person answers it: one of the core developers who maintain this repository and have
been with the network since its inception.

## Contributing

Contributions are welcome, and fixes of any size are useful. Fork the repository,
make your change, and open a pull request against `main`. For anything
substantial, open an issue first so the approach can be discussed.

[CONTRIBUTING.md](.github/CONTRIBUTING.md) has the coding guidelines, the build
and test commands, and the branch policy.

**If your fix applies to code shared with
[go-ethereum](https://github.com/ethereum/go-ethereum) rather than to Ethereum
Classic specifically, please send it upstream as well, or instead.** It reaches
every client built from that source, and it arrives here through the regular
merges.

Security vulnerabilities are not reported as issues. See [SECURITY.md](SECURITY.md).

## License

```
Copyright 2013-present The go-ethereum Authors
Copyright 2019-present The multi-geth Authors
Copyright 2020-present The core-geth Authors
```

The core-geth library (i.e. all code outside of the `cmd` directory) is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html),
also included in our repository in the `COPYING.LESSER` file.

The core-geth binaries (i.e. all code inside of the `cmd` directory) is licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also
included in our repository in the `COPYING` file.

Each source file carries a header stating which of the two applies to it. Where a
header and this summary disagree, the header governs. This program is distributed
in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the
implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

**Core-Geth is a modified version of
[go-ethereum](https://github.com/ethereum/go-ethereum), by way of
[multi-geth](https://github.com/multi-geth/multi-geth), and is a derivative work
of both.** Modification began in March 2018 as multi-geth and continues to the
present. The great majority of this source is upstream's, and each file's
copyright header records which project it belongs to.

Copyright is held by the individual contributors listed in `AUTHORS` and by the
authors of every commit in this repository's history. "The core-geth Authors"
refers to that body collectively, not to any company or organization.

## Project history

The code and the maintaining team have separate lineages, and neither is a
straight line.

**The code.** go-ethereum began in December 2013.
[multi-geth](https://github.com/multi-geth/multi-geth) forked it in March 2018 to
support several networks from a single client, holding chain rules as
configuration rather than as code branches. Core-Geth forked from multi-geth, a
separate client, in 2020 and kept that model, which is why a further network here is
a chain configuration rather than a new client.

**The team.** ETC Labs Core formed in December 2018, many of its developers having
previously been part of [ETCDEV](https://web.archive.org/web/20190330063218/https://www.etcdevteam.com/),
which supported the Classic Geth client. As Classic Geth was retired the team supported
multi-geth, and then Core-Geth from 2020, publishing as
[ETC Labs](https://web.archive.org/web/20200425081322/https://etclabs.org/) and
[ETC Core](https://web.archive.org/web/20200426174445/https://etccore.io/). ETC Labs left the
Ethereum Classic ecosystem in 2021, and from January 2022 the work was funded by the
[ETC Cooperative](https://web.archive.org/web/20250205192722/https://etccooperative.org/),
which entered maintenance mode at the end of 2024 and stopped developing Core-Geth in 2025.
The [Ethereum Classic DAO](https://ethereumclassicdao.org), a Wyoming DAO LLC launched in May
2025, succeeds it. [The ETC Cooperative transition](docs/etc-cooperative-transition.md) lists
the services moving to the `ethereumclassic` organization as the Cooperative winds down.

**This repository.** Maintenance moved to
[`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth), the Ethereum
Classic community repository, created on 2024-12-21 from the preceding repository at commit
`7ef3ecd7a` (2024-12-16). That commit is preserved on the `archive-etclabscore-2024-12` branch
and tagged
[`archive/etclabscore-2024-12`](https://github.com/ethereumclassic/core-geth/releases/tag/archive%2Fetclabscore-2024-12).
It is maintained by long-time Ethereum Classic core developers, with the network since its
inception, as a public good for the network.

**2026.** After the move the client went unfunded and unmaintained for 21 months, as the
[March 2026 security audit](docs/audits/2026-03-security-audit.md) counts, until
[White B0x](https://whiteb0x.com) took up its development, security work and disclosures for the
Ethereum Classic DAO in February 2026. That work was reported privately to the ETC Cooperative
in February and March 2026, a live attack on Ethereum Classic bootnodes in March 2026 confirmed
the exposure, and the fixes were cut into the emergency v1.12.21 and v1.12.22 releases in the
previous repository while v1.13 was tested from February to September 2026. The same fixes were
submitted here as 27 pull requests, [#10](https://github.com/ethereumclassic/core-geth/pull/10)
through [#36](https://github.com/ethereumclassic/core-geth/pull/36). v1.13.0, published on
14 September 2026, is the first release from this repository.

**The full account** of both lineages, the organizations behind them and the records for each
step is in [Project history](docs/about/project-history.md). Claims made about the release after
it shipped are answered from the record in
[v1.13.0: the record behind the release](docs/release-reports/v1.13.0-record.md).

**What comes next.** v1.13 is the last Core-Geth release line, maintained through the
transition. Core-Geth is scheduled to sunset gradually in favor of two efforts. Ethereum Classic
network extensions, overlays or plugins for Ethereum clients are in development, pending
release, and none is offered here as an option. [Fukuii](https://fukuii.org), a client native
to Ethereum Classic, removes upstream, third-party client dependencies from Ethereum Classic's
core software; moving to it starts with its first release, and the migration guide's
[Fukuii section](docs/tutorials/v1.13.0-migration.md#migrating-to-fukuii) says when. Fukuii is
inspired by the work of earlier Ethereum Classic development teams: ETCDEV's
[Classic Geth](https://github.com/ethereumproject/go-ethereum) and Orbita vision, and
[IOHK](https://iohk.io/)'s [Mantis](https://web.archive.org/web/20211026113958/https://mantisclient.io/).

**Supporting the work.** Maintaining Core-Geth since it moved to the community repository has been
unfunded public-goods work. Mining pools, centralized exchanges, issuers of Ethereum Classic
financial products, Etchash mining hardware manufacturers and large holders all depend on an
Ethereum Classic client that follows the network correctly and receives security fixes promptly. If
your operation relies on Ethereum Classic, please help fund that work with a donation or a
retroactive grant. Contact <donations@ethereumclassic.com>, or donate on any EVM chain.

**Retroactive Core-Geth Development Donation Address:**

```
0x86FE8d331A4B984B57d3e92C6F4cb9C881eC9B04
```
