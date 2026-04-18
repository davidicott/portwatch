package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeDDEvent(kind, proto string, port uint16) alert.Event {
	return alert.Event{Kind: alert.EventKind(kind), Port: scanner.Port{Port: port, Proto: proto}}
}

func TestDatadogNotifier_SkipsEmptyEvents(t *testing.T) {
	n := NewDatadogNotifier("key", "host")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestDatadogNotifier_PostsPayload(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("DD-API-KEY") == "" {
			t.Error("missing DD-API-KEY header")
		}
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	n := NewDatadogNotifier("testkey", "myhost")
	n.client = ts.Client()
	// Override URL via a simple wrapper trick — patch via direct field for test.
	// Since datadogEventsURL is a const, we use a real server and swap the client.
	// Instead, we test via a real request path by pointing to the test server.
	// Re-implement with overrideable URL field for testability:
	oldURL := datadogEventsURL
	_ = oldURL // const; test server used indirectly via client transport

	// Use a notifier with overridden URL by embedding test server URL.
	n2 := &DatadogNotifier{apiKey: "testkey", host: "myhost", client: &http.Client{
		Transport: &urlOverrideTransport{base: http.DefaultTransport, target: ts.URL},
	}}
	if err := n2.Notify([]alert.Event{makeDDEvent("opened", "tcp", 8080)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["title"] == nil {
		t.Error("expected title in payload")
	}
}

func TestDatadogNotifier_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()
	n := &DatadogNotifier{apiKey: "k", host: "h", client: &http.Client{
		Transport: &urlOverrideTransport{base: http.DefaultTransport, target: ts.URL},
	}}
	if err := n.Notify([]alert.Event{makeDDEvent("closed", "udp", 53)}); err == nil {
		t.Fatal("expected error on non-2xx status")
	}
}

type urlOverrideTransport struct {
	base   http.RoundTripper
	target string
}

func (u *urlOverrideTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := req.Clone(req.Context())
	req2.URL.Host = req.URL.Host
	parsed, _ := http.NewRequest(req.Method, u.target, req.Body)
	parsed.Header = req.Header
	return u.base.RoundTrip(parsed)
}
