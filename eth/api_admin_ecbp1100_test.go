package eth

import (
	"encoding/json"
	"math"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/rpc"
)

// A block tag is a negative rpc.BlockNumber. Converted straight to a height, "latest"
// became math.MaxUint64-1: once heights stopped wrapping, admin.messActivate("latest") would
// have scheduled MESS for a block no chain reaches instead of the head it names.
func TestEcbp1100ActivationBlock(t *testing.T) {
	const head = 20_000_000
	for _, c := range []struct {
		in      rpc.BlockNumber
		want    uint64
		refused bool
	}{
		{in: rpc.EarliestBlockNumber, want: 0},
		{in: rpc.BlockNumber(11_380_000), want: 11_380_000},
		{in: rpc.LatestBlockNumber, want: head},
		{in: rpc.PendingBlockNumber, want: head},
		{in: rpc.FinalizedBlockNumber, refused: true},
		{in: rpc.SafeBlockNumber, refused: true},
	} {
		got, err := ecbp1100ActivationBlock(c.in, head)
		switch {
		case c.refused && err == nil:
			t.Errorf("%v: resolved to %d, want refused", c.in, got)
		case !c.refused && err != nil:
			t.Errorf("%v: refused: %v", c.in, err)
		case !c.refused && got != c.want:
			t.Errorf("%v: resolved to %d, want %d", c.in, got, c.want)
		}
	}
}

// Both ends of the window resolve through the same helper, so the tag trap above is
// handled once. This pins that the deactivation end really does share it, and that its
// refusal names the end it was asked about rather than the other one.
func TestEcbp1100DeactivationBlock(t *testing.T) {
	const head = 20_000_000
	for _, c := range []struct {
		in      rpc.BlockNumber
		want    uint64
		refused bool
	}{
		{in: rpc.EarliestBlockNumber, want: 0},
		{in: rpc.BlockNumber(19_250_000), want: 19_250_000},
		{in: rpc.LatestBlockNumber, want: head},
		{in: rpc.PendingBlockNumber, want: head},
		{in: rpc.FinalizedBlockNumber, refused: true},
		{in: rpc.SafeBlockNumber, refused: true},
	} {
		got, err := ecbp1100DeactivationBlock(c.in, head)
		switch {
		case c.refused && err == nil:
			t.Errorf("%v: resolved to %d, want refused", c.in, got)
		case c.refused && !strings.Contains(err.Error(), "deactivate"):
			t.Errorf("%v: refused with %q, which does not name the deactivation end", c.in, err)
		case !c.refused && err != nil:
			t.Errorf("%v: refused: %v", c.in, err)
		case !c.refused && got != c.want:
			t.Errorf("%v: resolved to %d, want %d", c.in, got, c.want)
		}
	}
}

