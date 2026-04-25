package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeAmplitudeEvent(kind, proto string, port int) alert.Event {
	return alert.Event{
		Kind:      kind,
		Timestamp: time.Now(),
		Port: scanner.Port{
			Number:   port,
			Protocol: proto,
			Address:  "0.0.0.0",
		},
	}
}

func TestAmplitudeNotifier_SkipsEmptyEvents(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewAmplitudeNotifier("key123", ts.URL)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatal("expected no HTTP call for empty events")
	}
}

func TestAmplitudeNotifier_PostsPayload(t *testing.T) {
	var received amplitudePayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewAmplitudeNotifier("testkey", ts.URL)
	events := []alert.Event{
		makeAmplitudeEvent("opened", "tcp", 8080),
		makeAmplitudeEvent("closed", "udp", 53),
	}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.APIKey != "testkey" {
		t.Errorf("expected api_key 'testkey', got %q", received.APIKey)
	}
	if len(received.Events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(received.Events))
	}
	if received.Events[0].EventType != "port_opened" {
		t.Errorf("expected event_type 'port_opened', got %q", received.Events[0].EventType)
	}
	if received.Events[1].EventType != "port_closed" {
		t.Errorf("expected event_type 'port_closed', got %q", received.Events[1].EventType)
	}
}

func TestAmplitudeNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := NewAmplitudeNotifier("badkey", ts.URL)
	err := n.Notify([]alert.Event{makeAmplitudeEvent("opened", "tcp", 443)})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestAmplitudeNotifier_DefaultEndpoint(t *testing.T) {
	n := NewAmplitudeNotifier("key", "")
	if n.endpoint != defaultAmplitudeEndpoint {
		t.Errorf("expected default endpoint %q, got %q", defaultAmplitudeEndpoint, n.endpoint)
	}
}
