---
title: Build from Source
description: "Build Core-Geth from source on Linux, macOS or Windows: the tools and the Go release main pins, the build, a running node, and the test suites."
---



## Dependencies

A build needs four tools: Git, to clone the code and record the commit in the binary; Go; a C compiler, for the parts
written in C; and `make`, outside Windows. These steps start from a machine with none of them, install each at its
current release, and check them.

**Go is the release `main` pins.** `main` pins Go 1.26.8 in its `build/checksums.txt`, the Go its release builds use.
The steps install Go 1.26.8 from the Go project's downloads at go.dev, and check the download before installing it. A
package manager's Go is whatever release it carries, older than 1.26 on Debian 13 and Ubuntu 24.04.

=== "Linux"

    **Install Git, a C compiler and `make`.** On Debian or Ubuntu:

    ```shell
    sudo apt-get update && sudo apt-get install -y build-essential git curl
    ```

    On Fedora or RHEL:

    ```shell
    sudo dnf install -y --refresh gcc make git && sudo dnf upgrade -y gcc make git
    ```

    On Arch, where `-Syu` also brings the rest of the system current, as Arch requires:

    ```shell
    sudo pacman -Syu --needed base-devel git
    ```

    Each line also updates a tool that is already installed. `build-essential` and `base-devel` carry the C compiler
    and `make`, and the next step downloads with `curl`.

    **Install Go 1.26.8.** This checks the download against the SHA-256 `build/checksums.txt` pins for it:

    ```shell
    curl -LO https://go.dev/dl/go1.26.8.linux-amd64.tar.gz
    echo "d0f743b33e8d8945e6b1f432edd15785c70507121d6e2a723b21285eddf8b57b  go1.26.8.linux-amd64.tar.gz" | sha256sum -c &&
      sudo rm -rf /usr/local/go &&
      sudo tar -C /usr/local -xzf go1.26.8.linux-amd64.tar.gz &&
      sudo ln -sf /usr/local/go/bin/go /usr/local/bin/go &&
      hash -r &&
      go version
    ```

    `sha256sum -c` must print `OK`: if it prints `FAILED`, delete the file and download it again rather than installing
    it. The block replaces any Go in `/usr/local/go`, and `go version` must then print `go1.26.8`. On a 64-bit Arm
    machine, where `uname -m` prints `aarch64`, put `go1.26.8.linux-arm64.tar.gz` in place of the file name and
    `211ffced9dcb9633a55eac6364816ec0ddd951389a740e88fa8b3337971bdda0` in place of the SHA-256. On 32-bit Arm, use
    `go1.26.8.linux-armv6l.tar.gz` and `eab440beabf395870752021fa74cedf04f97b61e43ebbe005ecbb98b94e55697`.

    **Check the tools:**

    ```shell
    git --version
    go version
    gcc --version
    make --version
    ```

    Each line must print a version, and `go version` must print `go1.26.8`. A line that prints an error names a tool
    the steps above did not install. If `go version` prints another release, a Go earlier on `PATH` than
    `/usr/local/bin` runs in its place: uninstall that Go.

=== "macOS"

    **Install the Xcode Command Line Tools:**

    ```shell
    xcode-select --install
    ```

    A window asks whether to install the tools: choose Install, and wait until it says the software was installed. The
    tools carry Git, the C compiler and `make`. If they are already installed, the command says so: install any update
    for them that Software Update lists in System Settings, then go on to the next step.

    **Install Go 1.26.8.** This checks the download against the SHA-256 `build/checksums.txt` pins for it. `uname -m`
    prints `arm64` on Apple Silicon and `x86_64` on Intel. On an Apple Silicon Mac:

    ```shell
    curl -LO https://go.dev/dl/go1.26.8.darwin-arm64.tar.gz
    echo "a012b25b571bd0138a03dcd25375ceba866fe5ca822f426d2c66a4de56fd3f4b  go1.26.8.darwin-arm64.tar.gz" | shasum -a 256 -c &&
      sudo rm -rf /usr/local/go &&
      sudo tar -C /usr/local -xzf go1.26.8.darwin-arm64.tar.gz &&
      sudo mkdir -p /usr/local/bin &&
      sudo ln -sf /usr/local/go/bin/go /usr/local/bin/go &&
      hash -r &&
      go version
    ```

    On an Intel Mac:

    ```shell
    curl -LO https://go.dev/dl/go1.26.8.darwin-amd64.tar.gz
    echo "186be014105aa6542b767d2c6ed5cca10a0214bdff809ef1724022a8c7894150  go1.26.8.darwin-amd64.tar.gz" | shasum -a 256 -c &&
      sudo rm -rf /usr/local/go &&
      sudo tar -C /usr/local -xzf go1.26.8.darwin-amd64.tar.gz &&
      sudo mkdir -p /usr/local/bin &&
      sudo ln -sf /usr/local/go/bin/go /usr/local/bin/go &&
      hash -r &&
      go version
    ```

    `shasum -a 256 -c` must print `OK`: if it prints `FAILED`, delete the file and download it again rather than
    installing it. The block replaces any Go in `/usr/local/go`, and `go version` must then print `go1.26.8`.

    **Check the tools:**

    ```shell
    git --version
    go version
    gcc --version
    make --version
    ```

    Each line must print a version, and `go version` must print `go1.26.8`. `gcc --version` prints Apple clang, the C
    compiler on macOS. A line that prints an error names a tool the steps above did not install. If `go version` prints
    another release, a Go earlier on `PATH`, such as one from Homebrew, runs in its place: uninstall that Go.

