package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// MattermostNotifier sends alerts to a Mattermost incoming webhook.
type MattermostNotifier struct {
	webhookURL string
	channel    string
	client     *http.Client
}

type mattermostPayload struct {
	Channel string `json:"channel,omitempty"`
	Text    string `json:"text"`
}

// NewMattermostNotifier creates a notifier that posts to Mattermost.
func NewMattermostNotifier(webhookURL, channel string) *MattermostNotifier {
	return &MattermostNotifier{
		webhookURL: webhookURL,
		channel:    channel,
		client:     &http.Client{},
	}
}

// Notify sends port change events to Mattermost.
func (m *MattermostNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("**portwatch alert** — %d change(s) detected:\n", len(events)))
	for _, e := range events {
		buf.WriteString(fmt.Sprintf("- `%s` %s\n", e.Port, e.Kind))
	}

	payload := mattermostPayload{
		Channel: m.channel,
		Text:    buf.String(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("mattermost: marshal: %w", err)
	}

	resp, err := m.client.Post(m.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("mattermost: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("mattermost: unexpected status %d", resp.StatusCode)
	}
	return nil
}
