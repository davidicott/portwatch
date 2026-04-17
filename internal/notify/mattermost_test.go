package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeMattermostEvent(kind, proto string, port uint16) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Proto: proto, Port: port},
	}
}

func TestMattermostNotifier_PostsPayload(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewMattermostNotifier(ts.URL, "#alerts")
	events := []alert.Event{makeMattermostEvent("opened", "tcp", 8080)}

	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["channel"] != "#alerts" {
		t.Errorf("expected channel #alerts, got %v", received["channel"])
	}
	text, _ := received["text"].(string)
	if text == "" {
		t.Error("expected non-empty text")
	}
}

func TestMattermostNotifier_SkipsEmptyEvents(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewMattermostNotifier(ts.URL, "")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty events")
	}
}

func TestMattermostNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewMattermostNotifier(ts.URL, "")
	events := []alert.Event{makeMattermostEvent("closed", "tcp", 443)}
	if err := n.Notify(events); err == nil {
		t.Error("expected error on non-2xx status")
	}
}
