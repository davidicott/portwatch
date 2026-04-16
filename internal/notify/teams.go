package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// TeamsNotifier sends alerts to a Microsoft Teams channel via an incoming webhook.
type TeamsNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewTeamsNotifier creates a new TeamsNotifier.
func NewTeamsNotifier(webhookURL string) *TeamsNotifier {
	return &TeamsNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

type teamsPayload struct {
	Text string `json:"text"`
}

// Notify sends alert events to Microsoft Teams.
func (t *TeamsNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var buf bytes.Buffer
	for _, e := range events {
		fmt.Fprintf(&buf, "**[%s]** %s %s (pid %d)\n", e.Kind, e.Port.Proto, e.Port.Addr, e.Port.PID)
	}

	payload := teamsPayload{Text: buf.String()}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("teams: marshal payload: %w", err)
	}

	resp, err := t.client.Post(t.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("teams: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("teams: unexpected status %d", resp.StatusCode)
	}
	return nil
}
