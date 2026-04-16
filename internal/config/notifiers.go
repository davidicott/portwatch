package config

// NotifierConfig holds configuration for all supported notifiers.
type NotifierConfig struct {
	Stdout    StdoutConfig    `yaml:"stdout"`
	Webhook   WebhookConfig   `yaml:"webhook"`
	Slack     SlackConfig     `yaml:"slack"`
	Email     EmailConfig     `yaml:"email"`
	PagerDuty PagerDutyConfig `yaml:"pagerduty"`
	OpsGenie  OpsGenieConfig  `yaml:"opsgenie"`
	Teams     TeamsConfig     `yaml:"teams"`
	Discord   DiscordConfig   `yaml:"discord"`
	VictorOps VictorOpsConfig `yaml:"victorops"`
	Syslog    SyslogConfig    `yaml:"syslog"`
	Telegram  TelegramConfig  `yaml:"telegram"`
	SNS       SNSConfig       `yaml:"sns"`
}

type StdoutConfig struct {
	Enabled bool `yaml:"enabled"`
}

type WebhookConfig struct {
	Enabled bool   `yaml:"enabled"`
	URL     string `yaml:"url"`
}

type SlackConfig struct {
	Enabled    bool   `yaml:"enabled"`
	WebhookURL string `yaml:"webhook_url"`
}

type EmailConfig struct {
	Enabled  bool   `yaml:"enabled"`
	SMTPHost string `yaml:"smtp_host"`
	SMTPPort int    `yaml:"smtp_port"`
	From     string `yaml:"from"`
	To       string `yaml:"to"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type PagerDutyConfig struct {
	Enabled    bool   `yaml:"enabled"`
	RoutingKey string `yaml:"routing_key"`
}

type OpsGenieConfig struct {
	Enabled bool   `yaml:"enabled"`
	APIKey  string `yaml:"api_key"`
}

type TeamsConfig struct {
	Enabled    bool   `yaml:"enabled"`
	WebhookURL string `yaml:"webhook_url"`
}

type DiscordConfig struct {
	Enabled    bool   `yaml:"enabled"`
	WebhookURL string `yaml:"webhook_url"`
}

type VictorOpsConfig struct {
	Enabled    bool   `yaml:"enabled"`
	URL        string `yaml:"url"`
	RoutingKey string `yaml:"routing_key"`
}

type SyslogConfig struct {
	Enabled bool   `yaml:"enabled"`
	Tag     string `yaml:"tag"`
}

type TelegramConfig struct {
	Enabled  bool   `yaml:"enabled"`
	BotToken string `yaml:"bot_token"`
	ChatID   string `yaml:"chat_id"`
}

// SNSConfig holds AWS SNS notifier settings.
type SNSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	TopicARN string `yaml:"topic_arn"`
	Region   string `yaml:"region"`
}
