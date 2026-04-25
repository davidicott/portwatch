package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// NewChimeNotifier returns a Notifier that posts messages to an AWS Chime
// incoming webhook URL.
func NewChimeNotifier(webhookURL string) *ChimeNotifier {
	return &ChimeNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

// ChimeNotifier sends port-change events to an AWS Chime room via webhook.
type ChimeNotifier struct {
	webhookURL string
	client     *http.Client
}

type chimePayload struct {
	Content string `json:"Content"`
}

// Notify implements the Notifier interface.
func (n *ChimeNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var buf bytes.Buffer
	for _, e := range events {
		fmt.Fprintf(&buf, "[portwatch] %s port %d/%s\n", e.Kind, e.Port.Port, e.Port.Proto)
	}

	payload := chimePayload{Content: buf.String()}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("chime: marshal payload: %w", err)
	}

	resp, err := n.client.Post(n.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("chime: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("chime: unexpected status %d", resp.StatusCode)
	}
	return nil
}
