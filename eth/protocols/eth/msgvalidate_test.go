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
	"bytes"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
)

// makeRequestIDWrapped builds an RLP payload with the structure [requestId, [items...]].
// Each item is an empty string (RLP: 0x80), which is the minimal valid RLP element.
func makeRequestIDWrapped(requestID uint64, itemCount int) []byte {
	items := make([]interface{}, itemCount)
	for i := range items {
		items[i] = ""
	}
	payload := []interface{}{requestID, items}
	data, err := rlp.EncodeToBytes(payload)
	if err != nil {
		panic(err)
	}
	return data
}

// makeFlatList builds an RLP payload with the structure [items...].
// Each item is an empty string (RLP: 0x80).
func makeFlatList(itemCount int) []byte {
	items := make([]interface{}, itemCount)
	for i := range items {
		items[i] = ""
	}
	data, err := rlp.EncodeToBytes(items)
	if err != nil {
		panic(err)
	}
	return data
}

func TestValidateMessageItems_UnvalidatedType(t *testing.T) {
	// Messages not in responseItemLimits should pass through without reading.
	data, err := validateMessageItems(nil, 0, StatusMsg)
	if err != nil {
		t.Fatalf("unexpected error for unvalidated message: %v", err)
	}
	if data != nil {
		t.Fatal("expected nil data for unvalidated message")
	}
}

func TestValidateMessageItems_RequestIDWrapped_BelowLimit(t *testing.T) {
	payload := makeRequestIDWrapped(42, 100)
	data, err := validateMessageItems(bytes.NewReader(payload), uint32(len(payload)), BlockHeadersMsg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(data, payload) {
		t.Fatal("returned data does not match original payload")
	}
}

func TestValidateMessageItems_RequestIDWrapped_AtLimit(t *testing.T) {
	payload := makeRequestIDWrapped(42, maxResponseItems)
	data, err := validateMessageItems(bytes.NewReader(payload), uint32(len(payload)), BlockHeadersMsg)
	if err != nil {
		t.Fatalf("unexpected error at limit: %v", err)
	}
	if data == nil {
		t.Fatal("expected non-nil data at limit")
	}
}

func TestValidateMessageItems_RequestIDWrapped_AboveLimit(t *testing.T) {
	payload := makeRequestIDWrapped(42, maxResponseItems+1)
	_, err := validateMessageItems(bytes.NewReader(payload), uint32(len(payload)), BlockHeadersMsg)
	if err == nil {
		t.Fatal("expected error for oversized message")
	}
	if !errors.Is(err, errTooManyItems) {
		t.Fatalf("expected errTooManyItems, got: %v", err)
	}
}

func TestValidateMessageItems_RequestIDWrapped_AllTypes(t *testing.T) {
	// All request-ID wrapped message types should reject oversized payloads.
	types := []uint64{BlockHeadersMsg, BlockBodiesMsg, ReceiptsMsg, PooledTransactionsMsg}
	for _, code := range types {
		payload := makeRequestIDWrapped(1, maxResponseItems+1)
		_, err := validateMessageItems(bytes.NewReader(payload), uint32(len(payload)), code)
		if err == nil {
			t.Errorf("expected error for message code %#02x", code)
		}
		if !errors.Is(err, errTooManyItems) {
			t.Errorf("message code %#02x: expected errTooManyItems, got: %v", code, err)
		}
	}
}

func TestValidateMessageItems_FlatList_BelowLimit(t *testing.T) {
	payload := makeFlatList(100)
	data, err := validateMessageItems(bytes.NewReader(payload), uint32(len(payload)), TransactionsMsg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(data, payload) {
		t.Fatal("returned data does not match original payload")
	}
}

func TestValidateMessageItems_FlatList_AboveLimit(t *testing.T) {
	payload := makeFlatList(maxTransactionItems + 1)
	_, err := validateMessageItems(bytes.NewReader(payload), uint32(len(payload)), TransactionsMsg)
	if err == nil {
		t.Fatal("expected error for oversized transaction broadcast")
	}
	if !errors.Is(err, errTooManyItems) {
		t.Fatalf("expected errTooManyItems, got: %v", err)
	}
}

func TestValidateMessageItems_NewBlockHashes_AboveLimit(t *testing.T) {
	payload := makeFlatList(maxResponseItems + 1)
	_, err := validateMessageItems(bytes.NewReader(payload), uint32(len(payload)), NewBlockHashesMsg)
	if err == nil {
		t.Fatal("expected error for oversized NewBlockHashes")
	}
	if !errors.Is(err, errTooManyItems) {
		t.Fatalf("expected errTooManyItems, got: %v", err)
	}
}

func TestValidateMessageItems_MalformedRLP(t *testing.T) {
	// Malformed RLP should be passed through to the normal decoder.
	malformed := []byte{0xff, 0xff, 0xff}
	data, err := validateMessageItems(bytes.NewReader(malformed), uint32(len(malformed)), BlockHeadersMsg)
	if err != nil {
		t.Fatalf("malformed RLP should not return error: %v", err)
	}
	if !bytes.Equal(data, malformed) {
		t.Fatal("malformed RLP should return original data")
	}
}

func TestValidateMessageItems_EmptyList(t *testing.T) {
	// An empty response list should pass validation.
	payload := makeRequestIDWrapped(1, 0)
	data, err := validateMessageItems(bytes.NewReader(payload), uint32(len(payload)), BlockHeadersMsg)
	if err != nil {
		t.Fatalf("unexpected error for empty list: %v", err)
	}
	if data == nil {
		t.Fatal("expected non-nil data for validated message")
	}
}
