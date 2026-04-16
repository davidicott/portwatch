package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// SlackNotifier sends alert events to a Slack incoming webhook.
type SlackNotifier struct {
	webhookURL string
	client     *http.Client
}

type slackPayload struct {
	Text string `json:"text"`
}

// NewSlackNotifier creates a SlackNotifier that posts to the given Slack webhook URL.
func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{
		webhookURL: webhookURL,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Notify formats events and posts them to the configured Slack webhook.
func (s *SlackNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var buf bytes.Buffer
	for _, e := range events {
		buf.WriteString(fmt.Sprintf("[%s] %s %s/%d\n",
			e.Kind, e.Port.Addr, e.Port.Proto, e.Port.Port))
	}

	payload := slackPayload{Text: buf.String()}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("slack: marshal payload: %w", err)
	}

	resp, err := s.client.Post(s.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("slack: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 256))
		return fmt.Errorf("slack: unexpected status %d: %s", resp.StatusCode, bytes.TrimSpace(respBody))
	}
	return nil
}
