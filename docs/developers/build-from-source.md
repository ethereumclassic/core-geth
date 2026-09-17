---
title: Build from Source
description: "Build Core-Geth from source: the Go version required, the make targets, and the test suites that prove an Ethereum Classic build."
---

## Hardware Requirements

The disk, memory and time a node needs are measured on
[Hardware requirements](../getting-started/hardware-requirements.md).

## Dependencies

- **Go 1.26 or later.** <https://go.dev/doc/install>. The module declares
  `go 1.26.0`, so an older toolchain refuses to build it rather than producing a
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
$ make geth
$ ./build/bin/geth --help
```

`make geth` builds the node alone. `make all` builds every executable in
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
$ docker build -t core-geth:local .
```

Or with all tools:

```shell
$ docker build -t core-geth-alltools:local -f Dockerfile.alltools .
```

`core-geth:local` is just the tag you are giving the image; name it whatever you
like. The image's entry point is the `geth` binary, so flags go straight to the
node: `docker run core-geth:local --classic`. Running a published image instead
is covered under [Installation](../getting-started/installation.md#docker).
