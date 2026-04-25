package config

// PagerDutyV2Config holds settings for the PagerDuty Events API v2 notifier.
type PagerDutyV2Config struct {
	Enabled    bool   `yaml:"enabled"`
	RoutingKey string `yaml:"routing_key"`
	Endpoint   string `yaml:"endpoint"`
}

var pagerDutyV2Defaults = PagerDutyV2Config{
	Enabled:  false,
	Endpoint: "https://events.pagerduty.com/v2/enqueue",
}

func init() {
	registerNotifierDefaults("pagerduty_v2", pagerDutyV2Defaults)
}

// pagerDutyV2FromConfig returns a PagerDutyV2Config merged with defaults.
func pagerDutyV2FromConfig(raw map[string]interface{}) PagerDutyV2Config {
	cfg := pagerDutyV2Defaults
	if v, ok := raw["enabled"].(bool); ok {
		cfg.Enabled = v
	}
	if v, ok := raw["routing_key"].(string); ok {
		cfg.RoutingKey = v
	}
	if v, ok := raw["endpoint"].(string); ok && v != "" {
		cfg.Endpoint = v
	}
	return cfg
}
