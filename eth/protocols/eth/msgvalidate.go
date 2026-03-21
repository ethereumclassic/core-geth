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

package eth

import (
	"errors"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/rlp"
)

// maxResponseItems defines the maximum number of items allowed in response
// messages. This prevents a malicious peer from sending a response claiming
// to contain millions of items, which would cause the node to allocate
// excessive memory during RLP decoding before the response is validated.
//
// This is a mitigation for CVE-2026-26313.
const maxResponseItems = 2048

// maxTransactionItems defines the maximum number of transactions allowed in
// a broadcast TransactionsMsg. This is set to maxKnownTxs since a peer should
// never send more unique transactions than can be tracked.
const maxTransactionItems = 32768

// errTooManyItems is returned when a peer sends a message containing more
// items than the protocol allows.
var errTooManyItems = errors.New("too many items in message")

// responseItemLimits maps message codes to their maximum allowed item counts.
// Messages not in this map are not validated.
var responseItemLimits = map[uint64]int{
	BlockHeadersMsg:       maxResponseItems,
	BlockBodiesMsg:        maxResponseItems,
	ReceiptsMsg:           maxResponseItems,
	PooledTransactionsMsg: maxResponseItems,
	TransactionsMsg:       maxTransactionItems,
	NewBlockHashesMsg:     maxResponseItems,
}

// flatMessages is the set of message codes that are flat lists (no RequestId
// prefix). All other validated messages have the structure [RequestId, [items...]].
var flatMessages = map[uint64]bool{
	TransactionsMsg:   true,
	NewBlockHashesMsg: true,
}

// validateMessageItems reads the raw payload of a message and counts the
// number of items in the response list without performing full RLP decoding.
// For request-response messages (with a RequestId prefix), the count is of
// items in the second field (the response body). For broadcast messages like
// TransactionsMsg, the count is of the top-level list items.
//
// The function reads the entire payload into memory and returns it so the
// caller can replace the consumed io.Reader. The payload has already been
// read into memory by the transport layer, so this does not increase peak
// memory usage.
//
// Note: This validates only the top-level item count in each message. Inner
// lists (e.g., transactions within a BlockBody) are not counted. The 10 MiB
// maxMessageSize naturally limits inner list damage since each inner item has
// a minimum encoding size. A future full backport of upstream's rlp.RawList
// delayed decoding would address inner lists as well.
func validateMessageItems(payload io.Reader, size uint32, code uint64) ([]byte, error) {
	limit, ok := responseItemLimits[code]
	if !ok {
		// No validation needed for this message type.
		return nil, nil
	}

	// Read the raw payload.
	data := make([]byte, size)
	if _, err := io.ReadFull(payload, data); err != nil {
		return nil, fmt.Errorf("failed to read message payload: %w", err)
	}

	// Parse the outer RLP list.
	_, content, _, err := rlp.Split(data)
	if err != nil {
		return data, nil // Let the normal decoder report the error.
	}

	var itemData []byte
	if flatMessages[code] {
		// Flat list messages have no RequestId prefix.
		itemData = content
	} else {
		// Response messages have the structure [RequestId, [items...]].
		// Skip the RequestId (first element) and parse the response list.
		_, _, rest, err := rlp.Split(content)
		if err != nil {
			return data, nil // Let the normal decoder report the error.
		}
		// The remaining content should be the response list.
		_, responseContent, _, err := rlp.Split(rest)
		if err != nil {
			return data, nil // Let the normal decoder report the error.
		}
		itemData = responseContent
	}

	exceeded, err := countValuesExceedsLimit(itemData, limit)
	if err != nil {
		return data, nil // Let the normal decoder report the error.
	}
	if exceeded {
		return nil, fmt.Errorf("%w: message %#02x exceeds item limit %d",
			errTooManyItems, code, limit)
	}
	return data, nil
}

// countValuesExceedsLimit counts the number of RLP values in b and returns
// true if the count exceeds the given limit. Unlike rlp.CountValues, this
// exits early as soon as the limit is exceeded, avoiding 10M iterations
// for attack payloads with millions of tiny items.
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
