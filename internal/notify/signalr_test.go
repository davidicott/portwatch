package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func makeSignalREvent(t string) alert.Event {
	return alert.Event{Type: t, Host: "localhost", Port: 8080, Protocol: "tcp"}
}

func TestSignalRNotifier_SkipsEmptyEvents(t *testing.T) {
	n := NewSignalRNotifier("http://example.com", "testhub", "key")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestSignalRNotifier_PostsPayload(t *testing.T) {
	var received signalRPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode: %v", err)
		}
		if r.Header.Get("Authorization") != "Bearer testkey" {
			t.Errorf("missing auth header")
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	n := NewSignalRNotifier(ts.URL, "myhub", "testkey")
	if err := n.Notify([]alert.Event{makeSignalREvent("opened")}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Target != "portwatch" {
		t.Errorf("expected target portwatch, got %s", received.Target)
	}
	if len(received.Arguments) != 1 {
		t.Errorf("expected 1 argument, got %d", len(received.Arguments))
	}
}

func TestSignalRNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewSignalRNotifier(ts.URL, "hub", "")
	err := n.Notify([]alert.Event{makeSignalREvent("closed")})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}
