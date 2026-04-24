package config

import (
	"testing"
)

func TestGooglePubSubConfig_Defaults(t *testing.T) {
	defaults := googlePubSubDefaults()

	if defaults.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if defaults.ProjectID != "" {
		t.Errorf("expected empty ProjectID, got %q", defaults.ProjectID)
	}
	if defaults.TopicID != "portwatch-events" {
		t.Errorf("expected TopicID %q, got %q", "portwatch-events", defaults.TopicID)
	}
}

func TestGooglePubSubConfig_Fields(t *testing.T) {
	cfg := &NotifierConfig{
		GooglePubSub: &GooglePubSubConfig{
			Enabled:   true,
			ProjectID: "my-gcp-project",
			TopicID:   "alert-topic",
		},
	}

	if !cfg.GooglePubSub.Enabled {
		t.Error("expected Enabled to be true")
	}
	if cfg.GooglePubSub.ProjectID != "my-gcp-project" {
		t.Errorf("unexpected ProjectID: %q", cfg.GooglePubSub.ProjectID)
	}
	if cfg.GooglePubSub.TopicID != "alert-topic" {
		t.Errorf("unexpected TopicID: %q", cfg.GooglePubSub.TopicID)
	}
}

func TestGooglePubSubConfig_DefaultTopicIDWhenEmpty(t *testing.T) {
	n := &NotifierConfig{
		GooglePubSub: &GooglePubSubConfig{
			Enabled:   true,
			ProjectID: "proj",
			TopicID:   "",
		},
	}
	applyNotifierDefaults("google_pubsub", n)

	if n.GooglePubSub.TopicID != "portwatch-events" {
		t.Errorf("expected default TopicID, got %q", n.GooglePubSub.TopicID)
	}
}
