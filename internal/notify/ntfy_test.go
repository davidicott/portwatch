package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func makeNtfyEvent(msg string) alert.Event {
	return alert.Event{Message: msg}
}

func TestNtfyNotifier_SkipsEmptyEvents(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewNtfyNotifier(ts.URL, "portwatch")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatal("expected no HTTP call for empty events")
	}
}

func TestNtfyNotifier_PostsPayload(t *testing.T) {
	var gotTitle string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotTitle = r.Header.Get("Title")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewNtfyNotifier(ts.URL, "portwatch")
	events := []alert.Event{makeNtfyEvent("port 8080 opened"), makeNtfyEvent("port 9090 closed")}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotTitle == "" {
		t.Fatal("expected Title header to be set")
	}
}

func TestNtfyNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := NewNtfyNotifier(ts.URL, "portwatch")
	err := n.Notify([]alert.Event{makeNtfyEvent("port 22 opened")})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestNtfyNotifier_DefaultServer(t *testing.T) {
	n := NewNtfyNotifier("", "alerts")
	if n.serverURL != "https://ntfy.sh" {
		t.Fatalf("expected default server, got %s", n.serverURL)
	}
}
