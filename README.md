# Core-Geth

> An Ethereum Classic execution client, maintained as a public-good Go reference client.

| | |
|---|---|
| **Website** | [coregeth.com](https://coregeth.com), documentation at [docs.coregeth.com](https://docs.coregeth.com) |
| **Repository** | [`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth) |
| **Current release** | the v1.13 line, [latest release](https://github.com/ethereumclassic/core-geth/releases/latest) |
| **Maintainers** | the **@ethereumclassic/core-developers** team, in the [`ethereumclassic`](https://github.com/ethereumclassic) organization |

Core-Geth implements every Ethereum Classic fork, and the consensus rules are maintained here by the Ethereum Classic core developers: several admins and contributors in the community organization, not any single person or company.

## Supported networks

| Network | Chain ID | Consensus | Flag |
|---|---|---|---|
| Ethereum Classic | 61 | proof of work (Etchash) | `--classic`, `--mainnet`, or no flag |
| Mordor testnet | 63 | proof of work (Etchash) | `--mordor` |
| MintMe.com Coin | 24734 | proof of work | `--mintme` |
| Private chains | configurable | PoW / PoA | genesis config |

Started with no network flag, Core-Geth runs Ethereum Classic. **`--ethereum`, `--sepolia` and `--holesky` refuse to start:** this client implements Ethereum upgrades only through Cancun and is not a way to follow Ethereum mainnet.

Read the current fork schedule from [`params/`](params/) rather than restating it; it is data there, not code branches. **MESS** (ECBP-1100), the chain-selection defense against deep reorganizations, is covered on the [MESS page](https://docs.coregeth.com/operate/mess/).

## Quick start

Download `geth` for your platform from the [latest release](https://github.com/ethereumclassic/core-geth/releases/latest), verify it, and put it on your `PATH`. Then run an Ethereum Classic node:

```shell
geth --classic
```

[Installation](https://docs.coregeth.com/getting-started/installation/) has the per-platform download and verification steps, Docker images, and building from source. [Running a node](https://docs.coregeth.com/getting-started/run-a-node/) covers networks, running as a service, and the flags each kind of node needs; [Choose your role](https://docs.coregeth.com/guides/) routes operators, miners, mining pools, exchanges, RPC providers and developers to the right guide.

## Documentation

Published from this repository at [docs.coregeth.com](https://docs.coregeth.com/):

- [Installation](https://docs.coregeth.com/getting-started/installation/) and [running a node](https://docs.coregeth.com/getting-started/run-a-node/)
- [JSON-RPC API](https://docs.coregeth.com/JSON-RPC-API/)
- [Operating a node](https://docs.coregeth.com/operate/mess/): MESS, monitoring, maintenance, sync modes
- [Developers](https://docs.coregeth.com/developers/build-from-source/): build from source, add a network, testing
- [Upgrading](https://docs.coregeth.com/upgrading/) and [Releases](https://docs.coregeth.com/releases/)

## Building and testing

Building from source needs Go 1.26 or newer and a C compiler. With the repository cloned:

```shell
make geth                                 # build the node into ./build/bin/geth
make all                                  # build every executable
git submodule update --init --recursive   # fetch the fixture submodules the tests need
make test                                 # build, then run the test suite
make test-coregeth                        # the Ethereum Classic consensus suites
make lint                                 # the linter gate
```

[Build from source](https://docs.coregeth.com/developers/build-from-source/) has the toolchain details and the Docker build; [CONTRIBUTING.md](.github/CONTRIBUTING.md) has the coding guidelines and how a change reaches `main`.

## Security

Report a security issue privately, through a [GitHub advisory](https://github.com/ethereumclassic/core-geth/security/advisories) or to <security@ethereumclassic.net>, never a public issue. Mining pools, exchanges and service providers can use the same address to reach the maintainers. [SECURITY.md](SECURITY.md) has the policy, and the [security audits](https://docs.coregeth.com/audits/) establish which release carries which issue.

## Contributing

Fork the repository, make your change, and open a pull request against `main`; for anything substantial, open an issue first. [CONTRIBUTING.md](.github/CONTRIBUTING.md) has the guidelines. To discuss a change, the Ethereum Classic Core Developers run a [Discord](https://coregeth.com/discord). A fix that applies to code shared with [go-ethereum](https://github.com/ethereum/go-ethereum) is worth sending upstream as well.

## Project history

Core-Geth descends from go-ethereum by way of [multi-geth](https://github.com/multi-geth/multi-geth), and has been maintained in the [`ethereumclassic`](https://github.com/ethereumclassic) community organization since 2024. [Project history](https://docs.coregeth.com/about/project-history/) has the full account of the code and the maintaining team.

## Supporting the work

Core-Geth is public-goods infrastructure, funded by sponsorship. If your operation relies on Ethereum Classic, [Support this work](https://docs.coregeth.com/support/) sets out the ways to help.

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
