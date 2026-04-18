package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// SplunkNotifier sends events to a Splunk HTTP Event Collector (HEC) endpoint.
type SplunkNotifier struct {
	url    string
	token  string
	source string
	client *http.Client
}

type splunkEvent struct {
	Event  map[string]string `json:"event"`
	Source string            `json:"source,omitempty"`
}

// NewSplunkNotifier creates a notifier that posts to a Splunk HEC endpoint.
func NewSplunkNotifier(url, token, source string) *SplunkNotifier {
	return &SplunkNotifier{
		url:    url,
		token:  token,
		source: source,
		client: &http.Client{},
	}
}

// Notify sends each alert event to Splunk as a separate HEC event.
func (s *SplunkNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	for _, e := range events {
		payload := splunkEvent{
			Source: s.source,
			Event: map[string]string{
				"kind":     e.Kind,
				"protocol": e.Port.Protocol,
				"port":     fmt.Sprintf("%d", e.Port.Number),
				"process":  e.Port.Process,
			},
		}
		if err := enc.Encode(payload); err != nil {
			return fmt.Errorf("splunk: encode: %w", err)
		}
	}

	req, err := http.NewRequest(http.MethodPost, s.url, &buf)
	if err != nil {
		return fmt.Errorf("splunk: request: %w", err)
	}
	req.Header.Set("Authorization", "Splunk "+s.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("splunk: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("splunk: unexpected status %d", resp.StatusCode)
	}
	return nil
}
