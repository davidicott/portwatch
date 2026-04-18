package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func makeRCEvent(kind, port string) alert.Event {
	return alert.Event{Kind: kind, Port: port}
}

func TestRocketChatNotifier_SkipsEmptyEvents(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewRocketChatNotifier(ts.URL)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatal("expected no request for empty events")
	}
}

func TestRocketChatNotifier_PostsPayload(t *testing.T) {
	var got map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &got)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewRocketChatNotifier(ts.URL)
	events := []alert.Event{makeRCEvent("opened", "tcp:8080")}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["text"] == "" {
		t.Fatal("expected non-empty text in payload")
	}
}

func TestRocketChatNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewRocketChatNotifier(ts.URL)
	events := []alert.Event{makeRCEvent("closed", "tcp:9090")}
	if err := n.Notify(events); err == nil {
		t.Fatal("expected error on non-2xx status")
	}
}

func TestRocketChatNotifier_PayloadContainsPortAndKind(t *testing.T) {
	var got map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &got)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewRocketChatNotifier(ts.URL)
	events := []alert.Event{makeRCEvent("opened", "tcp:4444")}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	text := got["text"]
	if text == "" {
		t.Fatal("expected non-empty text in payload")
	}
	for _, substr := range []string{"tcp:4444", "opened"} {
		if !contains(text, substr) {
			t.Errorf("expected payload text to contain %q, got: %s", substr, text)
		}
	}
}

// contains is a simple substring helper for test assertions.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		})())
}
