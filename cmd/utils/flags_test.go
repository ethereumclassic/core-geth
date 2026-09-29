// Copyright 2019 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

// Package utils contains internal helper functions for go-ethereum commands.
package utils

import (
	"flag"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/miner"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/vars"
	"github.com/urfave/cli/v2"
)

// networkContext returns a command-line context holding the network flags,
// parsed from args, and the flags SetNodeConfig and SetEthConfig need a value for.
func networkContext(t *testing.T, args ...string) *cli.Context {
	t.Helper()
	app := cli.NewApp()
	app.Flags = slices.Concat(NetworkFlags, []cli.Flag{DeveloperFlag, DeveloperPoWFlag, NetworkIdFlag,
		USBFlag, GCModeFlag, EthProtocolsFlag, CryptoKZGFlag})
	set := flag.NewFlagSet("test", flag.ContinueOnError)
	for _, f := range app.Flags {
		if err := f.Apply(set); err != nil {
			t.Fatal(err)
		}
	}
	if err := set.Parse(args); err != nil {
		t.Fatal(err)
	}
	return cli.NewContext(app, set, nil)
}

// configContext is networkContext with id recorded as the network ID a
// configuration file sets, as geth records it once the file is loaded; 0
// stands for a file that sets none, or no file.
func configContext(t *testing.T, id uint64, args ...string) *cli.Context {
	t.Helper()
	ctx := networkContext(t, args...)
	if id != 0 {
		SetConfigNetworkID(ctx, id)
	}
	return ctx
}

// classicArgs are the command lines that must run Ethereum Classic mainnet.
var classicArgs = [][]string{{}, {"--mainnet"}, {"--classic"}, {"--networkid", "1"}}

// privateArgs are command lines for a private network: no network flag, and a --networkid other
// than Ethereum Classic's. They keep the defaults a node without a network flag had before
// Ethereum Classic became the default, rather than taking Ethereum Classic's. Mordor's network ID
// is among them because only --mordor selects Mordor.
var privateArgs = [][]string{{"--networkid", "42"}, {"--networkid", "7"}}

func TestNetworkDataDir(t *testing.T) {
	base := vars.DefaultDataDir()
	tests := []struct {
		args []string
		want string
	}{
		{[]string{"--mordor"}, filepath.Join(base, "mordor")},
		{[]string{"--mintme"}, filepath.Join(base, "mintme")},
		{[]string{"--dev"}, ""},
	}
	for _, args := range classicArgs {
		tests = append(tests, struct {
			args []string
			want string
		}{args, filepath.Join(base, "classic")})
	}
	for _, args := range privateArgs {
		tests = append(tests, struct {
			args []string
			want string
		}{args, base})
	}
	for _, test := range tests {
		cfg := node.Config{DataDir: base}
		SetDataDir(networkContext(t, test.args...), &cfg)
		if cfg.DataDir != test.want {
			t.Errorf("%v: datadir %s, want %s", test.args, cfg.DataDir, test.want)
		}
	}
}

func TestNetworkGenesis(t *testing.T) {
	for _, args := range [][]string{{"--mainnet"}, {"--classic"}} {
		genesis := genesisForCtxChainConfig(networkContext(t, args...))
		if genesis == nil || !reflect.DeepEqual(genesis.Config, params.ClassicChainConfig) {
			t.Errorf("%v: want the Ethereum Classic genesis", args)
		}
	}
	// With no network flag the genesis is settled against the database, which
	// holds a private network's genesis when it was initialized with one.
	if genesis := genesisForCtxChainConfig(networkContext(t)); genesis != nil {
		t.Errorf("no flag: chain ID %v, want nil", genesis.Config.GetChainID())
	}
	if genesis := genesisForCtxChainConfig(networkContext(t, "--mordor")); genesis == nil || !reflect.DeepEqual(genesis.Config, params.MordorChainConfig) {
		t.Error("--mordor: want the Mordor genesis")
	}
}

