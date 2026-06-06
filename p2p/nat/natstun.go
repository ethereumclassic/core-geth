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

package nat

import (
	"encoding/binary"
	"errors"
	"net"
	"time"

	"github.com/ethereum/go-ethereum/log"
)

// stunServers are the same canonical STUN servers used by Besu and Nethermind.
var stunServers = []string{
	"stun.l.google.com:19302",
	"stun1.l.google.com:19302",
	"stun.cloudflare.com:3478",
}

const stunProbeTimeout = 3 * time.Second

// natSTUN implements Interface via raw UDP STUN Binding Request (RFC 5389).
// No external dependencies. AddMapping and DeleteMapping are no-ops.
type natSTUN struct{}

func (natSTUN) String() string                                                     { return "STUN" }
func (natSTUN) AddMapping(string, int, int, string, time.Duration) (uint16, error) { return 0, nil }
func (natSTUN) DeleteMapping(string, int, int) error                               { return nil }

func (natSTUN) ExternalIP() (net.IP, error) {
	for _, server := range stunServers {
		ip, err := stunProbe(server)
		if err != nil {
			log.Debug("STUN probe failed", "server", server, "err", err)
			continue
		}
		return ip, nil
	}
	return nil, errors.New("all STUN probes failed")
}

func stunProbe(server string) (net.IP, error) {
	conn, err := net.DialTimeout("udp", server, stunProbeTimeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(stunProbeTimeout))

	// RFC 5389 §6: Binding Request — 20-byte fixed header, no body attributes
	var req [20]byte
	req[0], req[1] = 0x00, 0x01                             // Message Type: Binding Request
	req[2], req[3] = 0x00, 0x00                             // Message Length: 0 (no body)
	req[4], req[5], req[6], req[7] = 0x21, 0x12, 0xA4, 0x42 // Magic Cookie (RFC 5389)
	// bytes 8–19: Transaction ID — all zeros (we accept any response)

	if _, err := conn.Write(req[:]); err != nil {
		return nil, err
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return parseXorMappedAddress(buf[:n])
}

// parseXorMappedAddress extracts the IPv4 address from a STUN Binding Success Response.
// Attribute type 0x0020 (XOR-MAPPED-ADDRESS) layout (RFC 5389 §15.2):
//
//	1 byte  reserved
//	1 byte  family (0x01 = IPv4)
//	2 bytes XOR-mapped port (unused — IP only)
//	4 bytes XOR-mapped address (IPv4 XOR'd with magic cookie 0x2112A442)
func parseXorMappedAddress(buf []byte) (net.IP, error) {
	if len(buf) < 20 {
		return nil, errors.New("STUN response too short")
	}
	if binary.BigEndian.Uint16(buf[0:2]) != 0x0101 {
		return nil, errors.New("STUN response is not a Binding Success Response")
	}
	msgLen := int(binary.BigEndian.Uint16(buf[2:4]))
	pos := 20
	bodyEnd := 20 + msgLen

	for pos+4 <= bodyEnd && pos+4 <= len(buf) {
		attrType := binary.BigEndian.Uint16(buf[pos : pos+2])
		attrLen := int(binary.BigEndian.Uint16(buf[pos+2 : pos+4]))
		valStart := pos + 4

		if attrType == 0x0020 && attrLen >= 8 && valStart+8 <= len(buf) {
			if buf[valStart+1] == 0x01 { // IPv4 family
				xorAddr := binary.BigEndian.Uint32(buf[valStart+4 : valStart+8])
				addr := xorAddr ^ 0x2112A442
				ip := make(net.IP, 4)
				binary.BigEndian.PutUint32(ip, addr)
				return ip, nil
			}
		}
		pos = valStart + ((attrLen + 3) &^ 3) // advance past padded attribute
	}
	return nil, errors.New("STUN response contains no XOR-MAPPED-ADDRESS for IPv4")
}

// discoverSTUN tries the STUN probe and returns ExtIP on success, nil otherwise.
func discoverSTUN() Interface {
	ip, err := natSTUN{}.ExternalIP()
	if err != nil {
		return nil
	}
	return ExtIP(ip)
}
