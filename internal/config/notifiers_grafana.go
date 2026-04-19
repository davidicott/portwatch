package config

func init() {
	registerNotifierDefaults("grafana", grafanaDefaults)
}

func grafanaDefaults(n *NotifierConfig) {
	if n.Title == "" {
		n.Title = "portwatch alert"
	}
}

// grafana-specific fields are stored in the shared NotifierConfig.
// Required fields: WebhookURL.
// Optional fields: Title (defaults to "portwatch alert").
