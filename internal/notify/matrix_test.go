package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeMatrixEvent(t alert.EventType, port int) alert.Event {
	return alert.Event{
		Type: t,
		Port: scanner.Port{Port: port, Protocol: "tcp"},
	}
}

func TestMatrixNotifier_SkipsEmptyEvents(t *testing.T) {
	n := NewMatrixNotifier("http://localhost", "token", "!room:localhost")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestMatrixNotifier_PostsPayload(t *testing.T) {
	var gotBody map[string]string
	var gotAuth string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"event_id":"$abc"}`))
	}))
	defer ts.Close()

	n := NewMatrixNotifier(ts.URL, "mytoken", "!room:example.org")
	events := []alert.Event{
		makeMatrixEvent("opened", 8080),
		makeMatrixEvent("closed", 22),
	}

	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotBody["msgtype"] != "m.text" {
		t.Errorf("expected msgtype m.text, got %q", gotBody["msgtype"])
	}
	if gotAuth != "Bearer mytoken" {
		t.Errorf("expected Bearer token, got %q", gotAuth)
	}
	if gotBody["body"] == "" {
		t.Error("expected non-empty body")
	}
}

func TestMatrixNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := NewMatrixNotifier(ts.URL, "bad", "!room:example.org")
	err := n.Notify([]alert.Event{makeMatrixEvent("opened", 443)})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}
