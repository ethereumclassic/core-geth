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
	"bytes"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
)

// makeSnapResponse builds an RLP payload with the structure [requestId, [items...]]
// mimicking a snap response. Each item is an empty string (RLP: 0x80).
func makeSnapResponse(requestID uint64, itemCount int) []byte {
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

// makeSnapResponseWithProof builds an RLP payload with the structure
// [requestId, [items...], [proofs...]] mimicking AccountRange/StorageRanges.
func makeSnapResponseWithProof(requestID uint64, itemCount, proofCount int) []byte {
	items := make([]interface{}, itemCount)
	for i := range items {
		items[i] = ""
	}
	proofs := make([]interface{}, proofCount)
	for i := range proofs {
		proofs[i] = ""
	}
	payload := []interface{}{requestID, items, proofs}
	data, err := rlp.EncodeToBytes(payload)
	if err != nil {
		panic(err)
	}
	return data
}

func TestValidateSnapMessageItems_UnvalidatedType(t *testing.T) {
	// Request messages should pass through without reading.
	data, err := validateSnapMessageItems(nil, 0, GetAccountRangeMsg)
	if err != nil {
		t.Fatalf("unexpected error for unvalidated message: %v", err)
	}
	if data != nil {
		t.Fatal("expected nil data for unvalidated message")
	}
}

func TestValidateSnapMessageItems_BelowLimit(t *testing.T) {
	payload := makeSnapResponse(42, 100)
	data, err := validateSnapMessageItems(bytes.NewReader(payload), uint32(len(payload)), ByteCodesMsg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(data, payload) {
		t.Fatal("returned data does not match original payload")
	}
}

func TestValidateSnapMessageItems_AtLimit(t *testing.T) {
	payload := makeSnapResponse(42, maxSnapResponseItems)
	data, err := validateSnapMessageItems(bytes.NewReader(payload), uint32(len(payload)), TrieNodesMsg)
	if err != nil {
		t.Fatalf("unexpected error at limit: %v", err)
	}
	if data == nil {
		t.Fatal("expected non-nil data at limit")
	}
}

func TestValidateSnapMessageItems_AboveLimit(t *testing.T) {
	payload := makeSnapResponse(42, maxSnapResponseItems+1)
	_, err := validateSnapMessageItems(bytes.NewReader(payload), uint32(len(payload)), ByteCodesMsg)
	if err == nil {
		t.Fatal("expected error for oversized message")
	}
	if !errors.Is(err, errTooManySnapItems) {
		t.Fatalf("expected errTooManySnapItems, got: %v", err)
	}
}

func TestValidateSnapMessageItems_AllTypes(t *testing.T) {
	// All snap response types should reject oversized payloads.
	types := []uint64{AccountRangeMsg, StorageRangesMsg, ByteCodesMsg, TrieNodesMsg}
	for _, code := range types {
		payload := makeSnapResponse(1, maxSnapResponseItems+1)
		_, err := validateSnapMessageItems(bytes.NewReader(payload), uint32(len(payload)), code)
		if err == nil {
			t.Errorf("expected error for message code %#02x", code)
		}
		if !errors.Is(err, errTooManySnapItems) {
			t.Errorf("message code %#02x: expected errTooManySnapItems, got: %v", code, err)
		}
	}
}

func TestValidateSnapMessageItems_WithProof_AboveLimit(t *testing.T) {
	// AccountRange has [id, [accounts...], [proofs...]] — we count the accounts list.
	payload := makeSnapResponseWithProof(42, maxSnapResponseItems+1, 5)
	_, err := validateSnapMessageItems(bytes.NewReader(payload), uint32(len(payload)), AccountRangeMsg)
	if err == nil {
		t.Fatal("expected error for oversized AccountRange")
	}
	if !errors.Is(err, errTooManySnapItems) {
		t.Fatalf("expected errTooManySnapItems, got: %v", err)
	}
}

func TestValidateSnapMessageItems_WithProof_BelowLimit(t *testing.T) {
	// Small accounts list with proof should pass.
	payload := makeSnapResponseWithProof(42, 100, 50)
	data, err := validateSnapMessageItems(bytes.NewReader(payload), uint32(len(payload)), AccountRangeMsg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(data, payload) {
		t.Fatal("returned data does not match original payload")
	}
}

func TestValidateSnapMessageItems_MalformedRLP(t *testing.T) {
	malformed := []byte{0xff, 0xff, 0xff}
	data, err := validateSnapMessageItems(bytes.NewReader(malformed), uint32(len(malformed)), ByteCodesMsg)
	if err != nil {
		t.Fatalf("malformed RLP should not return error: %v", err)
	}
	if !bytes.Equal(data, malformed) {
		t.Fatal("malformed RLP should return original data")
	}
}

func TestValidateSnapMessageItems_EmptyList(t *testing.T) {
	payload := makeSnapResponse(1, 0)
	data, err := validateSnapMessageItems(bytes.NewReader(payload), uint32(len(payload)), TrieNodesMsg)
	if err != nil {
		t.Fatalf("unexpected error for empty list: %v", err)
	}
	if data == nil {
		t.Fatal("expected non-nil data for validated message")
	}
}
