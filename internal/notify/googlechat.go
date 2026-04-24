package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// GoogleChatNotifier sends port change alerts to a Google Chat webhook.
type GoogleChatNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewGoogleChatNotifier creates a GoogleChatNotifier that posts to the given webhook URL.
func NewGoogleChatNotifier(webhookURL string) *GoogleChatNotifier {
	return &GoogleChatNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

// Notify sends a Google Chat card message for each alert event.
func (n *GoogleChatNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var text string
	for _, e := range events {
		text += fmt.Sprintf("*[%s]* %s/%d\n", e.Kind, e.Port.Protocol, e.Port.Port)
	}

	payload := map[string]interface{}{
		"text": fmt.Sprintf("🔍 *portwatch alert* — %d event(s)\n%s", len(events), text),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("googlechat: marshal payload: %w", err)
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
