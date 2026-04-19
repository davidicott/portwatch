package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

const flowdockDefaultAPI = "https://api.flowdock.com/messages"

// FlowdockNotifier sends alerts to a Flowdock flow.
type FlowdockNotifier struct {
	token   string
	flowID  string
	apiURL  string
	client  *http.Client
}

// NewFlowdockNotifier creates a new FlowdockNotifier.
func NewFlowdockNotifier(token, flowID string) *FlowdockNotifier {
	return &FlowdockNotifier{
		token:  token,
		flowID: flowID,
		apiURL: flowdockDefaultAPI,
		client: &http.Client{},
	}
}

func (f *FlowdockNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	content := formatFlowdockMessage(events)

	payload := map[string]interface{}{
		"flow_token": f.token,
		"event":      "message",
		"content":    content,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("flowdock: marshal payload: %w", err)
	}

	resp, err := f.client.Post(f.apiURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("flowdock: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("flowdock: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func formatFlowdockMessage(events []alert.Event) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("portwatch: %d port change(s) detected\n", len(events)))
	for _, e := range events {
		buf.WriteString(fmt.Sprintf("  [%s] %s/%d\n", e.Kind, e.Port.Protocol, e.Port.Port))
	}
	return buf.String()
}
