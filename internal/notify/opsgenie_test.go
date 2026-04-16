package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeOGEvent(kind, proto string, port int) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Protocol: proto, Number: port, Address: "127.0.0.1"},
	}
}

func TestOpsGenieNotifier_PostsPayload(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "GenieKey test-key" {
			t.Errorf("missing or wrong auth header")
		}
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	n := NewOpsGenieNotifier("test-key")
	n.url = ts.URL

	err := n.Notify([]alert.Event{makeOGEvent("opened", "tcp", 8080)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["message"] == "" {
		t.Error("expected message in payload")
	}
}

func TestOpsGenieNotifier_SkipsEmptyEvents(t *testing.T) {
	n := NewOpsGenieNotifier("key")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestOpsGenieNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := NewOpsGenieNotifier("bad-key")
	n.url = ts.URL

	err := n.Notify([]alert.Event{makeOGEvent("closed", "udp", 53)})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}
