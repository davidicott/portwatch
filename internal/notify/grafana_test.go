package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/patrickdappollonio/portwatch/internal/alert"
	"github.com/patrickdappollonio/portwatch/internal/scanner"
)

func makeGrafanaEvent(kind, proto string, port int) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Proto: proto, Port: port},
	}
}

func TestGrafanaNotifier_SkipsEmptyEvents(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewGrafanaNotifier(ts.URL, "test", nil)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatal("expected no request for empty events")
	}
}

func TestGrafanaNotifier_PostsPayload(t *testing.T) {
	var got map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewGrafanaNotifier(ts.URL, "portwatch", nil)
	events := []alert.Event{makeGrafanaEvent("opened", "tcp", 8080)}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got["title"] != "portwatch" {
		t.Errorf("expected title 'portwatch', got %q", got["title"])
	}
	if got["state"] != "alerting" {
		t.Errorf("expected state 'alerting', got %q", got["state"])
	}
	if got["message"] == "" {
		t.Error("expected non-empty message")
	}
}

func TestGrafanaNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewGrafanaNotifier(ts.URL, "", nil)
	events := []alert.Event{makeGrafanaEvent("closed", "tcp", 9090)}
	if err := n.Notify(events); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}
