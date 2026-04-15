package history_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(kind alert.EventKind, port uint16) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Proto: "tcp", Port: port, PID: 0},
	}
}

func TestNew_DefaultCapacity(t *testing.T) {
	r := history.New(0)
	if r == nil {
		t.Fatal("expected non-nil Ring")
	}
}

func TestRecord_AndLen(t *testing.T) {
	r := history.New(10)
	r.Record([]alert.Event{
		makeEvent(alert.Opened, 8080),
		makeEvent(alert.Closed, 443),
	})
	if got := r.Len(); got != 2 {
		t.Fatalf("expected Len 2, got %d", got)
	}
}

func TestLatest_Order(t *testing.T) {
	r := history.New(10)
	r.Record([]alert.Event{makeEvent(alert.Opened, 1)})
	r.Record([]alert.Event{makeEvent(alert.Opened, 2)})
	r.Record([]alert.Event{makeEvent(alert.Opened, 3)})

	entries := r.Latest(0)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Event.Port.Port != 1 || entries[2].Event.Port.Port != 3 {
		t.Errorf("entries not in chronological order: %+v", entries)
	}
}

func TestLatest_LimitN(t *testing.T) {
	r := history.New(10)
	for i := uint16(1); i <= 5; i++ {
		r.Record([]alert.Event{makeEvent(alert.Opened, i)})
	}
	entries := r.Latest(3)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Event.Port.Port != 3 {
		t.Errorf("expected oldest of last-3 to be port 3, got %d", entries[0].Event.Port.Port)
	}
}

func TestRing_Eviction(t *testing.T) {
	r := history.New(3)
	for i := uint16(1); i <= 5; i++ {
		r.Record([]alert.Event{makeEvent(alert.Opened, i)})
	}
	if got := r.Len(); got != 3 {
		t.Fatalf("expected capacity-capped Len 3, got %d", got)
	}
	entries := r.Latest(0)
	if entries[0].Event.Port.Port != 3 {
		t.Errorf("expected evicted oldest, first entry port=3, got %d", entries[0].Event.Port.Port)
	}
}

func TestRecord_Empty(t *testing.T) {
	r := history.New(5)
	r.Record(nil)
	if r.Len() != 0 {
		t.Errorf("expected Len 0 after empty record, got %d", r.Len())
	}
}
