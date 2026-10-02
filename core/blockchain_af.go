package core

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

// errReorgFinality represents an error caused by artificial finality mechanisms.
var errReorgFinality = errors.New("finality-enforced invalid new chain")

// ArtificialFinalityNoDisable overrides toggling of AF features, forcing it on.
// n  = 1 : ON
// n != 1 : OFF
//
// It can be called while the chain runs, as the admin API does, since the setting is read
// atomically wherever the node would switch MESS off.
func (bc *BlockChain) ArtificialFinalityNoDisable(n int32) {
	bc.artificialFinalityNoDisable.Store(n)

	if n == 1 {
		log.Warn("Deactivating ECBP1100 (MESS) safety mechanisms", "always on", true)
		deactivateTransition := bc.chainConfig.GetECBP1100DeactivateTransition()
		logActivationBlockRaw := bc.chainConfig.GetECBP1100Transition()
		if notice := ecbp1100NoDisableNotice(logActivationBlockRaw, deactivateTransition); notice != "" {
			// Log the activation block as well as the deactivation block.
			// Context is nice to have for the user.
			var logActivationBlock uint64
			if logActivationBlockRaw == nil {
				// panic("impossible")
				logActivationBlock = *deactivateTransition
			} else {
				logActivationBlock = *logActivationBlockRaw
			}
			logCtx := []interface{}{"ECBP1100 activation block", logActivationBlock}
			if *deactivateTransition < ECBP1100Unreachable {
				logCtx = append(logCtx, "ECBP1100 deactivation block", *deactivateTransition)
			}
			log.Warn(notice, logCtx...)
		}
	} else {
		log.Info("ECBP1100 (MESS) safety mechanisms apply again")
	}
}

// IsArtificialFinalityNoDisable reports whether the setting ArtificialFinalityNoDisable sets is
// on: once MESS is on, the low-peer-count and stale-head safeguards no longer switch it off.
func (bc *BlockChain) IsArtificialFinalityNoDisable() bool {
	return bc.artificialFinalityNoDisable.Load() == 1
}

// ECBP1100Unreachable is the block a MESS setting moves one end of the window to, so that end
// never takes effect: --mess=false moves the activation there, and --mess, --mess.activate and
// --mess.nodisable move the deactivation there when none is given. No chain reaches it.
const ECBP1100Unreachable uint64 = math.MaxUint64 - 1

// ecbp1100NoDisableNotice is what --mess.nodisable adds to the log about the MESS window, or ""
// when there is nothing to add. A deactivation block after the activation ends the window the
// flag keeps on; an unreachable one, or one before the activation, does not.
func ecbp1100NoDisableNotice(activation, deactivation *uint64) string {
	switch {
	case deactivation == nil:
		return ""
	case *deactivation >= ECBP1100Unreachable:
		return "MESS is set to on, with its deactivation moved to an unreachable block. " +
			"--mess.nodisable keeps it on once enabled, bypassing its low-peer-count and stale-head safeguards."
	case activation != nil && *activation >= *deactivation:
		// The later block decides, so MESS runs from the activation with no end, or, at the
		// same block, never runs; the node warns about that case when it reads the flags.
		return ""
	default:
		return `An ECBP1100 (MESS) deactivation block is set together with --mess.nodisable.
--mess.nodisable keeps MESS on once enabled, bypassing its low-peer-count and stale-head safeguards, but MESS still stops applying at the deactivation block.
`
	}
}

