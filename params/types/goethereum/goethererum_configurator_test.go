package goethereum

import (
	"encoding/json"
	"math"
	"math/big"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/params/types/ctypes"
)

func TestChainConfig_converting(t *testing.T) {
	var c interface{} = &ChainConfig{}
	fromChainer := c.(ctypes.ChainConfigurator)

	if _, ok := reflect.TypeOf(fromChainer).Elem().FieldByName("Converting"); ok {
		reflect.ValueOf(fromChainer).Elem().Field(0).SetBool(true)
	}
}

// This configurator carries the same uint64-to-*big.Int boundary as coregeth's, and so the
// same wrap above math.MaxInt64; a chain configured in this format reaches it the same way.

func TestChainConfig_TransitionSettersAboveMaxInt64(t *testing.T) {
	// A getter answers only for its own consensus engine, so each setter is tried against a
	// configuration of each engine until one stores a small height as given.
	engines := []func() *ChainConfig{
		func() *ChainConfig { return &ChainConfig{Ethash: new(ctypes.EthashConfig)} },
		func() *ChainConfig { return &ChainConfig{Lyra2: new(ctypes.Lyra2Config)} },
	}
	uint64Ptr := reflect.TypeOf((*uint64)(nil))
	typ := reflect.TypeOf(&ChainConfig{})
	var checked int
	var skipped []string
	for i := 0; i < typ.NumMethod(); i++ {
		set := typ.Method(i)
		if !strings.HasPrefix(set.Name, "Set") || !strings.HasSuffix(set.Name, "Transition") {
			continue
		}
		get, ok := typ.MethodByName("Get" + strings.TrimPrefix(set.Name, "Set"))
		if !ok || set.Type.NumIn() != 2 || set.Type.In(1) != uint64Ptr ||
			get.Type.NumOut() != 1 || get.Type.Out(0) != uint64Ptr {
			continue
		}
		var engine func() *ChainConfig
		for _, e := range engines {
			if got := setThenGet(e(), set, get, 12_345); got != nil && *got == 12_345 {
				engine = e
				break
			}
		}
		if engine == nil {
			skipped = append(skipped, set.Name)
			continue
		}
		// math.MaxInt64+1 is not among these: through the int64 conversion it happens to read
		// back unchanged, so it cannot detect the wrap.
		for _, height := range []uint64{math.MaxInt64 + 2, math.MaxUint64 - 1, math.MaxUint64} {
			if got := setThenGet(engine(), set, get, height); got == nil || *got != height {
				t.Errorf("%s(%d): %s reads back %v", set.Name, height, get.Name, deref(got))
			}
		}
		checked++
	}
	// A setter that does not store even a small height as given is outside the loop above:
	// this format does not carry that setting, or derives it from others. Each is pinned
	// here, so none drops out silently.
	if want := []string{
		"SetEIP1153Transition", "SetEIP2200DisableTransition", "SetEIP2315Transition", "SetEIP2537Transition",
		"SetEIP3651Transition", "SetEIP3855Transition", "SetEIP3860Transition", "SetEIP4399Transition",
		"SetEIP4788Transition", "SetEIP4844Transition", "SetEIP4895Transition", "SetEIP5656Transition",
		"SetEIP6049Transition", "SetEIP6780Transition", "SetEIP7516Transition",
		"SetEthashECIP1010ContinueTransition", "SetEthashECIP1010PauseTransition", "SetEthashECIP1017Transition",
		"SetEthashECIP1041Transition", "SetEthashECIP1099Transition", "SetVerkleTransition",
	}; !slices.Equal(skipped, want) {
		t.Errorf("setters outside the round-trip check are %q, want %q", skipped, want)
	}
	t.Logf("%d transition setters checked above math.MaxInt64", checked)
}

func setThenGet(c *ChainConfig, set, get reflect.Method, height uint64) *uint64 {
	v := reflect.ValueOf(c)
	if err, _ := set.Func.Call([]reflect.Value{v, reflect.ValueOf(&height)})[0].Interface().(error); err != nil {
		return nil
	}
	return get.Func.Call([]reflect.Value{v})[0].Interface().(*uint64)
}

func deref(p *uint64) interface{} {
	if p == nil {
		return nil
	}
	return *p
}

