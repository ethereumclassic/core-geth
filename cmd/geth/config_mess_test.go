package main

import (
	"math"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/urfave/cli/v2"
)

// messConfig parses args against the MESS flags and returns the Ethereum service
// configuration geth would start with.
func messConfig(t *testing.T, args ...string) *ethconfig.Config {
	t.Helper()
	return messConfigFrom(t, new(ethconfig.Config), args...)
}

// messConfigFrom is messConfig starting from cfg, the way geth applies the flags over a
// loaded config file.
func messConfigFrom(t *testing.T, cfg *ethconfig.Config, args ...string) *ethconfig.Config {
	t.Helper()
	app := &cli.App{
		Flags: []cli.Flag{utils.MESSFlag, utils.MESSActivateFlag, utils.MESSDeactivateFlag, utils.MESSNoDisableFlag},
		Action: func(ctx *cli.Context) error {
			applyMESSFlags(ctx, cfg)
			return nil
		},
	}
	if err := app.Run(append([]string{"geth"}, args...)); err != nil {
		t.Fatalf("parsing %v: %v", args, err)
	}
	return cfg
}

// messEnabledAt reports whether MESS applies at block n on network once cfg's overrides
// are applied to the chain configuration the way eth.New applies them.
func messEnabledAt(t *testing.T, network *coregeth.CoreGethChainConfig, cfg *ethconfig.Config, n uint64) bool {
	t.Helper()
	c := *network
	if v := cfg.OverrideECBP1100; v != nil {
		if err := c.SetECBP1100Transition(v); err != nil {
			t.Fatal(err)
		}
	}
	if v := cfg.OverrideECBP1100Deactivate; v != nil {
		if err := c.SetECBP1100DeactivateTransition(v); err != nil {
			t.Fatal(err)
		}
	}
	return c.IsEnabled(c.GetECBP1100Transition, new(big.Int).SetUint64(n))
}

