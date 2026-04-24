package config

// KafkaNotifierConfig holds configuration for the Kafka notifier.
type KafkaNotifierConfig struct {
	Enabled bool   `yaml:"enabled"`
	Broker  string `yaml:"broker"`
	Topic   string `yaml:"topic"`
}

var kafkaDefaults = KafkaNotifierConfig{
	Enabled: false,
	Broker:  "localhost:9092",
	Topic:   "portwatch-events",
}

func init() {
	registerNotifierDefaults("kafka", kafkaDefaults)
}

// kafkaFromConfig returns a KafkaNotifierConfig populated with defaults
// for any zero-value fields.
func kafkaFromConfig(raw KafkaNotifierConfig) KafkaNotifierConfig {
	if raw.Broker == "" {
		raw.Broker = kafkaDefaults.Broker
	}
	if raw.Topic == "" {
		raw.Topic = kafkaDefaults.Topic
	}
	return raw
}
