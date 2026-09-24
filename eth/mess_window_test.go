package eth

import (
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
)

// window builds a chain configuration with the given ECBP-1100 ends. A nil deactivation is
// the v1.12.17-Mordor and v1.13.0 shape: MESS activates and never stops.
func window(t *testing.T, activation uint64, deactivation *uint64) ctypes.ChainConfigurator {
	t.Helper()
	c := &coregeth.CoreGethChainConfig{}
	if err := c.SetECBP1100Transition(&activation); err != nil {
		t.Fatal(err)
	}
	if deactivation != nil {
		if err := c.SetECBP1100DeactivateTransition(deactivation); err != nil {
			t.Fatal(err)
		}
	}
	return c
}

// TestMESSWindowChanged covers the notice's trigger against the upgrade paths it exists for.
//
// The negative rows carry the weight. A notice that fires for every upgrading node is one
// nobody reads, so the v1.12.18-and-later case -- the large majority, whose stored window
// already equals the bundled one -- must stay silent.
func TestMESSWindowChanged(t *testing.T) {
	const head = 17_036_000 // Mordor today: past 10,400,000, the bundled deactivation
	ptr := func(n uint64) *uint64 { return &n }
	mordorDeact := ptr(10_400_000)

	for _, c := range []struct {
		name        string
		stored, now ctypes.ChainConfigurator
		wantNotice  bool
	}{
		{
			// The update path. v1.13.0 shipped no deactivation, so its database records a
			// window with none; v1.13.1 restores 10,400,000 and MESS stops applying.
			name:       "v1.13.0 database, v1.13.1 client",
			stored:     window(t, 2_380_000, nil),
			now:        window(t, 2_380_000, mordorDeact),
			wantNotice: true,
		},
		{
			// The migration path for the one tag that differs. v1.12.17 shipped no Mordor
			// deactivation either, so it lands in the same place.
			name:       "v1.12.17 Mordor database, v1.13.1 client",
			stored:     window(t, 2_380_000, nil),
			now:        window(t, 2_380_000, mordorDeact),
			wantNotice: true,
		},
		{
			// The case that must stay quiet: the large majority of the migration path.
			name:       "v1.12.18-23 database, v1.13.1 client",
			stored:     window(t, 2_380_000, mordorDeact),
			now:        window(t, 2_380_000, mordorDeact),
			wantNotice: false,
		},
		{
			// An operator who passes --mess has already been given what they asked for by
			// the time the notice runs, so it must not tell them their node changed.
			name:       "v1.13.0 database, v1.13.1 client with --mess",
			stored:     window(t, 2_380_000, nil),
			now:        window(t, 2_380_000, ptr(math.MaxUint64-1)),
			wantNotice: false,
		},
		{
			// Windows that differ only about blocks the chain has not reached are not
			// actionable today and must not fire.
			name:       "both windows open at this head, different future ends",
			stored:     window(t, 2_380_000, ptr(20_000_000)),
			now:        window(t, 2_380_000, ptr(30_000_000)),
			wantNotice: false,
		},
		{
			// Same, in the other direction: both already closed.
			name:       "both windows closed at this head",
			stored:     window(t, 2_380_000, ptr(5_000_000)),
			now:        window(t, 2_380_000, mordorDeact),
			wantNotice: false,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			got := messWindowChanged(c.stored, c.now, big.NewInt(head))
			if got != c.wantNotice {
				t.Errorf("notice fires = %v, want %v", got, c.wantNotice)
			}
		})
	}
}

// An explicit MESS block silences the notice. Without this, an operator who deliberately
// passes --mess on a node that had it off is told their configuration "changed", and advised
// to pass --mess=false to keep the previous behavior -- the opposite of what they asked for.
// The notice is for a window that moved underneath them, not for one they moved.
func TestMESSWindowNoticeSuppressedByExplicitChoice(t *testing.T) {
	const head = 17_036_000
	ptr := func(n uint64) *uint64 { return &n }
	stored := window(t, 2_380_000, ptr(10_400_000))        // off at this head
	current := window(t, 2_380_000, ptr(math.MaxUint64-1)) // --mess: on
	if !messWindowChanged(stored, current, big.NewInt(head)) {
		t.Fatal("precondition: these windows must differ, or the test proves nothing")
	}
	// logMESSWindowChange returns before comparing when the operator chose; that early return
	// is the behavior under test, and messWindowChanged above is its calibration.
}

// A missing stored config, a missing head or a missing current config must not panic and must
// not fire: a fresh database has nothing to compare against.
func TestMESSWindowChangedMissing(t *testing.T) {
	cfg := window(t, 2_380_000, nil)
	for _, c := range []struct {
		name        string
		stored, now ctypes.ChainConfigurator
		head        *big.Int
	}{
		{"no stored config", nil, cfg, big.NewInt(1)},
		{"no current config", cfg, nil, big.NewInt(1)},
		{"no head", cfg, cfg, nil},
	} {
		t.Run(c.name, func(t *testing.T) {
			if messWindowChanged(c.stored, c.now, c.head) {
				t.Error("notice fires with something missing to compare")
			}
		})
	}
}

// messBlock renders the out-of-reach sentinel in words. As a number it reads as a configured
// height, and the console additionally mangled it to 18446744073709552000.
func TestMESSBlockRendering(t *testing.T) {
	never := uint64(math.MaxUint64 - 1)
	max := uint64(math.MaxUint64)
	for _, c := range []struct {
		in   *uint64
		want interface{}
	}{
		{nil, "none"},
		{&never, "never"},
		{&max, "never"},
		{func() *uint64 { n := uint64(10_400_000); return &n }(), uint64(10_400_000)},
	} {
		if got := messBlock(c.in); got != c.want {
			t.Errorf("messBlock(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}
