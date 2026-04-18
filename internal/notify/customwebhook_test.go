package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeCustomEvent(proto, addr string, port uint16) alert.Event {
	return alert.Event{
		Type: alert.Opened,
		Port: scanner.Port{Protocol: proto, Address: addr, Port: port},
	}
}

const defaultTmpl = `{"count":{{len .}},"events":[{{range $i,$e := .}}{{if $i}},{{end}}{"port":{{$e.Port.Port}}}{{end}}]}`

func TestCustomWebhookNotifier_SkipsEmptyEvents(t *testing.T) {
	n, err := NewCustomWebhookNotifier("http://localhost", "", defaultTmpl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := n.Notify(nil); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestCustomWebhookNotifier_PostsPayload(t *testing.T) {
	var got map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, err := NewCustomWebhookNotifier(ts.URL, http.MethodPost, defaultTmpl, map[string]string{"X-Token": "abc"})
	if err != nil {
		t.Fatal(err)
	}
	events := []alert.Event{makeCustomEvent("tcp", "0.0.0.0", 8080)}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["count"] != float64(1) {
		t.Errorf("expected count=1, got %v", got["count"])
	}
}

func TestCustomWebhookNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n, _ := NewCustomWebhookNotifier(ts.URL, "", defaultTmpl, nil)
	events := []alert.Event{makeCustomEvent("tcp", "0.0.0.0", 9090)}
	if err := n.Notify(events); err == nil {
		t.Error("expected error for non-2xx status")
	}
}

func TestCustomWebhookNotifier_InvalidTemplate(t *testing.T) {
	_, err := NewCustomWebhookNotifier("http://localhost", "", `{{.Unclosed`, nil)
	if err == nil {
		t.Error("expected error for invalid template")
	}
}
