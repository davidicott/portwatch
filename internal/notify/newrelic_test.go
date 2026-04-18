package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeNREvent(proto string, port uint16, kind alert.Kind) alert.Event {
	return alert.Event{
		Port: scanner.Port{Protocol: proto, Port: port},
		Kind: kind,
	}
}

func TestNewRelicNotifier_SkipsEmptyEvents(t *testing.T) {
	n := NewNewRelicNotifier("key", "")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestNewRelicNotifier_PostsPayload(t *testing.T) {
	var received []map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Api-Key") == "" {
			t.Error("expected Api-Key header")
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	n := NewNewRelicNotifier("test-key", ts.URL)
	events := []alert.Event{
		makeNREvent("tcp", 8080, alert.Opened),
	}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(received) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(received))
	}
	if received[0]["kind"] != "opened" {
		t.Errorf("expected kind=opened, got %v", received[0]["kind"])
	}
}

func TestNewRelicNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := NewNewRelicNotifier("bad-key", ts.URL)
	events := []alert.Event{makeNREvent("tcp", 443, alert.Closed)}
	if err := n.Notify(events); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}
