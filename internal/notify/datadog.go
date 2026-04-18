package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

const datadogEventsURL = "https://api.datadoghq.com/api/v1/events"

// NewDatadogNotifier sends port change events to the Datadog Events API.
func NewDatadogNotifier(apiKey, host string) *DatadogNotifier {
	return &DatadogNotifier{apiKey: apiKey, host: host, client: &http.Client{}}
}

// DatadogNotifier posts events to Datadog.
type DatadogNotifier struct {
	apiKey string
	host   string
	client *http.Client
}

type datadogPayload struct {
	Title string   `json:"title"`
	Text  string   `json:"text"`
	Host  string   `json:"host,omitempty"`
	Tags  []string `json:"tags"`
}

func (d *DatadogNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}
	var body bytes.Buffer
	for _, e := range events {
		payload := datadogPayload{
			Title: fmt.Sprintf("portwatch: port %s", e.Kind),
			Text:  fmt.Sprintf("%s %d/%s", e.Kind, e.Port.Port, e.Port.Proto),
			Host:  d.host,
			Tags:  []string{"portwatch", "port:" + fmt.Sprintf("%d", e.Port.Port), "proto:" + e.Port.Proto},
		}
		body.Reset()
		if err := json.NewEncoder(&body).Encode(payload); err != nil {
			return fmt.Errorf("datadog: encode: %w", err)
		}
		req, err := http.NewRequest(http.MethodPost, datadogEventsURL, &body)
		if err != nil {
			return fmt.Errorf("datadog: request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("DD-API-KEY", d.apiKey)
		resp, err := d.client.Do(req)
		if err != nil {
			return fmt.Errorf("datadog: post: %w", err)
		}
		resp.Body.Close()
		if resp.StatusCode >= 300 {
			return fmt.Errorf("datadog: unexpected status %d", resp.StatusCode)
		}
	}
	return nil
}
