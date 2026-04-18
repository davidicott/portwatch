package notify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// TwilioNotifier sends SMS alerts via the Twilio API.
type TwilioNotifier struct {
	accountSID string
	authToken  string
	from       string
	to         string
	client     *http.Client
}

// NewTwilioNotifier creates a new TwilioNotifier.
func NewTwilioNotifier(accountSID, authToken, from, to string) *TwilioNotifier {
	return &TwilioNotifier{
		accountSID: accountSID,
		authToken:  authToken,
		from:       from,
		to:         to,
		client:     &http.Client{},
	}
}

// Notify sends an SMS for each alert event.
func (t *TwilioNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	body := formatTwilioMessage(events)
	endpoint := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", t.accountSID)

	form := url.Values{}
	form.Set("From", t.from)
	form.Set("To", t.to)
	form.Set("Body", body)

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("twilio: build request: %w", err)
	}
	req.SetBasicAuth(t.accountSID, t.authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("twilio: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var result map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&result)
		return fmt.Errorf("twilio: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func formatTwilioMessage(events []alert.Event) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("portwatch: %d port change(s)\n", len(events)))
	for _, e := range events {
		sb.WriteString(fmt.Sprintf("[%s] %s\n", e.Kind, e.Port))
	}
	return sb.String()
}
