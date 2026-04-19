package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

const defaultClickUpAPI = "https://api.clickup.com/api/v2"

// NewClickUpNotifier creates a notifier that creates ClickUp tasks for port events.
func NewClickUpNotifier(apiToken, listID string) *ClickUpNotifier {
	return &ClickUpNotifier{
		apiToken: apiToken,
		listID:   listID,
		apiBase:  defaultClickUpAPI,
		client:   &http.Client{},
	}
}

// ClickUpNotifier posts port change events as tasks to a ClickUp list.
type ClickUpNotifier struct {
	apiToken string
	listID   string
	apiBase  string
	client   *http.Client
}

type clickUpTask struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Notify sends each event as a separate ClickUp task.
func (n *ClickUpNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}
	for _, e := range events {
		task := clickUpTask{
			Name:        fmt.Sprintf("[portwatch] %s port %d/%s", e.Kind, e.Port.Port, e.Port.Proto),
			Description: fmt.Sprintf("Port %d (%s) was %s on %s.", e.Port.Port, e.Port.Proto, e.Kind, e.At.Format("2006-01-02 15:04:05")),
		}
		body, err := json.Marshal(task)
		if err != nil {
			return fmt.Errorf("clickup: marshal: %w", err)
		}
		url := fmt.Sprintf("%s/list/%s/task", n.apiBase, n.listID)
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("clickup: request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", n.apiToken)
		resp, err := n.client.Do(req)
		if err != nil {
			return fmt.Errorf("clickup: send: %w", err)
		}
		resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("clickup: unexpected status %d", resp.StatusCode)
		}
	}
	return nil
}
