package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeSplunkEvent(kind string) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Number: 9200, Protocol: "tcp", Process: "elasticsearch"},
	}
}

func TestSplunkNotifier_SkipsEmptyEvents(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewSplunkNotifier(ts.URL, "token", "portwatch")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatal("expected no HTTP call for empty events")
	}
}

func TestSplunkNotifier_PostsPayload(t *testing.T) {
	var received []map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Splunk mytoken" {
			t.Errorf("missing or wrong Authorization header")
		}
		dec := json.NewDecoder(r.Body)
		for dec.More() {
			var m map[string]interface{}
			_ = dec.Decode(&m)
			received = append(received, m)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewSplunkNotifier(ts.URL, "mytoken", "portwatch")
	events := []alert.Event{makeSplunkEvent("opened"), makeSplunkEvent("closed")}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(received) != 2 {
		t.Fatalf("expected 2 events, got %d", len(received))
	}
}

func TestSplunkNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := NewSplunkNotifier(ts.URL, "bad", "portwatch")
	err := n.Notify([]alert.Event{makeSplunkEvent("opened")})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}
