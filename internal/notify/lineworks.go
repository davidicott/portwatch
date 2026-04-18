package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// LineWorksNotifier sends alerts to a LINE WORKS bot webhook.
type LineWorksNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewLineWorksNotifier creates a new LineWorksNotifier.
func NewLineWorksNotifier(webhookURL string) *LineWorksNotifier {
	return &LineWorksNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

func (n *LineWorksNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("portwatch: %d port change(s) detected\n", len(events)))
	for _, e := range events {
		buf.WriteString(fmt.Sprintf("  [%s] %s:%d (%s)\n", e.Type, e.Host, e.Port, e.Protocol))
	}

	payload := map[string]string{"content": buf.String()}
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("lineworks: marshal: %w", err)
	}

	resp, err := n.client.Post(n.webhookURL, "application/json", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("lineworks: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("lineworks: unexpected status %d", resp.StatusCode)
	}
	return nil
}