func TestNetworkBootnodes(t *testing.T) {
	urls := func(nodes []*enode.Node) (out []string) {
		for _, n := range nodes {
			out = append(out, n.URLv4())
		}
		return out
	}
	classic := urls(mustParseBootnodes(params.ClassicBootnodes))
	for _, args := range classicArgs {
		var cfg p2p.Config
		ctx := networkContext(t, args...)
		setBootstrapNodes(ctx, &cfg)
		if have := urls(cfg.BootstrapNodes); !slices.Equal(have, classic) {
			t.Errorf("%v: discv4 bootnodes %v, want Ethereum Classic's", args, have)
		}
		setBootstrapNodesV5(ctx, &cfg)
		if have := urls(cfg.BootstrapNodesV5); !slices.Equal(have, classic) {
			t.Errorf("%v: discv5 bootnodes %v, want Ethereum Classic's", args, have)
		}
	}
	for _, args := range privateArgs {
		var cfg p2p.Config
		ctx := networkContext(t, args...)
		setBootstrapNodes(ctx, &cfg)
		setBootstrapNodesV5(ctx, &cfg)
		if len(cfg.BootstrapNodes) != 0 || len(cfg.BootstrapNodesV5) != 0 {
			t.Errorf("%v: bootnodes %v and %v, want none", args, urls(cfg.BootstrapNodes), urls(cfg.BootstrapNodesV5))
		}
	}
}

func TestUnmaintainedNetworkRefused(t *testing.T) {
	for _, name := range []string{"ethereum", "sepolia", "holesky"} {
		err := checkUnmaintainedNetwork(networkContext(t, "--"+name))
		if err == nil || !strings.Contains(err.Error(), "--"+name+" is deprecated") {
			t.Errorf("--%s: error %v, want the deprecation refusal", name, err)
		}
	}
	for _, args := range slices.Concat(classicArgs, privateArgs, [][]string{{"--mordor"}, {"--mintme"}, {"--dev"}}) {
		if err := checkUnmaintainedNetwork(networkContext(t, args...)); err != nil {
			t.Errorf("%v: refused with %v", args, err)
		}
	}
}

func TestNetworkMinerDefaults(t *testing.T) {
	for _, args := range slices.Concat(classicArgs, [][]string{{"--mordor"}}) {
		ctx := networkContext(t, args...)
		var mcfg miner.Config
		setMiner(ctx, &mcfg)
		if mcfg.GasCeil != 8_000_000 {
			t.Errorf("%v: gas ceiling %d, want 8000000", args, mcfg.GasCeil)
		}
		ecfg := ethconfig.Defaults
		setEthashCacheDir(ctx, &ecfg)
		if ecfg.Ethash.CacheDir != "etchash" {
			t.Errorf("%v: ethash cache dir %q, want etchash", args, ecfg.Ethash.CacheDir)
		}
	}
	for _, args := range privateArgs {
		ctx := networkContext(t, args...)
		var mcfg miner.Config
		setMiner(ctx, &mcfg)
		if mcfg.GasCeil != 0 {
			t.Errorf("%v: gas ceiling %d, want the generic default left in place", args, mcfg.GasCeil)
		}
		ecfg := ethconfig.Defaults
		setEthashCacheDir(ctx, &ecfg)
		if ecfg.Ethash.CacheDir != ethconfig.Defaults.Ethash.CacheDir {
			t.Errorf("%v: ethash cache dir %q, want %q", args, ecfg.Ethash.CacheDir, ethconfig.Defaults.Ethash.CacheDir)
		}
	}
}

