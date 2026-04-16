package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func makePDEvent(kind, proto, addr string, port int) alert.Event {
	return makeEvent(kind, proto, addr, port)
}

func TestPagerDutyNotifier_PostsPayload(t *testing.T) {
	var received pagerDutyPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	n := NewPagerDutyNotifier("test-key", "portwatch")
	n.url = ts.URL

	events := []alert.Event{makePDEvent("opened", "tcp", "0.0.0.0", 8080)}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.RoutingKey != "test-key" {
		t.Errorf("routing key = %q, want %q", received.RoutingKey, "test-key")
	}
	if received.EventAction != "trigger" {
		t.Errorf("event_action = %q, want trigger", received.EventAction)
	}
	if received.Payload.Source != "portwatch" {
		t.Errorf("source = %q, want portwatch", received.Payload.Source)
	}
	if received.Payload.Severity != "warning" {
		t.Errorf("severity = %q, want warning", received.Payload.Severity)
	}
}

func TestPagerDutyNotifier_SkipsEmptyEvents(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewPagerDutyNotifier("key", "src")
	n.url = ts.URL
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty events")
	}
}

func TestPagerDutyNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	n := NewPagerDutyNotifier("key", "src")
	n.url = ts.URL
	events := []alert.Event{makePDEvent("closed", "tcp", "127.0.0.1", 22)}
	if err := n.Notify(events); err == nil {
		t.Error("expected error for non-2xx status")
	}
}
