# Tutorial: Migrating from `geth` to `core-geth`

The Core-Geth executable is now named `core-geth`. It was previously named `geth`, which
collided with go-ethereum's binary of the same name on any machine running both.

**No data migration is required, and no resync.** Your chaindata, node key and keystore
stay exactly where they are. This procedure changes how the client is *invoked*, nothing
about what it reads.

Budget five minutes. The node has to be stopped only for the moment it takes to restart it
under the new name.

### What changed, and what deliberately did not

Only the executable was renamed. The client's on-disk and on-wire identities are set from
constants in the source, not from the executable's filename, so they are unaffected.

| Identity | Before | After |
|---|---|---|
| Executable | `geth` | **`core-geth`** |
| Datadir instance directory | `<datadir>/geth/` | `<datadir>/geth/` — unchanged |
| Chaindata | `<datadir>/geth/chaindata` | unchanged |
| Node key | `<datadir>/geth/nodekey` | unchanged |
| IPC socket | `<datadir>/geth.ipc` | unchanged |
| Default datadir | `~/.ethereum` | unchanged |
| devp2p advertisement | `CoreGeth` | unchanged |

**Do not rename anything in the second column to match the first row.** The client looks
for its instance directory under `geth` regardless of what the executable is called. A node
pointed at `<datadir>/core-geth/` finds no chaindata and resyncs from genesis.

If you have already renamed a datadir in anticipation, move it back:

```bash
mv <datadir>/core-geth <datadir>/geth
```

### Before you start

```bash
# Confirm which client your current `geth` actually is.
# Core-Geth prints "Core-Geth <version>"; go-ethereum prints "Geth".
geth version | head -1
```

If that prints `Geth`, you are looking at go-ethereum and this guide does not apply to it.
If both are installed, the collision this rename removes is exactly your situation.

### 1. Install the renamed binary

From a release archive, the binary inside is now `core-geth`:

```bash
unzip core-geth-linux-v1.13.0.zip
sudo install -m 0755 core-geth /usr/local/bin/core-geth
```

Or from source:

```bash
make core-geth              # writes ./build/bin/core-geth
sudo install -m 0755 ./build/bin/core-geth /usr/local/bin/core-geth
```

### 2. Find everything that invokes `geth`

The repository ships a script that reports every reference and, on request, rewrites it:

```bash
# report only -- writes nothing
./build/migrate-geth-to-core-geth.sh

# rewrite, leaving a .bak beside each edited file
./build/migrate-geth-to-core-geth.sh --apply
```

It scans systemd units, shell scripts, compose files and your crontab. It rewrites a
reference only where the position makes it unambiguously an executable: the first token
after `ExecStart=`, a path under a `bin/` or `sbin/` directory, or a bare `geth` word.
**Anything else that mentions `geth` — a `--datadir` argument, an `ipcpath`, a log
directory — is reported for you to check, never rewritten**, because an absolute path
ending in `/geth` is far more often a datadir than a binary.

Typical systemd change:

```ini
# before
ExecStart=/usr/local/bin/geth --classic --datadir /var/lib/etc/geth

# after -- note the datadir is untouched
ExecStart=/usr/local/bin/core-geth --classic --datadir /var/lib/etc/geth
```

```bash
sudo systemctl daemon-reload
```

### 3. Remove the old binary

Only once nothing references it:

```bash
sudo rm /usr/local/bin/geth        # if it was Core-Geth
```

If your `geth` is go-ethereum, leave it. Keeping both is the point of the rename.

### 4. Start, and verify

```bash
sudo systemctl start core-geth     # or whatever your unit is called

core-geth attach --exec 'eth.blockNumber' <datadir>/geth.ipc
```

Two things to confirm, and the second is the one that matters:

```bash
# the client identifies itself and its version on one line:
#   Core-Geth 1.13.0-unstable
core-geth version | head -1

# it resumed from your existing chain data rather than starting over.
# A number near chain head means the datadir was found.
core-geth attach --exec 'eth.syncing' <datadir>/geth.ipc
```

If `eth.blockNumber` returns `0` and the node begins syncing from genesis, stop it
immediately — it is pointed at the wrong datadir. Check that `--datadir` still names the
same directory as before and that `<datadir>/geth/chaindata` exists.

### Docker

The image entrypoint is now `core-geth`:

```bash
# before
docker run etclabscore/core-geth geth --classic

# after
docker run etclabscore/core-geth core-geth --classic
```

A bare `docker run <image>` with no command needs no change; the entrypoint handles it.
Volume mounts are unaffected.

### For AI agents performing this migration

This section is the executable form of the procedure above. Follow it literally.

**Invariants. Violating any of these causes a full resync from genesis.**

1. **Never rename a datadir, or any directory or file under one.** `<datadir>/geth/`,
   `chaindata`, `nodekey`, `nodes` and `geth.ipc` keep their names permanently. The
   rename applies to the executable only.
2. **Never rewrite the value of `--datadir`, `--ipcpath`, `--keystore`, or any argument
   that names a path.** Only the executable token changes.
3. **Never rewrite a path merely because it ends in `/geth`.** Decide by position:
   is this token being *executed*, or is it an *argument*? Only the former changes.
4. **Do not add a `geth` symlink** to restore the old name. That re-creates the collision
   this rename exists to remove.

**Procedure.**

```bash
# 1. Identify. Core-Geth prints "Core-Geth <version>"; go-ethereum prints "Geth". Do not touch go-ethereum.
geth version | head -1

# 2. Record the pre-migration block height, to prove later that data survived.
geth attach --exec 'eth.blockNumber' <datadir>/geth.ipc

# 3. Report every reference before changing anything.
./build/migrate-geth-to-core-geth.sh

# 4. Apply, then read the NOT REWRITTEN section of the output and resolve each entry by hand.
./build/migrate-geth-to-core-geth.sh --apply

# 5. Reload and restart.
sudo systemctl daemon-reload && sudo systemctl restart <unit>

# 6. Verify data survived: this must be at or above the height recorded in step 2.
core-geth attach --exec 'eth.blockNumber' <datadir>/geth.ipc
```

**Stop and report to a human if:** step 6 returns a height far below step 2's; the script
reports a reference you cannot classify as executable or data; or `geth version` prints
neither `Core` nor `Geth`.

**Do not** attempt to "complete" the rename by making the datadir consistent with the
binary name. The asymmetry is deliberate and is documented in `AGENTS.md`.

### Rolling back

The rename carries no data change, so a rollback is just restoring the old invocation:

```bash
# restore edited config from the backups the script left
mv /etc/systemd/system/<unit>.service.bak /etc/systemd/system/<unit>.service
sudo systemctl daemon-reload

# put a binary named geth back
sudo install -m 0755 ./build/bin/core-geth /usr/local/bin/geth
```

Your datadir never changed, so there is nothing to restore there.

### Related

- [Migrating to v1.13.0](v1.13.0-migration.md) — the release migration this rename ships in
- [Build from source](../developers/build-from-source.md)