// TestNetworkDefaults follows each way of choosing a network, on the command line or in a
// configuration file, through the two functions geth applies its flags with, and checks every
// default that follows the network. A network ID set only in a configuration file once went
// unread, so a private network defined there took all of Ethereum Classic's defaults and
// announced itself to Ethereum Classic's discovery network.
func TestNetworkDefaults(t *testing.T) {
	coinType := accounts.BIP0044CoinType
	t.Cleanup(func() { accounts.SetCoinTypeConfiguration(coinType) })

	urls := func(nodes []*enode.Node) (out []string) {
		for _, n := range nodes {
			out = append(out, n.URLv4())
		}
		return out
	}
	base := vars.DefaultDataDir()
	type network struct {
		datadir   string
		bootnodes []string
		coinType  uint32
		trees     []string
		gasCeil   uint64
		etchash   bool
	}
	var (
		classic = network{
			datadir:   filepath.Join(base, "classic"),
			bootnodes: urls(mustParseBootnodes(params.ClassicBootnodes)),
			coinType:  accounts.BIP0044CoinTypeEtherClassic,
			trees:     []string{params.ClassicDNSNetwork1, params.ClassicDNSNetwork2, params.ClassicDNSNetwork3},
			gasCeil:   8_000_000,
			etchash:   true,
		}
		mordor = network{
			datadir:   filepath.Join(base, "mordor"),
			bootnodes: urls(mustParseBootnodes(params.MordorBootnodes)),
			coinType:  accounts.BIP0044CoinTypeTestnet,
			trees:     []string{params.MordorDNSNetwork1, params.MordorDNSNetwork2, params.MordorDNSNetwork3},
			gasCeil:   8_000_000,
			etchash:   true,
		}
		// A private network keeps the defaults a node without a network flag had before
		// Ethereum Classic became the default.
		private = network{
			datadir:  base,
			coinType: accounts.BIP0044CoinTypeEther,
			gasCeil:  ethconfig.Defaults.Miner.GasCeil,
		}
	)
	tests := []struct {
		name     string
		args     []string
		configID uint64 // the NetworkId a configuration file sets, 0 for none
		want     network
	}{
		{"no flags, no config", nil, 0, classic},
		{"--classic", []string{"--classic"}, 0, classic},
		{"--mainnet", []string{"--mainnet"}, 0, classic},
		{"--networkid 1", []string{"--networkid", "1"}, 0, classic},
		{"config NetworkId 1", nil, 1, classic},
		{"--classic over config NetworkId 12345", []string{"--classic"}, 12345, classic},
		{"--networkid 1 over config NetworkId 12345", []string{"--networkid", "1"}, 12345, classic},
		{"--mordor", []string{"--mordor"}, 0, mordor},
		{"--mordor with config NetworkId 12345", []string{"--mordor"}, 12345, mordor},
		{"config NetworkId 12345", nil, 12345, private},
		{"config NetworkId 7, which only --mordor selects", nil, 7, private},
		{"--networkid 12345", []string{"--networkid", "12345"}, 0, private},
		{"--networkid 12345 over config NetworkId 1", []string{"--networkid", "12345"}, 1, private},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// --usb makes SetNodeConfig choose the HD derivation path, which it sets globally.
			accounts.SetCoinTypeConfiguration(accounts.BIP0044CoinTypeEther)
			ctx := configContext(t, test.configID, append(slices.Clone(test.args), "--usb")...)

			ncfg := node.Config{DataDir: base}
			SetNodeConfig(ctx, &ncfg)
			if ncfg.DataDir != test.want.datadir {
				t.Errorf("datadir %s, want %s", ncfg.DataDir, test.want.datadir)
			}
			if have := urls(ncfg.P2P.BootstrapNodes); !slices.Equal(have, test.want.bootnodes) {
				t.Errorf("discv4 bootnodes %v, want %v", have, test.want.bootnodes)
			}
			if have := urls(ncfg.P2P.BootstrapNodesV5); !slices.Equal(have, test.want.bootnodes) {
				t.Errorf("discv5 bootnodes %v, want %v", have, test.want.bootnodes)
			}
			if accounts.BIP0044CoinType != test.want.coinType {
				t.Errorf("HD derivation path coin type %d, want %d", accounts.BIP0044CoinType, test.want.coinType)
			}

			ecfg := ethconfig.Defaults
			ecfg.NetworkId = test.configID // where loading the file puts it
			SetEthConfig(ctx, nil, &ecfg)
			if !slices.Equal(ecfg.EthDiscoveryURLs, test.want.trees) || !slices.Equal(ecfg.SnapDiscoveryURLs, test.want.trees) {
				t.Errorf("discovery trees %v and %v, want %v", ecfg.EthDiscoveryURLs, ecfg.SnapDiscoveryURLs, test.want.trees)
			}
			if ecfg.Miner.GasCeil != test.want.gasCeil {
				t.Errorf("gas ceiling %d, want %d", ecfg.Miner.GasCeil, test.want.gasCeil)
			}
			if etchash := ecfg.Ethash.CacheDir == "etchash"; etchash != test.want.etchash {
				t.Errorf("ethash cache dir %q, want etchash's: %v", ecfg.Ethash.CacheDir, test.want.etchash)
			}
			if etchash := ecfg.Ethash.DatasetDir != ethconfig.Defaults.Ethash.DatasetDir; etchash != test.want.etchash {
				t.Errorf("ethash dataset dir %q, want etchash's: %v", ecfg.Ethash.DatasetDir, test.want.etchash)
			}
		})
	}
}

