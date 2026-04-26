package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func makeGCEvent(kind, port string) alert.Event {
	return alert.Event{Kind: kind, Port: port}
}

func TestGoogleChatNotifier_SkipsEmptyEvents(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewGoogleChatNotifier(ts.URL)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatal("expected no HTTP call for empty events")
	}
}

func TestGoogleChatNotifier_PostsPayload(t *testing.T) {
	var received []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewGoogleChatNotifier(ts.URL)
	events := []alert.Event{
		makeGCEvent("opened", "tcp:8080"),
		makeGCEvent("closed", "tcp:9090"),
	}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload map[string]string
	if err := json.Unmarshal(received, &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	text := payload["text"]
	if !containsString(text, "tcp:8080") {
		t.Errorf("expected text to contain port, got: %s", text)
	}
	if !containsString(text, "2 port change") {
		t.Errorf("expected text to contain count, got: %s", text)
	}
}

func TestGoogleChatNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewGoogleChatNotifier(ts.URL)
	err := n.Notify([]alert.Event{makeGCEvent("opened", "tcp:443")})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func containsString(s, sub string) bool {
	return strings.Contains(s, sub)
}