// TestApplyMESSSwitch drives the decision behind admin_mess against a chain configuration,
// the way eth.New applies one, without standing up a node. The switch has to work whatever the
// two blocks are, including past a deactivation block, where ECBP-1110's bundled window leaves
// every Classic and Mordor node.
func TestApplyMESSSwitch(t *testing.T) {
	const head = 20_000_000
	never := uint64(math.MaxUint64 - 1)
	ptr := func(n uint64) *uint64 { return &n }

	bundled := ptr(11_380_000) // the activation block the network ships
	for _, c := range []struct {
		name               string
		activate, deactive *uint64
		bundled            *uint64
		enable             bool
		wantEnabled        bool
		wantActivation     *uint64 // where an activation that kept MESS off moves to
	}{
		// The window closed behind the head.
		{name: "on, past a deactivation", activate: ptr(11_380_000), deactive: ptr(19_250_000), bundled: bundled, enable: true, wantEnabled: true},
		// An off switch goes back to the network's own block, as --mess clears one.
		{name: "on, from an off switch", activate: ptr(never), bundled: bundled, enable: true, wantEnabled: true, wantActivation: bundled},
		{name: "on, from an off switch, no network block", activate: ptr(never), enable: true, wantEnabled: true, wantActivation: ptr(0)},
		{name: "on, from an off switch, network block beyond the head", activate: ptr(never), bundled: ptr(30_000_000), enable: true, wantEnabled: true, wantActivation: ptr(0)},
		{name: "on, window already open", activate: ptr(11_380_000), bundled: bundled, enable: true, wantEnabled: true},
		{name: "on, nothing configured", enable: true, wantEnabled: true, wantActivation: ptr(0)},
		{name: "off, from on", activate: ptr(11_380_000), bundled: bundled, enable: false},
		{name: "off, past a deactivation", activate: ptr(11_380_000), deactive: ptr(19_250_000), bundled: bundled, enable: false},
		{name: "off, already off", activate: ptr(never), bundled: bundled, enable: false},
	} {
		t.Run(c.name, func(t *testing.T) {
			config := &coregeth.CoreGethChainConfig{}
			if c.activate != nil {
				if err := config.SetECBP1100Transition(c.activate); err != nil {
					t.Fatal(err)
				}
			}
			if c.deactive != nil {
				if err := config.SetECBP1100DeactivateTransition(c.deactive); err != nil {
					t.Fatal(err)
				}
			}
			if err := applyMESSSwitch(config, head, c.bundled, c.enable); err != nil {
				t.Fatal(err)
			}
			got := ecbp1100Status(config, new(big.Int).SetUint64(head), true, false).Enabled
			if got != c.wantEnabled {
				t.Errorf("enabled=%v, want %v (activation %v, deactivation %v)",
					got, c.wantEnabled,
					config.GetECBP1100Transition(), config.GetECBP1100DeactivateTransition())
			}
			if c.wantActivation != nil {
				if a := config.GetECBP1100Transition(); a == nil || *a != *c.wantActivation {
					t.Errorf("activation moved to %v, want %d", messBlock(a), *c.wantActivation)
				}
			}
			// Turning it on must not strand an activation the operator chose.
			if c.enable && c.activate != nil && *c.activate <= head {
				if a := config.GetECBP1100Transition(); a == nil || *a != *c.activate {
					t.Errorf("activation moved from %d to %v, but it was already at or below the head",
						*c.activate, a)
				}
			}
			// Switching MESS on must keep it on while the local head moves backwards. Every
			// site that consults the window consults it against the local head, and a reorg
			// lowers that head. An activation pinned at the block the switch was thrown at
			// would therefore turn MESS off for the depth of any reorg, which is the event it
			// exists for. One block below the head is the shallowest case.
			if c.enable && c.wantEnabled {
				if !ecbp1100Status(config, new(big.Int).SetUint64(head-1), true, false).Enabled {
					t.Errorf("MESS is off one block below the head, so a reorg would disable it "+
						"(activation %v)", config.GetECBP1100Transition())
				}
			}
		})
	}
}

// Block numbers in admin_messStatus are hex quantities. As JSON numbers they lost
// precision in the console, which showed the off switch's block, 18446744073709551614, as
// 18446744073709552000.
func TestEcbp1100StatusEncoding(t *testing.T) {
	config := &coregeth.CoreGethChainConfig{}
	off := uint64(math.MaxUint64 - 1)
	if err := config.SetECBP1100Transition(&off); err != nil {
		t.Fatal(err)
	}
	got, err := json.Marshal(ecbp1100Status(config, big.NewInt(20_000_000), true, false))
	if err != nil {
		t.Fatal(err)
	}
	want := `{"enabled":false,"nodeSwitch":true,"noDisable":false,"activatedAtBlock":"0xfffffffffffffffe","defaultDisabledAtBlock":null,"head":"0x1312d00"}`
	if string(got) != want {
		t.Errorf("status encodes as\n%s\nwant\n%s", got, want)
	}
}

