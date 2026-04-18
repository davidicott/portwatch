package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/patrickdappollonio/portwatch/internal/alert"
)

func makeBCEvent(kind, msg string) alert.Event {
	return alert.Event{Kind: kind, Message: msg}
}

func TestBearyChat_SkipsEmptyEvents(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewBearyChat(ts.URL)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatal("expected no HTTP call for empty events")
	}
}

func TestBearyChat_PostsPayload(t *testing.T) {
	var got bearyChatPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewBearyChat(ts.URL)
	events := []alert.Event{
		makeBCEvent("opened", "port 8080/tcp opened"),
		makeBCEvent("closed", "port 9090/tcp closed"),
	}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Text == "" {
		t.Error("expected non-empty text in payload")
	}
	if got.Notification == "" {
		t.Error("expected non-empty notification field")
	}
}

func TestBearyChat_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewBearyChat(ts.URL)
	err := n.Notify([]alert.Event{makeBCEvent("opened", "port 22/tcp")})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}
