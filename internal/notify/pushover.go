package notify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/user/portwatch/internal/alert"
)

const pushoverAPI = "https://api.pushover.net/1/messages.json"

// PushoverNotifier sends alerts via the Pushover API.
type PushoverNotifier struct {
	token   string
	userKey string
	client  *http.Client
}

// NewPushoverNotifier creates a new PushoverNotifier.
func NewPushoverNotifier(token, userKey string) *PushoverNotifier {
	return &PushoverNotifier{
		token:   token,
		userKey: userKey,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// Notify sends port change events to Pushover.
func (p *PushoverNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var sb strings.Builder
	for _, e := range events {
		sb.WriteString(fmt.Sprintf("%s: %s\n", e.Kind, e.Port))
	}

	payload := map[string]string{
		"token":   p.token,
		"user":    p.userKey,
		"title":   fmt.Sprintf("portwatch: %d port change(s)", len(events)),
		"message": sb.String(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("pushover: marshal: %w", err)
	}

	resp, err := p.client.Post(pushoverAPI, "application/json", strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("pushover: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pushover: unexpected status %d", resp.StatusCode)
	}
	return nil
}
