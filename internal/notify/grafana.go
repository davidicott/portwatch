package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/patrickdappollonio/portwatch/internal/alert"
)

// NewGrafanaNotifier creates a notifier that sends alerts to a Grafana webhook (e.g. Grafana OnCall).
func NewGrafanaNotifier(webhookURL, title string, client *http.Client) *GrafanaNotifier {
	if client == nil {
		client = &http.Client{}
	}
	return &GrafanaNotifier{url: webhookURL, title: title, client: client}
}

// GrafanaNotifier sends port change events to a Grafana webhook endpoint.
type GrafanaNotifier struct {
	url    string
	title  string
	client *http.Client
}

type grafanaPayload struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	State   string `json:"state"`
}

// Notify posts events to the configured Grafana webhook URL.
func (g *GrafanaNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var buf bytes.Buffer
	for _, e := range events {
		fmt.Fprintf(&buf, "[%s] %s\n", e.Kind, e.Port)
	}

	title := g.title
	if title == "" {
		title = "portwatch alert"
	}

	p := grafanaPayload{
		Title:   title,
		Message: buf.String(),
		State:   "alerting",
	}

	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("grafana: marshal payload: %w", err)
	}

	resp, err := g.client.Post(g.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("grafana: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("grafana: unexpected status %d", resp.StatusCode)
	}
	return nil
}
