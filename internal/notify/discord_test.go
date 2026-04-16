package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeDiscordEvent(kind, proto, addr string, pid int) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Proto: proto, Addr: addr, PID: pid},
	}
}

func TestDiscordNotifier_PostsPayload(t *testing.T) {
	var received []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	n := NewDiscordNotifier(ts.URL)
	events := []alert.Event{makeDiscordEvent("opened", "tcp", "0.0.0.0:9090", 555)}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload map[string]string
	if err := json.Unmarshal(received, &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !strings.Contains(payload["content"], "9090") {
		t.Errorf("expected content to contain port, got: %s", payload["content"])
	}
}

func TestDiscordNotifier_SkipsEmptyEvents(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewDiscordNotifier(ts.URL)
	if err := n.Notify([]alert.Event{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty events")
	}
}

func TestDiscordNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	n := NewDiscordNotifier(ts.URL)
	events := []alert.Event{makeDiscordEvent("closed", "tcp", "127.0.0.1:22", 1)}
	if err := n.Notify(events); err == nil {
		t.Error("expected error for non-2xx status")
	}
}
