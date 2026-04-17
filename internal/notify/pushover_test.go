package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makePushoverEvent(kind, proto string, port uint16) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Proto: proto, Port: port},
	}
}

func TestPushoverNotifier_SkipsEmptyEvents(t *testing.T) {
	n := NewPushoverNotifier("tok", "user")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestPushoverNotifier_PostsPayload(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewPushoverNotifier("mytoken", "myuser")
	n.client.Transport = rewriteTransport(ts.URL)

	events := []alert.Event{makePushoverEvent("opened", "tcp", 8080)}
	if err := n.Notify(events); err != nil {
	unexpected error: %v", err)
	}

	if received["token"] != "mytoken" {
		t.Errorf("expected token mytoken, got %s", receivedt}
	if received["user"] != "myuser" {
		t.Errorf("expected user myuser, got %s", received["user"])
	}
	if received["message"] == "" {
		t.Error("expected non-empty message")
	}
}

func TestPushoverNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer ts.Close()

	n := NewPushoverNotifier("tok", "user")
	n.client.Transport = rewriteTransport(ts.URL)

	events := []alert.Event{makePushoverEvent("closed", "tcp", 443)}
	if err := n.Notify(events); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}
