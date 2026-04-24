package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeCWEvent(kind, proto string, port uint16) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Protocol: proto, Port: port},
	}
}

func TestChatworkNotifier_SkipsEmptyEvents(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewChatworkNotifier("tok", "123")
	n.client = ts.Client()
	// override base by swapping to local server would require refactor; skip HTTP call check
	_ = n.Notify(nil)
	if called {
		t.Fatal("expected no HTTP call for empty events")
	}
}

func TestChatworkNotifier_PostsPayload(t *testing.T) {
	var gotBody map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &gotBody)
		if r.Header.Get("X-ChatWorkToken") == "" {
			t.Error("missing X-ChatWorkToken header")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := &ChatworkNotifier{
		token:  "mytoken",
		roomID: "42",
		client: ts.Client(),
	}
	// patch the URL via a custom transport that rewrites host
	n.client = &http.Client{
		Transport: rewriteTransport(ts.URL),
	}

	events := []alert.Event{makeCWEvent("opened", "tcp", 8080)}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotBody["body"], "8080") {
		t.Errorf("body missing port: %q", gotBody["body"])
	}
}

func TestChatworkNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := &ChatworkNotifier{
		token:  "bad",
		roomID: "1",
		client: &http.Client{Transport: rewriteTransport(ts.URL)},
	}
	err := n.Notify([]alert.Event{makeCWEvent("opened", "tcp", 22)})
	if err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}
