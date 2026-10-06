// Copyright 2026 The core-geth Authors
// This file is part of the core-geth library.
//
// The core-geth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The core-geth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the core-geth library. If not, see <http://www.gnu.org/licenses/>.

package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
)

// messBlock renders an ECBP-1100 (MESS) window end for a log line. The block MESS settings move
// out of reach prints as never: as 18446744073709551614 it reads as a configured height, and it is
// not one.
func messBlock(n *uint64) interface{} {
	switch {
	case n == nil:
		return "none"
	case *n >= core.ECBP1100Unreachable:
		return "never"
	default:
		return *n
	}
}

// readStoredMESSWindow captures the chain configuration the database already holds.
//
// It must be called BEFORE the chain is opened. SetupGenesisBlock writes the current
// configuration over the stored one as part of opening, so a read afterwards returns what this
// client just wrote, and every comparison against it is trivially equal. A unit test cannot see a
// read in the wrong place, because it compares configurations handed to it; an upgrade on a real
// database does.
func readStoredMESSWindow(db ethdb.Reader) ctypes.ChainConfigurator {
	genesis := rawdb.ReadCanonicalHash(db, 0)
	if genesis == (common.Hash{}) {
		return nil // a database with no genesis has nothing to compare against
	}
	return rawdb.ReadChainConfig(db, genesis)
}

// logMESSWindowChange reports when the ECBP-1100 (MESS) window now in force differs, at the
// current head, from the one readStoredMESSWindow captured before the chain opened.
//
// Nothing else says so. The chain configuration is logged twice while the chain opens, the
// stored one and the one in force, and an upgrade that moves this field leaves those two lines
// disagreeing a few lines apart with nothing between them. A release that changes the bundled
// default therefore changes which of two competing chains a node prefers during a deep
// reorganization, and the only evidence is a disagreement nobody reads.
//
// It reports a change in EFFECT at the head rather than a change in the numbers. A stored
// window and a current window that disagree only about blocks the chain has not reached are
// not something an operator can act on today, and a notice that fires for everyone is one
// nobody reads either.
//
// chose reports whether the operator set a MESS block on the command line or in a config file,
// and silences the notice. It is for a window that moved underneath them, not one they moved:
// an operator who typed --mess does not need to be told MESS changed, and advising them to
// pass --mess=false to "keep the previous behavior" is advice against what they asked for.
//
// It fires once. Opening the database writes the new configuration, so the next start finds
// stored and current agreeing. That is deliberate: a line repeated at every start is one nobody
// reads. It is also why the release notes must carry the same information.
//
// This is not consensus. MESS selects between competing chains and never decides whether a
// block is valid, so a node whose setting changed is not on a different chain: it is one that
// would resist a deep reorganization differently from its fleet, which is why a silent change
// is worth a line at all.
func logMESSWindowChange(stored ctypes.ChainConfigurator, bc *core.BlockChain, chose bool) {
	if bc == nil {
		return
	}
	current := bc.Config()
	head := bc.CurrentBlock().Number
	if !messNoticeDue(stored, current, head, chose) {
		return
	}
	was := stored.IsEnabled(stored.GetECBP1100Transition, head)

	// The flag that keeps the previous setting, named so the operator does not have to work
	// out which end of the window moved.
	keep := "--mess"
	if !was {
		keep = "--mess=false"
	}
	log.Warn("ECBP1100 (MESS) is configured differently from the client that last used this database",
		"was", onOff(was), "now", onOff(!was), "head", head,
		"stored.activation", messBlock(stored.GetECBP1100Transition()),
		"stored.deactivation", messBlock(stored.GetECBP1100DeactivateTransition()),
		"current.activation", messBlock(current.GetECBP1100Transition()),
		"current.deactivation", messBlock(current.GetECBP1100DeactivateTransition()),
		"keep_previous", keep,
		"note", "chain selection only, not block validity; give every node in a fleet the same setting")
}

func onOff(b bool) string {
	if b {
		return "on"
	}
	return "off"
}

// messNoticeDue is logMESSWindowChange's decision, separated so it can be tested against two
// chain configurations without opening a database or starting a node. The notice is due when
// the operator did not choose the window themselves and it applies differently at the head.
func messNoticeDue(stored, current ctypes.ChainConfigurator, head *big.Int, chose bool) bool {
	return !chose && messWindowChanged(stored, current, head)
}

// messWindowChanged reports whether the two windows apply differently at head.
func messWindowChanged(stored, current ctypes.ChainConfigurator, head *big.Int) bool {
	if stored == nil || current == nil || head == nil {
		return false
	}
	return stored.IsEnabled(stored.GetECBP1100Transition, head) !=
		current.IsEnabled(current.GetECBP1100Transition, head)
}