// TestMESSFlags follows each MESS flag from the command line to whether MESS applies at a
// block. It parses the flags with applyMESSFlags and applies the resulting overrides to a
// bundled chain configuration the way eth.New does, without starting a node. The off switch
// once parsed correctly and left MESS enabled from block 2.
func TestMESSFlags(t *testing.T) {
	classic := params.ClassicChainConfig
	activation := *classic.GetECBP1100Transition()

	// The bundled configuration follows ECIP-1110: MESS activates and then stops applying
	// at the deactivation block. With no flags the node takes that window as it stands.
	deactivation := *classic.GetECBP1100DeactivateTransition()

	t.Run("no flags take the bundled window", func(t *testing.T) {
		cfg := messConfig(t)
		if cfg.OverrideECBP1100 != nil || cfg.OverrideECBP1100Deactivate != nil || cfg.ECBP1100NoDisable != nil {
			t.Error("no flags: sets an override")
		}
		if messEnabledAt(t, classic, cfg, activation-1) {
			t.Errorf("no flags: MESS applies before the bundled activation at block %d", activation)
		}
		if !messEnabledAt(t, classic, cfg, activation) {
			t.Errorf("no flags: MESS does not apply at the bundled activation block %d", activation)
		}
		if messEnabledAt(t, classic, cfg, deactivation) || messEnabledAt(t, classic, cfg, 23_000_000) {
			t.Errorf("no flags: MESS still applies at or past the bundled deactivation block %d", deactivation)
		}
	})

	// --mess has to reach both ends of that window. Moving only the activation would leave
	// MESS off past the deactivation block while the flag reported success.
	t.Run("--mess applies past the bundled deactivation", func(t *testing.T) {
		for _, args := range [][]string{{"--mess"}, {"--mess=true"}} {
			cfg := messConfig(t, args...)
			if messEnabledAt(t, classic, cfg, activation-1) {
				t.Errorf("%v: MESS applies before the bundled activation at block %d", args, activation)
			}
			if !messEnabledAt(t, classic, cfg, activation) ||
				!messEnabledAt(t, classic, cfg, deactivation) || !messEnabledAt(t, classic, cfg, 23_000_000) {
				t.Errorf("%v: MESS is not on from the bundled activation at block %d", args, activation)
			}
		}
	})

	t.Run("off with --mess=false", func(t *testing.T) {
		// With a deactivation block set as well, the activation block is compared a second
		// time, inside the deactivation window check.
		for _, args := range [][]string{{"--mess=false"}, {"--mess=false", "--mess.deactivate=20000000"}} {
			cfg := messConfig(t, args...)
			for _, network := range []*coregeth.CoreGethChainConfig{params.ClassicChainConfig, params.MordorChainConfig} {
				for n := range map[uint64]bool{*network.GetECBP1100Transition(): true, 11_380_000: true, 19_999_999: true, 23_000_000: true} {
					if messEnabledAt(t, network, cfg, n) {
						t.Errorf("%v, chain %v: MESS applies at block %d", args, network.GetChainID(), n)
					}
				}
			}
		}
	})

	t.Run("activation", func(t *testing.T) {
		for _, args := range [][]string{
			{"--mess.activate=15000000"},
			{"--ecbp1100=15000000"},
			{"--mess=false", "--mess.activate=15000000"}, // an explicit activation wins over the off switch
		} {
			cfg := messConfig(t, args...)
			if messEnabledAt(t, classic, cfg, 14_999_999) || !messEnabledAt(t, classic, cfg, 15_000_000) {
				t.Errorf("%v: MESS does not activate at exactly block 15000000", args)
			}
		}
		// A block no chain reaches keeps MESS off, math.MaxUint64 included.
		for _, arg := range []string{"--mess.activate=18446744073709551614", "--mess.activate=18446744073709551615"} {
			if cfg := messConfig(t, arg); messEnabledAt(t, classic, cfg, 23_000_000) {
				t.Errorf("%s: MESS applies at block 23000000", arg)
			}
		}
	})

	t.Run("deactivation", func(t *testing.T) {
		for _, args := range [][]string{{"--mess.deactivate=20000000"}, {"--override.ecbp1100.deactivate=20000000"}} {
			cfg := messConfig(t, args...)
			if !messEnabledAt(t, classic, cfg, 19_999_999) || messEnabledAt(t, classic, cfg, 20_000_000) {
				t.Errorf("%v: MESS does not deactivate at exactly block 20000000", args)
			}
		}
		// A deactivation block no chain reaches keeps MESS on, math.MaxUint64 included, even on
		// a chain whose configuration deactivates it.
		scheduled := *classic
		deactivation := uint64(20_000_000)
		if err := scheduled.SetECBP1100DeactivateTransition(&deactivation); err != nil {
			t.Fatal(err)
		}
		for _, arg := range []string{"--mess.deactivate=18446744073709551614", "--mess.deactivate=18446744073709551615"} {
			if cfg := messConfig(t, arg); !messEnabledAt(t, &scheduled, cfg, 23_000_000) {
				t.Errorf("%s: MESS does not apply at block 23000000 on a chain deactivating it at block %d", arg, deactivation)
			}
		}
	})

	t.Run("no-disable", func(t *testing.T) {
		for _, args := range [][]string{{"--mess.nodisable"}, {"--ecbp1100.nodisable"}} {
			if cfg := messConfig(t, args...); cfg.ECBP1100NoDisable == nil || !*cfg.ECBP1100NoDisable {
				t.Errorf("%v: no-disable is not set", args)
			}
		}
	})

	t.Run("over a config file", func(t *testing.T) {
		// A config file dumped with --mess=false or --mess.nodisable carries those settings, and geth
		// applies the flags over it. The matching flag given explicitly once left them in place.
		off := func() *ethconfig.Config {
			never := uint64(math.MaxUint64 - 1)
			return &ethconfig.Config{OverrideECBP1100: &never}
		}
		if cfg := messConfigFrom(t, off()); messEnabledAt(t, classic, cfg, 23_000_000) {
			t.Error("no flags: MESS applies at block 23000000 under the config file's off switch")
		}
		for _, args := range [][]string{{"--mess"}, {"--mess=true"}} {
			if cfg := messConfigFrom(t, off(), args...); !messEnabledAt(t, classic, cfg, 23_000_000) {
				t.Errorf("%v: MESS does not apply at block 23000000 over the config file's off switch", args)
			}
		}
		// An activation block set in the config file is not the off switch, and --mess keeps it.
		at := uint64(15_000_000)
		cfg := messConfigFrom(t, &ethconfig.Config{OverrideECBP1100: &at}, "--mess")
		if messEnabledAt(t, classic, cfg, 14_999_999) || !messEnabledAt(t, classic, cfg, 15_000_000) {
			t.Error("--mess: the config file's activation at block 15000000 is not kept")
		}
		yes := true
		if cfg := messConfigFrom(t, &ethconfig.Config{ECBP1100NoDisable: &yes}, "--mess.nodisable=false"); cfg.ECBP1100NoDisable != nil {
			t.Error("--mess.nodisable=false: no-disable stays set from the config file")
		}
		if cfg := messConfigFrom(t, &ethconfig.Config{ECBP1100NoDisable: &yes}); cfg.ECBP1100NoDisable == nil || !*cfg.ECBP1100NoDisable {
			t.Error("no flags: the config file's no-disable setting is dropped")
		}
	})
}

// TestMESSFlagsDumpConfig checks that dumpconfig writes the MESS settings the flags make, and
// that they survive being read back from the dumped file. The dump once left --mess=false
// out, so a node started from it ran with MESS on.
func TestMESSFlagsDumpConfig(t *testing.T) {
	dump := func(args ...string) string {
		t.Helper()
		out := filepath.Join(t.TempDir(), "config.toml")
		geth := runGeth(t, append(args, "dumpconfig", out)...)
		geth.WaitExit()
		if status := geth.ExitStatus(); status != 0 {
			t.Fatalf("%v dumpconfig: exit status %d\n%s", args, status, geth.StderrText())
		}
		b, err := os.ReadFile(out)
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	for _, c := range []struct {
		args []string
		want []string
	}{
		{[]string{"--mess=false"}, []string{"OverrideECBP1100 = 18446744073709551614"}},
		{
			[]string{"--mess.activate=15000000", "--mess.deactivate=20000000", "--mess.nodisable"},
			[]string{"OverrideECBP1100 = 15000000", "OverrideECBP1100Deactivate = 20000000", "ECBP1100NoDisable = true"},
		},
	} {
		dumped := dump(append([]string{"--mordor"}, c.args...)...)
		file := filepath.Join(t.TempDir(), "dumped.toml")
		if err := os.WriteFile(file, []byte(dumped), 0o600); err != nil {
			t.Fatal(err)
		}
		readBack := dump("--mordor", "--config", file)
		for _, want := range c.want {
			if !strings.Contains(dumped, want) {
				t.Errorf("%v: dumped config lacks %q", c.args, want)
			}
			if !strings.Contains(readBack, want) {
				t.Errorf("%v: config read back from the dump lacks %q", c.args, want)
			}
		}
	}
}
