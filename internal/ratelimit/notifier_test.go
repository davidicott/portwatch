package ratelimit_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/ratelimit"
	"github.com/user/portwatch/internal/scanner"
)

type captureNotifier struct {
	calls [][]alert.Event
}

func (c *captureNotifier) Notify(events []alert.Event) error {
	c.calls = append(c.calls, events)
	return nil
}

func makeRLEvent(proto string, port uint16) alert.Event {
	return alert.Event{
		Kind: "opened",
		Port: scanner.Port{Proto: proto, Port: port},
	}
}

func TestNotifierWrapper_AllowsUnderLimit(t *testing.T) {
	cap := &captureNotifier{}
	l := ratelimit.New(2, time.Minute)
	w := ratelimit.NewNotifierWrapper(cap, l)

	_ = w.Notify([]alert.Event{makeRLEvent("tcp", 80)})
	_ = w.Notify([]alert.Event{makeRLEvent("tcp", 80)})

	if len(cap.calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(cap.calls))
	}
}

func TestNotifierWrapper_BlocksOverLimit(t *testing.T) {
	cap := &captureNotifier{}
	l := ratelimit.New(1, time.Minute)
	w := ratelimit.NewNotifierWrapper(cap, l)

	_ = w.Notify([]alert.Event{makeRLEvent("tcp", 443)})
	_ = w.Notify([]alert.Event{makeRLEvent("tcp", 443)})

	if len(cap.calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(cap.calls))
	}
}

func TestNotifierWrapper_SkipsEmptyEvents(t *testing.T) {
	cap := &captureNotifier{}
	l := ratelimit.New(5, time.Minute)
	w := ratelimit.NewNotifierWrapper(cap, l)

	_ = w.Notify(nil)

	if len(cap.calls) != 0 {
		t.Fatalf("expected 0 calls, got %d", len(cap.calls))
	}
}
