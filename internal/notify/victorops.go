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

// NewVictorOpsNotifier returns a Notifier that sends alerts to a VictorOps REST endpoint.
func NewVictorOpsNotifier(routingKey, restEndpointURL string) *VictorOpsNotifier {
	return &VictorOpsNotifier{
		routingKey:      routingKey,
		restEndpointURL: restEndpointURL,
		client:          &http.Client{Timeout: 10 * time.Second},
	}
}

// VictorOpsNotifier sends port-change events to VictorOps.
type VictorOpsNotifier struct {
	routingKey      string
	restEndpointURL string
	client          *http.Client
}

type victorOpsPayload struct {
	MessageType       string `json:"message_type"`
	EntityID          string `json:"entity_id"`
	EntityDisplayName string `json:"entity_display_name"`
	StateMessage      string `json:"state_message"`
	RoutingKey        string `json:"routing_key"`
}

// Notify sends one VictorOps alert per event.
func (v *VictorOpsNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}
	for _, e := range events {
		p := victorOpsPayload{
			MessageType:       "CRITICAL",
			EntityID:          fmt.Sprintf("portwatch-%s-%d", e.Port.Protocol, e.Port.Number),
			EntityDisplayName: fmt.Sprintf("Port %s/%d %s", e.Port.Protocol, e.Port.Number, e.Kind),
			StateMessage:      fmt.Sprintf("portwatch detected port %s/%d was %s", e.Port.Protocol, e.Port.Number, e.Kind),
			RoutingKey:        v.routingKey,
		}
		body, err := json.Marshal(p)
		if err != nil {
			return fmt.Errorf("victorops: marshal: %w", err)
		}
		url := fmt.Sprintf("%s/%s", v.restEndpointURL, v.routingKey)
		resp, err := v.client.Post(url, "application/json", bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("victorops: post: %w", err)
		}
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("victorops: unexpected status %d: %s", resp.StatusCode, string(respBody))
		}
	}
	return nil
}
