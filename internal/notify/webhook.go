package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// WebhookNotifier POSTs JSON-encoded events to a URL.
type WebhookNotifier struct {
	url    string
	client *http.Client
}

// NewWebhookNotifier creates a WebhookNotifier targeting url.
func NewWebhookNotifier(url string, timeout time.Duration) *WebhookNotifier {
	if timeout <=t	timeout = 5 * time.Second
	}
	return &WebhookNotifier{
		url:    url,
		client: &http.Client{Timeout: timeout},
	}
}

type webhookPayload struct {
	Events []alert.Event `json:"events"`
}

// Notify POSTs events to the configured webhook URL.
// It returns an error if the server responds with a non-2xx status.
func (w *WebhookNotifier) Notify(ctx context.Context, events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}
	b, err := json.Marshal(webhookPayload{Events: events})
	if err != nil {
		return fmt.Errorf("webhook notifier: marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.url, bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("webhook notifier: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook notifier: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook notifier: unexpected status %d", resp.StatusCode)
	}
	return nil
}
