package config

import "time"

// RabbitMQNotifierConfig holds configuration for the RabbitMQ notifier.
type RabbitMQNotifierConfig struct {
	Enabled    bool          `yaml:"enabled"`
	URL        string        `yaml:"url"`
	Exchange   string        `yaml:"exchange"`
	RoutingKey string        `yaml:"routing_key"`
	Timeout    time.Duration `yaml:"timeout"`
}

var rabbitMQDefaults = RabbitMQNotifierConfig{
	Enabled:    false,
	URL:        "amqp://guest:guest@localhost:5672/",
	Exchange:   "portwatch",
	RoutingKey: "port.events",
	Timeout:    5 * time.Second,
}

func init() {
	registerNotifierDefaults("rabbitmq", rabbitMQDefaults)
}

// rabbitmqFromConfig returns a RabbitMQNotifierConfig merged with defaults.
func rabbitmqFromConfig(raw map[string]interface{}) RabbitMQNotifierConfig {
	cfg := rabbitMQDefaults
	if v, ok := raw["enabled"].(bool); ok {
		cfg.Enabled = v
	}
	if v, ok := raw["url"].(string); ok && v != "" {
		cfg.URL = v
	}
	if v, ok := raw["exchange"].(string); ok && v != "" {
		cfg.Exchange = v
	}
	if v, ok := raw["routing_key"].(string); ok && v != "" {
		cfg.RoutingKey = v
	}
	if v, ok := raw["timeout"].(time.Duration); ok && v > 0 {
		cfg.Timeout = v
	}
	return cfg
}
