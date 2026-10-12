package core

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/triedb"
)

func runMESSTest2(t *testing.T, enableMess bool, easyL, hardL, caN int, easyT, hardT int64) (hardHead bool, err error, hard, easy []*types.Block) {
	// Generate the original common chain segment and the two competing forks
	engine := ethash.NewFaker()

	db := rawdb.NewMemoryDatabase()
	genesis := params.DefaultMessNetGenesisBlock()
	genesisB := MustCommitGenesis(db, triedb.NewDatabase(db, nil), genesis)

	chain, err := NewBlockChain(db, nil, genesis, nil, engine, vm.Config{}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer chain.Stop()
	chain.EnableArtificialFinality(enableMess)

	easy, _ = GenerateChain(genesis.Config, genesisB, engine, db, easyL, func(i int, b *BlockGen) {
		b.SetNonce(types.EncodeNonce(uint64(rand.Int63n(math.MaxInt64))))
		b.OffsetTime(easyT)
	})
	commonAncestor := easy[caN-1]
	hard, _ = GenerateChain(genesis.Config, commonAncestor, engine, db, hardL, func(i int, b *BlockGen) {
		b.SetNonce(types.EncodeNonce(uint64(rand.Int63n(math.MaxInt64))))
		b.OffsetTime(hardT)
	})

	if _, err := chain.InsertChain(easy); err != nil {
		t.Fatal(err)
	}
	_, err = chain.InsertChain(hard)
	if err != nil {
		t.Logf("insert hard chain error = %v", err)
	}
	hardHead = chain.CurrentBlock().Hash() == hard[len(hard)-1].Hash()
	return
}

func TestBlockChain_AF_ECBP1100_2(t *testing.T) {
	offsetGreaterDifficulty := int64(-2) // 1..8 = -9..-2
	offsetSameDifficulty := int64(0)     // 9..17 = -1..8
	offsetWorseDifficulty := int64(8)    // 18..

	cases := []struct {
		easyLen, hardLen, commonAncestorN int
		easyOffset, hardOffset            int64
		hardGetsHead, accepted            bool
	}{
		// NOTE: Random coin tosses involved for equivalent difficulty.
		// Short trials for those are skipped.

		{
			1000, 30, 970,
			0, offsetSameDifficulty, // same difficulty
			false, true,
		},

		{
			1000, 1, 999,
			0, offsetWorseDifficulty, // worse! difficulty
			false, true,
		},
		{
			1000, 1, 999,
			0, offsetGreaterDifficulty, // better difficulty
			true, true,
		},
		{
			1000, 5, 995,
			0, offsetGreaterDifficulty,
			true, true,
		},
		{
			1000, 25, 975,
			0, offsetGreaterDifficulty,
			false, true,
		},
		{
			1000, 30, 970,
			0, offsetGreaterDifficulty,
			false, true,
		},
		{
			1000, 50, 950,
			0, offsetGreaterDifficulty,
			false, true,
		},
		{
			1000, 50, 950,
			0, offsetGreaterDifficulty,
			false, true,
		},
		{
			1000, 1000, 900,
			0, offsetGreaterDifficulty,
			true, true,
		},
		{
			1000, 2000, 800,
			0, offsetGreaterDifficulty,
			true, true,
		},
		{
			1000, 2000, 700,
			0, offsetGreaterDifficulty,
			true, true,
		},
		{
			1000, 2000, 700,
			0, offsetGreaterDifficulty,
			true, true,
		},
		{
			1000, 999, 1,
			0, offsetGreaterDifficulty,
			false, true,
		},
		{
			1000, 999, 1,
			0, offsetGreaterDifficulty,
			false, true,
		},
		{
			1000, 500, 500,
			0, offsetGreaterDifficulty,
			false, true,
		},
		{
			1000, 500, 500,
			0, offsetGreaterDifficulty,
			false, true,
		},
		{
			1000, 300, 700,
			0, offsetGreaterDifficulty,
			false, true,
		},
		{
			1000, 600, 700,
			0, offsetGreaterDifficulty,
			true, true,
		},
		// Will pass, takes a long time.
		// {
		// 	5000, 4000, 1000,
		// 	0, -2,
		// 	true, true,
		// },
	}

	for i, c := range cases {
		hardHead, err, hard, easy := runMESSTest2(t, true, c.easyLen, c.hardLen, c.commonAncestorN, c.easyOffset, c.hardOffset)

		ee, hh := easy[len(easy)-1], hard[len(hard)-1]
		rat, _ := new(big.Float).Quo(
			new(big.Float).SetInt(hh.Difficulty()),
			new(big.Float).SetInt(ee.Difficulty()),
		).Float64()

		logf := fmt.Sprintf("case=%d [easy=%d hard=%d ca=%d eo=%d ho=%d] drat=%0.6f span=%v hardHead(w|g)=%v|%v err=%v",
			i,
			c.easyLen, c.hardLen, c.commonAncestorN, c.easyOffset, c.hardOffset,
			rat,
			common.PrettyDuration(time.Second*time.Duration(10*(c.easyLen-c.commonAncestorN))),
			c.hardGetsHead, hardHead, err)

		if (err != nil && c.accepted) || (err == nil && !c.accepted) || (hardHead != c.hardGetsHead) {
			t.Error("FAIL", logf)
		} else {
			t.Log("PASS", logf)
		}
	}
}

/*
TestAFKnownBlock tests that AF functionality works for chain re-insertions.

Chain re-insertions use BlockChain.writeKnownBlockAsHead, where first-pass insertions
will hit writeBlockWithState.

AF needs to be implemented at both sites to prevent re-proposed chains from sidestepping
the AF criteria.
*/
func TestAFKnownBlock(t *testing.T) {
	engine := ethash.NewFaker()

	db := rawdb.NewMemoryDatabase()
	genesis := params.DefaultMessNetGenesisBlock()
	// genesis.Timestamp = 1
	genesisB := MustCommitGenesis(db, triedb.NewDatabase(db, nil), genesis)

	chain, err := NewBlockChain(db, nil, genesis, nil, engine, vm.Config{}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer chain.Stop()
	chain.EnableArtificialFinality(true)

	easy, _ := GenerateChain(genesis.Config, genesisB, engine, db, 1000, func(i int, gen *BlockGen) {
		gen.OffsetTime(0)
	})
	easyN, err := chain.InsertChain(easy)
	if err != nil {
		t.Fatal(err)
	}
	hard, _ := GenerateChain(genesis.Config, easy[easyN-300], engine, db, 300, func(i int, gen *BlockGen) {
		gen.OffsetTime(-7)
	})
	// writeBlockWithState
	if _, err := chain.InsertChain(hard); err != nil {
		t.Error("hard 1 not inserted (should be side)")
	}
	// writeKnownBlockAsHead
	if _, err := chain.InsertChain(hard); err != nil {
		t.Error("hard 2 inserted (will have 'ignored' known blocks, and never tried a reorg)")
	}
	hardHeadHash := hard[len(hard)-1].Hash()
	if chain.CurrentBlock().Hash() == hardHeadHash {
		t.Fatal("hard block got chain head, should be side")
	}
	if h := chain.GetHeaderByHash(hardHeadHash); h == nil {
		t.Fatal("missing hard block (should be imported as side, but still available)")
	}
}

// TestEcbp1100PolynomialV tests the general shape and return values of the ECBP1100 polynomial curve.
// It makes sure domain values above the 'cap' do indeed get limited, as well
// as sanity check some normal domain values.
func TestEcbp1100PolynomialV(t *testing.T) {
	cases := []struct {
		block, ag int64
	}{
		{100, 1},
		{300, 2},
		{500, 5},
		{1000, 16},
		{2000, 31},
		{10000, 31},
		{1e9, 31},
	}
	for i, c := range cases {
		y := ecbp1100PolynomialV(big.NewInt(c.block * 13))
		y.Div(y, ecbp1100PolynomialVCurveFunctionDenominator)
		if c.ag != y.Int64() {
			t.Fatal("mismatch", i)
		}
	}
}

// TestECBP1100NoDisableNotice checks that --mess.nodisable's log line says what is true of the MESS
// window, and never names block 18446744073709551614, which no chain reaches.
func TestECBP1100NoDisableNotice(t *testing.T) {
	u := func(n uint64) *uint64 { return &n }
	const (
		unreachable = "MESS is set to on, with its deactivation moved to an unreachable block"
		stops       = "MESS still stops applying at the deactivation block"
	)
	for _, c := range []struct {
		name                     string
		activation, deactivation *uint64
		want                     string // a substring, or "" for no notice
	}{
		{"no deactivation", u(11_380_000), nil, ""},
		{"deactivation out of reach", u(11_380_000), u(math.MaxUint64 - 1), unreachable},
		{"deactivation at math.MaxUint64", u(11_380_000), u(math.MaxUint64), unreachable},
		{"deactivation after the activation", u(11_380_000), u(19_250_000), stops},
		{"activation after the deactivation", u(25_000_000), u(20_000_000), ""},
		{"both at one block", u(20_000_000), u(20_000_000), ""},
	} {
		got := ecbp1100NoDisableNotice(c.activation, c.deactivation)
		if (c.want == "") != (got == "") || !strings.Contains(got, c.want) {
			t.Errorf("%s: notice %q, want one containing %q", c.name, got, c.want)
		}
		if strings.Contains(got, "18446744073709551614") {
			t.Errorf("%s: notice names the unreachable block: %q", c.name, got)
		}
	}
}

// TestArtificialFinalityNoDisable checks the no-disable setting on a running chain. On, the
// safeguards no longer switch MESS off once it is on, and the setting does not switch it on
// either; off, they switch it off again. It is changed from one goroutine while another switches
// MESS, as the admin API and the sync loop do, which the race detector checks.
func TestArtificialFinalityNoDisable(t *testing.T) {
	chain, err := NewBlockChain(rawdb.NewMemoryDatabase(), nil, params.DefaultMessNetGenesisBlock(), nil, ethash.NewFaker(), vm.Config{}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer chain.Stop()

	chain.ArtificialFinalityNoDisable(1)
	if !chain.IsArtificialFinalityNoDisable() {
		t.Fatal("no-disable is off after ArtificialFinalityNoDisable(1)")
	}
	if chain.IsArtificialFinalityEnabled() {
		t.Fatal("no-disable switched MESS on by itself")
	}
	chain.EnableArtificialFinality(true)
	chain.EnableArtificialFinality(false, "reason", "low peers")
	if !chain.IsArtificialFinalityEnabled() {
		t.Error("no-disable on: a safeguard switched MESS off")
	}
	chain.ArtificialFinalityNoDisable(0)
	if chain.IsArtificialFinalityNoDisable() {
		t.Error("no-disable is on after ArtificialFinalityNoDisable(0)")
	}
	chain.EnableArtificialFinality(false, "reason", "low peers")
	if chain.IsArtificialFinalityEnabled() {
		t.Error("no-disable off: a safeguard did not switch MESS off")
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 1000; i++ {
			chain.EnableArtificialFinality(i%2 == 0, "reason", "test")
		}
	}()
	for i := 0; i < 1000; i++ {
		chain.ArtificialFinalityNoDisable(int32(i % 2))
	}
	<-done
}