func Test_SplitTagsFlag(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		args string
		want map[string]string
	}{
		{
			"2 tags case",
			"host=localhost,bzzkey=123",
			map[string]string{
				"host":   "localhost",
				"bzzkey": "123",
			},
		},
		{
			"1 tag case",
			"host=localhost123",
			map[string]string{
				"host": "localhost123",
			},
		},
		{
			"empty case",
			"",
			map[string]string{},
		},
		{
			"garbage",
			"smth=smthelse=123",
			map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := SplitTagsFlag(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitTagsFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestParseMESSBlock checks how a MESS block flag reads its value. A leading zero is refused,
// because Go would read the number as octal: 010400000 would be block 2,228,224.
func TestParseMESSBlock(t *testing.T) {
	for _, c := range []struct {
		in   string
		want uint64
		err  string // a substring of the error, or "" for none
	}{
		{"0", 0, ""},
		{"15000000", 15_000_000, ""},
		{"0x9eb100", 10_400_000, ""},
		{"18446744073709551615", 18446744073709551615, ""},
		{"010400000", 0, "cannot start with 0"},
		{"00", 0, "cannot start with 0"},
		{"0_1", 0, "cannot start with 0"},
		{"", 0, "parse error"},
		{"abc", 0, "parse error"},
		{"-1", 0, "parse error"},
		{"18446744073709551616", 0, "out of range"},
	} {
		got, err := parseMESSBlock(c.in)
		switch {
		case c.err == "" && err != nil:
			t.Errorf("%q: %v, want %d", c.in, err, c.want)
		case c.err == "" && got != c.want:
			t.Errorf("%q: read as %d, want %d", c.in, got, c.want)
		case c.err != "" && (err == nil || !strings.Contains(err.Error(), c.err)):
			t.Errorf("%q: read as %d with error %v, want an error containing %q", c.in, got, err, c.err)
		}
	}
}

// TestMESSDefaultText checks the default --help gives for --mess: whether the bundled
// configuration leaves MESS on past its MESS blocks, then the release.
func TestMESSDefaultText(t *testing.T) {
	block := func(n uint64) *uint64 { return &n }
	for _, c := range []struct {
		name                     string
		activation, deactivation *uint64
		want                     string
	}{
		{"a deactivation after the activation", block(11_380_000), block(19_250_000), "off in Core-Geth 1.2.3-stable"},
		{"no deactivation", block(11_380_000), nil, "on in Core-Geth 1.2.3-stable"},
		{"an activation after the deactivation", block(20_000_000), block(19_250_000), "on in Core-Geth 1.2.3-stable"},
		{"no activation", nil, nil, "off in Core-Geth 1.2.3-stable"},
	} {
		config := *params.ClassicChainConfig
		config.SetECBP1100Transition(c.activation)
		config.SetECBP1100DeactivateTransition(c.deactivation)
		if got := messDefaultText(&config, "1.2.3-stable"); got != c.want {
			t.Errorf("%s: %q, want %q", c.name, got, c.want)
		}
	}

	// --help gives one default for --mess on every network, so the two bundled networks that
	// carry MESS have to agree on it.
	want := messDefaultText(params.ClassicChainConfig, params.VersionWithMeta)
	if got := messDefaultText(params.MordorChainConfig, params.VersionWithMeta); got != want {
		t.Errorf("Mordor gives %q and Classic %q, so no one default is true of both", got, want)
	}
	if MESSFlag.DefaultText != want {
		t.Errorf("--help gives %q, want %q", MESSFlag.DefaultText, want)
	}
}
