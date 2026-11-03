---
description: "Build and test Core-Geth: linting with golangci-lint, the Go test suites, and the Ethereum Classic consensus and chainspec equivalence tests."
---

# Testing

## Lint

- `make lint` runs the module-standard linter config. Under the hood, this uses `golangci-lint` and its configuration at `/.golangci.yml`.

## Run tests

The `Makefile` at the root of the project includes many commands that can be useful for testing.

- `make test` runs all package tests, most of them shared with go-ethereum.
- `make test-coregeth` runs the tests specific to Core-Geth.

`go test` and `go build` also work, as in any Go project.


## Test generation

Core-Geth is capable of generating some sets of tests used in the `tests` package, which are originally (and still largely)
driven by the [ethereum/tests](https://github.com/ethereum/tests) suite.

- `make tests-generate` runs test(s) generation for the `state` and `difficulty` subsections of this suite, extending the ethereum/tests version
  of the controls to include configurations for Ethereum Classic chain configurations at various points in Ethereum Classic hardfork history.

## Flaky (spuriously erroring) tests

Especially when run in CI environments, some tests can be expected to fail more-or-less randomly.