// On a chain past its MESS deactivation block, admin.messActivate does what --mess.activate does.
// With no deactivation given, MESS applies from the block with no end, whatever the block. With
// one given, the later block decides: a block after it turns MESS on again, and one at or before
// it leaves MESS off.
func TestEcbp1100AfterTheDeactivation(t *testing.T) {
	const head = 25_400_000
	for _, c := range []struct {
		in    rpc.BlockNumber
		given bool
		want  bool
	}{
		{rpc.LatestBlockNumber, false, true},
		{rpc.BlockNumber(20_000_000), false, true},
		{rpc.BlockNumber(19_250_000), false, true},
		{rpc.EarliestBlockNumber, false, true},     // before the deactivation, which moves out of reach
		{rpc.LatestBlockNumber, true, true},        // the head, after the deactivation
		{rpc.BlockNumber(20_000_000), true, true},  // after the deactivation
		{rpc.BlockNumber(19_250_000), true, false}, // at it, so neither block is later
		{rpc.EarliestBlockNumber, true, false},     // before it, so the deactivation ends MESS
	} {
		config := &coregeth.CoreGethChainConfig{}
		activation, deactivation := uint64(11_380_000), uint64(19_250_000)
		if err := config.SetECBP1100Transition(&activation); err != nil {
			t.Fatal(err)
		}
		if err := config.SetECBP1100DeactivateTransition(&deactivation); err != nil {
			t.Fatal(err)
		}
		block, err := ecbp1100ActivationBlock(c.in, head)
		if err != nil {
			t.Fatal(err)
		}
		if err := applyMESSActivation(config, block, c.given); err != nil {
			t.Fatal(err)
		}
		if got := ecbp1100Status(config, big.NewInt(head), true, false).Enabled; got != c.want {
			t.Errorf("admin.messActivate(%v) at head %d, deactivation %d given %t: enabled is %v, want %v",
				c.in, head, deactivation, c.given, got, c.want)
		}
		// A deactivation that was given stays where it was; a bundled one moves out of reach.
		want := deactivation
		if !c.given {
			want = math.MaxUint64 - 1
		}
		if got := config.GetECBP1100DeactivateTransition(); got == nil || *got != want {
			t.Errorf("admin.messActivate(%v), deactivation given %t: the deactivation is %v, want %d",
				c.in, c.given, messBlock(got), want)
		}
	}
}

