package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/internal/cmdtest"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/params"
)

// TestConfigFileNetworkDefaults checks that a network ID set in a configuration file chooses
// the network's defaults the way --networkid does. With no network flag, a private network
// defined in a file once took Ethereum Classic's bootnodes, discovery trees, gas ceiling,
// Ethash directories and data directory, and so announced itself to Ethereum Classic's
// discovery network.
func TestConfigFileNetworkDefaults(t *testing.T) {
	// The data directory is left at its default, which is under the home directory, so the
	// home directory is a temporary one.
	home := t.TempDir()
	t.Setenv("HOME", home)
	base := filepath.Join(home, ".ethereum")

	urls := func(nodes []*enode.Node) (out []string) {
		for _, n := range nodes {
			out = append(out, n.URLv4())
		}
		return out
	}
	var classicBootnodes []string
	for _, url := range params.ClassicBootnodes {
		classicBootnodes = append(classicBootnodes, enode.MustParse(url).URLv4())
	}
	classicTrees := []string{params.ClassicDNSNetwork1, params.ClassicDNSNetwork2, params.ClassicDNSNetwork3}

	for _, c := range []struct {
		networkID string // NetworkId in the file's [Eth] section, "" for none
		args      []string
		classic   bool
	}{
		{"", nil, true},
		{"1", nil, true},
		{"12345", nil, false},
		{"7", nil, false}, // Mordor's network ID: only --mordor selects Mordor
		{"", []string{"--networkid", "12345"}, false},
		{"12345", []string{"--networkid", "1"}, true}, // the flag wins over the file
		{"1", []string{"--networkid", "12345"}, false},
		{"12345", []string{"--classic"}, true},
	} {
		file := filepath.Join(t.TempDir(), "config.toml")
		toml := "[Eth]\n"
		if c.networkID != "" {
			toml += "NetworkId = " + c.networkID + "\n"
		}
		if err := os.WriteFile(file, []byte(toml), 0o600); err != nil {
			t.Fatal(err)
		}
		out := filepath.Join(t.TempDir(), "dump.toml")
		args := append(slices.Clone(c.args), "--config", file, "dumpconfig", out)

		// runGeth would add a --datadir, so start the command directly.
		geth := &testgeth{}
		geth.TestCmd = cmdtest.NewTestCmd(t, geth)
		geth.Run("core-geth-test", args...)
		geth.WaitExit()
		if status := geth.ExitStatus(); status != 0 {
			t.Fatalf("NetworkId %q, %v: exit status %d\n%s", c.networkID, c.args, status, geth.StderrText())
		}
		var cfg gethConfig
		if err := loadConfig(out, &cfg); err != nil {
			t.Fatal(err)
		}

		name := "NetworkId " + c.networkID
		if c.networkID == "" {
			name = "no NetworkId"
		}
		if c.classic {
			if want := filepath.Join(base, "classic"); cfg.Node.DataDir != want {
				t.Errorf("%s, %v: datadir %s, want %s", name, c.args, cfg.Node.DataDir, want)
			}
			if have := urls(cfg.Node.P2P.BootstrapNodes); !slices.Equal(have, classicBootnodes) {
				t.Errorf("%s, %v: discv4 bootnodes %v, want Ethereum Classic's", name, c.args, have)
			}
			if have := urls(cfg.Node.P2P.BootstrapNodesV5); !slices.Equal(have, classicBootnodes) {
				t.Errorf("%s, %v: discv5 bootnodes %v, want Ethereum Classic's", name, c.args, have)
			}
			if !slices.Equal(cfg.Eth.EthDiscoveryURLs, classicTrees) {
				t.Errorf("%s, %v: discovery trees %v, want Ethereum Classic's", name, c.args, cfg.Eth.EthDiscoveryURLs)
			}
			if cfg.Eth.Miner.GasCeil != 8_000_000 {
				t.Errorf("%s, %v: gas ceiling %d, want 8000000", name, c.args, cfg.Eth.Miner.GasCeil)
			}
			if cfg.Eth.Ethash.CacheDir != "etchash" {
				t.Errorf("%s, %v: ethash cache dir %q, want etchash", name, c.args, cfg.Eth.Ethash.CacheDir)
			}
			continue
		}
		if cfg.Node.DataDir != base {
			t.Errorf("%s, %v: datadir %s, want %s", name, c.args, cfg.Node.DataDir, base)
		}
		if len(cfg.Node.P2P.BootstrapNodes) != 0 || len(cfg.Node.P2P.BootstrapNodesV5) != 0 {
			t.Errorf("%s, %v: bootnodes %v and %v, want none", name, c.args, urls(cfg.Node.P2P.BootstrapNodes), urls(cfg.Node.P2P.BootstrapNodesV5))
		}
		if len(cfg.Eth.EthDiscoveryURLs) != 0 {
			t.Errorf("%s, %v: discovery trees %v, want none", name, c.args, cfg.Eth.EthDiscoveryURLs)
		}
		if cfg.Eth.Miner.GasCeil != ethconfig.Defaults.Miner.GasCeil {
			t.Errorf("%s, %v: gas ceiling %d, want the generic %d", name, c.args, cfg.Eth.Miner.GasCeil, ethconfig.Defaults.Miner.GasCeil)
		}
		if cfg.Eth.Ethash.CacheDir != ethconfig.Defaults.Ethash.CacheDir {
			t.Errorf("%s, %v: ethash cache dir %q, want %q", name, c.args, cfg.Eth.Ethash.CacheDir, ethconfig.Defaults.Ethash.CacheDir)
		}
	}
}

