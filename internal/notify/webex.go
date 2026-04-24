package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

const defaultWebexAPIURL = "https://webexapis.com/v1/messages"

// WebexNotifier sends port change alerts to a Cisco Webex space.
type WebexNotifier struct {
	token  string
	roomID string
	apiURL string
	client *http.Client
}

// NewWebexNotifier creates a WebexNotifier that posts to the given Webex room.
func NewWebexNotifier(token, roomID string) *WebexNotifier {
	return &WebexNotifier{
		token:  token,
		roomID: roomID,
		apiURL: defaultWebexAPIURL,
		client: &http.Client{},
	}
}

// Notify sends a Webex message for each port change event.
func (n *WebexNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	text := formatWebexMessage(events)

	payload := map[string]string{
		"roomId": n.roomID,
		"text":   text,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webex: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, n.apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webex: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+n.token)

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("webex: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webex: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func formatWebexMessage(events []alert.Event) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("portwatch: %d port change(s) detected\n", len(events)))
	for _, e := range events {
		buf.WriteString(fmt.Sprintf("  [%s] %s/%d\n", e.Kind, e.Port.Proto, e.Port.Port))
	}
	return buf.String()
}
