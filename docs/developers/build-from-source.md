---
title: Build from Source
description: "Build Core-Geth from source: the Go version required, the make targets, and the test suites that prove an Ethereum Classic build."
---

## Hardware Requirements

The disk, memory and time a node needs are measured on
[Hardware requirements](../getting-started/hardware-requirements.md).

## Dependencies

Building Core-Geth needs three things:

- **Go 1.26 or later.** The module declares `go 1.26.0`, so an older toolchain
  refuses to build it rather than producing a broken binary. Check yours with
  `go version`.
- **A C compiler.** Parts of the client are cgo, so a working C toolchain
  (`gcc` or `clang`) is required.
- **git**, to clone the repository.

Install all three for your platform:

=== "Debian / Ubuntu"

    ```shell
    sudo apt-get update && sudo apt-get install -y build-essential git
    ```

    `build-essential` is the C compiler. The `apt` Go package often lags 1.26, so
    install Go from [go.dev/doc/install](https://go.dev/doc/install) or with snap:

    ```shell
    sudo snap install go --classic
    go version
    ```

=== "Fedora / RHEL"

    ```shell
    sudo dnf install -y gcc git
    ```

    Then install Go 1.26 or later from
    [go.dev/doc/install](https://go.dev/doc/install). Fedora's own `golang` package
    works when it is 1.26 or newer, which `go version` confirms:

    ```shell
    sudo dnf install -y golang
    go version
    ```

=== "Arch"

    ```shell
    sudo pacman -S --needed base-devel git go
    go version
    ```

    `base-devel` provides the C compiler, and Arch packages a current Go; confirm
    it is 1.26 or newer with `go version`.

=== "macOS"

    The Xcode command line tools provide `clang` and `git`:

    ```shell
    xcode-select --install
    ```

    Then install Go 1.26 or later with [Homebrew](https://brew.sh), or from the
    macOS installer at [go.dev/doc/install](https://go.dev/doc/install):

    ```shell
    brew install go
    go version
    ```

=== "Windows"

    With [Chocolatey](https://chocolatey.org/install), in an elevated PowerShell:

    ```powershell
    choco install golang mingw git
    go version
    ```

    `mingw` is the MinGW-w64 C compiler that cgo needs on Windows. Without
    Chocolatey, install Go from the MSI at
    [go.dev/doc/install](https://go.dev/doc/install), install
    [MinGW-w64](https://www.mingw-w64.org/) separately, and make sure its `gcc.exe`
    is on your `PATH`.

!!! tip "Building without installing Go yourself"
    `go run build/ci.go install -dlgo` downloads the exact Go toolchain this
    project pins in `build/checksums.txt` and builds with that, rather than with
    whatever Go is on your `PATH`. This is what the release builds use, so it is
    the closest you can get to reproducing one locally.

## Source

Once the dependencies have been installed, clone and build.

**To run a node, build the latest release:**

=== "Linux and macOS"

    ```shell
    VERSION=$(curl -sLo /dev/null -w '%{url_effective}' https://github.com/ethereumclassic/core-geth/releases/latest | sed 's|.*/||')
    git clone --branch $VERSION --single-branch https://github.com/ethereumclassic/core-geth.git core-geth-$VERSION &&
      cd core-geth-$VERSION &&
      make geth &&
      ./build/bin/geth version
    ```

=== "Windows"

    In PowerShell. `make geth` runs this same `go run` command on other systems:

    ```powershell
    $VERSION = (Invoke-RestMethod "https://api.github.com/repos/ethereumclassic/core-geth/releases/latest").tag_name
    git clone --branch $VERSION --single-branch https://github.com/ethereumclassic/core-geth.git "core-geth-$VERSION"
    cd "core-geth-$VERSION"
    go run build/ci.go install ./cmd/geth
    .\build\bin\geth.exe version
    ```

**To work on the client or try what is in development, build `main`.** It keeps moving, and it is
not a release:

=== "Linux and macOS"

    ```shell
    git clone https://github.com/ethereumclassic/core-geth.git &&
      cd core-geth &&
      make geth &&
      ./build/bin/geth version
    ```

=== "Windows"

    ```powershell
    git clone https://github.com/ethereumclassic/core-geth.git
    cd core-geth
    go run build/ci.go install ./cmd/geth
    .\build\bin\geth.exe version
    ```

`make geth` builds the node alone. `make all` builds every executable in
`cmd/`, which is what the `alltools` release archives contain; on Windows,
`go run build/ci.go install` does the same. Run `make help` to list the available
targets.

!!! note "Running the test suites needs the fixture submodules"
    The consensus tests read from three git submodules that a plain clone does
    not populate. Without them the suites fail in ways that look like consensus
    errors rather than missing files:

    ```shell
    git submodule update --init --recursive
    make test-coregeth
    ```

## Build docker image

You can build a local docker image of the latest release directly from the source:

```shell
VERSION=$(curl -sLo /dev/null -w '%{url_effective}' https://github.com/ethereumclassic/core-geth/releases/latest | sed 's|.*/||')
git clone --branch $VERSION --single-branch https://github.com/ethereumclassic/core-geth.git core-geth-$VERSION &&
  cd core-geth-$VERSION &&
  docker build -t core-geth:local .
```

Or with all tools:

```shell
docker build -t core-geth-alltools:local -f Dockerfile.alltools .
```

`core-geth:local` is just the tag you are giving the image; name it whatever you
like. The image's entry point is the `geth` binary, so flags go straight to the
node: `docker run core-geth:local --classic`. Running a published image instead
is covered under [Installation](../getting-started/installation.md#docker).
