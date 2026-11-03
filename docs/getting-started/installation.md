---
title: Installation
description: "Download, verify and install a Core-Geth release on Linux, macOS, Windows or Docker, and check an archive against its published checksum."
---

!!! danger "Security advisory: the v1.12.x line"
        Every `v1.12.x` release is affected by [the flaws the audits found](../index.md#what-was-audited). `v1.13.0` and
    every release since fix them. Archives and images under the previous `etclabscore` namespace are not built from this
    source and do not carry those fixes. [Upgrading](../upgrading.md#leaving-the-v112x-line) covers the move.

## Choose your route

Pick the row that matches your machine, and follow only that section.

| Your machine | Route |
| --- | --- |
| Linux, x86_64: most servers and desktops | [Linux, x86_64](#linux-x86_64) |
| Linux, 64-bit ARM: Raspberry Pi on a 64-bit OS, AWS Graviton, Ampere | [Linux, ARM](#linux-arm) |
| Linux, 32-bit ARM: older Raspberry Pi and embedded boards | [Linux, ARM](#linux-arm) |
| macOS, Apple Silicon | [macOS](#macos) |
| macOS, Intel | [macOS](#macos) |
| Windows, x86_64 | [Windows](#windows) |
| Anything, in a container | [Docker](#docker) |
| Anything, compiled yourself | [Build from source](#build-from-source) |

On Linux and macOS, `uname -sm` prints both.

**The executable is named `geth`.** `core-geth` is the project, the repository, the container image and the archive
names.

The Linux and macOS routes end with `geth` on your `PATH`, and the Windows route with `geth.exe` in `C:\core-geth`. Each
archive on the [releases page](https://github.com/ethereumclassic/core-geth/releases) is a `.zip` holding the binary and
`COPYING`, the license, with a `.sha256` file beside it.

**Each route reads the latest release's tag into `VERSION`.** To install a different release, set `VERSION` to its tag
instead of running that first line.

## Linux, x86_64

```shell
VERSION=$(curl -sLo /dev/null -w '%{url_effective}' https://github.com/ethereumclassic/core-geth/releases/latest | sed 's|.*/||')
BASE=https://github.com/ethereumclassic/core-geth/releases/download/$VERSION
curl -LO $BASE/core-geth-linux-$VERSION.zip
curl -LO $BASE/core-geth-linux-$VERSION.zip.sha256
sha256sum -c core-geth-linux-$VERSION.zip.sha256 &&
  unzip -o core-geth-linux-$VERSION.zip &&
  sudo install -m 0755 geth /usr/local/bin/geth &&
  geth version
```

`sha256sum -c` must print `OK`. If it prints `FAILED`, the lines after it do not run:
delete the file and download it again.

**This archive needs glibc 2.17 or newer**: RHEL 7, CentOS 7 and anything more
recent. `ldd --version` prints what you have. On an older distribution
the binary will not start; [build from source](#build-from-source) there instead.

Before running it on a node that holds value, also
[verify where the archive came from](#verify-where-the-archive-came-from).

## Linux, ARM

Four ARM archives are published. `uname -m` tells you which one you need:

| `uname -m` prints | Download | Needs glibc |
| --- | --- | --- |
| `aarch64` or `arm64` | `core-geth-arm64-<tag>.zip` | 2.17 or newer |
| `armv7l` | `core-geth-arm7-<tag>.zip` | 2.28 or newer |
| `armv6l` | `core-geth-arm6-<tag>.zip` | 2.28 or newer |
| `armv5tel` and other `armv5` | `core-geth-arm5-<tag>.zip` | 2.28 or newer |

Follow what `uname -m` prints, not the board: Raspberry Pi OS reports `aarch64` in its 64-bit edition and `armv7l` in
its 32-bit edition on the same hardware.

Then, substituting your architecture for `arm64` below:

```shell
VERSION=$(curl -sLo /dev/null -w '%{url_effective}' https://github.com/ethereumclassic/core-geth/releases/latest | sed 's|.*/||')
BASE=https://github.com/ethereumclassic/core-geth/releases/download/$VERSION
curl -LO $BASE/core-geth-arm64-$VERSION.zip
curl -LO $BASE/core-geth-arm64-$VERSION.zip.sha256
sha256sum -c core-geth-arm64-$VERSION.zip.sha256 &&
  unzip -o core-geth-arm64-$VERSION.zip &&
  sudo install -m 0755 geth /usr/local/bin/geth &&
  geth version
```

`core-geth-arm-<tag>.zip` is the same file as `arm5`, kept for existing scripts. Use `arm5`.

Before running it on a node that holds value, also
[verify where the archive came from](#verify-where-the-archive-came-from).

## macOS

Two archives, chosen by processor:

| `uname -m` prints | Download |
| --- | --- |
| `arm64`: Apple Silicon, M1 and later | `core-geth-osx-arm64-<tag>.zip` |
| `x86_64`: Intel | `core-geth-osx-<tag>.zip` |

macOS has no `sha256sum`; use `shasum -a 256 -c`, which reads the same file. On an Intel Mac, put `osx` where the
file names below say `osx-arm64`:

```shell
VERSION=$(curl -sLo /dev/null -w '%{url_effective}' https://github.com/ethereumclassic/core-geth/releases/latest | sed 's|.*/||')
BASE=https://github.com/ethereumclassic/core-geth/releases/download/$VERSION
curl -LO $BASE/core-geth-osx-arm64-$VERSION.zip
curl -LO $BASE/core-geth-osx-arm64-$VERSION.zip.sha256
shasum -a 256 -c core-geth-osx-arm64-$VERSION.zip.sha256 &&
  unzip -o core-geth-osx-arm64-$VERSION.zip &&
  sudo install -m 0755 geth /usr/local/bin/geth &&
  geth version
```

!!! tip "If macOS refuses to run it"
    These binaries are not code-signed, so Gatekeeper blocks one that carries the
    quarantine attribute with *"cannot be opened because the developer cannot be
    verified"*. Browsers set that attribute; `curl` does not, so the commands above
    avoid it. If you downloaded through a browser, clear it:

    ```shell
    xattr -d com.apple.quarantine geth
    ```

    Once macOS has blocked it, **Open Anyway** in System Settings, under Privacy & Security,
    also lets it run ([Apple's instructions](https://support.apple.com/en-us/102445)).

Before running it on a node that holds value, also
[verify where the archive came from](#verify-where-the-archive-came-from).

## Windows

Download `core-geth-win64-<tag>.zip` and its `.sha256`. In PowerShell:

```powershell
$VERSION = (Invoke-RestMethod "https://api.github.com/repos/ethereumclassic/core-geth/releases/latest").tag_name
$BASE = "https://github.com/ethereumclassic/core-geth/releases/download/$VERSION"
Invoke-WebRequest "$BASE/core-geth-win64-$VERSION.zip" -OutFile "core-geth-win64-$VERSION.zip"
Invoke-WebRequest "$BASE/core-geth-win64-$VERSION.zip.sha256" -OutFile "core-geth-win64-$VERSION.zip.sha256"
(Get-FileHash "core-geth-win64-$VERSION.zip" -Algorithm SHA256).Hash
Get-Content "core-geth-win64-$VERSION.zip.sha256"
```

Compare the last two lines: the letters must match, and the case does not matter. If they differ, delete the file and
download it again.

If they match, unpack it:

```powershell
Expand-Archive "core-geth-win64-$VERSION.zip" -DestinationPath C:\core-geth -Force
C:\core-geth\geth.exe version
```

To keep `geth.exe` in another folder, unpack it there, and add that folder to `PATH` to run it by name.

!!! tip "If Windows warns about the publisher"
        These binaries are not code-signed, so Windows shows their publisher as unknown, in the firewall alert on the first
    start and in the *"Windows protected your PC"* screen. Check the archive first, against its checksum and
    [where it came from](#verify-where-the-archive-came-from). Then allow the firewall access, and select **More info**,
    then **Run anyway**. To avoid that screen for a file a browser downloaded, remove its mark before you extract it:

    ```powershell
    Unblock-File "core-geth-win64-$VERSION.zip"
    ```

Before running it on a node that holds value, also
[verify where the archive came from](#verify-where-the-archive-came-from).

## Verify where the archive came from

**A checksum proves the download was not corrupted. It does not prove what built the
file**, because whoever publishes the archive publishes the checksum beside it. From
v1.13.0 each archive also carries a build attestation, which records the workflow,
commit and ref that produced it. In the same shell as your route above, naming the archive you downloaded:

=== "Linux, x86_64"

    ```shell
    gh attestation verify core-geth-linux-$VERSION.zip --repo ethereumclassic/core-geth
    ```

=== "Linux, ARM"

    ```shell
    gh attestation verify core-geth-arm64-$VERSION.zip --repo ethereumclassic/core-geth
    ```

=== "macOS"

    ```shell
    gh attestation verify core-geth-osx-arm64-$VERSION.zip --repo ethereumclassic/core-geth
    ```

=== "Windows"

    ```powershell
    gh attestation verify "core-geth-win64-$VERSION.zip" --repo ethereumclassic/core-geth
    ```

That needs the [GitHub CLI](https://cli.github.com/), signed in with `gh auth login`. It succeeds only for a file built
by this repository's release workflow: a file from anywhere else fails, as does one whose bytes have changed.

## Docker

Images are published to the GitHub Container Registry for each tagged release, built for
`linux/amd64` and `linux/arm64`.

**Pulling needs no account and no token:**

```shell
docker pull ghcr.io/ethereumclassic/core-geth:latest
```

To check where the pulled image came from, verify its build attestation:

```shell
gh attestation verify oci://ghcr.io/ethereumclassic/core-geth:latest --repo ethereumclassic/core-geth
```

That needs the [GitHub CLI](https://cli.github.com/), signed in. It succeeds only for an image built by this
repository's workflows.

With no route to the registry, [load the image from the release](#load-the-image-from-the-release) instead.

Tags mirror the release tags. `ghcr.io/ethereumclassic/core-geth:<tag>` is a
specific release; `latest` follows the most recent non-prerelease. The image
built from `Dockerfile.alltools`, containing the full tool set rather than the
node alone, is published under the same name with an `alltools-` prefix, as
`ghcr.io/ethereumclassic/core-geth:alltools-latest`.




### Load the image from the release

Each release also carries its images as files, one per architecture:
`core-geth-docker-amd64-<tag>.tar.gz` and `core-geth-docker-arm64-<tag>.tar.gz`, and
`core-geth-alltools-docker-amd64-<tag>.tar.gz` and `core-geth-alltools-docker-arm64-<tag>.tar.gz`
for the full tool set. Each is the published image saved from the registry, with a `.sha256` and
a build attestation. Download the one for your machine, check it and load it. On a 64-bit ARM host, put `arm64` where
the file names below say `amd64`:

```shell
VERSION=$(curl -sLo /dev/null -w '%{url_effective}' https://github.com/ethereumclassic/core-geth/releases/latest | sed 's|.*/||')
BASE=https://github.com/ethereumclassic/core-geth/releases/download/$VERSION
curl -LO $BASE/core-geth-docker-amd64-$VERSION.tar.gz
curl -LO $BASE/core-geth-docker-amd64-$VERSION.tar.gz.sha256
sha256sum -c core-geth-docker-amd64-$VERSION.tar.gz.sha256
```

On macOS, use `shasum -a 256 -c` in place of `sha256sum -c`. It must print `OK`: if it prints `FAILED`, delete the
file and download it again rather than loading it. On a node that holds value, also check where it came from, which
needs the [GitHub CLI](https://cli.github.com/), signed in with `gh auth login`:

```shell
gh attestation verify core-geth-docker-amd64-$VERSION.tar.gz --repo ethereumclassic/core-geth
```

Then load it:

```shell
docker load -i core-geth-docker-amd64-$VERSION.tar.gz
```

The loaded image keeps its published name with its release tag, `ghcr.io/ethereumclassic/core-geth:<tag>`, and
carries no `latest` tag. Use that name in the `docker run` command below: a command that names `latest` makes
Docker pull from the registry.

### Build an image from the release archive

This packages the binary the release attested, rather than compiling one. Save this `Dockerfile` in an empty folder:

```dockerfile
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*
COPY geth /usr/local/bin/geth
EXPOSE 8545 8546 30303 30303/udp
ENTRYPOINT ["geth"]
```

Then, in that folder, download the latest release for `linux/amd64`, check it and build. For `linux/arm64`, put
`arm64` where the file names say `linux`. On macOS, use `shasum -a 256 -c` in place of `sha256sum -c`. On a node that
holds value, also [verify where the archive came from](#verify-where-the-archive-came-from):

```shell
VERSION=$(curl -sLo /dev/null -w '%{url_effective}' https://github.com/ethereumclassic/core-geth/releases/latest | sed 's|.*/||')
BASE=https://github.com/ethereumclassic/core-geth/releases/download/$VERSION
curl -LO $BASE/core-geth-linux-$VERSION.zip
curl -LO $BASE/core-geth-linux-$VERSION.zip.sha256
sha256sum -c core-geth-linux-$VERSION.zip.sha256 &&
  unzip -o core-geth-linux-$VERSION.zip &&
  docker build -t core-geth:local . &&
  docker run --rm core-geth:local version
```

Use a glibc-based image such as Debian or Ubuntu. The release binaries link against glibc, so an
Alpine base cannot run them.

### Build an image from source

The `Dockerfile` produces an image containing `geth`, and `Dockerfile.alltools` one containing
the full tool set. These commands build the latest release:

```shell
VERSION=$(curl -sLo /dev/null -w '%{url_effective}' https://github.com/ethereumclassic/core-geth/releases/latest | sed 's|.*/||')
git clone --branch $VERSION --single-branch https://github.com/ethereumclassic/core-geth.git core-geth-$VERSION &&
  cd core-geth-$VERSION &&
  docker build -t core-geth:local .
```

### Run the container

Run the image the same way whichever method produced it. Flags go straight to the node:

```shell
docker run -d \
    --name core-geth \
    -v $LOCAL_DATADIR:/root \
    -p 30303:30303 -p 30303:30303/udp \
    -p 127.0.0.1:8545:8545 \
    ghcr.io/ethereumclassic/core-geth:latest \
    --classic \
    --http --http.addr 0.0.0.0 --http.port 8545
```

If you built the image yourself, use your own tag (`core-geth:local`,
above) in place of `ghcr.io/ethereumclassic/core-geth:latest`.

That maps the devp2p port over both TCP and UDP, keeps chain data in
`$LOCAL_DATADIR` on the host so it survives the container, and reaches the
JSON-RPC endpoint from the host and nowhere else.

!!! warning "Keep both halves of that RPC line"
    `--http` alone binds to `localhost` inside the container, so a published port reaches nothing, with no error.
    `--http.addr 0.0.0.0` fixes that, and `-p 127.0.0.1:8545:8545` keeps the port off the host's public interfaces.
    **Publishing as `-p 8545:8545` instead exposes your RPC on every interface of the host, which on a VPS is the open
    internet.**

    If you do not need RPC from the host, drop both the `-p 127.0.0.1:8545:8545` and the `--http*` flags.

The image also exposes 8546 for the WebSocket endpoint, which needs
`--ws --ws.addr 0.0.0.0` to be served.

## Build from source

Building needs **Go 1.26** and a C compiler; the developer guide has the
[per-platform toolchain setup](../developers/build-from-source.md#dependencies)
for both. This is also the route for a platform no archive covers, or a
distribution older than the glibc floors above. These commands build the latest release:

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

Do not clone with `--recursive`: the submodules are several gigabytes, and `make geth` does not read them.

To check a build, compare the commit you checked out with the one the binary was built from. In the clone, both lines
must name the same commit:

```shell
git rev-parse HEAD
go version -m ./build/bin/geth | grep vcs.revision
```

On Windows, run `go version -m .\build\bin\geth.exe` and read its `vcs.revision` line.

`make all` builds every tool rather than the node alone.
[Build from source](../developers/build-from-source.md) covers the full set of targets,
building what is in development, and building inside Docker.

## What else is in a release

Each platform publishes two archives:

| Archive | Contains |
| --- | --- |
| `core-geth-<platform>-<tag>.zip` | the `geth` node binary, and `COPYING` |
| `core-geth-alltools-<platform>-<tag>.zip` | `geth` plus the other tools built from this source, and `COPYING` |

Every route above uses the first. The `alltools` archive adds every other executable
built from this source (`abidump`, `abigen`, `ancient-store-mem`, `bootnode`,
`clef`, `devp2p`, `echainspec`, `era`, `ethkey`, `evm`, `p2psim` and `rlpdump`) and
installs the same way.

Once `geth` is installed, [Running a node](run-a-node.md) takes it to a synced node that keeps
running, with a guide for each platform and the flags for each kind of node.
