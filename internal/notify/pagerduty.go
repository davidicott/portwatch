package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

const pagerDutyEventURL = "https://events.pagerduty.com/v2/enqueue"

type pagerDutyPayload struct {
	RoutingKey  string         `json:"routing_key"`
	EventAction string         `json:"event_action"`
	Payload     pdInnerPayload `json:"payload"`
}

type pdInnerPayload struct {
	Summary   string `json:"summary"`
	Source    string `json:"source"`
	Severity  string `json:"severity"`
	Timestamp string `json:"timestamp"`
}

// PagerDutyNotifier sends alerts to PagerDuty via the Events API v2.
type PagerDutyNotifier struct {
	routingKey string
	source     string
	client     *http.Client
	url        string
}

// NewPagerDutyNotifier creates a PagerDutyNotifier with the given routing key.
func NewPagerDutyNotifier(routingKey, source string) *PagerDutyNotifier {
	return &PagerDutyNotifier{
		routingKey: routingKey,
		source:     source,
		client:     &http.Client{Timeout: 10 * time.Second},
		url:        pagerDutyEventURL,
	}
}

// Notify sends each event as a PagerDuty trigger.
func (p *PagerDutyNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}
	for _, e := range events {
		payload := pagerDutyPayload{
			RoutingKey:  p.routingKey,
			EventAction: "trigger",
			Payload: pdInnerPayload{
				Summary:   e.String(),
				Source:    p.source,
				Severity:  "warning",
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			},
		}
		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("pagerduty: marshal: %w", err)
		}
		resp, err := p.client.Post(p.url, "application/json", bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("pagerduty: post: %w", err)
		}
		resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("pagerduty: unexpected status %d", resp.StatusCode)
		}
	}
	return nil
}
