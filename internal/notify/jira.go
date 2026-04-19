package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// JiraNotifier creates Jira issues for port change events.
type JiraNotifier struct {
	baseURL   string
	username  string
	token     string
	projectKey string
	issueType  string
	client    *http.Client
}

// NewJiraNotifier returns a JiraNotifier.
func NewJiraNotifier(baseURL, username, token, projectKey, issueType string) *JiraNotifier {
	if issueType == "" {
		issueType = "Task"
	}
	return &JiraNotifier{
		baseURL:    baseURL,
		username:   username,
		token:      token,
		projectKey: projectKey,
		issueType:  issueType,
		client:     &http.Client{},
	}
}

// Notify posts a Jira issue for each event.
func (j *JiraNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	summary := fmt.Sprintf("portwatch: %d port change(s) detected", len(events))
	description := ""
	for _, e := range events {
		description += fmt.Sprintf("- [%s] %s/%d\n", e.Kind, e.Port.Protocol, e.Port.Port)
	}

	payload := map[string]any{
		"fields": map[string]any{
			"project":   map[string]string{"key": j.projectKey},
			"summary":   summary,
			"description": description,
			"issuetype": map[string]string{"name": j.issueType},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("jira: marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, j.baseURL+"/rest/api/2/issue", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("jira: request: %w", err)
	}
	req.SetBasicAuth(j.username, j.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := j.client.Do(req)
	if err != nil {
		return fmt.Errorf("jira: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("jira: unexpected status %d", resp.StatusCode)
	}
	return nil
}
