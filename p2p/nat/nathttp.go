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
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/log"
)

// httpProbeURLs are the canonical HTTPS probes, aligned with Besu and Nethermind.
var httpProbeURLs = []string{
	"https://icanhazip.com",
	"https://checkip.amazonaws.com",
	"https://api4.ipify.org",
	"https://4.ident.me",
}

const httpProbeTimeout = 3 * time.Second

// natHTTP implements Interface by probing a set of HTTPS endpoints for the external IP.
// AddMapping and DeleteMapping are no-ops; only ExternalIP is meaningful.
type natHTTP struct{}

func (natHTTP) String() string                                                     { return "HTTP" }
func (natHTTP) AddMapping(string, int, int, string, time.Duration) (uint16, error) { return 0, nil }
func (natHTTP) DeleteMapping(string, int, int) error                               { return nil }

func (natHTTP) ExternalIP() (net.IP, error) {
	client := &http.Client{Timeout: httpProbeTimeout}
	for _, url := range httpProbeURLs {
		ip, err := httpProbe(client, url)
		if err != nil {
			log.Debug("HTTP NAT probe failed", "url", url, "err", err)
			continue
		}
		return ip, nil
	}
	return nil, errors.New("all HTTP NAT probes failed")
}

func httpProbe(client *http.Client, url string) (net.IP, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP probe %s: status %d", url, resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 64))
	if err != nil {
		return nil, err
	}
	ipStr := strings.TrimSpace(string(body))
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP from %s: %q", url, ipStr)
	}
	return ip, nil
}

// discoverHTTP tries the HTTP probe and returns ExtIP on success, nil otherwise.
func discoverHTTP() Interface {
	ip, err := natHTTP{}.ExternalIP()
	if err != nil {
		return nil
	}
	return ExtIP(ip)
}
