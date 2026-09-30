package eth

import (
	"encoding/json"
	"math"
	"math/big"
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

// On a chain past its MESS deactivation block, admin.ecbp1100 with a block after that
// deactivation turns MESS on again, since the later block decides.
func TestEcbp1100AfterTheDeactivation(t *testing.T) {
	const head = 25_400_000
	for _, c := range []struct {
		in   rpc.BlockNumber
		want bool
	}{
		{rpc.LatestBlockNumber, true},        // the head, after the deactivation
		{rpc.BlockNumber(20_000_000), true},  // after the deactivation
		{rpc.BlockNumber(19_250_000), false}, // at it, so neither block is later
		{rpc.EarliestBlockNumber, false},     // before it, so the deactivation ends MESS
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
		if err := config.SetECBP1100Transition(&block); err != nil {
			t.Fatal(err)
		}
		if got := ecbp1100Status(config, big.NewInt(head), true).Enabled; got != c.want {
			t.Errorf("admin.ecbp1100(%v) at head %d, deactivation %d: enabled is %v, want %v", c.in, head, deactivation, got, c.want)
		}
	}
}
