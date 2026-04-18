package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/patrickdappollonio/portwatch/internal/alert"
)

// BearyChat incoming webhook payload.
type bearyChatPayload struct {
	Text        string `json:"text"`
	Notification string `json:"notification,omitempty"`
}

// BearyChat notifier sends alerts to a BearyChat incoming webhook.
type bearyChat struct {
	webhookURL string
	client     *http.Client
}

// NewBearyChat returns a Notifier that posts to a BearyChat webhook URL.
func NewBearyChat(webhookURL string) Notifier {
	return &bearyChat{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

func (b *bearyChat) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var buf bytes.Buffer
	for _, e := range events {
		buf.WriteString(fmt.Sprintf("[%s] %s\n", e.Kind, e.Message))
	}

	payload := bearyChatPayload{
		Text:        buf.String(),
		Notification: fmt.Sprintf("%d port event(s) detected", len(events)),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("bearychat: marshal payload: %w", err)
	}

	resp, err := b.client.Post(b.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("bearychat: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("bearychat: unexpected status %d", resp.StatusCode)
	}

	return nil
}
