package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

const telegramAPIBase = "https://api.telegram.org/bot"

// TelegramNotifier sends alerts to a Telegram chat via the Bot API.
type TelegramNotifier struct {
	token  string
	chatID string
	client *http.Client
}

// NewTelegramNotifier creates a notifier that posts messages to a Telegram chat.
func NewTelegramNotifier(token, chatID string) *TelegramNotifier {
	return &TelegramNotifier{
		token:  token,
		chatID: chatID,
		client: &http.Client{},
	}
}

// Notify sends port change events as a Telegram message.
func (t *TelegramNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	text := formatTelegramMessage(events)

	payload := map[string]string{
		"chat_id": t.chatID,
		"text":    text,
		"parse_mode": "Markdown",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("telegram: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s%s/sendMessage", telegramAPIBase, t.token)
	resp, err := t.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("telegram: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("telegram: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func formatTelegramMessage(events []alert.Event) string {
	var buf bytes.Buffer
	buf.WriteString("*portwatch alert*\n")
	for _, e := range events {
		buf.WriteString(fmt.Sprintf("`%s` — %s\n", e.Port, e.Kind))
	}
	return buf.String()
}
