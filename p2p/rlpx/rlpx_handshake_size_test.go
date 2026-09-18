package rlpx

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/ecies"
)

// sizePrefix returns a handshake message made of nothing but the two-byte size prefix.
// A reader over it runs out of data before the packet body, so a readMsg that reaches
// the body read fails on the short read. That is what separates "refused the declared
// size" from "tried to read the packet".
func sizePrefix(size uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, size)
	return b
}

// TestHandshakeReadMsgSizeLimit covers the 2KB bound on handshake messages, ported from
// go-ethereum 27654d302 (#30029). That commit shipped no test, so this one is ours.
//
// readMsg runs on the inbound path before the peer has authenticated, and readBuffer.grow
// reserves whatever size the prefix declares. The bound is what stops an unauthenticated
// peer from choosing that allocation.
func TestHandshakeReadMsgSizeLimit(t *testing.T) {
	key := newkey()

	t.Run("refuses an oversized message before allocating", func(t *testing.T) {
		var h handshakeState
		_, err := h.readMsg(new(authMsgV4), key, bytes.NewReader(sizePrefix(2049)))
		if err == nil {
			t.Fatal("readMsg accepted a size prefix above the bound")
		}
		if err.Error() != "message too big" {
			t.Fatalf("readMsg error = %q, want %q", err, "message too big")
		}
		// readMsg grows the read buffer to 512 bytes before it reads the prefix. Any
		// capacity past that was reserved for the declared size, which is the
		// allocation the bound exists to refuse.
		if got := cap(h.rbuf.data); got > 512 {
			t.Fatalf("read buffer grew to %d bytes for a message that was refused", got)
		}
	})

	t.Run("admits the largest allowed message", func(t *testing.T) {
		var h handshakeState
		_, err := h.readMsg(new(authMsgV4), key, bytes.NewReader(sizePrefix(2048)))
		if err == nil {
			t.Fatal("readMsg returned no error for a packet that is only a prefix")
		}
		if err.Error() == "message too big" {
			t.Fatal("readMsg refused a size prefix of 2048, which is within the bound")
		}
	})

	t.Run("a real handshake packet still decrypts", func(t *testing.T) {
		initKey, respKey := newkey(), newkey()
		init := handshakeState{
			initiator: true,
			remote:    ecies.ImportECDSAPublic(&respKey.PublicKey),
		}
		authMsg, err := init.makeAuthMsg(initKey)
		if err != nil {
			t.Fatal(err)
		}
		packet, err := init.sealEIP8(authMsg)
		if err != nil {
			t.Fatal(err)
		}
		if size := binary.BigEndian.Uint16(packet[:2]); size > 2048 {
			t.Fatalf("sealed auth packet declares %d bytes, above the bound", size)
		}
		var recv handshakeState
		if _, err := recv.readMsg(new(authMsgV4), respKey, bytes.NewReader(packet)); err != nil {
			t.Fatalf("readMsg refused a valid handshake packet: %v", err)
		}
	})
}
