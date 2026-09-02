# Contributing

Thank you for considering a contribution. Fixes of any size are welcome.

Fork the repository, make your change, and open a pull request. For anything
substantial, open an issue first so the approach can be discussed before you
spend time on it.

## Send upstream what belongs upstream

Core-Geth is a downstream of
[go-ethereum](https://github.com/ethereum/go-ethereum). If your fix applies to
code shared with upstream rather than to Ethereum Classic specifically, please
send it there as well, or instead — a change that lands upstream reaches every
client built from it, and it arrives here through the regular merges.

Ethereum Classic chain configuration, the ECIP implementations, and this
client's own multi-network model belong here.

## Which branch

Pull requests are opened against `master`, this repository's default branch.

A `main` branch exists and carries in-progress modernization work. It is not the
target for contributions yet. When that work is ready for release, `main` becomes
the default branch and this document will say so.

## Coding guidelines

- Code is `gofmt`-formatted and documented per the Go
  [commentary](https://go.dev/doc/effective_go#commentary) conventions.
- `make lint` is the gate. It runs the linter set enabled in `.golangci.yml`,
  which is narrower than the golangci-lint defaults — a change that passes
  `golangci-lint run` with default settings has not been checked against this
  project's configuration.
- Commit messages are prefixed with the package or area they modify, for example
  `eth, rpc: make trace configs optional`.
- Comments carry reasoning rather than description.

## Building and testing

```bash
git submodule update --init --recursive   # the test suites need the fixture submodules
make geth                                 # build cmd/geth into ./build/bin/geth
make all                                  # build every executable
make test                                 # make all, then build/ci.go test
make lint                                 # the linter gate
make test-coregeth                        # the Ethereum Classic specific suites
```

`make help` lists the annotated targets. `make test-coregeth` is what proves this
client's chain configuration work and is worth running for any change under
`params/`, `core/` or `consensus/`.

Note that `make tests-generate` **overwrites** generated fixtures under
`tests/testdata-etc/`. Read the target before running it.

## Reporting a vulnerability

Do not open a public issue. See [SECURITY.md](../SECURITY.md).
