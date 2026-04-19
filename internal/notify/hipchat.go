package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// HipChatNotifier sends alerts to a HipChat room via the v2 API.
type HipChatNotifier struct {
	serverURL string
	roomID    string
	token     string
	client    *http.Client
}

type hipChatPayload struct {
	Message       string `json:"message"`
	MessageFormat string `json:"message_format"`
	Color         string `json:"color"`
	Notify        bool   `json:"notify"`
}

// NewHipChatNotifier creates a notifier that posts to the given HipChat room.
func NewHipChatNotifier(serverURL, roomID, token string) *HipChatNotifier {
	return &HipChatNotifier{
		serverURL: serverURL,
		roomID:    roomID,
		token:     token,
		client:    &http.Client{},
	}
}

// Notify sends port change events to HipChat.
func (h *HipChatNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var buf bytes.Buffer
	for _, e := range events {
		fmt.Fprintf(&buf, "[%s] %s\n", e.Kind, e.Port)
	}

	payload := hipChatPayload{
		Message:       buf.String(),
		MessageFormat: "text",
		Color:         "red",
		Notify:        true,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("hipchat: marshal: %w", err)
	}

	url := fmt.Sprintf("%s/v2/room/%s/notification", h.serverURL, h.roomID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("hipchat: request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.token)

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("hipchat: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("hipchat: unexpected status %d", resp.StatusCode)
	}
	return nil
}
