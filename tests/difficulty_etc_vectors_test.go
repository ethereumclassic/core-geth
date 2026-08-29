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

package tests

import (
	"encoding/json"
	"math/big"
	"os"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
)

// Difficulty-bomb conformance vectors for Ethereum Classic: the bomb's introduction at
// 200,000, its pause at 3,000,000 by ECIP-1010, and its removal at 5,900,000 by ECIP-1041,
// across every ETC rule set.
//
// The values come from three independent implementations read as source rather than run.
// Provenance is in etcvectors/README.md.
//
// THE FILE IS LABEL-KEYED. Each ETC_* section states expectations under that upgrade's rule
// set with currentBlockNumber as a free parameter, so a configuration must be built per
// label. Running every label against one params.ClassicChainConfig resolves the rules from
// the block number instead and scores 65 of 192 -- and this repository once reported the
// fixture as defective on exactly that mistake.
//
// The failure is one adjustment unit because EIP-100 changes the interval divisor from 10 to
// 9: a nine-second interval gives max(1-9//9, -99) = 0 where Frontier and EIP-2 give +1.

// etcBombConfig builds the rule set a label names.
//
// The two halves are not symmetric. The label's own adjustment rule is IN FORCE, so its
// transition is 0 rather than its mainnet height. Rules the label merely carries keep their
// REAL parameters -- ECIP-1010's pause window and ECIP-1041's removal height are parameters
// of those rules, not label selectors.
func etcBombConfig(eip2, eip100, ecip1010, ecip1041 bool) ctypes.ChainConfigurator {
	c := &coregeth.CoreGethChainConfig{Ethash: new(ctypes.EthashConfig)}
	z := big.NewInt(0)
	if eip2 {
		c.EIP2FBlock = z
	}
	if eip100 {
		c.EIP100FBlock = z
	}
	if ecip1010 {
		c.ECIP1010PauseBlock = big.NewInt(3000000)
		c.ECIP1010Length = big.NewInt(2000000)
	}
	if ecip1041 {
		c.DisposalBlock = big.NewInt(5900000)
	}
	return c
}

var etcBombLabels = map[string]ctypes.ChainConfigurator{
	"ETC_Frontier":             etcBombConfig(false, false, false, false),
	"ETC_Homestead":            etcBombConfig(true, false, false, false),
	"ETC_GasReprice":           etcBombConfig(true, false, false, false),
	"ETC_DieHard":              etcBombConfig(true, false, true, false),
	"ETC_Gotham":               etcBombConfig(true, false, true, false),
	"ETC_DefuseDifficultyBomb": etcBombConfig(true, false, true, true),
	"ETC_Atlantis":             etcBombConfig(true, true, true, true),
	"ETC_Agharta":              etcBombConfig(true, true, true, true),
	"ETC_Phoenix":              etcBombConfig(true, true, true, true),
	"ETC_Magneto":              etcBombConfig(true, true, true, true),
	"ETC_Mystique":             etcBombConfig(true, true, true, true),
	"ETC_Spiral":               etcBombConfig(true, true, true, true),
}

func TestETCDifficultyBombVectors(t *testing.T) {
	const path = "etcvectors/etc_bomb_pause_and_removal.json"

	blob, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	var outer map[string]map[string]json.RawMessage
	if err := json.Unmarshal(blob, &outer); err != nil {
		t.Fatalf("parsing %s: %v", path, err)
	}
	root, ok := outer["difficultyBombPauseAndRemoval"]
	if !ok {
		t.Fatalf("%s: expected root key difficultyBombPauseAndRemoval", path)
	}

	labels := make([]string, 0, len(root))
	for k := range root {
		if k != "_info" {
			labels = append(labels, k)
		}
	}
	sort.Strings(labels)

	var ran int
	for _, label := range labels {
		cfg, known := etcBombLabels[label]
		if !known {
			// Never skip silently: an unmapped label is a gap in this map, not a
			// reason to pass over cases the fixture provides.
			t.Errorf("%s: no configuration mapped for this label", label)
			continue
		}
		var cases map[string]DifficultyTest
		if err := json.Unmarshal(root[label], &cases); err != nil {
			t.Errorf("%s: %v", label, err)
			continue
		}
		if len(cases) == 0 {
			t.Errorf("%s: no cases -- a label that runs nothing would pass silently", label)
			continue
		}
		names := make([]string, 0, len(cases))
		for n := range cases {
			names = append(names, n)
		}
		sort.Strings(names)
		for _, name := range names {
			c := cases[name]
			t.Run(label+"/"+name, func(t *testing.T) {
				if err := c.Run(cfg); err != nil {
					t.Errorf("block %d: %v", c.CurrentBlockNumber, err)
				}
			})
			ran++
		}
	}
	if ran == 0 {
		t.Fatal("no cases ran -- the vectors loaded but nothing was asserted")
	}
	t.Logf("%d cases asserted across %d rule sets", ran, len(labels))
}
