package config

import "testing"

func TestKafkaConfig_Defaults(t *testing.T) {
	cfg := kafkaFromConfig(KafkaNotifierConfig{})
	if cfg.Broker != "localhost:9092" {
		t.Errorf("expected default broker, got %q", cfg.Broker)
	}
	if cfg.Topic != "portwatch-events" {
		t.Errorf("expected default topic, got %q", cfg.Topic)
	}
	if cfg.Enabled {
		t.Error("expected enabled=false by default")
	}
}

func TestKafkaConfig_Fields(t *testing.T) {
	cfg := kafkaFromConfig(KafkaNotifierConfig{
		Enabled: true,
		Broker:  "kafka.example.com:9092",
		Topic:   "my-topic",
	})
	if !cfg.Enabled {
		t.Error("expected enabled=true")
	}
	if cfg.Broker != "kafka.example.com:9092" {
		t.Errorf("unexpected broker: %q", cfg.Broker)
	}
	if cfg.Topic != "my-topic" {
		t.Errorf("unexpected topic: %q", cfg.Topic)
	}
}

func TestKafkaConfig_DefaultsPreservedWhenPartial(t *testing.T) {
	cfg := kafkaFromConfig(KafkaNotifierConfig{Enabled: true})
	if cfg.Broker != kafkaDefaults.Broker {
		t.Errorf("expected default broker fallback, got %q", cfg.Broker)
	}
	if cfg.Topic != kafkaDefaults.Topic {
		t.Errorf("expected default topic fallback, got %q", cfg.Topic)
	}
}
