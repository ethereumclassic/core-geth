## Core-Geth: Ethereum Classic Execution Client

> A [go-ethereum](https://github.com/ethereum/go-ethereum) fork providing the production Ethereum Classic (ETC) execution client.

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

**MESS** (Modified Exponential Subjective Scoring, ECBP-1100) is a chain-selection defense
against deep reorganizations. Being a client-side policy rather than a state transition, it
appears twice above: ECBP-1100 turns it on by default, ECBP-1110 recommends clients ship it
off. Both are Best Practices rather than consensus rules, so each client sets its own
default. `--mess` turns it on, `--mess=false` turns it off.
[MESS](docs/operate/mess.md) has this client's default and the trade-off; the
[confirmation calculator](docs/guides/mess-calculator.md) shows the hashrate a
reorganization of a given age needs.

**Through Spiral (the head configuration this client implements), Ethereum Classic has
not adopted EIP-1559**, so transactions are the legacy and EIP-2930 access-list types.
That is the state of the fork schedule in `params/`, not a permanent property of the
network: adopting it is a protocol decision for Ethereum Classic to make, and a client
question after that.

### Wire protocol

**This client speaks `eth/68`.** `eth/69` (EIP-7642) removes Total Difficulty from the
handshake, which this client's proof-of-work chain selection reads, so adopting it would
mean reworking that path. **That is a scoping decision about this client, not a limitation
of Ethereum Classic.**

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

- Current release notes and upgrade steps:
  [`docs/tutorials/v1.13.0-migration.md`](docs/tutorials/v1.13.0-migration.md).
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

**To report a security issue privately**, use a GitHub private advisory or email
<security@ethereumclassic.net>. [`SECURITY.md`](SECURITY.md) has the disclosure policy and
the details of each channel. Security issues are never reported as public issues.

Ethereum Classic stakeholders such as mining pools, exchanges and service providers can use
the same address to reach the core developers who maintain this repository.

**Each security fix is its own commit, and the repository is the record.** `git log --grep
CVE-` lists them, including the upstream commit every fix derives from. The
[security audits](docs/audits/) establish which release carries which issue, and each
release's notes state what that release closed.

**Releases are cut from
[`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth/releases)**, which
is where to watch for security updates. Upgrading from an earlier release can carry required
operator steps beyond replacing the binary; the release notes and the
[migration guides](https://docs.coregeth.com/tutorials/) state them per release.

## Contributing

Contributions are welcome, and fixes of any size are useful. Fork the repository,
make your change, and open a pull request against `main`. For anything
substantial, open an issue first so the approach can be discussed.

[CONTRIBUTING.md](.github/CONTRIBUTING.md) has the coding guidelines, the build
and test commands, and the branch policy. To discuss a change before writing it, the Ethereum
Classic Core Developers run a [Discord](https://ethereumclassic.com/discord). Security issues go through the private routes in
[Security](#security) instead, never a public channel.

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

**The team.** Development has changed hands four times:

| Development lead | Client | When |
| --- | --- | --- |
| [Sorpaas](https://github.com/sorpaas) | multi-geth, from its start | 2018 to 2020 |
| [ETC Labs](https://web.archive.org/web/20200425081322/https://etclabs.org/), publishing also as [ETC Core](https://web.archive.org/web/20200426174445/https://etccore.io/) | Core-Geth, through the fork and rebrand from multi-geth | 2020 to 2021 |
| [ETC Cooperative](https://web.archive.org/web/20250205192722/https://etccooperative.org/) | Core-Geth | 2021 to 2024 |
| Community maintained, in the [`ethereumclassic`](https://github.com/ethereumclassic) organization | Core-Geth | 2024 onward |

[The ETC Cooperative transition](docs/etc-cooperative-transition.md) lists the services moving
to that organization as the Cooperative winds down.

**This repository.** Maintenance moved to
[`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth), the Ethereum
Classic community repository, created on 2024-12-21 from the preceding repository at commit
`7ef3ecd7a` (2024-12-16). That commit is preserved on the `archive-etclabscore-2024-12` branch
and tagged
[`archive/etclabscore-2024-12`](https://github.com/ethereumclassic/core-geth/releases/tag/archive%2Fetclabscore-2024-12).

Development then stayed dormant across every Core-Geth repository for 21 months.
[White B0x](https://whiteb0x.com) has led it since February 2026, cutting the v1.13.0
modernization release and clearing years of accumulated technical debt inherited from the
previous repositories. That leaves a modernized client for community maintenance in the
[`ethereumclassic`](https://github.com/ethereumclassic) organization going forward.

[Project history](docs/about/project-history.md) has the full account of both lineages, the
organizations behind them and the records for each step. Claims made about the client are
answered from the record in
[v1.13.0: the record behind the release](docs/release-reports/v1.13.0-record.md).

**Supporting the work.** Maintaining Core-Geth since it moved to the community repository has been
unfunded public-goods work. Mining pools, centralized exchanges, issuers of Ethereum Classic
financial products, Etchash mining hardware manufacturers and large holders all depend on an
Ethereum Classic client that follows the network correctly and receives security fixes promptly. If
your operation relies on Ethereum Classic, please help fund that work.
[Support this work](https://docs.coregeth.com/support/) carries the routes: sponsorship on GitHub,
an invoiced maintenance agreement for organizations, or a direct transfer to an address provided on
request.
