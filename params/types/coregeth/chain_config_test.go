package coregeth

import (
	"encoding/json"
	"math"
	"math/big"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
)

// nolint:unused
var testConfig = &CoreGethChainConfig{
	NetworkID:  1,
	Ethash:     new(ctypes.EthashConfig),
	ChainID:    big.NewInt(61),
	EIP2FBlock: big.NewInt(1150000),
	EIP7FBlock: big.NewInt(1150000),
	// DAOForkBlock:        big.NewInt(1920000),
	EIP150Block:        big.NewInt(2500000),
	EIP155Block:        big.NewInt(3000000),
	EIP160FBlock:       big.NewInt(3000000),
	EIP161FBlock:       big.NewInt(8772000),
	EIP170FBlock:       big.NewInt(8772000),
	EIP100FBlock:       big.NewInt(8772000),
	EIP140FBlock:       big.NewInt(8772000),
	EIP198FBlock:       big.NewInt(8772000),
	EIP211FBlock:       big.NewInt(8772000),
	EIP212FBlock:       big.NewInt(8772000),
	EIP213FBlock:       big.NewInt(8772000),
	EIP214FBlock:       big.NewInt(8772000),
	EIP658FBlock:       big.NewInt(8772000),
	EIP145FBlock:       big.NewInt(9573000),
	EIP1014FBlock:      big.NewInt(9573000),
	EIP1052FBlock:      big.NewInt(9573000),
	EIP1283FBlock:      nil,
	PetersburgBlock:    nil, // Un1283
	EIP2200FBlock:      nil, // RePetersburg (== re-1283)
	DisposalBlock:      big.NewInt(5900000),
	ECIP1017FBlock:     big.NewInt(5000000),
	ECIP1017EraRounds:  big.NewInt(5000000),
	ECIP1010PauseBlock: big.NewInt(3000000),
	ECIP1010Length:     big.NewInt(2000000),
	RequireBlockHashes: map[uint64]common.Hash{
		2500000: common.HexToHash("0xca12c63534f565899681965528d536c52cb05b7c48e269c2a6cb77ad864d878a"),
	},
}

func TestCoreGethChainConfig_String(t *testing.T) {
	t.Skip("(noop) development use only")
	t.Log(testConfig.String())
}

func TestCoreGethChainConfig_ECBP1100Deactivate(t *testing.T) {
	var _testConfig = &CoreGethChainConfig{}
	*_testConfig = *testConfig

	activate := uint64(100)
	deactivate := uint64(200)
	_testConfig.SetECBP1100Transition(&activate)
	_testConfig.SetECBP1100DeactivateTransition(&deactivate)

	n := uint64(10)
	bigN := new(big.Int).SetUint64(n)
	if _testConfig.IsEnabled(_testConfig.GetECBP1100Transition, bigN) {
		t.Errorf("ECBP1100 should be not yet be activated at block %d", n)
	}

	n = uint64(100)
	bigN = new(big.Int).SetUint64(n)
	if !_testConfig.IsEnabled(_testConfig.GetECBP1100Transition, bigN) {
		t.Errorf("ECBP1100 should be activated at block %d", n)
	}

	n = uint64(110)
	bigN = new(big.Int).SetUint64(n)
	if !_testConfig.IsEnabled(_testConfig.GetECBP1100Transition, bigN) {
		t.Errorf("ECBP1100 should be activated at block %d", n)
	}

	n = uint64(200)
	bigN = new(big.Int).SetUint64(n)
	if _testConfig.IsEnabled(_testConfig.GetECBP1100Transition, bigN) {
		t.Errorf("ECBP1100 should be deactivated at block %d", n)
	}

	n = uint64(210)
	bigN = new(big.Int).SetUint64(n)
	if _testConfig.IsEnabled(_testConfig.GetECBP1100Transition, bigN) {
		t.Errorf("ECBP1100 should be deactivated at block %d", n)
	}
}

// Transition heights cross the configurator as uint64 and are stored as *big.Int. A
// conversion through int64 at that boundary wraps every height above math.MaxInt64: the
// height --mess=false assigns to put MESS out of reach read back as block 2, and left MESS
// enabled from there.

