package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeFlowdockEvent(kind, proto string, port uint16) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Protocol: proto, Port: port},
	}
}

func TestFlowdockNotifier_SkipsEmptyEvents(t *testing.T) {
	n := NewFlowdockNotifier("tok", "flow1")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestFlowdockNotifier_PostsPayload(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	n := NewFlowdockNotifier("mytoken", "myflow")
	n.apiURL = ts.URL

	events := []alert.Event{makeFlowdockEvent("opened", "tcp", 8080)}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["flow_token"] != "mytoken" {
		t.Errorf("expected token mytoken, got %v", received["flow_token"])
	}
	if received["event"] != "message" {
		t.Errorf("expected event=message, got %v", received["event"])
	}
	if received["content"] == "" {
		t.Error("expected non-empty content")
	}
}

func TestFlowdockNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := NewFlowdockNotifier("bad", "flow")
	n.apiURL = ts.URL

	events := []alert.Event{makeFlowdockEvent("closed", "udp", 53)}
	if err := n.Notify(events); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}
