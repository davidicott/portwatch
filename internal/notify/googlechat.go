package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// GoogleChatNotifier sends alerts to a Google Chat webhook.
type GoogleChatNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewGoogleChatNotifier creates a new GoogleChatNotifier.
func NewGoogleChatNotifier(webhookURL string) *GoogleChatNotifier {
	return &GoogleChatNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

type googleChatPayload struct {
	Text string `json:"text"`
}

// Notify sends events to the configured Google Chat webhook.
func (n *GoogleChatNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var buf bytes.Buffer
	for _, e := range events {
		buf.WriteString(fmt.Sprintf("[%s] %s:%d/%s\n", e.Kind, e.Port.Host, e.Port.Port, e.Port.Proto))
	}

	payload := googleChatPayload{Text: buf.String()}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("googlechat: marshal: %w", err)
	}

	resp, err := n.client.Post(n.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("googlechat: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("googlechat: unexpected status %d", resp.StatusCode)
	}
	return nil
}