func TestChainConfig_IsEnabledAboveMaxInt64(t *testing.T) {
	outOfReach := uint64(math.MaxUint64 - 1)

	c := &ChainConfig{}
	if err := c.SetECBP1100Transition(&outOfReach); err != nil {
		t.Fatal(err)
	}
	for _, h := range []uint64{11_380_000, 23_000_000} {
		if c.IsEnabled(c.GetECBP1100Transition, new(big.Int).SetUint64(h)) {
			t.Errorf("activation at %d: enabled at block %d", outOfReach, h)
		}
	}

	// Control: a height in range still activates from exactly that block.
	activation := uint64(11_380_000)
	c = &ChainConfig{}
	if err := c.SetECBP1100Transition(&activation); err != nil {
		t.Fatal(err)
	}
	for h, want := range map[uint64]bool{activation - 1: false, activation: true, 23_000_000: true} {
		if got := c.IsEnabled(c.GetECBP1100Transition, new(big.Int).SetUint64(h)); got != want {
			t.Errorf("activation at %d: enabled at block %d is %v, want %v", activation, h, got, want)
		}
	}

	// A deactivation out of reach leaves the window open.
	if err := c.SetECBP1100DeactivateTransition(&outOfReach); err != nil {
		t.Fatal(err)
	}
	if !c.IsEnabled(c.GetECBP1100Transition, new(big.Int).SetUint64(23_000_000)) {
		t.Errorf("deactivation at %d: disabled at block 23000000", outOfReach)
	}

	// With a deactivation set, the activation is compared in a second place; out of reach
	// there too, it keeps the window shut.
	deactivation := uint64(20_000_000)
	c = &ChainConfig{}
	if err := c.SetECBP1100Transition(&outOfReach); err != nil {
		t.Fatal(err)
	}
	if err := c.SetECBP1100DeactivateTransition(&deactivation); err != nil {
		t.Fatal(err)
	}
	for _, h := range []uint64{11_380_000, deactivation - 1} {
		if c.IsEnabled(c.GetECBP1100Transition, new(big.Int).SetUint64(h)) {
			t.Errorf("activation at %d, deactivation at %d: enabled at block %d", outOfReach, deactivation, h)
		}
	}

	// A chainspec decodes a height straight into the field, so the comparison is reached
	// without any setter.
	c = &ChainConfig{EIP155Block: new(big.Int).SetUint64(outOfReach)}
	if c.IsEnabled(c.GetEIP155Transition, new(big.Int).SetUint64(23_000_000)) {
		t.Errorf("EIP-155 decoded at %d: enabled at block 23000000", outOfReach)
	}
}

// A JSON chain configuration can hold a block number outside the uint64 range. Keeping only its
// low 64 bits read 2^64+2 as block 2 and -11380000 as block 11380000. Clamped, a number below
// zero is already reached and one of 2^64 or more is never reached, as a comparison against a
// real block has it.
func TestChainConfig_BlockNumbersOutsideUint64(t *testing.T) {
	var c ChainConfig
	if err := json.Unmarshal([]byte(`{"eip155Block": -11380000, "byzantiumBlock": 18446744073709551618}`), &c); err != nil {
		t.Fatal(err)
	}
	for _, x := range []struct {
		key  string
		get  func() *uint64
		want uint64
	}{
		{"eip155Block", c.GetEIP155Transition, 0},
		{"byzantiumBlock", c.GetEIP140Transition, math.MaxUint64},
	} {
		if got := x.get(); got == nil || *got != x.want {
			t.Errorf("%s reads back %v, want %d", x.key, deref(got), x.want)
		}
	}
	if !c.IsEnabled(c.GetEIP155Transition, big.NewInt(0)) {
		t.Error("eip155Block below zero: not enabled at block 0")
	}
	if c.IsEnabled(c.GetEIP140Transition, big.NewInt(23_000_000)) {
		t.Error("byzantiumBlock above math.MaxUint64: enabled at block 23000000")
	}
}

// TestChainConfig_ECBP1100LaterBlockWins pins how the MESS activation and deactivation blocks
// combine: the later of the two decides from that block on, so an activation after the deactivation
// turns MESS on again.
func TestChainConfig_ECBP1100LaterBlockWins(t *testing.T) {
	for _, c := range []struct {
		name                     string
		activation, deactivation uint64
		want                     map[uint64]bool
	}{
		{"activation first", 15_000_000, 20_000_000,
			map[uint64]bool{14_999_999: false, 15_000_000: true, 19_999_999: true, 20_000_000: false, 25_000_000: false}},
		{"activation after the deactivation", 25_000_000, 20_000_000,
			map[uint64]bool{15_000_000: false, 20_000_000: false, 24_999_999: false, 25_000_000: true, 30_000_000: true}},
		{"both at one block", 20_000_000, 20_000_000,
			map[uint64]bool{19_999_999: false, 20_000_000: false, 25_000_000: false}},
	} {
		var cfg ChainConfig
		activation, deactivation := c.activation, c.deactivation
		if err := cfg.SetECBP1100Transition(&activation); err != nil {
			t.Fatal(err)
		}
		if err := cfg.SetECBP1100DeactivateTransition(&deactivation); err != nil {
			t.Fatal(err)
		}
		for h, want := range c.want {
			if got := cfg.IsEnabled(cfg.GetECBP1100Transition, new(big.Int).SetUint64(h)); got != want {
				t.Errorf("%s (activation %d, deactivation %d): enabled at block %d is %v, want %v",
					c.name, activation, deactivation, h, got, want)
			}
		}
	}
}
