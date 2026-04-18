package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// NewZendutyNotifier returns a Notifier that sends alerts to Zenduty via its
// REST incidents API.
func NewZendutyNotifier(apiKey, serviceID, escalationPolicyID string) *ZendutyNotifier {
	return &ZendutyNotifier{
		apiKey:             apiKey,
		serviceID:          serviceID,
		escalationPolicyID: escalationPolicyID,
		client:             &http.Client{},
	}
}

// ZendutyNotifier sends port-change events to Zenduty.
type ZendutyNotifier struct {
	apiKey             string
	serviceID          string
	escalationPolicyID string
	client             *http.Client
}

type zendutyPayload struct {
	Title              string `json:"title"`
	Message            string `json:"message"`
	Service            string `json:"service"`
	EscalationPolicy   string `json:"escalation_policy"`
}

const zendutyEndpoint = "https://www.zenduty.com/api/v1/incidents/"

// Notify implements the Notifier interface.
func (z *ZendutyNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	title := fmt.Sprintf("portwatch: %d port change(s) detected", len(events))
	var body bytes.Buffer
	for _, e := range events {
		fmt.Fprintf(&body, "[%s] %s\n", e.Kind, e.Port)
	}

	p := zendutyPayload{
		Title:            title,
		Message:          body.String(),
		Service:          z.serviceID,
		EscalationPolicy: z.escalationPolicyID,
	}

	b, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("zenduty: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, zendutyEndpoint, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("zenduty: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+z.apiKey)

	resp, err := z.client.Do(req)
	if err != nil {
		return fmt.Errorf("zenduty: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("zenduty: unexpected status %d", resp.StatusCode)
	}
	return nil
}
