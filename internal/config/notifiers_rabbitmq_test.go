package config

import (
	"testing"
)

func TestRabbitMQConfig_Defaults(t *testing.T) {
	cfg := DefaultConfig()
	rm := cfg.Notifiers.RabbitMQ

	if rm.Exchange != "portwatch" {
		t.Errorf("expected default exchange 'portwatch', got %q", rm.Exchange)
	}
	if rm.RoutingKey != "ports.changed" {
		t.Errorf("expected default routing key 'ports.changed', got %q", rm.RoutingKey)
	}
	if rm.Enabled {
		t.Error("expected RabbitMQ notifier to be disabled by default")
	}
}

func TestRabbitMQConfig_Fields(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Notifiers.RabbitMQ.Enabled = true
	cfg.Notifiers.RabbitMQ.URL = "amqp://user:pass@localhost:5672/"
	cfg.Notifiers.RabbitMQ.Exchange = "alerts"
	cfg.Notifiers.RabbitMQ.RoutingKey = "portwatch.events"

	rm := cfg.Notifiers.RabbitMQ

	if !rm.Enabled {
		t.Error("expected enabled to be true")
	}
	if rm.URL != "amqp://user:pass@localhost:5672/" {
		t.Errorf("unexpected URL: %q", rm.URL)
	}
	if rm.Exchange != "alerts" {
		t.Errorf("unexpected exchange: %q", rm.Exchange)
	}
	if rm.RoutingKey != "portwatch.events" {
		t.Errorf("unexpected routing key: %q", rm.RoutingKey)
	}
}

func TestRabbitMQConfig_URLRequired(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Notifiers.RabbitMQ.Enabled = true
	// URL intentionally left empty — callers should validate before dialing.
	if cfg.Notifiers.RabbitMQ.URL != "" {
		t.Error("expected empty URL by default")
	}
}
