package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// MatrixNotifier sends alerts to a Matrix room via the Client-Server API.
type MatrixNotifier struct {
	homeserver string
	token      string
	roomID     string
	client     *http.Client
}

// NewMatrixNotifier creates a notifier that posts to the given Matrix room.
func NewMatrixNotifier(homeserver, token, roomID string) *MatrixNotifier {
	return &MatrixNotifier{
		homeserver: homeserver,
		token:      token,
		roomID:     roomID,
		client:     &http.Client{},
	}
}

// Notify sends alert events as a Matrix m.room.message.
func (n *MatrixNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	body := formatMatrixMessage(events)
	payload, err := json.Marshal(map[string]string{
		"msgtype": "m.text",
		"body":    body,
	})
	if err != nil {
		return fmt.Errorf("matrix: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/_matrix/client/v3/rooms/%s/send/m.room.message",
		n.homeserver, n.roomID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("matrix: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+n.token)

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("matrix: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("matrix: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func formatMatrixMessage(events []alert.Event) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "portwatch: %d port change(s)\n", len(events))
	for _, e := range events {
		fmt.Fprintf(&buf, "  [%s] %s/%d\n", e.Type, e.Port.Protocol, e.Port.Port)
	}
	return buf.String()
}
