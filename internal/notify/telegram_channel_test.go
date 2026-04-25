package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func makeTCEvent(kind, proto string, port uint16) alert.Event {
	return alert.Event{
		Kind:     kind,
		Protocol: proto,
		Port:     port,
	}
}

func TestTelegramChannelNotifier_SkipsEmptyEvents(t *testing.T) {
	n := NewTelegramChannelNotifier("token", "@channel")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("expected no error for empty events, got %v", err)
	}
}

func TestTelegramChannelNotifier_PostsPayload(t *testing.T) {
	var received map[string]interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer ts.Close()

	n := NewTelegramChannelNotifier("mytoken", "-100123456")
	n.baseURL = ts.URL

	events := []alert.Event{makeTCEvent("opened", "tcp", 8080)}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["chat_id"] != "-100123456" {
		t.Errorf("expected chat_id -100123456, got %v", received["chat_id"])
	}
	if received["parse_mode"] != "Markdown" {
		t.Errorf("expected parse_mode Markdown, got %v", received["parse_mode"])
	}
	if received["text"] == "" {
		t.Error("expected non-empty text")
	}
}

func TestTelegramChannelNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := NewTelegramChannelNotifier("badtoken", "@channel")
	n.baseURL = ts.URL

	events := []alert.Event{makeTCEvent("closed", "udp", 53)}
	if err := n.Notify(events); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}
