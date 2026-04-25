package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makePDV2Event(kind alert.EventKind, port uint16) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Port: port, Proto: "tcp"},
	}
}

func TestPagerDutyV2Notifier_SkipsEmptyEvents(t *testing.T) {
	n := NewPagerDutyV2Notifier("key", "")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestPagerDutyV2Notifier_PostsPayload(t *testing.T) {
	var received pdV2Payload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	n := NewPagerDutyV2Notifier("test-routing-key", srv.URL)
	events := []alert.Event{
		makePDV2Event(alert.EventOpened, 8080),
		makePDV2Event(alert.EventClosed, 22),
	}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.RoutingKey != "test-routing-key" {
		t.Errorf("routing key: got %q, want %q", received.RoutingKey, "test-routing-key")
	}
	if received.EventAction != "trigger" {
		t.Errorf("event_action: got %q, want trigger", received.EventAction)
	}
	if received.Payload.Severity != "warning" {
		t.Errorf("severity: got %q, want warning", received.Payload.Severity)
	}
}

func TestPagerDutyV2Notifier_NonSuccessStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	n := NewPagerDutyV2Notifier("key", srv.URL)
	err := n.Notify([]alert.Event{makePDV2Event(alert.EventOpened, 443)})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestPagerDutyV2Notifier_DefaultEndpoint(t *testing.T) {
	n := NewPagerDutyV2Notifier("key", "")
	if n.endpoint != defaultPagerDutyV2URL {
		t.Errorf("endpoint: got %q, want %q", n.endpoint, defaultPagerDutyV2URL)
	}
}
