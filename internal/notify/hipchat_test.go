package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func makeHCEvent(kind, port string) alert.Event {
	return alert.Event{Kind: kind, Port: port}
}

func TestHipChatNotifier_SkipsEmptyEvents(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewHipChatNotifier(ts.URL, "42", "token")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatal("expected no HTTP call for empty events")
	}
}

func TestHipChatNotifier_PostsPayload(t *testing.T) {
	var got hipChatPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Errorf("decode: %v", err)
		}
		auth := r.Header.Get("Authorization")
		if auth != "Bearer mytoken" {
			t.Errorf("unexpected auth header: %s", auth)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	n := NewHipChatNotifier(ts.URL, "42", "mytoken")
	events := []alert.Event{makeHCEvent("opened", "tcp:8080")}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Message == "" {
		t.Error("expected non-empty message")
	}
	if !got.Notify {
		t.Error("expected notify=true")
	}
}

func TestHipChatNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := NewHipChatNotifier(ts.URL, "42", "bad")
	err := n.Notify([]alert.Event{makeHCEvent("opened", "tcp:9090")})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}
