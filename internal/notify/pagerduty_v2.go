// Package notify provides notifier implementations for portwatch alerts.
package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

const defaultPagerDutyV2URL = "https://events.pagerduty.com/v2/enqueue"

// NewPagerDutyV2Notifier creates a notifier that sends alerts via the
// PagerDuty Events API v2 using a routing key.
func NewPagerDutyV2Notifier(routingKey, endpoint string) *pagerDutyV2Notifier {
	if endpoint == "" {
		endpoint = defaultPagerDutyV2URL
	}
	return &pagerDutyV2Notifier{
		routingKey: routingKey,
		endpoint:   endpoint,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

type pagerDutyV2Notifier struct {
	routingKey string
	endpoint   string
	client     *http.Client
}

type pdV2Payload struct {
	RoutingKey  string        `json:"routing_key"`
	EventAction string        `json:"event_action"`
	Payload     pdV2EventBody `json:"payload"`
}

type pdV2EventBody struct {
	Summary  string `json:"summary"`
	Severity string `json:"severity"`
	Source   string `json:"source"`
	Timestamp string `json:"timestamp"`
}

func (n *pagerDutyV2Notifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}
	summary := fmt.Sprintf("portwatch: %d port change(s) detected", len(events))
	body := pdV2Payload{
		RoutingKey:  n.routingKey,
		EventAction: "trigger",
		Payload: pdV2EventBody{
			Summary:   summary,
			Severity:  "warning",
			Source:    "portwatch",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("pagerduty v2: marshal: %w", err)
	}
	resp, err := n.client.Post(n.endpoint, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("pagerduty v2: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pagerduty v2: unexpected status %d", resp.StatusCode)
	}
	return nil
}
