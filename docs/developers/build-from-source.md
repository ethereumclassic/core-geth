---
title: Build from Source
---

## Hardware Requirements

Minimum:

* CPU with 2+ cores
* 4GB RAM
* 500GB free storage space to sync the Mainnet
* 8 MBit/sec download Internet service

Recommended:

* Fast CPU with 4+ cores
* 16GB+ RAM
* High Performance SSD with at least 500GB free space
* 25+ MBit/sec download Internet service

## Dependencies

- **Go 1.26 or later.** <https://go.dev/doc/install> — the module declares
  `go 1.26`, so an older toolchain refuses to build it rather than producing a
  broken binary.
- **A C compiler.** Parts of the client are cgo, so a working toolchain is
  required. On Debian or Ubuntu:

```shell
$ sudo apt-get install -y build-essential
```

!!! tip "Building without installing Go yourself"
    `go run build/ci.go install -dlgo` downloads the exact Go toolchain this
    project pins in `build/checksums.txt` and builds with that, rather than with
    whatever Go is on your `PATH`. This is what the release builds use, so it is
    the closest you can get to reproducing one locally.

## Source

Once the dependencies have been installed, clone and build:

```shell
$ git clone https://github.com/ethereumclassic/core-geth.git
$ cd core-geth
$ make core-geth
$ ./build/bin/core-geth --help
```

`make core-geth` builds the node alone. `make all` builds every executable in
`cmd/`, which is what the `alltools` release archives contain. Run `make help`
to list the available targets.

!!! note "Running the test suites needs the fixture submodules"
    The consensus tests read from three git submodules that a plain clone does
    not populate. Without them the suites fail in ways that look like consensus
    errors rather than missing files:

    ```shell
    $ git submodule update --init --recursive
    $ make test-coregeth
    ```

## Build docker image

You can build a local docker image directly from the source:

```shell
$ git clone https://github.com/ethereumclassic/core-geth.git
$ cd core-geth
$ docker build -t=core-geth .
```

Or with all tools:

```shell
$ docker build -t core-geth-alltools -f Dockerfile.alltools .
```
