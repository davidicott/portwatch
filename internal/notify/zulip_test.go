package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeZulipEvent(kind, addr string, port int) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Address: addr, Port: port, Protocol: "tcp"},
	}
}

func TestZulipNotifier_PostsPayload(t *testing.T) {
	var got map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewZulipNotifier(ts.URL, "bot@example.com", "secret", "alerts", "portwatch")
	events := []alert.Event{makeZulipEvent("opened", "0.0.0.0", 9090)}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got["type"] != "stream" {
		t.Errorf("expected type=stream, got %q", got["type"])
	}
	if got["to"] != "alerts" {
		t.Errorf("expected to=alerts, got %q", got["to"])
	}
	if got["topic"] != "portwatch" {
		t.Errorf("expected topic=portwatch, got %q", got["topic"])
	}
}

func TestZulipNotifier_SkipsEmptyEvents(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewZulipNotifier(ts.URL, "bot@example.com", "secret", "alerts", "portwatch")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty events")
	}
}

func TestZulipNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := NewZulipNotifier(ts.URL, "bot@example.com", "bad-key", "alerts", "portwatch")
	events := []alert.Event{makeZulipEvent("opened", "0.0.0.0", 22)}
	if err := n.Notify(events); err == nil {
		t.Error("expected error on non-2xx status")
	}
}
