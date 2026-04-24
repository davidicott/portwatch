package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

const chatworkAPIBase = "https://api.chatwork.com/v2"

// ChatworkNotifier sends alerts to a Chatwork room via the REST API.
type ChatworkNotifier struct {
	token  string
	roomID string
	client *http.Client
}

// NewChatworkNotifier creates a notifier that posts messages to the given
// Chatwork room. token is a personal API token; roomID is the numeric room ID.
func NewChatworkNotifier(token, roomID string) *ChatworkNotifier {
	return &ChatworkNotifier{
		token:  token,
		roomID: roomID,
		client: &http.Client{},
	}
}

// Notify sends each alert event as a message to the configured Chatwork room.
func (n *ChatworkNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	body := formatChatworkMessage(events)

	payload, _ := json.Marshal(map[string]string{"body": body})

	url := fmt.Sprintf("%s/rooms/%s/messages", chatworkAPIBase, n.roomID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("chatwork: build request: %w", err)
	}
	req.Header.Set("X-ChatWorkToken", n.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("chatwork: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("chatwork: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func formatChatworkMessage(events []alert.Event) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("[info][title]portwatch – %d port change(s)[/title]\n", len(events)))
	for _, e := range events {
		buf.WriteString(fmt.Sprintf("%s %s/%d\n", e.Kind, e.Port.Protocol, e.Port.Port))
	}
	buf.WriteString("[/info]")
	return buf.String()
}