// EnableArtificialFinality enables and disable artificial finality features for the blockchain.
// Currently toggled features include:
// - ECBP1100-MESS: modified exponential subject scoring
//
// This level of activation works BELOW the chain configuration for any of the
// potential features. eg. If ECBP1100 is not activated at the chain config x block number,
// then calling bc.EnableArtificialFinality(true) will be a noop.
// The method is idempotent.
func (bc *BlockChain) EnableArtificialFinality(enable bool, logValues ...interface{}) {
	// Short circuit if AF state is enabled and nodisable=true.
	if bc.IsArtificialFinalityNoDisable() && bc.IsArtificialFinalityEnabled() && !enable {
		log.Warn("Preventing disable artificial finality", "enabled", true, "nodisable", true)
		return
	}

	// Store enable/disable value regardless of config activation.
	var statusLog string
	if enable {
		statusLog = "Enabled"
		atomic.StoreInt32(&bc.artificialFinalityEnabledStatus, 1)
	} else {
		statusLog = "Disabled"
		atomic.StoreInt32(&bc.artificialFinalityEnabledStatus, 0)
	}
	if !bc.chainConfig.IsEnabled(bc.chainConfig.GetECBP1100Transition, bc.CurrentHeader().Number) {
		// Don't log anything if the config hasn't enabled it yet.
		return
	}
	logFn := log.Warn // Deactivated and enabled
	if enable {
		logFn = log.Info // Activated and enabled
	}
	logFn(fmt.Sprintf("%s artificial finality features", statusLog), logValues...)
}

// IsArtificialFinalityEnabled returns the status of the blockchain's artificial
// finality feature setting.
// This status is agnostic of feature activation by chain configuration.
func (bc *BlockChain) IsArtificialFinalityEnabled() bool {
	return atomic.LoadInt32(&bc.artificialFinalityEnabledStatus) == 1
}

// ecbp1100 implements the "MESS" artificial finality mechanism
// "Modified Exponential Subjective Scoring" used to prefer known chain segments
// over later-to-come counterparts, especially proposed segments stretching far into the past.
func ecbp1100(commonAncestor, current, proposed *types.Header, getTDFunc func(common.Hash, uint64) *big.Int) error {
	// Get the total difficulties of the proposed chain segment and the existing one.
	commonAncestorTD := getTDFunc(commonAncestor.Hash(), commonAncestor.Number.Uint64())
	proposedParentTD := getTDFunc(proposed.ParentHash, proposed.Number.Uint64()-1)
	proposedTD := new(big.Int).Add(proposed.Difficulty, proposedParentTD)
	localTD := getTDFunc(current.Hash(), current.Number.Uint64())

	// if proposed_subchain_td * CURVE_FUNCTION_DENOMINATOR < get_curve_function_numerator(current.Time - commonAncestor.Time) * local_subchain_td.
	//
	// The curve's x is the LOCAL head's age, not the proposed head's. This comment said
	// proposed.Time while the line below correctly used current.Time; ECIP-1100's own
	// pseudocode says current.Time. Corrected rather than left, because a reader
	// reimplementing from the comment would get the rule wrong in a way that only shows
	// up when the two heads' timestamps differ.
	proposedSubchainTD := new(big.Int).Sub(proposedTD, commonAncestorTD)
	localSubchainTD := new(big.Int).Sub(localTD, commonAncestorTD)

	xBig := big.NewInt(int64(current.Time - commonAncestor.Time))
	eq := ecbp1100PolynomialV(xBig)
	want := eq.Mul(eq, localSubchainTD)

	got := new(big.Int).Mul(proposedSubchainTD, ecbp1100PolynomialVCurveFunctionDenominator)

	if got.Cmp(want) < 0 {
		prettyRatio, _ := new(big.Float).Quo(
			new(big.Float).SetInt(got),
			new(big.Float).SetInt(want),
		).Float64()
		return fmt.Errorf(`%w: ECBP1100-MESS 🔒 status=rejected age=%v current.span=%v proposed.span=%v tdr/gravity=%0.6f common.bno=%d common.hash=%s current.bno=%d current.hash=%s proposed.bno=%d proposed.hash=%s`,
			errReorgFinality,
			common.PrettyAge(time.Unix(int64(commonAncestor.Time), 0)),
			common.PrettyDuration(time.Duration(current.Time-commonAncestor.Time)*time.Second),
			common.PrettyDuration(time.Duration(int32(xBig.Uint64()))*time.Second),
			prettyRatio,
			commonAncestor.Number.Uint64(), commonAncestor.Hash().Hex(),
			current.Number.Uint64(), current.Hash().Hex(),
			proposed.Number.Uint64(), proposed.Hash().Hex(),
		)
	}
	return nil
}

