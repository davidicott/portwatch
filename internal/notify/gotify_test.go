package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeGotifyEvent(kind, proto string, port int) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Host: "localhost", Port: port, Proto: proto},
	}
}

func TestGotifyNotifier_PostsPayload(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewGotifyNotifier(ts.URL, "testtoken", 7)
	events := []alert.Event{makeGotifyEvent("opened", "tcp", 8080)}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["title"] == "" {
		t.Error("expected non-empty title")
	}
	if received["priority"].(float64) != 7 {
		t.Errorf("expected priority 7, got %v", received["priority"])
	}
}

func TestGotifyNotifier_SkipsEmptyEvents(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewGotifyNotifier(ts.URL, "tok", 0)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty events")
	}
}

func TestGotifyNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := NewGotifyNotifier(ts.URL, "bad", 5)
	events := []alert.Event{makeGotifyEvent("closed", "tcp", 443)}
	if err := n.Notify(events); err == nil {
		t.Error("expected error for non-2xx status")
	}
}
