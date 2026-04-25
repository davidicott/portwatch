package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

const defaultTelegramChannelAPI = "https://api.telegram.org"

// TelegramChannelNotifier sends alerts to a Telegram channel via bot API.
type TelegramChannelNotifier struct {
	token      string
	channelID  string
	baseURL    string
	httpClient *http.Client
}

type telegramChannelPayload struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

// NewTelegramChannelNotifier creates a notifier that posts to a Telegram channel.
func NewTelegramChannelNotifier(token, channelID string) *TelegramChannelNotifier {
	return &TelegramChannelNotifier{
		token:      token,
		channelID:  channelID,
		baseURL:    defaultTelegramChannelAPI,
		httpClient: &http.Client{},
	}
}

// Notify sends port change events to the configured Telegram channel.
func (n *TelegramChannelNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	text := formatTelegramMessage(events)
	payload := telegramChannelPayload{
		ChatID:    n.channelID,
		Text:      text,
		ParseMode: "Markdown",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("telegram channel: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", n.baseURL, n.token)
	resp, err := n.httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("telegram channel: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("telegram channel: unexpected status %d", resp.StatusCode)
	}
	return nil
}
