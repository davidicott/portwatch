package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/patrickdappollonio/portwatch/internal/alert"
	"github.com/patrickdappollonio/portwatch/internal/scanner"
)

func makeSPEvent(kind, proto string, port int) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Protocol: proto, Port: port},
	}
}

func TestStatuspageNotifier_SkipsEmptyEvents(t *testing.T) {
	n := NewStatuspageNotifier("key", "page", "comp", "")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("expected no error for empty events, got %v", err)
	}
}

func TestStatuspageNotifier_PostsPayload(t *testing.T) {
	var received map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		if r.Header.Get("Authorization") == "" {
			t.Error("expected Authorization header")
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	n := NewStatuspageNotifier("mykey", "mypageid", "mycompid", srv.URL)
	events := []alert.Event{
		makeSPEvent("opened", "tcp", 8080),
		makeSPEvent("closed", "tcp", 22),
	}

	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	inc, ok := received["incident"].(map[string]interface{})
	if !ok {
		t.Fatal("expected incident key in payload")
	}
	name, _ := inc["name"].(string)
	if name == "" {
		t.Error("expected non-empty incident name")
	}
}

func TestStatuspageNotifier_NonSuccessStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	n := NewStatuspageNotifier("bad-key", "page", "comp", srv.URL)
	err := n.Notify([]alert.Event{makeSPEvent("opened", "tcp", 9090)})
	if err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}

func TestStatuspageNotifier_DefaultEndpoint(t *testing.T) {
	n := NewStatuspageNotifier("k", "p", "c", "")
	if n.endpoint != defaultStatuspageEndpoint {
		t.Errorf("expected default endpoint %q, got %q", defaultStatuspageEndpoint, n.endpoint)
	}
}
