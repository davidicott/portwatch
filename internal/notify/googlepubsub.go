package notify

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"

	"github.com/user/portwatch/internal/alert"
)

// GooglePubSubNotifier publishes port change events to a Google Cloud Pub/Sub topic.
type GooglePubSubNotifier struct {
	client    *pubsub.Client
	topic     *pubsub.Topic
	projectID string
	topicID   string
}

type pubSubPayload struct {
	Kind    string `json:"kind"`
	Proto   string `json:"proto"`
	Port    int    `json:"port"`
	Address string `json:"address"`
}

// NewGooglePubSubNotifier creates a notifier that publishes to the given Pub/Sub topic.
func NewGooglePubSubNotifier(ctx context.Context, projectID, topicID string) (*GooglePubSubNotifier, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub: create client: %w", err)
	}
	topic := client.Topic(topicID)
	return &GooglePubSubNotifier{
		client:    client,
		topic:     topic,
		projectID: projectID,
		topicID:   topicID,
	}, nil
}

// Notify publishes each event as a separate Pub/Sub message.
func (n *GooglePubSubNotifier) Notify(ctx context.Context, events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}
	for _, e := range events {
		payload := pubSubPayload{
			Kind:    string(e.Kind),
			Proto:   e.Port.Proto,
			Port:    e.Port.Port,
			Address: e.Port.Address,
		}
		data, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("pubsub: marshal event: %w", err)
		}
		result := n.topic.Publish(ctx, &pubsub.Message{Data: data})
		if _, err := result.Get(ctx); err != nil {
			return fmt.Errorf("pubsub: publish message: %w", err)
		}
	}
	return nil
}

// Close stops the topic and closes the underlying client.
func (n *GooglePubSubNotifier) Close() error {
	n.topic.Stop()
	return n.client.Close()
}
