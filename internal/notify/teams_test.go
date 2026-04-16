package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeTeamsEvent(kind, proto, addr string, pid int) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Proto: proto, Addr: addr, PID: pid},
	}
}

func TestTeamsNotifier_PostsPayload(t *testing.T) {
	var received []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewTeamsNotifier(ts.URL)
	events := []alert.Event{makeTeamsEvent("opened", "tcp", "0.0.0.0:8080", 1234)}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload map[string]string
	if err := json.Unmarshal(received, &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if payload["text"] == "" {
		t.Error("expected non-empty text field")
	}
}

func TestTeamsNotifier_SkipsEmptyEvents(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewTeamsNotifier(ts.URL)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty events")
	}
}

func TestTeamsNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewTeamsNotifier(ts.URL)
	events := []alert.Event{makeTeamsEvent("closed", "udp", "0.0.0.0:53", 99)}
	if err := n.Notify(events); err == nil {
		t.Error("expected error for non-2xx status")
	}
}
