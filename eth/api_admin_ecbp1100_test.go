package eth

import (
	"encoding/json"
	"math"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/rpc"
)

// A block tag is a negative rpc.BlockNumber. Converted straight to a height, "latest"
// became math.MaxUint64-1: once heights stopped wrapping, admin.ecbp1100("latest") would
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
// the way eth.New applies one, without standing up a node. The switch has to reach whichever
// end is holding the window shut: admin_ecbp1100 sets the activation alone and so could not
// turn MESS on at all once the head was past the deactivation block, which is the state the
// bundled ECIP-1110 configuration leaves every Classic and Mordor node in.
func TestApplyMESSSwitch(t *testing.T) {
	const head = 20_000_000
	never := uint64(math.MaxUint64 - 1)
	ptr := func(n uint64) *uint64 { return &n }

	for _, c := range []struct {
		name               string
		activate, deactive *uint64
		enable             bool
		wantEnabled        bool
	}{
		// The case admin_ecbp1100 alone cannot reach: the window closed behind the head.
		{name: "on, past a deactivation", activate: ptr(11_380_000), deactive: ptr(19_250_000), enable: true, wantEnabled: true},
		{name: "on, from an off switch", activate: ptr(never), enable: true, wantEnabled: true},
		{name: "on, window already open", activate: ptr(11_380_000), enable: true, wantEnabled: true},
		{name: "on, nothing configured", enable: true, wantEnabled: true},
		{name: "off, from on", activate: ptr(11_380_000), enable: false},
		{name: "off, past a deactivation", activate: ptr(11_380_000), deactive: ptr(19_250_000), enable: false},
		{name: "off, already off", activate: ptr(never), enable: false},
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
			if err := applyMESSSwitch(config, head, c.enable); err != nil {
				t.Fatal(err)
			}
			got := ecbp1100Status(config, new(big.Int).SetUint64(head), true).Enabled
			if got != c.wantEnabled {
				t.Errorf("enabled=%v, want %v (activation %v, deactivation %v)",
					got, c.wantEnabled,
					config.GetECBP1100Transition(), config.GetECBP1100DeactivateTransition())
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
			// lowers that head -- so an activation pinned at the block the switch was thrown
			// at would turn MESS off for the depth of any reorg, which is the event it
			// exists for. One block below the head is the shallowest case.
			if c.enable && c.wantEnabled {
				if !ecbp1100Status(config, new(big.Int).SetUint64(head-1), true).Enabled {
					t.Errorf("MESS is off one block below the head, so a reorg would disable it "+
						"(activation %v)", config.GetECBP1100Transition())
				}
			}
		})
	}
}

// Block numbers in admin_ecbp1100Status are hex quantities. As JSON numbers they lost
// precision in the console, which showed the off switch's block, 18446744073709551614, as
// 18446744073709552000.
func TestEcbp1100StatusEncoding(t *testing.T) {
	config := &coregeth.CoreGethChainConfig{}
	off := uint64(math.MaxUint64 - 1)
	if err := config.SetECBP1100Transition(&off); err != nil {
		t.Fatal(err)
	}
	got, err := json.Marshal(ecbp1100Status(config, big.NewInt(20_000_000), true))
	if err != nil {
		t.Fatal(err)
	}
	want := `{"enabled":false,"nodeSwitch":true,"activatedAtBlock":"0xfffffffffffffffe","defaultDisabledAtBlock":null,"head":"0x1312d00"}`
	if string(got) != want {
		t.Errorf("status encodes as\n%s\nwant\n%s", got, want)
	}
}
