package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeVOEvent(proto string, port int, kind string) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Protocol: proto, Number: port},
	}
}

func TestVictorOpsNotifier_PostsPayload(t *testing.T) {
	var got map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewVictorOpsNotifier("key123", ts.URL)
	err := n.Notify([]alert.Event{makeVOEvent("tcp", 8080, "opened")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["message_type"] != "CRITICAL" {
		t.Errorf("expected CRITICAL, got %s", got["message_type"])
	}
	if got["routing_key"] != "key123" {
		t.Errorf("expected key123, got %s", got["routing_key"])
	}
	if got["entity_id"] != "portwatch-tcp-8080" {
		t.Errorf("unexpected entity_id: %s", got["entity_id"])
	}
}

func TestVictorOpsNotifier_SkipsEmptyEvents(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewVictorOpsNotifier("key", ts.URL)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty events")
	}
}

func TestVictorOpsNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewVictorOpsNotifier("key", ts.URL)
	err := n.Notify([]alert.Event{makeVOEvent("udp", 53, "closed")})
	if err == nil {
		t.Error("expected error for non-2xx status")
	}
}
