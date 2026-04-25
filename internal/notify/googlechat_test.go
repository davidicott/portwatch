package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func makeGCEvent(port, kind string) alert.Event {
	return alert.Event{Port: port, Kind: kind}
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
		t.Error("expected no HTTP call for empty events")
	}
}

func TestGoogleChatNotifier_PostsPayload(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewGoogleChatNotifier(ts.URL)
	events := []alert.Event{
		makeGCEvent("tcp:8080", "opened"),
		makeGCEvent("tcp:9090", "closed"),
	}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["text"] == "" {
		t.Error("expected non-empty text in payload")
	}
	for _, e := range events {
		if !containsString(received["text"], e.Port) {
			t.Errorf("expected payload to contain port %q", e.Port)
		}
	}
}

func TestGoogleChatNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewGoogleChatNotifier(ts.URL)
	err := n.Notify([]alert.Event{makeGCEvent("tcp:22", "opened")})
	if err == nil {
		t.Error("expected error for non-2xx response")
	}
}

func containsString(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
