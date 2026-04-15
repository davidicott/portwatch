package notify_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/notify"
)

func TestWebhookNotifier_PostsPayload(t *testing.T) {
	var received []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		received = body
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := notify.NewWebhookNotifier(srv.URL, time.Second)
	events := []alert.Event{makeEvent("opened", "tcp:443")}
	if err := n.Notify(context.Background(), events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var payload struct {
		Events []alert.Event `json:"events"`
	}
	if err := json.Unmarshal(received, &payload); err != nil {
		t.Fatalf("could not unmarshal payload: %v", err)
	}
	if len(payload.Events) != 1 || payload.Events[0].Port != "tcp:443" {
		t.Errorf("unexpected payload: %+v", payload)
	}
}

func TestWebhookNotifier_NonSuccessStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	n := notify.NewWebhookNotifier(srv.URL, time.Second)
	err := n.Notify(context.Background(), []alert.Event{makeEvent("closed", "tcp:80")})
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestWebhookNotifier_SkipsEmptyEvents(t *testing.T) {
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer srv.Close()

	n := notify.NewWebhookNotifier(srv.URL, time.Second)
	if err := n.Notify(context.Background(), nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP request for empty event list")
	}
}
