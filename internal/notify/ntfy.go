package notify

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// NtfyNotifier sends alerts to an ntfy.sh topic.
type NtfyNotifier struct {
	serverURL string
	topic     string
	client    *http.Client
}

// NewNtfyNotifier creates a notifier that publishes to the given ntfy topic.
func NewNtfyNotifier(serverURL, topic string) *NtfyNotifier {
	if serverURL == "" {
		serverURL = "https://ntfy.sh"
	}
	return &NtfyNotifier{
		serverURL: serverURL,
		topic:     topic,
		client:    &http.Client{},
	}
}

// Notify sends each event as a message to the ntfy topic.
func (n *NtfyNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var buf bytes.Buffer
	for _, e := range events {
		fmt.Fprintf(&buf, "%s\n", e.String())
	}

	url := fmt.Sprintf("%s/%s", n.serverURL, n.topic)
	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return fmt.Errorf("ntfy: build request: %w", err)
	}
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Title", fmt.Sprintf("portwatch: %d event(s)", len(events)))

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("ntfy: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ntfy: unexpected status %d", resp.StatusCode)
	}
	return nil
}
