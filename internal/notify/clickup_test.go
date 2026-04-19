package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeCUEvent(kind alert.EventKind, port uint16, proto string) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Port: port, Proto: proto},
		At:   time.Now(),
	}
}

func TestClickUpNotifier_SkipsEmptyEvents(t *testing.T) {
	n := NewClickUpNotifier("token", "list123")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestClickUpNotifier_PostsPayload(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		if r.Header.Get("Authorization") != "mytoken" {
			t.Errorf("missing auth header")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewClickUpNotifier("mytoken", "list42")
	n.apiBase = ts.URL

	events := []alert.Event{makeCUEvent(alert.EventOpened, 8080, "tcp")}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["name"] == "" {
		t.Error("expected task name to be set")
	}
}

func TestClickUpNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := NewClickUpNotifier("bad", "list1")
	n.apiBase = ts.URL

	events := []alert.Event{makeCUEvent(alert.EventClosed, 443, "tcp")}
	if err := n.Notify(events); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}
