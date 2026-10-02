---
title: Troubleshooting
---

Each section starts from what you see and ends with what the node does once the problem is fixed.
Where another page owns a procedure, the section links to it.

A node started in the foreground logs to the terminal. A node run as a systemd service logs to the
system journal ([Linux users guide](../getting-started/run/linux.md#6-keep-it-running-a-systemd-service)), and a node in Docker to `docker logs`. The
commands below pass `--classic`; on a Mordor node, pass `--mordor` instead.

## Starting the node

### Why does `geth` refuse `--ethereum`, `--sepolia` or `--holesky`?

```
Fatal: --ethereum is deprecated: this client implements Ethereum upgrades only through Cancun, so it cannot follow Ethereum or its test networks. Use --classic, the default, or --mordor
```

- **Why:** this client cannot follow Ethereum or its test networks, so each of these flags stops it
  before it opens a data directory
  ([Networks this client does not support](../index.md#networks-this-client-does-not-support)).
- **Fix:** start the node with `--classic` or `--mordor`.
- **Once fixed:** the first log line names the network, as `Starting Core-Geth on Ethereum Classic...`
  or `Starting Core-Geth on Mordor testnet...`.

### Why does the node refuse its data directory, or start on the wrong network?

A data directory holds the chain of the network it was created for. Started with another network's
flag, the node stops:

```
Fatal: Failed to register the Ethereum service: database contains incompatible genesis (have a68ebde7932eccb177d38d55dcc6461a019dd795a681e59b5a3e4f3a7259a3f1, new d4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3)
```

With no network flag, a node started on a Mordor data directory names the wrong network:

```
INFO [..] Starting Core-Geth on Ethereum Classic...
INFO [..] Initialising Ethereum protocol           network=63 dbversion=8
```

- **Why:** `have` is the genesis block the data directory holds, and `new` is the one the flag
  selects. With no flag, the node does not recognize a Mordor data directory, and Mordor nodes refuse
  to peer with it
  ([Pass `--mordor` on every command](../getting-started/run-mordor-node.md#pass-mordor-on-every-command)).
- **Check:** the first log line names the network the node is running.
- **Fix:** start the node with the network flag the data directory was created with.
- **Once fixed:** a Mordor node logs `Starting Core-Geth on Mordor testnet...` and `network=7`.

### Why did the node start a new sync instead of using my chain?

```
INFO [..] Allocated cache and file handles         database=<datadir>/geth/chaindata cache=<size> handles=<count>
INFO [..] Initialising Ethereum protocol           network=1 dbversion=<nil>
INFO [..] Writing custom genesis block
```

- **Why:** the node opened a data directory that holds no chain and created a new database in it.
  `database=` names the directory it opened, and `dbversion=<nil>` means the database is new. That
  happens when:
    - a node upgraded from v1.12.x starts with neither a network flag nor `--datadir`
      ([If you start the node without a network flag](../tutorials/v1.13.0-migration.md#if-you-start-the-node-without-a-network-flag));
    - `--datadir` names a different path from the one that holds the chain, for example when a
      service unit and a foreground run use different paths;
    - a Mordor node starts without `--mordor` or `--datadir`, and opens Ethereum Classic's default
      directory instead of Mordor's.
- **Check:** compare the path after `database=` with the directory that holds your chain.
- **Fix:** stop the node ([Stop it safely](../getting-started/run-classic-node.md#stop-it-safely)),
  then start it with the network flag and `--datadir` that match the chain.
- **Once fixed:** the start logs no `Writing custom genesis block`, and `Loaded most recent local block`
  gives your chain's head:

    ```
    INFO [..] Loaded most recent local block           number=<block> hash=<hash> td=<td> age=<age>
    ```

### Why does the node stop with `incompatible state scheme`?

```
Fatal: Failed to register the Ethereum service: incompatible state scheme, stored: hash, provided: path
```

- **Why:** `--state.scheme` names a different scheme from the one the data directory was created
  with, and a data directory keeps its scheme ([State scheme](sync-modes.md#state-scheme)).
- **Check:** `stored` is the data directory's scheme, and `provided` is the one the flag asks for.
- **Fix:** remove `--state.scheme` from the command or the service unit. Using the other scheme takes
  a new data directory.
- **Once fixed:** the node logs the stored scheme:

    ```
    INFO [..] State scheme set to already existing     scheme=hash
    ```

### Why does the node stop with `invalid peer config: light peer count`?

```
Fatal: Error starting protocol stack: invalid peer config: light peer count (100) >= total peer count (50)
```

- **Why:** the configuration file sets `LightServ` above 0. The client has no light server, but it
  still compares `LightPeers` with the peer limit and stops when `LightPeers` is at least that limit.
  With a smaller `LightPeers`, the node starts, and accepts that many fewer peers.
- **Check:** the file has a `LightServ` line under `[Eth]`.
- **Fix:** delete the `LightServ` and `LightPeers` lines from the file.
- **Once fixed:** the node starts and keeps running.

### Why does the node stop with `bind: address already in use`?

```
Fatal: Error starting protocol stack: listen tcp :<port>: bind: address already in use
```

For the HTTP or WebSocket port, an error line comes first:

```
ERROR[..] Failed to open RPC endpoints             error="listen tcp 127.0.0.1:<port>: bind: address already in use"
Fatal: Error starting protocol stack: listen tcp 127.0.0.1:<port>: bind: address already in use
```

- **Why:** another process on the host holds the port, often a second node using the same peer, HTTP
  or WebSocket port.
- **Check:** on Linux, `ss -lntup` lists the process that holds each port
  ([Ports and listeners](security.md#ports-and-listeners)).
- **Fix:** give each node its own ports.
  [Mordor beside Ethereum Classic on one host](../getting-started/run-mordor-node.md#mordor-beside-ethereum-classic-on-one-host)
  lists the flag for each listener.
- **Once fixed:** the node starts, and `ss -lntup` shows it on its own ports.

### Why does `geth` say the data directory is already in use?

```
Fatal: Failed to create the protocol stack: datadir already used by another process
```

- **Why:** another `geth` process has the data directory open. It may be a node that is still
  running or still stopping, or a second node or database command given the same directory.
- **Check:** on Linux, `pgrep -a geth` lists every `geth` process with its command line.
- **Fix:** stop the other process and wait for it to exit
  ([Stop it safely](../getting-started/run-classic-node.md#stop-it-safely)). Each node needs its own
  data directory.
- **Once fixed:** the command runs.

### Why does the node stop with `ancient chain segments already extracted`?

```
Fatal: Failed to register the Ethereum service: ancient chain segments already extracted, please set --datadir.ancient to the correct path
```

- **Why:** the node's ancient store is in the directory `--datadir.ancient` names, and this start did
  not pass the flag
  ([The `ancient` store on another disk](../getting-started/hardware-requirements.md#the-ancient-store-on-another-disk)).
- **Fix:** pass the same `--datadir.ancient` on every start, in the service unit as well.
  [History on another disk](maintenance.md#history-on-another-disk) covers the empty store the refused
  start leaves behind.
- **Once fixed:** the node logs where it opened the store:

    ```
    INFO [..] Opened ancient database                  database=<ancient dir>/chain readonly=false
    ```

### Why does the node hang at start after `removedb`?

The log stops after these lines, and the node never opens its IPC endpoint:

```
WARN [..] Empty database, resetting chain
WARN [..] Rewinding blockchain to block            target=0
WARN [..] Empty database, resetting chain
```

- **Why:** `removedb` deleted the state data and kept the ancient chain, and a data directory in that
  state does not start ([Resync from scratch](maintenance.md#resync-from-scratch)).
- **Fix:** stop the node, run `removedb` again, and answer `y` to both parts.
- **Once fixed:** the next start logs `Writing custom genesis block`, and the node syncs the chain from
  its peers.

## Peers

### Why does the node have no peers?

[Bootnodes and peer discovery](../guides/bootnodes-and-discovery.md) explains where a node finds peers.

`net.peerCount` returns `0`
([How to tell it is done](../getting-started/run-classic-node.md#how-to-tell-it-is-done)). Check these
causes in order.

1. **A peer limit of 0.** The `Maximum peer count` line reads `ETH=0 total=0` when `--maxpeers 0` is
   given, or when the configuration file sets `MaxPeers = 0`:

    ```
    INFO [..] Maximum peer count                       ETH=0 total=0
    ```

    Remove the setting, or set a limit above 0.

2. **Discovery turned off.** With `--nodiscover`, the node finds no peers of its own, and the enode in
   `Started P2P networking` ends in `discport=0`:

    ```
    INFO [..] Started P2P networking                   self="enode://<node id>@<IP>:<port>?discport=0"
    ```

    Remove the flag, or list the peers to connect to in `StaticNodes`
    ([Why does the node ignore `static-nodes.json`?](#why-does-the-node-ignore-static-nodesjson)).

3. **The wrong network.** The first log line names the network the node is running
   ([Why does the node refuse its data directory, or start on the wrong network?](#why-does-the-node-refuse-its-data-directory-or-start-on-the-wrong-network)).

4. **Blocked ports.** Other nodes reach this one over TCP for connections and UDP for discovery
   ([Ports and listeners](security.md#ports-and-listeners)). Let both through the host's firewall
   ([A host firewall](security.md#a-host-firewall)), and forward both on any router in front of the
   host. The node opens at most a third of its `--maxpeers` connections itself, and the rest have to
   come from other nodes, so a peer count that stops at that third usually means other nodes cannot
   reach it.

5. **NAT without port mapping.** Behind a router, the default `--nat any` asks the router to forward
   the ports over UPnP or NAT-PMP ([Ports and listeners](security.md#ports-and-listeners)). Where the
   router supports neither, forward the ports yourself and give the node its public address with
   `--nat extip:<public IP>`. The enode in `Started P2P networking` then carries that address:

    ```
    INFO [..] Started P2P networking                   self=enode://<node id>@<public IP>:<port>
    ```

6. **The clock.** Discovery fails between two nodes whose clocks are too far apart
   ([Clock](security.md#clock)).

7. **A configuration file from another network.** The file's `BootstrapNodes` replace the network's
   own ([Configuration](../getting-started/run-cli.md#configuration)).

**Once fixed:** `net.peerCount` returns a number above 0.

### Why did nodes that list this one stop connecting after its key changed?

Once this node's key has changed, nodes that name it in `StaticNodes` no longer connect to it, and
neither node logs the failed connection at the default verbosity.

- **Why:** the enode ID comes from the node key ([The node key](security.md#the-node-key)). A node
  that dials the old ID reaches a node holding a different key, and the connection fails.
- **Check:** read this node's enode, and compare it with the entry on the other node:

    ```shell
    $ geth --classic attach --exec 'admin.nodeInfo.enode' <datadir>/geth.ipc
    ```

- **Fix:** give the new enode to every node that lists this one, in `StaticNodes` or `TrustedNodes`.
  [3. Rotate the P2P node key](../tutorials/v1.13.0-migration.md#rotate-the-p2p-node-key) lists the
  other places that name a node by its enode.
- **Once fixed:** on the other node, `admin.peers` includes the new enode:

    ```shell
    $ geth --classic attach --exec 'admin.peers.map(function(p){return p.enode})' <datadir>/geth.ipc
    ```

### Why does the node ignore `static-nodes.json`?

```
ERROR[..] The static-nodes.json file is deprecated and ignored. Use P2P.StaticNodes in config.toml instead.
ERROR[..] The trusted-nodes.json file is deprecated and ignored. Use P2P.TrustedNodes in config.toml instead.
```

- **Why:** the node reads static and trusted peers only from its configuration file. It logs these
  lines when it finds the old files in `<datadir>/geth/`, and does not read them.
- **Fix:** move each enode to the `StaticNodes` or `TrustedNodes` line under `[Node.P2P]`, in a
  configuration file written by `dumpconfig`
  ([Configuration](../getting-started/run-cli.md#configuration)), and start the node with `--config`.
  The lines belong under `[Node.P2P]`, and a `[P2P]` table of its own stops the node:

    ```
    Fatal: <config file>, line <line>: field 'P2P' is not defined in main.gethConfig
    ```

    Nothing reads the two JSON files, and removing them stops the error lines.

- **Once fixed:** the two error lines are gone, and `admin.peers` lists the static peers.

## Syncing

### Are the `WARN` lines during the first sync a problem?

```
WARN [..] Failed to load snapshot                  err="missing or corrupted snapshot"
WARN [..] Pivot seemingly stale, moving            old=<block> new=<block>
```

No. A first sync logs lines like these, and none of them stops it.
[Warnings a first sync logs](../getting-started/run-classic-node.md#warnings-a-first-sync-logs) says
why each one appears.

- **Check:** `Syncing:` progress lines go on appearing between them.

### Is the sync stuck?

A first sync can look stalled: its `eta` swings, `state healing` lines return, or
`Pivot seemingly stale, moving` appears more than once.

- **Why:** the chain moves on while the state downloads, so the node moves its sync target forward and
  heals the state that changed
  ([What a healthy first sync looks like](../getting-started/run-classic-node.md#what-a-healthy-first-sync-looks-like)).
  The chain download can also go on after state healing reaches `pending=0`.
- **Check:** read `eth.syncing` twice, a while apart
  ([How to tell it is done](../getting-started/run-classic-node.md#how-to-tell-it-is-done)). The sync
  is moving if `currentBlock` or any of the `synced` and `healed` counters has risen, or if `synced=`
  has risen in the `Syncing: chain download in progress` lines.
- **Fix:** if nothing has risen, look for
  [no peers](#why-does-the-node-have-no-peers) and
  [low disk space](#why-did-the-node-shut-itself-down-with-low-disk-space).
- **Once fixed:** the log shows `Snap sync complete, auto disabling`, and `eth.syncing` returns
  `false`.

### What do `Disabled artificial finality features` and `Reorg disallowed` mean?

`Disabled artificial finality features` means MESS
([MESS on this node](../getting-started/run-classic-node.md#mess-on-this-node)) has switched itself
off. The line's `reason` gives the cause: `low peers` when the node has fewer peers than MESS needs,
and `stale safety interval` when the node's newest block is older than MESS allows.
`Reorg disallowed` means MESS refused a reorganization to another chain.

- **Check:** while MESS is off, `admin.messStatus()` reports `nodeSwitch` as `false`.
- **Fix:** bring back the node's [peers](#why-does-the-node-have-no-peers) or its
  [sync](#is-the-sync-stuck). MESS switches itself back on once the node is in sync with enough peers.
  [MESS](mess.md) describes both conditions.
  After `Reorg disallowed`, compare the node's head with the network's, as step 4 of
  [How to tell it is done](../getting-started/run-classic-node.md#how-to-tell-it-is-done) does.
- **Once fixed:** the log shows MESS switching back on:

    ```
    INFO [..] Enabled artificial finality features     reason=synced peers=<count>
    ```

## Disk and memory

### Why did the node shut itself down with `Low disk space`?

```
ERROR[..] Low disk space. Gracefully shutting down Geth to prevent database corruption. available=<size> path=<datadir>/geth
INFO [..] Got interrupt, shutting down...
```

- **Why:** the free space in `<datadir>/geth` fell below the node's threshold, and the node stopped
  cleanly to protect its database
  ([Running out of disk](../getting-started/hardware-requirements.md#running-out-of-disk)). A service
  unit with `Restart=on-failure` leaves it stopped.
- **Check:** `df -h <datadir>` shows the space left on that file system.
- **Fix:** free space on that file system, or add space to it.
  [History on another disk](maintenance.md#history-on-another-disk) moves the ancient store off it.
  Then start the node.
- **Once fixed:** the node runs without logging `Disk space is running low`.

### Why was the node killed for memory?

The log ends without `Got interrupt, shutting down...`. For a node run as a service, the journal
shows:

```
core-geth.service: A process of this unit has been killed by the OOM killer.
core-geth.service: Main process exited, code=killed, status=9/KILL
```

- **Why:** the machine ran short of memory, and the kernel's out-of-memory killer ended the node.
- **Check:**

    ```shell
    $ systemctl show core-geth -p Result
    Result=oom-kill
    ```

- **Fix:** lower `--cache`, the memory the node gives its caches
  ([Memory and `--cache`](../getting-started/hardware-requirements.md#memory-and-cache)), or free
  memory that other programs on the machine use. `geth account new` can be killed the same way on a
  small machine ([Key derivation on a small machine](../getting-started/run-cli.md#key-derivation-on-a-small-machine)).
- **Once fixed:** started again, the node keeps running, and the same command prints `Result=success`.

## JSON-RPC

### Why does a request for the `finalized` or `safe` block fail?

```shell
$ curl -s -X POST -H 'Content-Type: application/json' \
    --data '{"jsonrpc":"2.0","id":1,"method":"eth_getBlockByNumber","params":["finalized",false]}' \
    http://127.0.0.1:8545
{"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":"finalized block not found"}}
```

A request for `safe` fails the same way, with `safe block not found`.

- **Why:** a node on Ethereum Classic or Mordor never has a finalized or safe block. The client sets
  them only from the Engine API's forkchoice update, which only an Ethereum consensus client sends.
  Any request that names either tag fails, `eth_getBalance` among them, and `eth_getLogs` fails with
  `finalized header not found` or `safe header not found`.
- **Fix:** have the application ask for `latest` or a block number instead. A gateway in front of
  the node judges finality by block depth
  ([What a gateway checks](../guides/public-rpc-endpoint.md#what-a-gateway-checks)).
- **Once fixed:** the request returns a block.

## Stopping the node

### Why doesn't Ctrl-C stop a node started with `geth console`?

- **Why:** at the console prompt, Ctrl-C discards the line being typed and leaves the node running
  ([A node on an Ethereum Classic network](../getting-started/run-cli.md#a-node-on-an-ethereum-classic-network)).
- **Fix:** type `exit` or press Ctrl-D. The node stops with the console.
- **Once fixed:** the log ends with `Blockchain stopped`.

### Why does pressing Ctrl-C again log `Already shutting down`?

```
WARN [..] Already shutting down, interrupt more to panic. times=9
```

- **Why:** the node is already stopping, and each further interrupt counts toward a panic, which is
  an unclean stop ([Stop it safely](../getting-started/run-classic-node.md#stop-it-safely)).
- **Fix:** send no more interrupts, and wait.
- **Once fixed:** the log shows `Blockchain stopped`, and then the process exits.

### Why does a stop during the first sync log `Failed to journal state snapshot`?

```
ERROR[..] Failed to journal state snapshot         err="snapshot [0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421] missing"
```

The error is expected when a node stops before its first sync finishes. Started again, the node
resumes the sync ([Stop it safely](../getting-started/run-classic-node.md#stop-it-safely)).

### Why does every start log `Unclean shutdown detected`?

```
WARN [..] Unclean shutdown detected                booted=<time> age=<age>
```

- **Why:** the node keeps a record of its recent unclean stops, and logs this line for each one on
  every start, including a start after a clean stop
  ([After an unclean stop](maintenance.md#after-an-unclean-stop)).
- **Check:** compare `booted` with the time the previous run started. Only a line whose `booted`
  matches that run is new.
- **Fix:** none for an old line. A new line means the previous run was killed or crashed, and
  [After an unclean stop](maintenance.md#after-an-unclean-stop) says what the node does about it.

## Collecting diagnostics

Before you open an issue, collect these:

1. **The release.** `geth version` prints the release, the commit it was built from, the Go version
   and the platform:

    ```shell
    $ geth version
    ```

2. **How the node runs:** the network, the command line or service unit, and the configuration file
   if there is one, with passwords and tokens removed. Say whether the data directory came from a
   v1.12.x release.

3. **The log, as text,** from the start of the run through the problem. `--log.file` also writes the
   log to a file, and `--verbosity 4` adds debug lines ([Logs](monitoring.md#logs)):

    ```shell
    $ geth --classic --datadir <datadir> --log.file geth.log --verbosity 4
    ```

    The log shows the node's IP address and enode ID, in its `New local node record` and
    `Started P2P networking` lines. Remove them first if you do not want them public.

4. **For a node that hangs, a dump of its goroutines.** Start the node with `--pprof`
   ([Profiling](monitoring.md#profiling)), and save the dump while it hangs:

    ```shell
    $ geth --classic --datadir <datadir> --pprof
    $ curl -s -o goroutines.txt 'http://127.0.0.1:6060/debug/pprof/goroutine?debug=2'
    ```

5. **What you expected, what happened, and the steps that reproduce it.**

Open the issue on [`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth/issues).
To report a vulnerability, follow
[`SECURITY.md`](https://github.com/ethereumclassic/core-geth/blob/main/SECURITY.md) instead.
`geth version-check` does not check this client for vulnerabilities
([Reporting a vulnerability](security.md#reporting-a-vulnerability)).
