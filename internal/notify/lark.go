package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// larkPayload represents the message body sent to a Lark (Feishu) webhook.
type larkPayload struct {
	MsgType string      `json:"msg_type"`
	Content larkContent `json:"content"`
}

type larkContent struct {
	Text string `json:"text"`
}

// LarkNotifier sends port-change alerts to a Lark (Feishu) incoming webhook.
type LarkNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewLarkNotifier creates a LarkNotifier that posts to the given webhook URL.
func NewLarkNotifier(webhookURL string) *LarkNotifier {
	return &LarkNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

// Notify sends a Lark message for each non-empty batch of events.
func (n *LarkNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("portwatch: %d port change(s) detected\n", len(events)))
	for _, e := range events {
		buf.WriteString(fmt.Sprintf("  [%s] %s\n", e.Kind, e.Port))
	}

	payload := larkPayload{
		MsgType: "text",
		Content: larkContent{Text: buf.String()},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("lark: marshal payload: %w", err)
	}

	resp, err := n.client.Post(n.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("lark: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("lark: unexpected status %d", resp.StatusCode)
	}
	return nil
}
