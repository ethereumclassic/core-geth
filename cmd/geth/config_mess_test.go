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
			// In the order geth applies them: the flags naming a block in
			// makeConfigNode, then the on/off switch in makeFullNode.
			applyMESSBlockFlags(ctx, cfg)
			applyMESSToggleFlag(ctx, cfg)
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

	// The v1.12.x line has no boolean switch: it names blocks, with these three flags. Every
	// node migrating off that line carries one of these spellings if it configured MESS at
	// all, so the migration guide's promise that nothing changes rests on them still parsing.
	// Nothing else asserts that, and a rename would break those operators silently.
	t.Run("the v1.12.x spellings still work", func(t *testing.T) {
		for _, c := range []struct {
			args []string
			want func(*ethconfig.Config) bool
			what string
		}{
			{[]string{"--ecbp1100=15000000"},
				func(c *ethconfig.Config) bool { return c.OverrideECBP1100 != nil && *c.OverrideECBP1100 == 15_000_000 },
				"activation block"},
			{[]string{"--override.ecbp1100.deactivate=20000000"},
				func(c *ethconfig.Config) bool {
					return c.OverrideECBP1100Deactivate != nil && *c.OverrideECBP1100Deactivate == 20_000_000
				},
				"deactivation block"},
			{[]string{"--ecbp1100.nodisable"},
				func(c *ethconfig.Config) bool { return c.ECBP1100NoDisable != nil && *c.ECBP1100NoDisable },
				"no-disable"},
		} {
			if cfg := messConfig(t, c.args...); !c.want(cfg) {
				t.Errorf("%v: the v1.12.x spelling no longer sets the %s", c.args, c.what)
			}
		}
		// A v1.12.x operator who disabled MESS the way ECIP-1110 documents, by pushing the
		// activation out of reach, must still get MESS off. On the v1.12.x line this value was
		// read as signed and inverted, leaving MESS on; reading it as written is the fix, and
		// this pins that the fix survives.
		cfg := messConfig(t, "--ecbp1100=18446744073709551615")
		if messEnabledAt(t, classic, cfg, 23_000_000) {
			t.Error("--ecbp1100=<max>: MESS applies at block 23000000, so the unreachable activation inverted")
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
		// --mess=false clears an on switch a config file carries, mirroring what --mess does
		// to an off switch. Without it the file ends up holding both sentinels: the right
		// answer, since an unreachable activation wins, but a self-contradictory file.
		never := uint64(math.MaxUint64 - 1)
		on := func() *ethconfig.Config {
			n := never
			return &ethconfig.Config{OverrideECBP1100Deactivate: &n}
		}
		if cfg := messConfigFrom(t, on(), "--mess=false"); cfg.OverrideECBP1100Deactivate != nil {
			t.Errorf("--mess=false: the config file's on switch is left set to %d",
				*cfg.OverrideECBP1100Deactivate)
		}
		if messEnabledAt(t, classic, messConfigFrom(t, on(), "--mess=false"), 23_000_000) {
			t.Error("--mess=false: MESS applies at block 23000000 over the config file's on switch")
		}
		// A real block an operator chose is not a sentinel and must survive.
		real := uint64(20_000_000)
		if cfg := messConfigFrom(t, &ethconfig.Config{OverrideECBP1100Deactivate: &real}, "--mess=false"); cfg.OverrideECBP1100Deactivate == nil || *cfg.OverrideECBP1100Deactivate != real {
			t.Error("--mess=false: a deactivation block the operator chose was cleared")
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

// TestMESSFlagsDumpConfig checks what dumpconfig writes, and the split is the point: a flag
// naming a block is a setting an operator chose and survives a round trip, while the on/off
// switch does not persist at all, so the flags a node runs with decide MESS and dropping one
// returns it to the bundled default.
//
// The switch was captured once, and this test asserted that. It was protecting a client whose
// bundled default was MESS on, where a dump silently dropping --mess=false left a node
// defenceless. Under the ECIP-1110 default that inverts: dropping the flag yields off, which is
// the published default. Revisit if the bundled default ever returns to on.
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
		// notWant is checked against every MESS key, so an empty want is a real
		// assertion rather than a test that checks nothing.
		notWant []string
	}{
		// The on/off switch writes nothing, in either position.
		{args: []string{"--mess"}, notWant: allMESSConfigKeys},
		{args: []string{"--mess=false"}, notWant: allMESSConfigKeys},
		// A block an operator named is kept, and the switch beside it adds nothing.
		{
			args:    []string{"--mess.activate=15000000", "--mess.deactivate=20000000", "--mess.nodisable"},
			want:    []string{"OverrideECBP1100 = 15000000", "OverrideECBP1100Deactivate = 20000000", "ECBP1100NoDisable = true"},
			notWant: []string{"18446744073709551614"},
		},
		{
			args:    []string{"--mess", "--mess.deactivate=20000000"},
			want:    []string{"OverrideECBP1100Deactivate = 20000000"},
			notWant: []string{"18446744073709551614", "OverrideECBP1100 ="},
		},
		{
			args:    []string{"--mess=false", "--mess.activate=15000000"},
			want:    []string{"OverrideECBP1100 = 15000000"},
			notWant: []string{"18446744073709551614", "OverrideECBP1100Deactivate"},
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
		for _, notWant := range c.notWant {
			if strings.Contains(dumped, notWant) {
				t.Errorf("%v: dumped config carries %q, which must not persist", c.args, notWant)
			}
			if strings.Contains(readBack, notWant) {
				t.Errorf("%v: config read back from the dump carries %q", c.args, notWant)
			}
		}
	}
}

// allMESSConfigKeys is every key a MESS flag can write. Asserting a dump carries none of them
// is what makes "the switch does not persist" a test rather than an absence nobody checks.
var allMESSConfigKeys = []string{"OverrideECBP1100", "OverrideECBP1100Deactivate", "ECBP1100NoDisable"}

// TestMESSFlagsDumpConfigCalibration is the negative control for the notWant assertions above.
// A flag naming a block must still reach the dump, so a dumpconfig that silently stopped
// writing MESS settings altogether would fail here rather than pass the whole suite.
func TestMESSFlagsDumpConfigCalibration(t *testing.T) {
	out := filepath.Join(t.TempDir(), "config.toml")
	geth := runGeth(t, "--mordor", "--mess.activate=15000000", "dumpconfig", out)
	geth.WaitExit()
	if status := geth.ExitStatus(); status != 0 {
		t.Fatalf("dumpconfig: exit status %d\n%s", status, geth.StderrText())
	}
	b, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "OverrideECBP1100 = 15000000") {
		t.Fatal("dumpconfig writes no MESS key at all, so the non-persistence assertions prove nothing")
	}
}
