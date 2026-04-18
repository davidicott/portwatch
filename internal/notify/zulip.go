package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// ZulipNotifier sends alerts to a Zulip stream via the Zulip REST API.
type ZulipNotifier struct {
	baseURL  string
	email    string
	apiKey   string
	stream   string
	topic    string
	client   *http.Client
}

// NewZulipNotifier creates a ZulipNotifier.
func NewZulipNotifier(baseURL, email, apiKey, stream, topic string) *ZulipNotifier {
	return &ZulipNotifier{
		baseURL: baseURL,
		email:   email,
		apiKey:  apiKey,
		stream:  stream,
		topic:   topic,
		client:  &http.Client{},
	}
}

func (z *ZulipNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var buf bytes.Buffer
	for _, e := range events {
		fmt.Fprintf(&buf, "**%s** — %s:%d (%s)\n", e.Kind, e.Port.Address, e.Port.Port, e.Port.Protocol)
	}

	payload := map[string]string{
		"type":    "stream",
		"to":      z.stream,
		"topic":   z.topic,
		"content": buf.String(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, z.baseURL+"/api/v1/messages", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.SetBasicAuth(z.email, z.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := z.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("zulip: unexpected status %d", resp.StatusCode)
	}
	return nil
}
