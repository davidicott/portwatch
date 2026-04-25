package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/patrickdappollonio/portwatch/internal/alert"
)

const defaultStatuspageEndpoint = "https://api.statuspage.io/v1"

// StatuspageNotifier posts port-change events as Statuspage.io incidents.
type StatuspageNotifier struct {
	apiKey     string
	pageID     string
	componentID string
	endpoint   string
	client     *http.Client
}

// NewStatuspageNotifier creates a new Statuspage.io notifier.
func NewStatuspageNotifier(apiKey, pageID, componentID, endpoint string) *StatuspageNotifier {
	if endpoint == "" {
		endpoint = defaultStatuspageEndpoint
	}
	return &StatuspageNotifier{
		apiKey:      apiKey,
		pageID:      pageID,
		componentID: componentID,
		endpoint:    endpoint,
		client:      &http.Client{Timeout: 10 * time.Second},
	}
}

// Notify sends an incident to Statuspage.io for each batch of events.
func (n *StatuspageNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	body := map[string]interface{}{
		"incident": map[string]interface{}{
			"name":                fmt.Sprintf("portwatch: %d port change(s) detected", len(events)),
			"status":             "investigating",
			"impact_override":    "minor",
			"body":               buildStatuspageBody(events),
			"component_ids":      []string{n.componentID},
			"deliver_notifications": true,
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("statuspage: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/pages/%s/incidents", n.endpoint, n.pageID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("statuspage: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "OAuth "+n.apiKey)

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("statuspage: post incident: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("statuspage: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func buildStatuspageBody(events []alert.Event) string {
	var buf bytes.Buffer
	for _, e := range events {
		fmt.Fprintf(&buf, "[%s] %s/%d\n", e.Kind, e.Port.Protocol, e.Port.Port)
	}
	return buf.String()
}
