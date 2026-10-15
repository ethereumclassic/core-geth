---
description: "Run a Core-Geth Ethereum Classic or Mordor node on Windows: from geth.exe to a node that starts when you sign in."
---

# Windows users guide

This guide runs a node on Windows, on Ethereum Classic or on the Mordor test network, from
`geth.exe` to a node that starts when you sign in. [Running a node](../run-a-node.md) covers what
every platform shares.

The commands are for PowerShell. They keep `geth.exe` in `C:\core-geth`; if you keep it in another
folder, use that one.

## 1. Install

Download and check the archive as [Installation, Windows](../installation.md#windows) shows, and
unpack `geth.exe` into `C:\core-geth`. Then:

```powershell
cd C:\core-geth
.\geth.exe version        # Core-Geth 1.13.0-stable
```

## 2. Choose the network

| Network | Flag | Data directory |
|---|---|---|
| Ethereum Classic | `--classic` | `%LOCALAPPDATA%\Ethereum\classic` |
| Mordor | `--mordor` | `%LOCALAPPDATA%\Ethereum\mordor` |

With no flag, the node runs Ethereum Classic. This guide passes `--classic` anyway, so that every
command says which network it is for. If an `AppData\Roaming\Ethereum` folder in your user folder
already exists and is not empty, the node uses the `classic` or `mordor` folder inside it instead.

## 3. Start the node

```powershell
.\geth.exe --classic      # Ethereum Classic
.\geth.exe --mordor       # Mordor
```

The node runs in the window and logs to it. Its first line names the network:
`Starting Core-Geth on Ethereum Classic...` or `Starting Core-Geth on Mordor testnet...`.
[What a healthy first sync looks like](../run-classic-node.md#what-a-healthy-first-sync-looks-like)
explains the lines that follow.

**The first start can raise a Windows Security alert** asking whether to allow `geth.exe` through the
firewall. It shows the publisher as unknown, because the release binaries are not code-signed
([If Windows warns about the publisher](../installation.md#windows)). Allow it, so other nodes can
connect to yours on port 30303. If you do not, the node still
syncs, through the connections it makes itself.

**A click in the window can pause the node.** In a classic console window, clicking starts a text
selection. While the selection lasts, the node waits to write its log and stops syncing. Press Esc to
end the selection.

**For a [configuration](../run-a-node.md#choose-your-configuration)**, add its flags after the
network flag, such as `.\geth.exe --classic --http --http.api eth,net,web3`. In step 6, add them to
the end of the command file's `geth.exe` line.

To keep the chain on another drive, add `--datadir`, such as `--datadir D:\core-geth\mordor`.

## 4. Check that it syncs

Open a second PowerShell window:

```powershell
cd C:\core-geth
.\geth.exe attach --exec 'eth.syncing'       # an object while it syncs, false once synced
.\geth.exe attach --exec 'net.peerCount'     # above 0 within minutes
.\geth.exe attach --exec 'eth.blockNumber'   # once synced, compare with a block explorer
```

**These commands are the same on both networks.** On Windows, `attach` connects to the named pipe
`\\.\pipe\geth.ipc`, which the running node opens whatever its network and data directory. A second
node on the same machine needs a pipe of its own: start it with `--ipcpath geth-mordor.ipc`, and
attach to it with `.\geth.exe attach \\.\pipe\geth-mordor.ipc`. It needs its own ports too
([Mordor beside Ethereum Classic on one host](../run-mordor-node.md#mordor-beside-ethereum-classic-on-one-host)).

[How to tell it is done](../run-classic-node.md#how-to-tell-it-is-done) lists every sign of a
finished sync.

## 5. Use it

`attach` without `--exec` opens an interactive console on the node. Type `exit` to leave it; the node
keeps running.

```powershell
.\geth.exe attach
```

For wallets and programs, start the node with `--http`, such as `.\geth.exe --mordor --http`. It
serves JSON-RPC on `127.0.0.1:8545`. From a second PowerShell window:

```powershell
Invoke-RestMethod http://127.0.0.1:8545 -Method Post -ContentType 'application/json' -Body '{"jsonrpc":"2.0","id":1,"method":"eth_chainId","params":[]}'
```

The `result` is `0x3d` on Ethereum Classic and `0x3f` on Mordor. Read
[RPC exposure](../../operate/security.md#rpc-exposure) before you serve it beyond `127.0.0.1`.

## 6. Keep it running: start it when you sign in

Core-Geth does not install itself as a Windows service. This step starts the node in its own window
each time you sign in, where Ctrl+C still stops it cleanly. It writes a command file to your Startup
folder:

```powershell
$startup = [Environment]::GetFolderPath('Startup')

# Ethereum Classic
Set-Content "$startup\core-geth-classic.cmd" '@title Core-Geth: Ethereum Classic', '"C:\core-geth\geth.exe" --classic'

# Mordor
Set-Content "$startup\core-geth-mordor.cmd" '@title Core-Geth: Mordor', '"C:\core-geth\geth.exe" --mordor'
```

To start it now without signing out, open the file, such as
`Start-Process "$startup\core-geth-mordor.cmd"`. To stop it starting at sign-in, delete the file:
`Remove-Item "$startup\core-geth-mordor.cmd"`.

The node runs while you are signed in. To run a node with no one signed in, run it under a service
wrapper that stops it by sending Ctrl+C and then waits for it to exit. Core-Geth ships no wrapper.

## 7. Stop it cleanly

Press Ctrl+C once in the node's window, then wait for `Blockchain stopped` and for the node to exit.
A window opened from the Startup folder then asks `Terminate batch job (Y/N)?`; type `Y`.

**Stop the node before you close its window, sign out, or shut Windows down.** Each of those gives the
node only a short time before Windows ends it, and a node still stopping at that point is ended
without finishing ([After an unclean stop](../../operate/maintenance.md#after-an-unclean-stop)).
[Stop it safely](../run-classic-node.md#stop-it-safely) shows what a clean stop logs, and why a stop
during the first sync takes longer.

## 8. Update it

Stop the node, replace `C:\core-geth\geth.exe` with the new release's, and start the node again.
[Upgrading within v1.13.x](../../operate/maintenance.md#upgrading-within-v113x) lists what to read
and check. A node on a v1.12.x release follows the
[Windows migration guide](../../tutorials/migration/windows.md) instead.
