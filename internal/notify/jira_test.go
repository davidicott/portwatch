package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeJiraEvent(kind, proto string, port int) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Protocol: proto, Port: port},
	}
}

func TestJiraNotifier_SkipsEmptyEvents(t *testing.T) {
	n := NewJiraNotifier("http://example.com", "user", "token", "OPS", "")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestJiraNotifier_PostsPayload(t *testing.T) {
	var received map[string]any
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/2/issue" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	n := NewJiraNotifier(ts.URL, "user", "token", "OPS", "Bug")
	events := []alert.Event{makeJiraEvent("opened", "tcp", 8080)}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fields, ok := received["fields"].(map[string]any)
	if !ok {
		t.Fatal("missing fields")
	}
	if fields["summary"] == "" {
		t.Error("expected non-empty summary")
	}
}

func TestJiraNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := NewJiraNotifier(ts.URL, "user", "token", "OPS", "")
	events := []alert.Event{makeJiraEvent("closed", "tcp", 443)}
	if err := n.Notify(events); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}