/*
ecbp1100PolynomialV is a cubic function that looks a lot like Option 3's sin function,
but adds the benefit that the calculation can be done with integers (instead of yucky floating points).
> https://github.com/ethereumclassic/ECIPs/issues/374#issuecomment-694156719

CURVE_FUNCTION_DENOMINATOR = 128

def get_curve_function_numerator(time_delta: int) -> int:
    xcap = 25132 # = floor(8000*pi)
    ampl = 15
    height = CURVE_FUNCTION_DENOMINATOR * (ampl * 2)
    if x > xcap:
        x = xcap
    # The sine approximator `y = 3*x**2 - 2*x**3` rescaled to the desired height and width
    return CURVE_FUNCTION_DENOMINATOR + (3 * x**2 - 2 * x**3 // xcap) * height // xcap ** 2


The if tdRatio < antiGravity check would then be

if proposed_subchain_td * CURVE_FUNCTION_DENOMINATOR < get_curve_function_numerator(current.Time - commonAncestor.Time) * local_subchain_td.
*/
// nolint:goimports
func ecbp1100PolynomialV(x *big.Int) *big.Int {
	// Make a copy; do not mutate argument value.

	// if x > xcap:
	//    x = xcap
	xA := new(big.Int).Set(x)
	if xA.Cmp(ecbp1100PolynomialVXCap) > 0 {
		xA.Set(ecbp1100PolynomialVXCap)
	}

	xB := new(big.Int).Set(x)
	if xB.Cmp(ecbp1100PolynomialVXCap) > 0 {
		xB.Set(ecbp1100PolynomialVXCap)
	}

	out := big.NewInt(0)

	// 3 * x**2
	xA.Exp(xA, big2, nil)
	xA.Mul(xA, big3)

	// 3 * x**2 // xcap
	xB.Exp(xB, big3, nil)
	xB.Mul(xB, big2)
	xB.Div(xB, ecbp1100PolynomialVXCap)

	// (3 * x**2 - 2 * x**3 // xcap)
	out.Sub(xA, xB)

	// // (3 * x**2 - 2 * x**3 // xcap) * height
	out.Mul(out, ecbp1100PolynomialVHeight)

	// xcap ** 2
	xcap2 := new(big.Int).Exp(ecbp1100PolynomialVXCap, big2, nil)

	// (3 * x**2 - 2 * x**3 // xcap) * height // xcap ** 2
	out.Div(out, xcap2)

	// CURVE_FUNCTION_DENOMINATOR + (3 * x**2 - 2 * x**3 // xcap) * height // xcap ** 2
	out.Add(out, ecbp1100PolynomialVCurveFunctionDenominator)
	return out
}

var big2 = big.NewInt(2)
var big3 = big.NewInt(3)

// ecbp1100PolynomialVCurveFunctionDenominator
// CURVE_FUNCTION_DENOMINATOR = 128
var ecbp1100PolynomialVCurveFunctionDenominator = big.NewInt(128)

// ecbp1100PolynomialVXCap
// xcap = 25132 # = floor(8000*pi)
var ecbp1100PolynomialVXCap = big.NewInt(25132)

// ecbp1100PolynomialVAmpl
// ampl = 15
var ecbp1100PolynomialVAmpl = big.NewInt(15)

// ecbp1100PolynomialVHeight
// height = CURVE_FUNCTION_DENOMINATOR * (ampl * 2)
var ecbp1100PolynomialVHeight = new(big.Int).Mul(new(big.Int).Mul(ecbp1100PolynomialVCurveFunctionDenominator, ecbp1100PolynomialVAmpl), big2)
