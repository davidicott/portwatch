package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// SignalRNotifier sends alerts to a SignalR-compatible HTTP endpoint (e.g. Azure SignalR REST API).
type SignalRNotifier struct {
	endpoint string
	hub      string
	apiKey   string
	client   *http.Client
}

type signalRPayload struct {
	Target    string   `json:"target"`
	Arguments []string `json:"arguments"`
}

// NewSignalRNotifier creates a new SignalRNotifier.
func NewSignalRNotifier(endpoint, hub, apiKey string) *SignalRNotifier {
	return &SignalRNotifier{
		endpoint: endpoint,
		hub:      hub,
		apiKey:   apiKey,
		client:   &http.Client{},
	}
}

// Notify sends port change events to the SignalR hub.
func (s *SignalRNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var lines []string
	for _, e := range events {
		lines = append(lines, fmt.Sprintf("[%s] %s:%d/%s", e.Type, e.Host, e.Port, e.Protocol))
	}

	body, err := json.Marshal(signalRPayload{
		Target:    "portwatch",
		Arguments: lines,
	})
	if err != nil {
		return fmt.Errorf("signalr: marshal: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/hubs/%s", s.endpoint, s.hub)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("signalr: request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("signalr: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("signalr: unexpected status %d", resp.StatusCode)
	}
	return nil
}
