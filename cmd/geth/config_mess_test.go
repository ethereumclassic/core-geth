package main

import (
	"fmt"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

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
			applyMESSRunSettings(ctx, cfg)
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
// block. It applies the flags as makeFullNode does, and the resulting overrides to a bundled
// chain configuration as eth.New does, without starting a node. The off switch
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
		if cfg := messConfigFrom(t, &ethconfig.Config{OverrideECBP1100Deactivate: &until}, "--mess.nodisable"); !messEnabledAt(t, classic, cfg, 23_000_000) {
			t.Error("--mess.nodisable: MESS does not apply at block 23000000, past the config file's deactivation at block 20000000")
		}
		// Dropping the file's no-disable line returns MESS to the bundled window.
		if cfg := messConfigFrom(t, &ethconfig.Config{ECBP1100NoDisable: &yes}, "--mess.nodisable=false"); messEnabledAt(t, classic, cfg, 23_000_000) {
			t.Error("--mess.nodisable=false: MESS applies at block 23000000, past the bundled deactivation")
		}
	})
}

// TestMESSFlagsLegacyEnvironment checks that the v1.12.x environment variable names of the MESS
// flags set them, and that the new name wins when both are set.
func TestMESSFlagsLegacyEnvironment(t *testing.T) {
	t.Setenv("GETH_ECBP1100", "15000000")
	t.Setenv("GETH_OVERRIDE_ECBP1100_DEACTIVATE", "20000000")
	t.Setenv("GETH_ECBP1100_NODISABLE", "true")
	cfg := messConfig(t)
	if v := cfg.OverrideECBP1100; v == nil || *v != 15_000_000 {
		t.Error("GETH_ECBP1100=15000000: the activation is not block 15000000")
	}
	if v := cfg.OverrideECBP1100Deactivate; v == nil || *v != 20_000_000 {
		t.Error("GETH_OVERRIDE_ECBP1100_DEACTIVATE=20000000: the deactivation is not block 20000000")
	}
	if v := cfg.ECBP1100NoDisable; v == nil || !*v {
		t.Error("GETH_ECBP1100_NODISABLE=true: no-disable is not set")
	}
	// The new name wins when both are set.
	t.Setenv("GETH_MESS_ACTIVATE", "16000000")
	if v := messConfig(t).OverrideECBP1100; v == nil || *v != 16_000_000 {
		t.Error("GETH_MESS_ACTIVATE=16000000 with GETH_ECBP1100=15000000: the activation is not block 16000000")
	}
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
				applyMESSRunSettings(ctx, cfg)
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

// TestMESSFlagsDumpConfig checks that dumpconfig writes the line of each MESS flag given and
// nothing else, and that the lines survive being read back from the dumped file. The end that
// --mess.activate and --mess.nodisable imply is applied when the node starts instead: written
// into the file, it would outlive the setting, and deleting the activation line later would
// leave MESS on. The dump once left --mess=false out, so a node started from it ran with MESS
// on.
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
	// check compares a dump's MESS lines with want exactly, so a line nobody set fails it.
	check := func(label, dumped string, want []string) {
		t.Helper()
		var got []string
		for _, line := range strings.Split(dumped, "\n") {
			if line = strings.TrimSpace(line); strings.Contains(line, "ECBP1100") {
				got = append(got, line)
			}
		}
		slices.Sort(got)
		want = slices.Clone(want)
		slices.Sort(want)
		if !slices.Equal(got, want) {
			t.Errorf("%s: MESS lines %q, want %q", label, got, want)
		}
	}
	never := "18446744073709551614"
	for _, c := range []struct {
		args []string
		want []string
	}{
		{nil, nil},
		{[]string{"--mess"}, []string{"OverrideECBP1100Deactivate = " + never}},
		{[]string{"--mess=false"}, []string{"OverrideECBP1100 = " + never}},
		{[]string{"--mess.activate=15000000"}, []string{"OverrideECBP1100 = 15000000"}},
		{[]string{"--mess.nodisable"}, []string{"ECBP1100NoDisable = true"}},
		{[]string{"--mess=false", "--mess.activate=15000000"}, []string{"OverrideECBP1100 = 15000000"}},
		{
			[]string{"--mess.activate=15000000", "--mess.deactivate=20000000", "--mess.nodisable"},
			[]string{"OverrideECBP1100 = 15000000", "OverrideECBP1100Deactivate = 20000000", "ECBP1100NoDisable = true"},
		},
	} {
		dumped := dump(append([]string{"--mordor"}, c.args...)...)
		check(fmt.Sprint(c.args), dumped, c.want)
		file := filepath.Join(t.TempDir(), "dumped.toml")
		if err := os.WriteFile(file, []byte(dumped), 0o600); err != nil {
			t.Fatal(err)
		}
		check(fmt.Sprint(c.args)+" read back", dump("--mordor", "--config", file), c.want)
	}
	// A config file holding only an activation or a no-disable line is dumped as it is.
	for _, line := range []string{"OverrideECBP1100 = 15000000", "ECBP1100NoDisable = true"} {
		file := filepath.Join(t.TempDir(), "config.toml")
		if err := os.WriteFile(file, []byte("[Eth]\n"+line+"\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		check(fmt.Sprintf("file %q", line), dump("--mordor", "--config", file), []string{line})
	}
}

// TestMESSFlagsRunningNode starts a Mordor node and reads the MESS window it runs with, so the
// end that --mess.activate and --mess.nodisable imply, which dumpconfig does not write, is
// shown to reach the running node.
func TestMESSFlagsRunningNode(t *testing.T) {
	config := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(config, []byte("[Eth]\nOverrideECBP1100 = 15000000\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		name string
		args []string
		want string // the activation and deactivation blocks admin.ecbp1100Status reports
	}{
		{"bundled", nil, "0x2450e0 0x9eb100"},
		{"activate", []string{"--mess.activate=15000000"}, "0xe4e1c0 0xfffffffffffffffe"},
		{"nodisable", []string{"--mess.nodisable"}, "0x2450e0 0xfffffffffffffffe"},
		{"config file", []string{"--config", config}, "0xe4e1c0 0xfffffffffffffffe"},
	} {
		t.Run(c.name, func(t *testing.T) {
			args := append(c.args, "--ipcdisable", "--exec",
				`var s = admin.ecbp1100Status(); s.activatedAtBlock + " " + s.defaultDisabledAtBlock`, "console")
			geth := runMinimalGeth(t, args...)
			geth.KillTimeout = 20 * time.Second
			geth.Expect(fmt.Sprintf("%q\n", c.want))
			geth.ExpectExit()
		})
	}
}

// TestMESSFlagsRefused checks that geth refuses a value a MESS flag cannot take: it exits with an
// error naming the flag and writes no config file, rather than going on with the bundled default.
func TestMESSFlagsRefused(t *testing.T) {
	run := func(args ...string) (status int, stderr string, wrote bool) {
		t.Helper()
		out := filepath.Join(t.TempDir(), "config.toml")
		geth := runGeth(t, append(append([]string{"--mordor"}, args...), "dumpconfig", out)...)
		geth.WaitExit()
		_, err := os.Stat(out)
		return geth.ExitStatus(), geth.StderrText(), err == nil
	}
	// The control: a value the flag takes is written, so the cases below fail because of their
	// values and not because this harness cannot pass.
	if status, stderr, wrote := run("--mess.activate=15000000"); status != 0 || !wrote {
		t.Fatalf("--mess.activate=15000000: exit status %d, config written %t\n%s", status, wrote, stderr)
	}
	for _, c := range []struct {
		args []string
		want string // in the error geth prints
	}{
		{[]string{"--mess.activate="}, "-mess.activate:"},
		{[]string{"--mess.activate=abc"}, "-mess.activate:"},
		{[]string{"--mess.activate=-1"}, "-mess.activate:"},
		{[]string{"--mess.activate=18446744073709551616"}, "-mess.activate:"},
		// Without a value, the flag takes the next argument as its value.
		{[]string{"--mess.activate", "--mess"}, "-mess.activate:"},
		{[]string{"--mess.deactivate=abc"}, "-mess.deactivate:"},
		{[]string{"--mess=abc"}, "-mess:"},
		{[]string{"--mess.nodisable=abc"}, "-mess.nodisable:"},
		// Go would read a leading zero as octal, a different block, so it is refused.
		{[]string{"--mess.activate=010400000"}, "-mess.activate: a block number cannot start with 0"},
		{[]string{"--mess.deactivate=010400000"}, "-mess.deactivate: a block number cannot start with 0"},
	} {
		status, stderr, wrote := run(c.args...)
		if status == 0 || wrote {
			t.Errorf("%v: exit status %d, config written %t, want refused", c.args, status, wrote)
		}
		if !strings.Contains(stderr, c.want) {
			t.Errorf("%v: error %q does not name the flag (%q)", c.args, strings.TrimSpace(stderr), c.want)
		}
	}
	// A value from an environment variable is refused the same way.
	t.Setenv("GETH_MESS_ACTIVATE", "abc")
	if status, stderr, wrote := run(); status == 0 || wrote || !strings.Contains(stderr, "mess.activate") {
		t.Errorf("GETH_MESS_ACTIVATE=abc: exit status %d, config written %t, error %q, want refused naming the flag",
			status, wrote, strings.TrimSpace(stderr))
	}
	// A value refused from an environment variable names the variable as well, under either name,
	// and a leading zero is refused there too.
	for _, c := range []struct{ name, value, want string }{
		{"GETH_MESS_ACTIVATE", "abc", "parse error"},
		{"GETH_MESS_ACTIVATE", "010400000", "cannot start with 0"},
		{"GETH_ECBP1100", "010400000", "cannot start with 0"},
	} {
		os.Unsetenv("GETH_MESS_ACTIVATE")
		t.Setenv(c.name, c.value)
		if status, stderr, wrote := run(); status == 0 || wrote || !strings.Contains(stderr, c.want) ||
			!strings.Contains(stderr, "(from "+c.name+")") {
			t.Errorf("%s=%s: exit status %d, config written %t, error %q, want refused naming %s",
				c.name, c.value, status, wrote, strings.TrimSpace(stderr), c.name)
		}
	}
	// A value given on the command line names no variable, even beside one that holds another.
	os.Unsetenv("GETH_ECBP1100")
	t.Setenv("GETH_MESS_ACTIVATE", "15000000")
	if status, stderr, _ := run("--mess.activate=010400000"); status == 0 || strings.Contains(stderr, "(from") {
		t.Errorf("--mess.activate=010400000 beside GETH_MESS_ACTIVATE=15000000: exit status %d, error %q, want refused naming no variable",
			status, strings.TrimSpace(stderr))
	}
}

// TestMESSBlockRefusedBeforeStart checks that geth refuses a MESS block value before any command
// runs: the data directory stays empty, rather than being opened and then abandoned.
func TestMESSBlockRefusedBeforeStart(t *testing.T) {
	for _, arg := range []string{"--mess.activate=010400000", "--mess.deactivate=abc"} {
		out := filepath.Join(t.TempDir(), "config.toml")
		geth := runGeth(t, "--mordor", arg, "dumpconfig", out)
		geth.WaitExit()
		if status := geth.ExitStatus(); status == 0 {
			t.Errorf("%s: exit status 0, want refused", arg)
		}
		if entries, err := os.ReadDir(geth.Datadir); err != nil || len(entries) != 0 {
			t.Errorf("%s: the data directory holds %d entries (err %v), want none", arg, len(entries), err)
		}
	}
}

// TestMESSEmptyEnvVarIgnored checks that a MESS block variable set to the empty string gives no
// block, as an empty variable does for every numeric flag: geth runs as if it were unset, and
// writes no MESS line. A block on the command line still applies beside one, and an empty value
// on the command line is still refused, as TestMESSFlagsRefused checks.
func TestMESSEmptyEnvVarIgnored(t *testing.T) {
	dump := func(t *testing.T, args ...string) string {
		t.Helper()
		out := filepath.Join(t.TempDir(), "config.toml")
		geth := runGeth(t, append(append([]string{"--mordor"}, args...), "dumpconfig", out)...)
		geth.WaitExit()
		if status := geth.ExitStatus(); status != 0 {
			t.Fatalf("exit status %d, want 0: %s", status, strings.TrimSpace(geth.StderrText()))
		}
		data, err := os.ReadFile(out)
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	// messLines returns the dump's MESS keys. The data directory's path names the test, and so
	// can contain ECBP1100 too, so only a line that starts with a MESS key counts.
	messLines := func(dump string) []string {
		var keys []string
		for _, line := range strings.Split(dump, "\n") {
			if strings.HasPrefix(line, "OverrideECBP1100") || strings.HasPrefix(line, "ECBP1100NoDisable") {
				keys = append(keys, line)
			}
		}
		return keys
	}
	for _, name := range []string{"GETH_MESS_ACTIVATE", "GETH_ECBP1100", "GETH_MESS_DEACTIVATE", "GETH_OVERRIDE_ECBP1100_DEACTIVATE"} {
		t.Run(name, func(t *testing.T) {
			t.Setenv(name, "")
			if keys := messLines(dump(t)); len(keys) != 0 {
				t.Errorf("%s set to the empty string: the dump carries %q", name, keys)
			}
		})
	}
	t.Run("a block on the command line beside it", func(t *testing.T) {
		t.Setenv("GETH_MESS_ACTIVATE", "")
		if keys := messLines(dump(t, "--mess.activate=15000000")); !slices.Equal(keys, []string{"OverrideECBP1100 = 15000000"}) {
			t.Errorf("GETH_MESS_ACTIVATE empty with --mess.activate=15000000: the dump carries %q, want the activation only", keys)
		}
	})
}
