// Package history maintains a rolling log of port change events
// so that users can review recent alerts after the fact.
package history

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Entry records a single alert event with a timestamp.
type Entry struct {
	OccurredAt time.Time
	Event      alert.Event
}

// Ring is a fixed-capacity circular buffer of history entries.
type Ring struct {
	mu       sync.RWMutex
	entries  []Entry
	cap      int
	head     int
	size     int
}

// New creates a Ring that retains at most capacity entries.
// If capacity is less than 1 it defaults to 100.
func New(capacity int) *Ring {
	if capacity < 1 {
		capacity = 100
	}
	return &Ring{
		entries: make([]Entry, capacity),
		cap:     capacity,
	}
}

// Record appends events to the ring, evicting the oldest entry when full.
func (r *Ring) Record(events []alert.Event) {
	if len(events) == 0 {
		return
	}
	now := time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, e := range events {
		r.entries[r.head] = Entry{OccurredAt: now, Event: e}
		r.head = (r.head + 1) % r.cap
		if r.size < r.cap {
			r.size++
		}
	}
}

// Latest returns up to n most-recent entries in chronological order.
// If n <= 0 all stored entries are returned.
func (r *Ring) Latest(n int) []Entry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.size == 0 {
		return nil
	}
	count := r.size
	if n > 0 && n < count {
		count = n
	}
	out := make([]Entry, count)
	// oldest index inside the ring
	start := (r.head - r.size + r.cap) % r.cap
	offset := r.size - count
	for i := 0; i < count; i++ {
		out[i] = r.entries[(start+offset+i)%r.cap]
	}
	return out
}

// Len returns the number of entries currently stored.
func (r *Ring) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.size
}
