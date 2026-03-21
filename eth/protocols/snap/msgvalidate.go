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
	"errors"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/rlp"
)

// maxSnapResponseItems defines the maximum number of items allowed in snap
// protocol response messages. Snap serve-side limits are maxCodeLookups=1024
// and maxTrieNodeLookups=1024, so 2048 provides headroom.
//
// This is a mitigation for CVE-2026-26313.
const maxSnapResponseItems = 2048

// errTooManySnapItems is returned when a peer sends a snap message containing
// more items than the protocol allows.
var errTooManySnapItems = errors.New("too many items in snap message")

// snapResponseItemLimits maps snap response message codes to their maximum
// allowed item counts. Messages not in this map are not validated.
var snapResponseItemLimits = map[uint64]int{
	AccountRangeMsg:  maxSnapResponseItems,
	StorageRangesMsg: maxSnapResponseItems,
	ByteCodesMsg:     maxSnapResponseItems,
	TrieNodesMsg:     maxSnapResponseItems,
}

// validateSnapMessageItems reads the raw payload of a snap message and counts
// the number of items in the primary response list without performing full RLP
// decoding. All snap response messages have the structure [RequestId, [items...], ...]
// where we count items in the first list after the RequestId.
//
// The function reads the entire payload into memory and returns it so the
// caller can replace the consumed io.Reader. The payload has already been
// read into memory by the transport layer, so this does not increase peak
// memory usage.
func validateSnapMessageItems(payload io.Reader, size uint32, code uint64) ([]byte, error) {
	limit, ok := snapResponseItemLimits[code]
	if !ok {
		// No validation needed for this message type.
		return nil, nil
	}

	// Read the raw payload.
	data := make([]byte, size)
	if _, err := io.ReadFull(payload, data); err != nil {
		return nil, fmt.Errorf("failed to read snap message payload: %w", err)
	}

	// Parse the outer RLP list.
	_, content, _, err := rlp.Split(data)
	if err != nil {
		return data, nil // Let the normal decoder report the error.
	}

	// All snap responses have [RequestId, [primary_list...], ...].
	// Skip the RequestId (first element).
	_, _, rest, err := rlp.Split(content)
	if err != nil {
		return data, nil // Let the normal decoder report the error.
	}

	// Parse the primary response list (accounts, slots, codes, or nodes).
	_, responseContent, _, err := rlp.Split(rest)
	if err != nil {
		return data, nil // Let the normal decoder report the error.
	}

	exceeded, err := countValuesExceedsLimit(responseContent, limit)
	if err != nil {
		return data, nil // Let the normal decoder report the error.
	}
	if exceeded {
		return nil, fmt.Errorf("%w: message %#02x exceeds item limit %d",
			errTooManySnapItems, code, limit)
	}
	return data, nil
}

// countValuesExceedsLimit counts the number of RLP values in b and returns
// true if the count exceeds the given limit. Exits early as soon as the
// limit is exceeded, avoiding excessive iterations for attack payloads.
func countValuesExceedsLimit(b []byte, limit int) (bool, error) {
	count := 0
	for len(b) > 0 {
		_, _, rest, err := rlp.Split(b)
		if err != nil {
			return false, err
		}
		b = rest
		count++
		if count > limit {
			return true, nil
		}
	}
	return false, nil
}