// TestNetworkBanner checks that the startup banner names the network the node took its defaults
// for, and only that one: Ethereum Classic exactly when IsClassic says so, a network ID set in a
// configuration file and a network flag given as false included. The banner once read only the
// command line, before the file was loaded, and tested whether a network flag was given rather
// than its value: it announced Ethereum Classic for a private network defined in a file, Mordor
// for --mordor=false, and nothing for --networkid 1.
func TestNetworkBanner(t *testing.T) {
	file := filepath.Join(t.TempDir(), "private.toml")
	if err := os.WriteFile(file, []byte("[Eth]\nNetworkId = 12345\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		args    []string
		network string // the network the banner names, "" for none
	}{
		{nil, "Ethereum Classic"},
		{[]string{"--classic"}, "Ethereum Classic"},
		{[]string{"--networkid", "1"}, "Ethereum Classic"},
		{[]string{"--classic=false"}, "Ethereum Classic"},
		{[]string{"--mainnet=false"}, "Ethereum Classic"},
		{[]string{"--mordor=false"}, "Ethereum Classic"},
		{[]string{"--networkid", "12345"}, ""},
		{[]string{"--config", file}, ""},
		{[]string{"--mordor=false", "--networkid", "12345"}, ""},
	} {
		// With no peers and no listener, the node starts, prints its name and exits.
		geth := runGeth(t, slices.Concat(c.args, []string{"--nodiscover", "--maxpeers", "0", "--port", "0",
			"--nat", "none", "--exec", "admin.nodeInfo.name", "console"})...)
		geth.KillTimeout = 30 * time.Second
		geth.Output()
		geth.WaitExit()
		stderr := geth.StderrText()
		if status := geth.ExitStatus(); status != 0 {
			t.Fatalf("%v: exit status %d\n%s", c.args, status, stderr)
		}
		var named, want []string
		for _, line := range strings.Split(stderr, "\n") {
			if _, network, ok := strings.Cut(line, "Starting Core-Geth on "); ok {
				named = append(named, strings.TrimSuffix(strings.TrimSpace(network), "..."))
			}
		}
		if c.network != "" {
			want = []string{c.network}
		}
		if !slices.Equal(named, want) {
			t.Errorf("%v: the banner names %q, want %q", c.args, named, want)
		}
	}
}

// TestEthereumDatabaseRefused checks that geth refuses a data directory holding Ethereum
// mainnet's chain when no network flag chooses a genesis, both to start a node and to export
// blocks, as it refuses --ethereum. Ethereum mainnet shares Ethereum Classic's genesis, so geth
// once ran Ethereum's rules on such a directory, which this client follows only through the
// Cancun upgrade, with Ethereum Classic's bootnodes.
func TestEthereumDatabaseRefused(t *testing.T) {
	datadir := t.TempDir()
	genesis := filepath.Join(t.TempDir(), "ethereum.json")
	b, err := json.Marshal(params.DefaultGenesisBlock())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(genesis, b, 0o600); err != nil {
		t.Fatal(err)
	}
	setup := runGeth(t, "--datadir", datadir, "init", genesis)
	setup.WaitExit()
	if status := setup.ExitStatus(); status != 0 {
		t.Fatalf("init: exit status %d\n%s", status, setup.StderrText())
	}
	console := []string{"--nodiscover", "--maxpeers", "0", "--port", "0", "--nat", "none", "--exec", "admin.nodeInfo.name", "console"}
	for _, args := range [][]string{
		console,
		slices.Concat([]string{"--networkid", "12345"}, console),
		{"export", filepath.Join(t.TempDir(), "chain.rlp")},
	} {
		geth := runGeth(t, slices.Concat([]string{"--datadir", datadir}, args)...)
		geth.KillTimeout = 30 * time.Second
		geth.Output()
		geth.WaitExit()
		if status, stderr := geth.ExitStatus(), geth.StderrText(); status == 0 || !strings.Contains(stderr, "holds Ethereum mainnet's chain") {
			t.Errorf("%v: exit status %d, want the directory refused\n%s", args, status, stderr)
		}
	}
}
