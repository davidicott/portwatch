package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeWebexEvent(kind, proto string, port int) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Proto: proto, Port: port},
	}
}

func TestWebexNotifier_SkipsEmptyEvents(t *testing.T) {
	n := NewWebexNotifier("tok", "room1")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestWebexNotifier_PostsPayload(t *testing.T) {
	var received map[string]string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		auth := r.Header.Get("Authorization")
		if auth != "Bearer mytoken" {
			t.Errorf("expected Bearer mytoken, got %q", auth)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := NewWebexNotifier("mytoken", "roomABC")
	n.apiURL = srv.URL

	events := []alert.Event{
		makeWebexEvent("opened", "tcp", 8080),
	}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["roomId"] != "roomABC" {
		t.Errorf("expected roomId=roomABC, got %q", received["roomId"])
	}
	if received["text"] == "" {
		t.Error("expected non-empty text")
	}
}

func TestWebexNotifier_NonSuccessStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	n := NewWebexNotifier("bad", "room")
	n.apiURL = srv.URL

	err := n.Notify([]alert.Event{makeWebexEvent("opened", "tcp", 443)})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}
