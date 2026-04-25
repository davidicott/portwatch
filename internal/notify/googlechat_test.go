package notify

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func makeGCEvent(kind, port string) alert.Event {
	return alert.Event{Kind: kind, Port: port}
}

func TestGoogleChatNotifier_SkipsEmptyEvents(t *testing.T) {
	n := NewGoogleChatNotifier("http://example.com")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGoogleChatNotifier_PostsPayload(t *testing.T) {
	var received string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf strings.Builder
		buf.ReadFrom(r.Body)
		received = buf.String()
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewGoogleChatNotifier(ts.URL)
	events := []alert.Event{
		makeGCEvent("opened", "tcp:8080"),
		makeGCEvent("closed", "tcp:9090"),
	}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsString(received, "tcp:8080") {
		t.Errorf("expected payload to contain tcp:8080, got: %s", received)
	}
	if !containsString(received, "2 change(s)") {
		t.Errorf("expected payload to mention change count, got: %s", received)
	}
}

func TestGoogleChatNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewGoogleChatNotifier(ts.URL)
	err := n.Notify([]alert.Event{makeGCEvent("opened", "tcp:80")})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func containsString(s, sub string) bool {
	return strings.Contains(s, sub)
}
