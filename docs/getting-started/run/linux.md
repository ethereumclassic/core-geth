---
description: "Run a Core-Geth Ethereum Classic or Mordor node on Linux: install geth, start it, and set up a systemd service that starts the node at boot."
---

# Linux users guide

This guide runs a node on Linux, on Ethereum Classic or on the Mordor test network, from an
installed `geth` to a service that starts at boot. [Running a node](../run-a-node.md) covers what
every platform shares.

Commands are shown for both networks. Use the lines for the network you chose, and keep its flag on
every command.

## 1. Install

Install `geth` as [Installation](../installation.md) shows for
[x86_64](../installation.md#linux-x86_64) or [ARM](../installation.md#linux-arm), then check it:

```bash
geth version        # Core-Geth 1.13.0-stable
```

## 2. Choose the network

| Network | Flag | Data directory |
|---|---|---|
| Ethereum Classic | `--classic` | `~/.ethereum/classic` |
| Mordor | `--mordor` | `~/.ethereum/mordor` |

With no flag, the node runs Ethereum Classic. This guide passes `--classic` anyway, so that every
command says which network it is for.

## 3. Start the node

Other nodes reach yours on port 30303, over TCP and UDP. Allow it through your firewall
([A host firewall](../../operate/security.md#a-host-firewall)), then start the node:

```bash
geth --classic      # Ethereum Classic
geth --mordor       # Mordor
```

The node runs in the foreground and logs to the terminal. Its first line names the network:
`Starting Core-Geth on Ethereum Classic...` or `Starting Core-Geth on Mordor testnet...`.
[What a healthy first sync looks like](../run-classic-node.md#what-a-healthy-first-sync-looks-like)
explains the lines that follow.

**For a [configuration](../run-a-node.md#choose-your-configuration)**, add its flags after the
network flag, such as `geth --classic --http --http.api eth,net,web3`, and add them to `ExecStart`
in step 6 as well.

To keep the chain on another disk, add `--datadir`, such as `--datadir /data/core-geth/mordor`. Then
give `attach` the socket inside that directory: `geth --mordor attach /data/core-geth/mordor/geth.ipc`.

## 4. Check that it syncs

In a second terminal:

```bash
# Ethereum Classic
geth --classic attach --exec 'eth.syncing'       # an object while it syncs, false once synced
geth --classic attach --exec 'net.peerCount'     # above 0 within minutes
geth --classic attach --exec 'eth.blockNumber'   # once synced, compare with a block explorer

# Mordor
geth --mordor attach --exec 'eth.syncing'
geth --mordor attach --exec 'net.peerCount'
geth --mordor attach --exec 'eth.blockNumber'
```

`attach` looks for the node's socket in the data directory of the network you name. Without
`--mordor`, it looks in `~/.ethereum/classic` and does not find a Mordor node:
`no such file or directory`. [How to tell it is done](../run-classic-node.md#how-to-tell-it-is-done)
lists every sign of a finished sync.

## 5. Use it

`attach` without `--exec` opens an interactive console on the node. Type `exit` to leave it; the node
keeps running.

```bash
geth --mordor attach
```

For wallets and programs, start the node with `--http`, such as `geth --mordor --http`. It serves
JSON-RPC on `127.0.0.1:8545`. From a second terminal:

```bash
curl -s -X POST -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","id":1,"method":"eth_chainId","params":[]}' http://127.0.0.1:8545
```

The `result` is `0x3d` on Ethereum Classic and `0x3f` on Mordor. Read
[RPC exposure](../../operate/security.md#rpc-exposure) before you serve it beyond `127.0.0.1`.

## 6. Keep it running: a systemd service

The project ships no service unit, so this one is an example to adapt. It runs the node as a
dedicated user, starts it at boot, restarts it if it fails, and gives it five minutes to stop.

**Create the user.** `--create-home` matters: `geth attach` needs a home directory it can write to,
even when you give it the socket's path.

```bash
sudo useradd --system --create-home --shell /usr/sbin/nologin geth
```

**Save the unit.** For Ethereum Classic, save this as `/etc/systemd/system/core-geth.service`:

```ini
[Unit]
Description=Core-Geth node on Ethereum Classic
Wants=network-online.target
After=network-online.target

[Service]
User=geth
StateDirectory=core-geth
ExecStart=/usr/local/bin/geth --classic --datadir /var/lib/core-geth
Restart=on-failure
TimeoutStopSec=300

[Install]
WantedBy=multi-user.target
```

For Mordor, save the same unit as `/etc/systemd/system/core-geth-mordor.service`, with these three
lines in place of the ones they match:

```ini
Description=Core-Geth node on Mordor
StateDirectory=core-geth-mordor
ExecStart=/usr/local/bin/geth --mordor --datadir /var/lib/core-geth-mordor
```

**Start it now and at every boot, then follow its log.** For Mordor, use `core-geth-mordor` in place
of `core-geth`. Reading the system journal takes membership of the `adm` or `systemd-journal` group,
or `sudo`:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now core-geth
journalctl -u core-geth -f
```

What the unit's settings do:

- **`StateDirectory`** has systemd create the data directory under `/var/lib`, owned by `geth`. To
  keep a data directory you already synced in the foreground, put its path in `--datadir` instead,
  and make it owned by `geth`.
- **`ExecStart`** uses the path [Installation](../installation.md) installs `geth` to. A configuration's
  flags go at the end of it.
- **`Restart=on-failure`** restarts the node after it exits with an error or dies from a signal such
  as SIGKILL. A node that shuts itself down because the disk is nearly full exits cleanly, so systemd
  leaves it stopped ([running out of disk](../hardware-requirements.md#running-out-of-disk)).
- **`TimeoutStopSec=300`** is how long systemd waits after its SIGTERM before it kills the node with
  SIGKILL. Five minutes is well beyond every clean stop measured, including one made during a first
  sync ([stop times](../hardware-requirements.md#stopping)). A killed node skips the shutdown that
  writes its state to disk.

The node's socket belongs to `geth`, and only `geth` and root can connect to it, so run `attach` as
`geth`. `-H` gives it `geth`'s home directory to write to:

```bash
sudo -u geth -H geth --classic attach --exec 'eth.syncing' /var/lib/core-geth/geth.ipc
sudo -u geth -H geth --mordor attach --exec 'eth.syncing' /var/lib/core-geth-mordor/geth.ipc
```

Both services on one host also need separate ports:
[Mordor beside Ethereum Classic on one host](../run-mordor-node.md#mordor-beside-ethereum-classic-on-one-host).

## 7. Stop it cleanly

- **In the foreground:** press Ctrl-C once, then wait for `Blockchain stopped` and the prompt.
- **As a service:** `sudo systemctl stop core-geth`, or `core-geth-mordor`. systemd waits up to
  `TimeoutStopSec` for the node to finish.

Pressing Ctrl-C again does not speed it up. [Stop it safely](../run-classic-node.md#stop-it-safely)
shows what a clean stop logs, and why a stop during the first sync takes longer.

## 8. Update it

Stop the node, install the new release's `geth` the same way as in step 1, and start the node again.
[Upgrading within v1.13.x](../../operate/maintenance.md#upgrading-within-v113x) lists what to read
and check. A node on a v1.12.x release follows the
[Linux migration guide](../../tutorials/migration/linux.md) instead.
