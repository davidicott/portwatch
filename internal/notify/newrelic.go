package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

const defaultNewRelicURL = "https://log-api.newrelic.com/log/v1"

// NewRelicNotifier sends port change events to New Relic Log API.
type NewRelicNotifier struct {
	apiKey  string
	url     string
	client  *http.Client
}

// NewNewRelicNotifier creates a notifier that forwards events to New Relic.
func NewNewRelicNotifier(apiKey, url string) *NewRelicNotifier {
	if url == "" {
		url = defaultNewRelicURL
	}
	return &NewRelicNotifier{
		apiKey: apiKey,
		url:    url,
		client: &http.Client{},
	}
}

func (n *NewRelicNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	type logEntry struct {
		Message  string `json:"message"`
		Kind     string `json:"kind"`
		Protocol string `json:"protocol"`
		Port     uint16 `json:"port"`
	}

	entries := make([]logEntry, 0, len(events))
	for _, e := range events {
		entries = append(entries, logEntry{
			Message:  fmt.Sprintf("portwatch: port %s %s", e.Port.String(), e.Kind),
			Kind:     string(e.Kind),
			Protocol: e.Port.Protocol,
			Port:     e.Port.Port,
		})
	}

	body, err := json.Marshal(entries)
	if err != nil {
		return fmt.Errorf("newrelic: marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, n.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("newrelic: request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", n.apiKey)

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("newrelic: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("newrelic: unexpected status %d", resp.StatusCode)
	}
	return nil
}
