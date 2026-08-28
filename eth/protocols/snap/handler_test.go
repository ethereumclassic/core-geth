// Copyright 2026 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package snap

import (
	"encoding/binary"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/holiman/uint256"
)

// slimEOABody returns the genuine slim-format encoding of an externally owned
// account (empty storage root, empty code hash), as served in AccountRange
// responses.
func slimEOABody(nonce uint64, balance *uint256.Int) rlp.RawValue {
	return types.SlimAccountRLP(types.StateAccount{
		Nonce:    nonce,
		Balance:  balance,
		Root:     types.EmptyRootHash,
		CodeHash: types.EmptyCodeHash[:],
	})
}

// accountHashAt returns a synthetic account hash that sorts in ascending order
// with the index, matching the monotonicity requirement on honest responses.
func accountHashAt(i int) common.Hash {
	var h common.Hash
	binary.BigEndian.PutUint64(h[:8], uint64(i)+1)
	h[31] = 0x7f // make it look less degenerate than a bare counter
	return h
}

// realisticAccounts generates n consecutive accounts carrying real slim EOA
// encodings of varying sizes (5 to ~15 bytes), in ascending hash order.
func realisticAccounts(n int) []*AccountData {
	accounts := make([]*AccountData, n)
	for i := 0; i < n; i++ {
		balance := new(uint256.Int).Mul(uint256.NewInt(uint64(i%5)), uint256.NewInt(1e18))
		accounts[i] = &AccountData{
			Hash: accountHashAt(i),
			Body: slimEOABody(uint64(i%7), balance),
		}
	}
	return accounts
}

// densestHonestResponse replays ServiceGetAccountRangeQuery's accounting
// (size += common.HashLength + len(body); append; stop once size > budget)
// against a maxRequestSize budget using the smallest real EOA body (a used,
// zero-balance account encodes to 5 bytes in slim format). This is the
// maximum item count an honest byte-filled response can carry.
func densestHonestResponse() []*AccountData {
	var (
		accounts []*AccountData
		size     uint64
	)
	for i := 0; ; i++ {
		body := slimEOABody(1, uint256.NewInt(0))
		size += uint64(common.HashLength + len(body))
		accounts = append(accounts, &AccountData{Hash: accountHashAt(i), Body: body})
		if size > maxRequestSize {
			break
		}
	}
	return accounts
}

// TestAccountRangeHonestMaximum pins that an AccountRange response filled to
// the byte budget the server is allowed to spend is decoded without error.
//
// This client carried an item-count ceiling of 2048 applied to every snap
// response type, which rejected such responses and disconnected the peer:
// AccountRange is bounded by bytes rather than by item count, so an honest
// reply holds as many entries as fit. The ceiling is gone and the protocol
// now decodes the list lazily, but the case is worth pinning, because the
// failure it produced was a torn-down connection to a peer behaving correctly.
func TestAccountRangeHonestMaximum(t *testing.T) {
	accounts := densestHonestResponse()
	if len(accounts) < 10000 {
		t.Fatalf("test lost its meaning: honest maximum fell to %d items", len(accounts))
	}
	blob, err := rlp.EncodeToBytes(&AccountRangePacket{ID: 1, Accounts: accounts})
	if err != nil {
		t.Fatalf("failed to encode response: %v", err)
	}
	var res AccountRangePacket
	if err := rlp.DecodeBytes(blob, &res); err != nil {
		t.Fatalf("honest byte-filled response (%d items) rejected: %v", len(accounts), err)
	}
	if got := len(res.Accounts); got != len(accounts) {
		t.Fatalf("decoded %d accounts, want %d", got, len(accounts))
	}
}