// newMESSTestAPI builds the admin API on the MESS test network at its genesis, with a
// deactivation at block 0, which the head has reached. startupDeactivation is the deactivation a
// flag or a config file line gave, if any.
func newMESSTestAPI(t *testing.T, startupDeactivation *uint64) (*AdminAPI, *core.BlockChain) {
	t.Helper()
	zero := uint64(0)
	genesis := params.DefaultMessNetGenesisBlock()
	config := *params.MessNetConfig // a copy, since the calls change it
	if err := config.SetECBP1100DeactivateTransition(&zero); err != nil {
		t.Fatal(err)
	}
	genesis.Config = &config
	chain, err := core.NewBlockChain(rawdb.NewMemoryDatabase(), nil, genesis, nil, ethash.NewFaker(), vm.Config{}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(chain.Stop)
	eth := &Ethereum{
		blockchain:            chain,
		config:                &ethconfig.Config{OverrideECBP1100Deactivate: startupDeactivation},
		messBundledActivation: chain.Config().GetECBP1100Transition(),
	}
	return NewAdminAPI(eth), chain
}

// TestMessActivateKeepsAGivenDeactivation drives the admin API on a chain at its genesis, where a
// deactivation at block 0 has been reached. MessActivate moves a bundled deactivation out of
// reach, as --mess.activate does, and keeps one given at startup or with MessDeactivate.
func TestMessActivateKeepsAGivenDeactivation(t *testing.T) {
	zero := uint64(0)
	for _, c := range []struct {
		name      string
		atStartup *uint64 // the deactivation a flag or a config file line gave
		viaRPC    bool    // the deactivation given with MessDeactivate first
		wantOn    bool
	}{
		{name: "bundled", wantOn: true},
		{name: "given at startup", atStartup: &zero},
		{name: "given with MessDeactivate", viaRPC: true},
	} {
		t.Run(c.name, func(t *testing.T) {
			api, chain := newMESSTestAPI(t, c.atStartup)
			if c.viaRPC {
				if _, err := api.MessDeactivate(rpc.BlockNumber(0)); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := api.MessActivate(rpc.BlockNumber(0)); err != nil {
				t.Fatal(err)
			}
			cc := chain.Config()
			if got := cc.IsEnabled(cc.GetECBP1100Transition, chain.CurrentBlock().Number); got != c.wantOn {
				t.Errorf("MESS at the head after admin.messActivate(0): %v, want %v (deactivation %v)",
					got, c.wantOn, messBlock(cc.GetECBP1100DeactivateTransition()))
			}
		})
	}
}

// TestApplyMESSNoDisable checks what MessNoDisable(true) does to the window: what
// --mess.nodisable does, MESS on with no end unless a deactivation was given, and an off switch
// back to the network's own activation block.
func TestApplyMESSNoDisable(t *testing.T) {
	const head = 20_000_000
	never := uint64(math.MaxUint64 - 1)
	bundled, deactivation := uint64(11_380_000), uint64(19_250_000)
	for _, c := range []struct {
		name           string
		activation     uint64
		given          bool
		wantActivation uint64
		wantDeact      uint64
		wantOn         bool
	}{
		{name: "the bundled window", activation: bundled, wantActivation: bundled, wantDeact: never, wantOn: true},
		{name: "a deactivation given", activation: bundled, given: true, wantActivation: bundled, wantDeact: deactivation},
		{name: "from an off switch", activation: never, wantActivation: bundled, wantDeact: never, wantOn: true},
	} {
		t.Run(c.name, func(t *testing.T) {
			config := &coregeth.CoreGethChainConfig{}
			a, d := c.activation, deactivation
			if err := config.SetECBP1100Transition(&a); err != nil {
				t.Fatal(err)
			}
			if err := config.SetECBP1100DeactivateTransition(&d); err != nil {
				t.Fatal(err)
			}
			if err := applyMESSNoDisable(config, head, &bundled, c.given); err != nil {
				t.Fatal(err)
			}
			if got := config.GetECBP1100Transition(); got == nil || *got != c.wantActivation {
				t.Errorf("activation %v, want %d", messBlock(got), c.wantActivation)
			}
			if got := config.GetECBP1100DeactivateTransition(); got == nil || *got != c.wantDeact {
				t.Errorf("deactivation %v, want %d", messBlock(got), c.wantDeact)
			}
			if got := ecbp1100Status(config, big.NewInt(head), true, true).Enabled; got != c.wantOn {
				t.Errorf("enabled at the head %v, want %v", got, c.wantOn)
			}
		})
	}
}

// TestMessNoDisable drives MessNoDisable on a running chain. True keeps MESS on through the
// safeguards and turns the window on, keeping a deactivation that was given; false lets the
// safeguards switch MESS off again, and leaves the window alone.
func TestMessNoDisable(t *testing.T) {
	zero := uint64(0)
	for _, c := range []struct {
		name      string
		atStartup *uint64
		wantOn    bool // the window at the head after MessNoDisable(true)
	}{
		{name: "bundled deactivation", wantOn: true},
		{name: "deactivation given at startup", atStartup: &zero},
	} {
		t.Run(c.name, func(t *testing.T) {
			api, chain := newMESSTestAPI(t, c.atStartup)
			status, err := api.MessNoDisable(true)
			if err != nil {
				t.Fatal(err)
			}
			if !status.NoDisable || !chain.IsArtificialFinalityNoDisable() {
				t.Error("MessNoDisable(true) left the safeguards able to switch MESS off")
			}
			cc := chain.Config()
			if got := cc.IsEnabled(cc.GetECBP1100Transition, chain.CurrentBlock().Number); got != c.wantOn {
				t.Errorf("MESS at the head after MessNoDisable(true): %v, want %v", got, c.wantOn)
			}
			act, deact := messBlock(cc.GetECBP1100Transition()), messBlock(cc.GetECBP1100DeactivateTransition())
			if status, err = api.MessNoDisable(false); err != nil {
				t.Fatal(err)
			}
			if status.NoDisable || chain.IsArtificialFinalityNoDisable() {
				t.Error("MessNoDisable(false) left the safeguards unable to switch MESS off")
			}
			if a, d := messBlock(cc.GetECBP1100Transition()), messBlock(cc.GetECBP1100DeactivateTransition()); a != act || d != deact {
				t.Errorf("MessNoDisable(false) moved the window from %v to %v, to %v to %v", act, deact, a, d)
			}
		})
	}
}

// TestMESSCallNames checks that each ECBP-1100 name does exactly what its MESS name does: the
// same call on two identical chains leaves the same state, read with each status name.
func TestMESSCallNames(t *testing.T) {
	block := rpc.BlockNumber(0)
	for _, c := range []struct {
		name       string
		mess, ecbp func(*AdminAPI) error
	}{
		{"activate",
			func(a *AdminAPI) error { _, err := a.MessActivate(block); return err },
			func(a *AdminAPI) error { _, err := a.Ecbp1100(block); return err }},
		{"deactivate",
			func(a *AdminAPI) error { _, err := a.MessDeactivate(block); return err },
			func(a *AdminAPI) error { _, err := a.Ecbp1100Deactivate(block); return err }},
		{"no-disable",
			func(a *AdminAPI) error { _, err := a.MessNoDisable(true); return err },
			func(a *AdminAPI) error { _, err := a.Ecbp1100NoDisable(true); return err }},
	} {
		t.Run(c.name, func(t *testing.T) {
			mess, _ := newMESSTestAPI(t, nil)
			ecbp, _ := newMESSTestAPI(t, nil)
			if err := c.mess(mess); err != nil {
				t.Fatal(err)
			}
			if err := c.ecbp(ecbp); err != nil {
				t.Fatal(err)
			}
			a, _ := json.Marshal(mess.MessStatus())
			b, _ := json.Marshal(ecbp.Ecbp1100Status())
			if string(a) != string(b) {
				t.Errorf("the MESS name leaves %s, the ECBP-1100 name %s", a, b)
			}
		})
	}
}

// TestMESSKeptOffBy checks when admin_messActivate names the deactivation block as what keeps MESS
// off: only when that block, being at or after the activation, decides and has been reached.
func TestMESSKeptOffBy(t *testing.T) {
	ptr := func(n uint64) *uint64 { return &n }
	for _, c := range []struct {
		activation   uint64
		deactivation *uint64
		head         uint64
		want         bool
	}{
		{11_380_000, ptr(19_250_000), 25_000_000, true},  // the deactivation is later, and reached
		{19_250_000, ptr(19_250_000), 25_000_000, true},  // the same block: MESS is off
		{25_000_000, ptr(19_250_000), 25_000_000, false}, // the activation is later: MESS is on
		{30_000_000, ptr(19_250_000), 25_000_000, false}, // the activation is later and ahead: MESS starts there
		{11_380_000, ptr(30_000_000), 25_000_000, false}, // the deactivation is ahead: MESS is on
		{11_380_000, nil, 25_000_000, false},             // no deactivation block
	} {
		if got := messKeptOffBy(c.activation, c.deactivation, c.head); got != c.want {
			t.Errorf("activation %d, deactivation %v, head %d: %v, want %v", c.activation, messBlock(c.deactivation), c.head, got, c.want)
		}
	}
}
