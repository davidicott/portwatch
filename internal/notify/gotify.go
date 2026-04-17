package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// GotifyNotifier sends alerts to a self-hosted Gotify server.
type GotifyNotifier struct {
	baseURL  string
	token    string
	priority int
	client   *http.Client
}

type gotifyPayload struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
}

// NewGotifyNotifier creates a notifier that posts to a Gotify server.
func NewGotifyNotifier(baseURL, token string, priority int) *GotifyNotifier {
	if priority <= 0 {
		priority = 5
	}
	return &GotifyNotifier{
		baseURL:  baseURL,
		token:    token,
		priority: priority,
		client:   &http.Client{},
	}
}

// Notify sends port change events to Gotify.
func (g *GotifyNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var body bytes.Buffer
	for _, e := range events {
		fmt.Fprintf(&body, "[%s] %s:%d/%s\n", e.Kind, e.Port.Host, e.Port.Port, e.Port.Proto)
	}

	payload := gotifyPayload{
		Title:    fmt.Sprintf("portwatch: %d port change(s)", len(events)),
		Message:  body.String(),
		Priority: g.priority,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("gotify: marshal: %w", err)
	}

	url := fmt.Sprintf("%s/message?token=%s", g.baseURL, g.token)
	resp, err := g.client.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("gotify: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("gotify: unexpected status %d", resp.StatusCode)
	}
	return nil
}
