package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

const defaultOpsGenieURL = "https://api.opsgenie.com/v2/alerts"

// OpsGenieNotifier sends alerts to OpsGenie.
type OpsGenieNotifier struct {
	apiKey  string
	url     string
	client  *http.Client
}

// NewOpsGenieNotifier creates an OpsGenieNotifier with the given API key.
func NewOpsGenieNotifier(apiKey string) *OpsGenieNotifier {
	return &OpsGenieNotifier{
		apiKey: apiKey,
		url:    defaultOpsGenieURL,
		client: &http.Client{},
	}
}

type opsGeniePayload struct {
	Message     string `json:"message"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

// Notify sends each event as a separate OpsGenie alert.
func (o *OpsGenieNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}
	for _, e := range events {
		payload := opsGeniePayload{
			Message:     fmt.Sprintf("portwatch: %s %s/%d", e.Kind, e.Port.Protocol, e.Port.Number),
			Description: fmt.Sprintf("Port %s/%d on %s", e.Port.Protocol, e.Port.Number, e.Port.Address),
			Priority:    "P3",
		}
		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("opsgenie: marshal: %w", err)
		}
		req, err := http.NewRequest(http.MethodPost, o.url, bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("opsgenie: request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "GenieKey "+o.apiKey)
		resp, err := o.client.Do(req)
		if err != nil {
			return fmt.Errorf("opsgenie: send: %w", err)
		}
		resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("opsgenie: unexpected status %d", resp.StatusCode)
		}
	}
	return nil
}
