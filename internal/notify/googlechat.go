package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// NewGoogleChatNotifier returns a Notifier that posts messages to a Google Chat webhook.
func NewGoogleChatNotifier(webhookURL string) *GoogleChatNotifier {
	return &GoogleChatNotifier{webhookURL: webhookURL, client: &http.Client{}}
}

// GoogleChatNotifier sends alerts to a Google Chat space via incoming webhook.
type GoogleChatNotifier struct {
	webhookURL string
	client     *http.Client
}

type googleChatPayload struct {
	Text string `json:"text"`
}

func (g *GoogleChatNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var buf bytes.Buffer
	for _, e := range events {
		buf.WriteString(fmt.Sprintf("[%s] %s port %d/%s\n", e.Kind, e.Action, e.Port.Number, e.Port.Protocol))
	}

	payload := googleChatPayload{Text: buf.String()}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("googlechat: marshal: %w", err)
	}

	resp, err := g.client.Post(g.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("googlechat: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("googlechat: unexpected status %d", resp.StatusCode)
	}
	return nil
}
