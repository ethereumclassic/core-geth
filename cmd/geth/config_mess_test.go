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
		Flags: messFlags(),
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

// messFlags returns fresh copies of the MESS flags, environment variables included. A flag set
// from an environment variable keeps that value in the flag itself, so tests sharing the
// package's flags would leak it into each other.
func messFlags() []cli.Flag {
	mess, activate, deactivate, noDisable := *utils.MESSFlag, *utils.MESSActivateFlag, *utils.MESSDeactivateFlag, *utils.MESSNoDisableFlag
	return []cli.Flag{&mess, &activate, &deactivate, &noDisable}
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

	// The bundled configuration activates MESS and then stops applying it at the
	// deactivation block. With no flags the node takes that window as it stands.
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
		// --mess.activate turns MESS on from its block, past the bundled deactivation too, with or
		// without --mess.
		for _, args := range [][]string{
			{"--mess.activate=15000000"},
			{"--ecbp1100=15000000"},
			{"--mess", "--mess.activate=15000000"},
			{"--mess=false", "--mess.activate=15000000"}, // an explicit activation wins over the off switch
		} {
			cfg := messConfig(t, args...)
			if messEnabledAt(t, classic, cfg, 14_999_999) || !messEnabledAt(t, classic, cfg, 15_000_000) {
				t.Errorf("%v: MESS does not activate at exactly block 15000000", args)
			}
			if !messEnabledAt(t, classic, cfg, deactivation) || !messEnabledAt(t, classic, cfg, 23_000_000) {
				t.Errorf("%v: MESS stops at the bundled deactivation block %d", args, deactivation)
			}
		}
		// A block no chain reaches keeps MESS off, math.MaxUint64 included. Block 15000000 is
		// inside the bundled window, where only the activation can keep MESS off.
		for _, arg := range []string{"--mess.activate=18446744073709551614", "--mess.activate=18446744073709551615"} {
			cfg := messConfig(t, arg)
			if messEnabledAt(t, classic, cfg, 15_000_000) || messEnabledAt(t, classic, cfg, 23_000_000) {
				t.Errorf("%s: MESS applies at block 15000000 or 23000000", arg)
			}
		}
	})

	t.Run("deactivation", func(t *testing.T) {
		for _, args := range [][]string{
			{"--mess.deactivate=20000000"},
			{"--override.ecbp1100.deactivate=20000000"},
			{"--mess", "--mess.deactivate=20000000"}, // an explicit deactivation wins over --mess
		} {
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

	t.Run("the later block wins", func(t *testing.T) {
		// Each block is where MESS turns on or off, and the later one decides from there, so an
		// activation after the deactivation turns MESS on again.
		cfg := messConfig(t, "--mess.deactivate=20000000", "--mess.activate=25000000")
		for n, want := range map[uint64]bool{20_000_000: false, 24_999_999: false, 25_000_000: true, 30_000_000: true} {
			if got := messEnabledAt(t, classic, cfg, n); got != want {
				t.Errorf("deactivation 20000000, activation 25000000: MESS applies at block %d is %v, want %v", n, got, want)
			}
		}
		// At one block neither is later, and MESS stays off.
		cfg = messConfig(t, "--mess.activate=20000000", "--mess.deactivate=20000000")
		for _, n := range []uint64{19_999_999, 20_000_000, 25_000_000} {
			if messEnabledAt(t, classic, cfg, n) {
				t.Errorf("activation and deactivation both 20000000: MESS applies at block %d", n)
			}
		}
	})

	t.Run("no-disable", func(t *testing.T) {
		// --mess.nodisable turns MESS on and keeps it on, past the bundled deactivation too, with
		// or without --mess.
		for _, args := range [][]string{{"--mess.nodisable"}, {"--ecbp1100.nodisable"}, {"--mess", "--mess.nodisable"}} {
			cfg := messConfig(t, args...)
			if cfg.ECBP1100NoDisable == nil || !*cfg.ECBP1100NoDisable {
				t.Errorf("%v: no-disable is not set", args)
			}
			if messEnabledAt(t, classic, cfg, activation-1) || !messEnabledAt(t, classic, cfg, activation) ||
				!messEnabledAt(t, classic, cfg, deactivation) || !messEnabledAt(t, classic, cfg, 23_000_000) {
				t.Errorf("%v: MESS is not on from the bundled activation at block %d with no end", args, activation)
			}
		}
		// An explicit deactivation ends it.
		cfg := messConfig(t, "--mess.nodisable", "--mess.deactivate=20000000")
		if !messEnabledAt(t, classic, cfg, 19_999_999) || messEnabledAt(t, classic, cfg, 20_000_000) {
			t.Error("--mess.nodisable --mess.deactivate=20000000: MESS does not stop at exactly block 20000000")
		}
		// --mess=false wins over it.
		cfg = messConfig(t, "--mess=false", "--mess.nodisable")
		if messEnabledAt(t, classic, cfg, 15_000_000) || messEnabledAt(t, classic, cfg, 23_000_000) {
			t.Error("--mess=false --mess.nodisable: MESS applies")
		}
	})

	t.Run("over a config file", func(t *testing.T) {
		// A config file dumped with --mess=false or --mess.nodisable carries those settings, and geth
		// applies the flags over it. The matching flag given explicitly once left them in place.
		off := func() *ethconfig.Config {
			never := uint64(math.MaxUint64 - 1)
			return &ethconfig.Config{OverrideECBP1100: &never}
		}
		// Block 15000000 is inside the bundled window, where only the off switch keeps MESS off.
		if cfg := messConfigFrom(t, off()); messEnabledAt(t, classic, cfg, 15_000_000) || messEnabledAt(t, classic, cfg, 23_000_000) {
			t.Error("no flags: MESS applies at block 15000000 or 23000000 under the config file's off switch")
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

		// A config file's activation or no-disable line means what its flag means: MESS on, with no
		// end unless a deactivation is given.
		cfg = messConfigFrom(t, &ethconfig.Config{OverrideECBP1100: &at})
		if messEnabledAt(t, classic, cfg, 14_999_999) || !messEnabledAt(t, classic, cfg, 15_000_000) ||
			!messEnabledAt(t, classic, cfg, 23_000_000) {
			t.Error("no flags: the config file's activation at block 15000000 does not turn MESS on with no end")
		}
		until := uint64(20_000_000)
		cfg = messConfigFrom(t, &ethconfig.Config{OverrideECBP1100: &at, OverrideECBP1100Deactivate: &until})
		if !messEnabledAt(t, classic, cfg, 19_999_999) || messEnabledAt(t, classic, cfg, 20_000_000) {
			t.Error("no flags: the config file's deactivation at block 20000000 does not end MESS")
		}
		cfg = messConfigFrom(t, &ethconfig.Config{ECBP1100NoDisable: &yes})
		if !messEnabledAt(t, classic, cfg, activation) || !messEnabledAt(t, classic, cfg, 23_000_000) {
			t.Error("no flags: the config file's no-disable line does not turn MESS on with no end")
		}
		// The file's off switch wins over its no-disable line, and a flag wins over the file.
		cfg = messConfigFrom(t, &ethconfig.Config{OverrideECBP1100: off().OverrideECBP1100, ECBP1100NoDisable: &yes})
		if messEnabledAt(t, classic, cfg, 15_000_000) || messEnabledAt(t, classic, cfg, 23_000_000) {
			t.Error("no flags: MESS applies under the config file's off switch and no-disable line")
		}
		if cfg := messConfigFrom(t, off(), "--mess.nodisable"); !messEnabledAt(t, classic, cfg, 23_000_000) {
			t.Error("--mess.nodisable: MESS does not apply at block 23000000 over the config file's off switch")
		}
		// Dropping the file's no-disable line returns MESS to the bundled window.
		if cfg := messConfigFrom(t, &ethconfig.Config{ECBP1100NoDisable: &yes}, "--mess.nodisable=false"); messEnabledAt(t, classic, cfg, 23_000_000) {
			t.Error("--mess.nodisable=false: MESS applies at block 23000000, past the bundled deactivation")
		}
	})
}

// TestMESSConflicts checks that MESS settings asking for opposite things are reported at startup,
// and that settings which agree are not.
func TestMESSConflicts(t *testing.T) {
	conflicts := func(args ...string) []string {
		t.Helper()
		var out []string
		cfg := new(ethconfig.Config)
		app := &cli.App{
			Flags: messFlags(),
			Action: func(ctx *cli.Context) error {
				applyMESSFlags(ctx, cfg)
				out = messConflicts(ctx, cfg)
				return nil
			},
		}
		if err := app.Run(append([]string{"geth"}, args...)); err != nil {
			t.Fatalf("parsing %v: %v", args, err)
		}
		return out
	}
	for _, c := range []struct {
		args []string
		want string // a substring of the one warning, or "" for none
	}{
		{[]string{"--mess=false", "--mess.activate=15000000"}, "--mess=false and --mess.activate=15000000; MESS applies from block 15000000"},
		{[]string{"--mess=false", "--mess.nodisable"}, "--mess=false and --mess.nodisable; MESS is off"},
		{[]string{"--mess.activate=20000000", "--mess.deactivate=20000000"}, "both block 20000000; MESS is off"},
		{nil, ""},
		{[]string{"--mess=false", "--mess.deactivate=20000000"}, ""}, // both say off
		{[]string{"--mess", "--mess.activate=15000000"}, ""},
		{[]string{"--mess", "--mess.deactivate=20000000"}, ""},
		{[]string{"--mess", "--mess.nodisable"}, ""},
		{[]string{"--mess.deactivate=20000000", "--mess.activate=25000000"}, ""}, // the later block wins
		{[]string{"--mess=false", "--mess.activate=18446744073709551614"}, ""},   // both keep MESS off
	} {
		got := conflicts(c.args...)
		switch {
		case c.want == "" && len(got) != 0:
			t.Errorf("%v: unexpected warnings %q", c.args, got)
		case c.want != "" && (len(got) != 1 || !strings.Contains(got[0], c.want)):
			t.Errorf("%v: warnings %q, want one containing %q", c.args, got, c.want)
		}
	}
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
		{[]string{"--mess"}, []string{"OverrideECBP1100Deactivate = 18446744073709551614"}},
		{[]string{"--mess=false"}, []string{"OverrideECBP1100 = 18446744073709551614"}},
		{[]string{"--mess.activate=15000000"}, []string{"OverrideECBP1100 = 15000000", "OverrideECBP1100Deactivate = 18446744073709551614"}},
		{[]string{"--mess.nodisable"}, []string{"ECBP1100NoDisable = true", "OverrideECBP1100Deactivate = 18446744073709551614"}},
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
