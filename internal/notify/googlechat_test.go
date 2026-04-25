package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeGCEvent(kind, host string, port int, proto string) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Host: host, Port: port, Proto: proto},
	}
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
	var got map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewGoogleChatNotifier(ts.URL)
	events := []alert.Event{makeGCEvent("opened", "localhost", 8080, "tcp")}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got["text"], "8080") {
		t.Errorf("expected payload to contain port, got: %v", got)
	}
}

func TestGoogleChatNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewGoogleChatNotifier(ts.URL)
	events := []alert.Event{makeGCEvent("opened", "localhost", 9090, "udp")}
	if err := n.Notify(events); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func containsString(s, sub string) bool {
	return strings.Contains(s, sub)
}