=== "Windows"

    Chocolatey installs Git and the C compiler. A Windows build needs no `make`.

    **Check for Chocolatey.** In a PowerShell window opened with **Run as administrator**, which the Chocolatey steps
    need:

    ```powershell
    choco --version
    ```

    If PowerShell says the term is not recognized, Chocolatey is not installed: take the next step. If the version starts
    with 1, skip the next step and take the one after it. If it starts with 2, go to "Update Chocolatey".

    **If Chocolatey is not installed, install it.** WinGet, which comes with Windows 11, installs Chocolatey from its
    publisher:

    ```powershell
    winget install --id Chocolatey.Chocolatey -e --source winget --accept-source-agreements --accept-package-agreements
    ```

    Approve the prompt the installer shows. Then open a new PowerShell window with **Run as administrator** for the next
    steps: a window opened before the install does not find Chocolatey.

    **If Chocolatey is version 1, move it to 1.4.0 first:**

    ```powershell
    choco upgrade chocolatey --version 1.4.0 --allow-downgrade -y
    ```

    Chocolatey's [guide to version 2](https://docs.chocolatey.org/en-us/guides/upgrading-to-chocolatey-v2-v6/) moves a
    version 1 to 1.4.0 before the upgrade to version 2, which the next step makes. Skip this step on version 2: it would
    downgrade Chocolatey.

    **Update Chocolatey, then install Git and a C compiler:**

    ```powershell
    choco upgrade chocolatey -y
    choco upgrade mingw git -y
    choco --version
    ```

    `choco upgrade` installs what is missing and moves what is installed to its current release, so an older Chocolatey,
    Git or MinGW-w64 is updated rather than kept. The last line must print a version starting with 2. If the first line
    asks for a restart to install .NET Framework 4.8, which Chocolatey 2 needs, restart and run this block again. `mingw`
    is MinGW-w64, the C compiler cgo needs on Windows.

    **Download Go 1.26.8 and check it.** `build/checksums.txt` pins Go's archives, not its installer, so this checks the
    installer against the SHA-256 go.dev publishes for it. An administrator window opens in `C:\Windows\System32`, so the
    first line moves to your home folder:

    ```powershell
    cd $HOME
    Invoke-WebRequest "https://go.dev/dl/go1.26.8.windows-amd64.msi" -OutFile "go1.26.8.windows-amd64.msi"
    (Get-FileHash "go1.26.8.windows-amd64.msi" -Algorithm SHA256).Hash -eq "0581e37d54c8f2876c447204737d4aaa3f65f983005c65a1d5186b5f18ae3eca"
    ```

    The last line must print `True`. If it prints `False`, delete the file and download it again rather than running it.

    **If the check printed `True`, run the Go installer:**

    ```powershell
    msiexec /i "$HOME\go1.26.8.windows-amd64.msi" /passive | Out-Null
    ```

    The installer shows a progress bar and takes its defaults, which put Go in `C:\Program Files\Go` and on `PATH`. The
    line returns once it has finished.

    **Check the tools,** in the same window:

    ```powershell
    $env:Path = [Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [Environment]::GetEnvironmentVariable("Path", "User")
    git --version
    go version
    gcc --version
    ```

    The first line reloads `PATH` from Windows' settings, so this window finds what the steps above installed. Each line
    after it must print a version, and `go version` must print `go1.26.8`. A line saying the term is not recognized
    names a tool the steps above did not install.

!!! tip "When your Go is not the pinned release"
    After `main` moves its pin past the Go you installed, `go run build/ci.go install -dlgo ./cmd/geth` downloads the
    release `build/checksums.txt` names, checks it against the SHA-256 that file lists, and builds `geth` with it. The
    Go you installed stays as it is, and runs the build script.

## Source

Once the tools are installed, clone `main` and build it:

=== "Linux and macOS"

    ```shell
    git clone --branch main --single-branch https://github.com/ethereumclassic/core-geth.git &&
      cd core-geth &&
      make geth &&
      ./build/bin/geth version
    ```

=== "Windows"

    In the PowerShell window you checked the tools in, one line at a time. `make geth` runs this same `go run` command
    on other systems:

    ```powershell
    cd $HOME
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

Then install the build where a release installs the client:

=== "Linux and macOS"

    ```shell
    ls build/bin/geth &&
      sudo mkdir -p /usr/local/bin &&
      sudo install -m 0755 build/bin/geth /usr/local/bin/geth &&
      geth version
    ```

    Run this in the `core-geth` folder the build left you in: the first line stops the block anywhere else. A new Mac
    can lack `/usr/local/bin`, which the second line creates. `geth version` must print the same first line and
    `Git Commit` line the build printed.

=== "Windows"

    ```powershell
    New-Item -ItemType Directory -Force C:\core-geth | Out-Null
    Copy-Item .\build\bin\geth.exe C:\core-geth\geth.exe
    C:\core-geth\geth.exe version
    ```

    Run these in the `core-geth` folder the build left you in. The last line must print the same first line and
    `Git Commit` line the build printed. `C:\core-geth` is where the Windows download puts the client.

Then start a node on Ethereum Classic:

=== "Linux and macOS"

    ```shell
    geth --classic
    ```

    For the Mordor test network, run `geth --mordor` instead.

=== "Windows"

    ```powershell
    cd C:\core-geth
    .\geth.exe --classic
    ```

    For the Mordor test network, run `.\geth.exe --mordor` instead.

The node runs in the foreground and logs to the terminal. [Running a node](../getting-started/run-a-node.md) has a guide
for each platform: how to check that it syncs, keep it running and stop it cleanly.

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

To build a local image of `main`, with Docker installed and running as
[Installation, Docker](../getting-started/installation.md#docker) checks:

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
