package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeGCEvent(kind, proto string, port int) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Protocol: proto, Port: port},
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
	var captured map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &captured)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewGoogleChatNotifier(ts.URL)
	events := []alert.Event{
		makeGCEvent("opened", "tcp", 8080),
		makeGCEvent("closed", "udp", 53),
	}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	text, ok := captured["text"].(string)
	if !ok || text == "" {
		t.Fatal("expected non-empty text field in payload")
	}
	if len(text) == 0 {
		t.Fatal("expected payload text to contain event info")
	}
}

func TestGoogleChatNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewGoogleChatNotifier(ts.URL)
	events := []alert.Event{makeGCEvent("opened", "tcp", 443)}
	if err := n.Notify(events); err == nil {
		t.Fatal("expected error on non-2xx status")
	}
}
