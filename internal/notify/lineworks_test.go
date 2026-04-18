package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func makeLWEvent(typ, host string, port int) alert.Event {
	return alert.Event{Type: typ, Host: host, Port: port, Protocol: "tcp"}
}

func TestLineWorksNotifier_SkipsEmptyEvents(t *testing.T) {
	n := NewLineWorksNotifier("http://example.com")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestLineWorksNotifier_PostsPayload(t *testing.T) {
	var body []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewLineWorksNotifier(ts.URL)
	events := []alert.Event{makeLWEvent("opened", "localhost", 8080)}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload map[string]string
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !strings.Contains(payload["content"], "8080") {
		t.Errorf("expected port in content, got: %s", payload["content"])
	}
}

func TestLineWorksNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewLineWorksNotifier(ts.URL)
	events := []alert.Event{makeLWEvent("closed", "localhost", 9090)}
	if err := n.Notify(events); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}
