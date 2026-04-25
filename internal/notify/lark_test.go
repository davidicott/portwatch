package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func makeLarkEvent(kind, port string) alert.Event {
	return alert.Event{Kind: kind, Port: port}
}

func TestLarkNotifier_SkipsEmptyEvents(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewLarkNotifier(ts.URL)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty events")
	}
}

func TestLarkNotifier_PostsPayload(t *testing.T) {
	var received larkPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewLarkNotifier(ts.URL)
	events := []alert.Event{
		makeLarkEvent("opened", "tcp:8080"),
		makeLarkEvent("closed", "tcp:9090"),
	}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.MsgType != "text" {
		t.Errorf("expected msg_type=text, got %q", received.MsgType)
	}
	if received.Content.Text == "" {
		t.Error("expected non-empty text content")
	}
}

func TestLarkNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewLarkNotifier(ts.URL)
	events := []alert.Event{makeLarkEvent("opened", "tcp:443")}
	err := n.Notify(events)
	if err == nil {
		t.Error("expected error on non-2xx response")
	}
}
