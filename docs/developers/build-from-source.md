---
title: Build from Source
description: "Build Core-Geth from source: the Go version required, the make targets, and the test suites that prove an Ethereum Classic build."
---



## Dependencies

Building Core-Geth needs three things:

- **Go 1.26 or later.** Check yours with `go version`.
- **A C compiler**, `gcc` or `clang`.
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

!!! tip "Building with the Go toolchain the project pins"
    `go run build/ci.go install -dlgo` downloads the exact Go toolchain this project pins in `build/checksums.txt` and
    builds with that. It still needs an installed Go to run the build script.

## Source

Once the dependencies have been installed, clone `main` and build:

=== "Linux and macOS"

    ```shell
    git clone --branch main --single-branch https://github.com/ethereumclassic/core-geth.git &&
      cd core-geth &&
      make geth &&
      ./build/bin/geth version
    ```

=== "Windows"

    In PowerShell. `make geth` runs this same `go run` command on other systems:

    ```powershell
    git clone --branch main --single-branch https://github.com/ethereumclassic/core-geth.git
    cd core-geth
    go run build/ci.go install ./cmd/geth
    .\build\bin\geth.exe version
    ```

If a `core-geth` folder is already there from an earlier clone, run `git -C core-geth pull` in place of the clone, here
and in [Build docker image](#build-docker-image).

**`main` is where changes land before a release.** Every change reaches it through a pull request, which needs an
approving review from someone other than its last pusher, and five passing checks, before it merges. So a build of
`main` carries the fixes and dependency updates merged since the latest release.

**A build of `main` is not a release, and `main` keeps moving:** two builds a day apart can be different code.
`geth version` prints the version `main` is working toward, marked `-unstable`, and the commit it was built from, for
example:

```text
Core-Geth 1.13.1-unstable
Git Commit: 50f00770daa3f9e31e4d6bd63e17f49a1d78a35f
Git Commit Date: 20261012
```

A release's checksum and build attestation verify a release file, so they do not apply to a build. Note the commit
instead, from that `Git Commit` line or from `git rev-parse --short HEAD` in the clone: it identifies the code you run.

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

To build a local image of `main`:

```shell
git clone --branch main --single-branch https://github.com/ethereumclassic/core-geth.git &&
  cd core-geth &&
  docker build -t core-geth:local .
```

Or with all tools:

```shell
docker build -t core-geth-alltools:local -f Dockerfile.alltools .
```

Flags go straight to the node: `docker run core-geth:local --classic`.
[Installation](../getting-started/installation.md#docker) covers running a published image.
