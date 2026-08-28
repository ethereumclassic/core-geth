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

package tracker

import (
	"errors"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/p2p"
)

// TestCleanAfterStop verifies that the clean method does not crash when called
// after Stop. This can happen because clean is scheduled via time.AfterFunc and
// may fire after Stop sets t.expire to nil.
func TestCleanAfterStop(t *testing.T) {
	cap := p2p.Cap{Name: "test", Version: 1}
	timeout := 50 * time.Millisecond
	tr := New(cap, "peer1", timeout)

	// Track a request to start the expiration timer.
	tr.Track(Request{ID: 1, ReqCode: 0x01, RespCode: 0x02, Size: 1})

	// Stop the tracker, then wait for the timer to fire.
	tr.Stop()
	time.Sleep(timeout + 50*time.Millisecond)

	// Also verify that calling clean directly after stop doesn't panic.
	tr.clean()
}

// This checks that metrics gauges for pending requests are be decremented when a
// Tracker is stopped.
func TestMetricsOnStop(t *testing.T) {
	metrics.Enabled = true

	cap := p2p.Cap{Name: "test", Version: 1}
	tr := New(cap, "peer1", time.Minute)

	// Track some requests with different ReqCodes.
	var id uint64
	for i := 0; i < 3; i++ {
		tr.Track(Request{ID: id, ReqCode: 0x01, RespCode: 0x02, Size: 1})
		id++
	}
	for i := 0; i < 5; i++ {
		tr.Track(Request{ID: id, ReqCode: 0x03, RespCode: 0x04, Size: 1})
		id++
	}

	gauge1 := tr.trackedGauge(0x01)
	gauge2 := tr.trackedGauge(0x03)

	if gauge1.Snapshot().Value() != 3 {
		t.Fatalf("gauge1 value mismatch: got %d, want 3", gauge1.Snapshot().Value())
	}
	if gauge2.Snapshot().Value() != 5 {
		t.Fatalf("gauge2 value mismatch: got %d, want 5", gauge2.Snapshot().Value())
	}

	tr.Stop()

	if gauge1.Snapshot().Value() != 0 {
		t.Fatalf("gauge1 value after stop: got %d, want 0", gauge1.Snapshot().Value())
	}
	if gauge2.Snapshot().Value() != 0 {
		t.Fatalf("gauge2 value after stop: got %d, want 0", gauge2.Snapshot().Value())
	}
}

// TestFulfilRejections covers every way Fulfil can refuse a response. Since the
// pre-decode item ceilings were removed from the snap protocol, this is the only
// bound standing between an oversized response and the materialization of its
// contents, so each refusal path is pinned here rather than left to the callers.
func TestFulfilRejections(t *testing.T) {
	cap := p2p.Cap{Name: "test", Version: 1}
	newTracker := func() *Tracker { return New(cap, "peer1", time.Minute) }

	t.Run("no matching request", func(t *testing.T) {
		tr := newTracker()
		defer tr.Stop()
		if err := tr.Fulfil(Response{ID: 7, MsgCode: 0x02, Size: 1}); !errors.Is(err, ErrNoMatchingRequest) {
			t.Fatalf("have %v, want %v", err, ErrNoMatchingRequest)
		}
	})

	t.Run("wrong response code", func(t *testing.T) {
		tr := newTracker()
		defer tr.Stop()
		if err := tr.Track(Request{ID: 1, ReqCode: 0x01, RespCode: 0x02, Size: 4}); err != nil {
			t.Fatalf("track: %v", err)
		}
		if err := tr.Fulfil(Response{ID: 1, MsgCode: 0x03, Size: 4}); !errors.Is(err, ErrCodeMismatch) {
			t.Fatalf("have %v, want %v", err, ErrCodeMismatch)
		}
	})

	t.Run("response larger than requested", func(t *testing.T) {
		tr := newTracker()
		defer tr.Stop()
		if err := tr.Track(Request{ID: 1, ReqCode: 0x01, RespCode: 0x02, Size: 4}); err != nil {
			t.Fatalf("track: %v", err)
		}
		if err := tr.Fulfil(Response{ID: 1, MsgCode: 0x02, Size: 5}); !errors.Is(err, ErrTooManyItems) {
			t.Fatalf("have %v, want %v", err, ErrTooManyItems)
		}
	})

	t.Run("response at the requested size is accepted", func(t *testing.T) {
		tr := newTracker()
		defer tr.Stop()
		if err := tr.Track(Request{ID: 1, ReqCode: 0x01, RespCode: 0x02, Size: 4}); err != nil {
			t.Fatalf("track: %v", err)
		}
		if err := tr.Fulfil(Response{ID: 1, MsgCode: 0x02, Size: 4}); err != nil {
			t.Fatalf("response at the requested size rejected: %v", err)
		}
	})

	t.Run("request id collision", func(t *testing.T) {
		tr := newTracker()
		defer tr.Stop()
		if err := tr.Track(Request{ID: 1, ReqCode: 0x01, RespCode: 0x02, Size: 4}); err != nil {
			t.Fatalf("track: %v", err)
		}
		if err := tr.Track(Request{ID: 1, ReqCode: 0x01, RespCode: 0x02, Size: 4}); !errors.Is(err, ErrCollision) {
			t.Fatalf("have %v, want %v", err, ErrCollision)
		}
	})

	t.Run("tracking limit reached", func(t *testing.T) {
		tr := newTracker()
		defer tr.Stop()
		for i := 0; i < maxTrackedPackets; i++ {
			if err := tr.Track(Request{ID: uint64(i), ReqCode: 0x01, RespCode: 0x02, Size: 1}); err != nil {
				t.Fatalf("track %d: %v", i, err)
			}
		}
		err := tr.Track(Request{ID: maxTrackedPackets, ReqCode: 0x01, RespCode: 0x02, Size: 1})
		if !errors.Is(err, ErrLimitReached) {
			t.Fatalf("have %v, want %v", err, ErrLimitReached)
		}
	})

	t.Run("tracking after stop", func(t *testing.T) {
		tr := newTracker()
		tr.Stop()
		if err := tr.Track(Request{ID: 1, ReqCode: 0x01, RespCode: 0x02, Size: 1}); !errors.Is(err, ErrStopped) {
			t.Fatalf("have %v, want %v", err, ErrStopped)
		}
	})
}
