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

func makeChimeEvent(kind alert.EventKind, port uint16, proto string) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Port: port, Proto: proto},
	}
}

func TestChimeNotifier_SkipsEmptyEvents(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewChimeNotifier(ts.URL)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatal("expected no HTTP call for empty events")
	}
}

func TestChimeNotifier_PostsPayload(t *testing.T) {
	var gotBody []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewChimeNotifier(ts.URL)
	events := []alert.Event{
		makeChimeEvent(alert.EventOpened, 8080, "tcp"),
	}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload map[string]string
	if err := json.Unmarshal(gotBody, &payload); err != nil {
		t.Fatalf("invalid JSON payload: %v", err)
	}
	content, ok := payload["Content"]
	if !ok {
		t.Fatal("payload missing Content field")
	}
	if !strings.Contains(content, "8080") {
		t.Errorf("expected port 8080 in content, got: %s", content)
	}
}

func TestChimeNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewChimeNotifier(ts.URL)
	events := []alert.Event{
		makeChimeEvent(alert.EventClosed, 443, "tcp"),
	}
	if err := n.Notify(events); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}
