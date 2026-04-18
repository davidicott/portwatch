package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// RocketChatNotifier sends alerts to a Rocket.Chat webhook.
type RocketChatNotifier struct {
	webhookURL string
	client     *http.Client
}

type rocketChatPayload struct {
	Text string `json:"text"`
}

// NewRocketChatNotifier creates a new RocketChatNotifier.
func NewRocketChatNotifier(webhookURL string) *RocketChatNotifier {
	return &RocketChatNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

// Notify sends port change events to Rocket.Chat.
func (r *RocketChatNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var buf bytes.Buffer
	for _, e := range events {
		buf.WriteString(fmt.Sprintf("[%s] %s\n", e.Kind, e.Port))
	}

	payload := rocketChatPayload{Text: buf.String()}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("rocketchat: marshal payload: %w", err)
	}

	resp, err := r.client.Post(r.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("rocketchat: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("rocketchat: unexpected status %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}
