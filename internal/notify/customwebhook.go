package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/user/portwatch/internal/alert"
)

// CustomWebhookNotifier sends alerts to an arbitrary HTTP endpoint
// with a user-defined JSON body template.
type CustomWebhookNotifier struct {
	url      string
	method   string
	headers  map[string]string
	tmpl     *template.Template
	client   *http.Client
}

// NewCustomWebhookNotifier creates a notifier that POSTs a templated JSON body.
// bodyTemplate is a Go text/template producing valid JSON; it receives []alert.Event.
func NewCustomWebhookNotifier(url, method, bodyTemplate string, headers map[string]string) (*CustomWebhookNotifier, error) {
	if method == "" {
		method = http.MethodPost
	}
	tmpl, err := template.New("body").Parse(bodyTemplate)
	if err != nil {
		return nil, fmt.Errorf("customwebhook: invalid template: %w", err)
	}
	return &CustomWebhookNotifier{
		url:     url,
		method:  method,
		headers: headers,
		tmpl:    tmpl,
		client:  &http.Client{},
	}, nil
}

// Notify renders the template with the given events and sends the HTTP request.
func (n *CustomWebhookNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var buf bytes.Buffer
	if err := n.tmpl.Execute(&buf, events); err != nil {
		return fmt.Errorf("customwebhook: template render: %w", err)
	}

	// Validate the rendered output is JSON.
	if !json.Valid(buf.Bytes()) {
		return fmt.Errorf("customwebhook: rendered body is not valid JSON")
	}

	req, err := http.NewRequest(n.method, n.url, &buf)
	if err != nil {
		return fmt.Errorf("customwebhook: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range n.headers {
		req.Header.Set(k, v)
	}

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("customwebhook: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("customwebhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}
