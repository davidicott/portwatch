package config

// GooglePubSubConfig holds configuration for the Google Cloud Pub/Sub notifier.
type GooglePubSubConfig struct {
	Enabled   bool   `yaml:"enabled"`
	ProjectID string `yaml:"project_id"`
	TopicID   string `yaml:"topic_id"`
}

func googlePubSubDefaults() GooglePubSubConfig {
	return GooglePubSubConfig{
		Enabled:   false,
		ProjectID: "",
		TopicID:   "portwatch-events",
	}
}

func init() {
	registerNotifierDefaults("google_pubsub", func(n *NotifierConfig) {
		if n.GooglePubSub == nil {
			defaults := googlePubSubDefaults()
			n.GooglePubSub = &defaults
		} else {
			if n.GooglePubSub.TopicID == "" {
				n.GooglePubSub.TopicID = googlePubSubDefaults().TopicID
			}
		}
	})
}
