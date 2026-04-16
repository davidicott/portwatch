package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeTGEvent(proto, addr string, port uint16, kind alert.EventKind) alert.Event {
	return alert.Event{
		Port: scanner.Port{Protocol: proto, Addr: addr, Port: port},
		Kind: kind,
	}
}

func TestTelegramNotifier_PostsPayload(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := &TelegramNotifier{
		token:  "testtoken",
		chatID: "12345",
		client: ts.Client(),
	}
	// Override base URL by patching via custom server — use direct field for test.
	n.client = &http.Client{
		Transport: rewriteTransport(ts.URL),
	}

	events := []alert.Event{makeTGEvent("tcp", "0.0.0.0", 8080, alert.Opened)}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTelegramNotifier_SkipsEmptyEvents(t *testing.T) {
	n := NewTelegramNotifier("tok", "cid")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestTelegramNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := &TelegramNotifier{
		token:  "bad",
		chatID: "cid",
		client: &http.Client{Transport: rewriteTransport(ts.URL)},
	}

	events := []alert.Event{makeTGEvent("tcp", "0.0.0.0", 9090, alert.Closed)}
	if err := n.Notify(events); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}
