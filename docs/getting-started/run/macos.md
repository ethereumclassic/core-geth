---
description: "Run a Core-Geth Ethereum Classic or Mordor node on macOS: install geth and set up a launchd agent that starts the node when you log in."
---

# Mac users guide

This guide runs a node on macOS, on Ethereum Classic or on the Mordor test network, from an
installed `geth` to a launchd agent that starts it when you log in.
[Running a node](../run-a-node.md) covers what every platform shares.

Commands are shown for both networks. Use the lines for the network you chose, and keep its flag on
every command.

## 1. Install

Install `geth` as [Installation, macOS](../installation.md#macos) shows, including clearing the
quarantine attribute if a browser downloaded it. Then check it:

```bash
geth version        # Core-Geth 1.13.0-stable
```

## 2. Choose the network

| Network | Flag | Data directory |
|---|---|---|
| Ethereum Classic | `--classic` | `~/Library/Ethereum/classic` |
| Mordor | `--mordor` | `~/Library/Ethereum/mordor` |

With no flag, the node runs Ethereum Classic. This guide passes `--classic` anyway, so that every
command says which network it is for.

## 3. Start the node

```bash
geth --classic      # Ethereum Classic
geth --mordor       # Mordor
```

The node runs in the foreground and logs to the terminal. Its first line names the network:
`Starting Core-Geth on Ethereum Classic...` or `Starting Core-Geth on Mordor testnet...`.
[What a healthy first sync looks like](../run-classic-node.md#what-a-healthy-first-sync-looks-like)
explains the lines that follow. Other nodes reach yours on port 30303, over TCP and UDP
([Ports and listeners](../../operate/security.md#ports-and-listeners)).

**For a [configuration](../run-a-node.md#choose-your-configuration)**, add its flags after the
network flag, such as `geth --classic --http --http.api eth,net,web3`. Step 6 shows how the agent
takes them.

To keep the chain on another disk, add `--datadir`, such as `--datadir /Volumes/Data/core-geth/mordor`.
Then give `attach` the socket inside that directory:
`geth --mordor attach /Volumes/Data/core-geth/mordor/geth.ipc`.

## 4. Check that it syncs

In a second terminal window:

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

`attach` looks for the node's socket in the data directory of the network you name, so a Mordor node
needs `--mordor` here too. [How to tell it is done](../run-classic-node.md#how-to-tell-it-is-done)
lists every sign of a finished sync.

## 5. Use it

`attach` without `--exec` opens an interactive console on the node. Type `exit` to leave it; the node
keeps running.

```bash
geth --mordor attach
```

For wallets and programs, start the node with `--http`, such as `geth --mordor --http`. It serves
JSON-RPC on `127.0.0.1:8545`. From a second terminal window:

```bash
curl -s -X POST -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","id":1,"method":"eth_chainId","params":[]}' http://127.0.0.1:8545
```

The `result` is `0x3d` on Ethereum Classic and `0x3f` on Mordor. Read
[RPC exposure](../../operate/security.md#rpc-exposure) before you serve it beyond `127.0.0.1`.

## 6. Keep it running: a launchd agent

A launchd agent starts the node when you log in, restarts it if it fails, and gives it five minutes to
stop. It runs while you are logged in, and a Mac that sleeps pauses the node until it wakes.

**Create the agent.** Set `NETWORK` to `classic` or `mordor`, then run the rest as it is. The file
it writes names the network in the agent's label, the flag and the log file:

```bash
NETWORK=classic      # or: NETWORK=mordor
mkdir -p ~/Library/LaunchAgents ~/Library/Logs
cat > ~/Library/LaunchAgents/com.coregeth.$NETWORK.plist <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>com.coregeth.$NETWORK</string>
  <key>ProgramArguments</key>
  <array>
    <string>/usr/local/bin/geth</string>
    <string>--$NETWORK</string>
  </array>
  <key>RunAtLoad</key>
  <true/>
  <key>KeepAlive</key>
  <dict>
    <key>SuccessfulExit</key>
    <false/>
  </dict>
  <key>ExitTimeOut</key>
  <integer>300</integer>
  <key>StandardErrorPath</key>
  <string>$HOME/Library/Logs/core-geth-$NETWORK.log</string>
</dict>
</plist>
EOF
```

**For a configuration**, add each of its flags to `ProgramArguments` as its own `<string>` line
after `<string>--$NETWORK</string>`, joining a flag to its value with `=`:
`<string>--http.api=eth,net,web3</string>`.

**Start it now, then follow its log.** It starts again each time you log in:

```bash
launchctl bootstrap gui/$(id -u) ~/Library/LaunchAgents/com.coregeth.$NETWORK.plist
tail -f ~/Library/Logs/core-geth-$NETWORK.log
```

What the agent's settings do:

- **`RunAtLoad`** starts the node when the agent loads: now, and at each login.
- **`KeepAlive`**, with `SuccessfulExit` set to false, restarts the node after it exits with an error.
  A node that stops cleanly exits without one, and launchd leaves it stopped.
- **`ExitTimeOut`** is how long launchd waits after its SIGTERM before it kills the node. The system
  default is shorter than a stop during the first sync can take
  ([stop times](../hardware-requirements.md#stopping)).
- **`StandardErrorPath`** is where the log goes, because the node logs to standard error.

The agent runs the node as you, with the default data directory, so the commands in steps 4 and 5
reach it unchanged. To run a node with no one logged in, install the same file as a launch daemon in
`/Library/LaunchDaemons` instead, with a `UserName` key naming the account to run it as;
`man launchd.plist` documents both.

## 7. Stop it cleanly

- **In the foreground:** press Ctrl-C once, then wait for `Blockchain stopped` and the prompt.
- **As an agent:** with `NETWORK` set as in step 6, unload it. launchd sends the node SIGTERM, and the
  log then ends with `Blockchain stopped`. The agent loads again at your next login; delete its file to
  keep it from starting.

```bash
launchctl bootout gui/$(id -u)/com.coregeth.$NETWORK
rm ~/Library/LaunchAgents/com.coregeth.$NETWORK.plist      # only to keep it from starting at login
```

Pressing Ctrl-C again does not speed up a stop. [Stop it safely](../run-classic-node.md#stop-it-safely)
shows what a clean stop logs, and why a stop during the first sync takes longer.

## 8. Update it

Stop the node, install the new release's `geth` the same way as in step 1, and start the node again.
[Upgrading within v1.13.x](../../operate/maintenance.md#upgrading-within-v113x) lists what to read
and check. A node on a v1.12.x release follows the
[Mac migration guide](../../tutorials/migration/macos.md) instead.
