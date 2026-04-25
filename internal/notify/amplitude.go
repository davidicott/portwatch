package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

const defaultAmplitudeEndpoint = "https://api2.amplitude.com/2/httpapi"

// AmplitudeNotifier sends port change events to Amplitude as track events.
type AmplitudeNotifier struct {
	apiKey   string
	endpoint string
	client   *http.Client
}

// NewAmplitudeNotifier creates a new AmplitudeNotifier with the given API key.
func NewAmplitudeNotifier(apiKey, endpoint string) *AmplitudeNotifier {
	if endpoint == "" {
		endpoint = defaultAmplitudeEndpoint
	}
	return &AmplitudeNotifier{
		apiKey:   apiKey,
		endpoint: endpoint,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

type amplitudeEvent struct {
	UserID    string                 `json:"user_id"`
	EventType string                 `json:"event_type"`
	EventProp map[string]interface{} `json:"event_properties"`
	Time      int64                  `json:"time"`
}

type amplitudePayload struct {
	APIKey string           `json:"api_key"`
	Events []amplitudeEvent `json:"events"`
}

// Notify sends all port events to Amplitude.
func (n *AmplitudeNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	amEvents := make([]amplitudeEvent, 0, len(events))
	for _, e := range events {
		amEvents = append(amEvents, amplitudeEvent{
			UserID:    "portwatch",
			EventType: fmt.Sprintf("port_%s", e.Kind),
			Time:      e.Timestamp.UnixMilli(),
			EventProp: map[string]interface{}{
				"port":     e.Port.Number,
				"protocol": e.Port.Protocol,
				"address":  e.Port.Address,
			},
		})
	}

	body, err := json.Marshal(amplitudePayload{APIKey: n.apiKey, Events: amEvents})
	if err != nil {
		return fmt.Errorf("amplitude: marshal payload: %w", err)
	}

	resp, err := n.client.Post(n.endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("amplitude: post events: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("amplitude: unexpected status %d", resp.StatusCode)
	}
	return nil
}
