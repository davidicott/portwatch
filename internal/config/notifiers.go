package config

// NotifierConfig holds configuration for all supported notification channels.
// Each field is optional; a nil/zero value means the notifier is disabled.
type NotifierConfig struct {
	Stdout    bool              `yaml:"stdout"`
	Webhook   *WebhookConfig    `yaml:"webhook,omitempty"`
	Slack     *SlackConfig      `yaml:"slack,omitempty"`
	Email     *EmailConfig      `yaml:"email,omitempty"`
	PagerDuty *PagerDutyConfig  `yaml:"pagerduty,omitempty"`
	OpsGenie  *OpsGenieConfig   `yaml:"opsgenie,omitempty"`
	Teams     *TeamsConfig      `yaml:"teams,omitempty"`
	Discord   *DiscordConfig    `yaml:"discord,omitempty"`
	VictorOps *VictorOpsConfig  `yaml:"victorops,omitempty"`
}

// WebhookConfig configures a generic webhook notifier.
type WebhookConfig struct {
	URL string `yaml:"url"`
}

// SlackConfig configures the Slack notifier.
type SlackConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

// EmailConfig configures the SMTP email notifier.
type EmailConfig struct {
	SMTPHost string   `yaml:"smtp_host"`
	SMTPPort int      `yaml:"smtp_port"`
	From     string   `yaml:"from"`
	To       []string `yaml:"to"`
}

// PagerDutyConfig configures the PagerDuty notifier.
type PagerDutyConfig struct {
	RoutingKey string `yaml:"routing_key"`
}

// OpsGenieConfig configures the OpsGenie notifier.
type OpsGenieConfig struct {
	APIKey string `yaml:"api_key"`
}

// TeamsConfig configures the Microsoft Teams notifier.
type TeamsConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

// DiscordConfig configures the Discord notifier.
type DiscordConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

// VictorOpsConfig configures the VictorOps notifier.
type VictorOpsConfig struct {
	RoutingKey      string `yaml:"routing_key"`
	RestEndpointURL string `yaml:"rest_endpoint_url"`
}