func TestCoreGethChainConfig_TransitionSettersAboveMaxInt64(t *testing.T) {
	// A getter answers only for its own consensus engine, so each setter is tried against a
	// configuration of each engine until one stores a small height as given.
	engines := []func() *CoreGethChainConfig{
		func() *CoreGethChainConfig { return &CoreGethChainConfig{Ethash: new(ctypes.EthashConfig)} },
		func() *CoreGethChainConfig { return &CoreGethChainConfig{Lyra2: new(ctypes.Lyra2Config)} },
	}
	uint64Ptr := reflect.TypeOf((*uint64)(nil))
	typ := reflect.TypeOf(&CoreGethChainConfig{})
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
		var engine func() *CoreGethChainConfig
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
	// A setter that does not store even a small height as given is outside the loop above.
	// Each is pinned here, so none drops out silently, and covered by its own test:
	// the continue transition needs a pause block to measure from.
	if want := []string{"SetEthashECIP1010ContinueTransition"}; !slices.Equal(skipped, want) {
		t.Errorf("setters outside the round-trip check are %q, want %q", skipped, want)
	}
	t.Logf("%d transition setters checked above math.MaxInt64", checked)
}

func setThenGet(c *CoreGethChainConfig, set, get reflect.Method, height uint64) *uint64 {
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

// The continue height is stored as a length past the pause, so it takes its own path.
func TestCoreGethChainConfig_ECIP1010ContinueAboveMaxInt64(t *testing.T) {
	c := &CoreGethChainConfig{Ethash: new(ctypes.EthashConfig), ECIP1010PauseBlock: big.NewInt(3_000_000)}
	height := uint64(math.MaxUint64 - 1)
	if err := c.SetEthashECIP1010ContinueTransition(&height); err != nil {
		t.Fatal(err)
	}
	if got := c.GetEthashECIP1010ContinueTransition(); got == nil || *got != height {
		t.Errorf("continue transition reads back %v, want %d", deref(got), height)
	}
}

func TestCoreGethChainConfig_IsEnabledAboveMaxInt64(t *testing.T) {
	outOfReach := uint64(math.MaxUint64 - 1)

	c := &CoreGethChainConfig{}
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
	c = &CoreGethChainConfig{}
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
	c = &CoreGethChainConfig{}
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
	c = &CoreGethChainConfig{EIP155Block: new(big.Int).SetUint64(outOfReach)}
	if c.IsEnabled(c.GetEIP155Transition, new(big.Int).SetUint64(23_000_000)) {
		t.Errorf("EIP-155 decoded at %d: enabled at block 23000000", outOfReach)
	}
}

// A JSON chain configuration can hold a block number outside the uint64 range. Keeping only its
// low 64 bits read 2^64+2 as block 2 and -11380000 as block 11380000. Clamped, a number below
// zero is already reached and one of 2^64 or more is never reached, as a comparison against a
// real block has it.
func TestCoreGethChainConfig_BlockNumbersOutsideUint64(t *testing.T) {
	var c CoreGethChainConfig
	in := `{"eip155Block": -11380000, "eip160Block": 18446744073709551618,
		"ecbp1100FBlock": 18446744073709551616, "ecbp1100DeactivateFBlockFBlock": -1}`
	if err := json.Unmarshal([]byte(in), &c); err != nil {
		t.Fatal(err)
	}
	for _, x := range []struct {
		key  string
		get  func() *uint64
		want uint64
	}{
		{"eip155Block", c.GetEIP155Transition, 0},
		{"eip160Block", c.GetEIP160Transition, math.MaxUint64},
		{"ecbp1100FBlock", c.GetECBP1100Transition, math.MaxUint64},
		{"ecbp1100DeactivateFBlockFBlock", c.GetECBP1100DeactivateTransition, 0},
	} {
		if got := x.get(); got == nil || *got != x.want {
			t.Errorf("%s reads back %v, want %d", x.key, deref(got), x.want)
		}
	}
	if !c.IsEnabled(c.GetEIP155Transition, big.NewInt(0)) {
		t.Error("eip155Block below zero: not enabled at block 0")
	}
	if c.IsEnabled(c.GetEIP160Transition, big.NewInt(23_000_000)) {
		t.Error("eip160Block above math.MaxUint64: enabled at block 23000000")
	}
}

// TestCoreGethChainConfig_ECBP1100LaterBlockWins pins how the MESS activation and deactivation
// blocks combine: the later of the two decides from that block on, so an activation after the
// deactivation turns MESS on again.
func TestCoreGethChainConfig_ECBP1100LaterBlockWins(t *testing.T) {
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
		var cfg CoreGethChainConfig
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
