package notify

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"

	"portwatch/internal/alert"
)

// SNSClient is the subset of the AWS SNS API we use.
type SNSClient interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

// SNSNotifier sends alerts to an AWS SNS topic.
type SNSNotifier struct {
	client   SNSClient
	topicARN string
}

// NewSNSNotifier creates an SNSNotifier using the default AWS credential chain.
func NewSNSNotifier(topicARN, region string) (*SNSNotifier, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("sns: load aws config: %w", err)
	}
	return &SNSNotifier{
		client:   sns.NewFromConfig(cfg),
		topicARN: topicARN,
	}, nil
}

// NewSNSNotifierWithClient creates an SNSNotifier with a custom client (useful for testing).
func NewSNSNotifierWithClient(client SNSClient, topicARN string) *SNSNotifier {
	return &SNSNotifier{client: client, topicARN: topicARN}
}

// Notify publishes port change events to the configured SNS topic.
func (n *SNSNotifier) Notify(ctx context.Context, events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}
	payload, err := json.Marshal(map[string]any{
		"source": "portwatch",
		"events": events,
	})
	if err != nil {
		return fmt.Errorf("sns: marshal payload: %w", err)
	}
	_, err = n.client.Publish(ctx, &sns.PublishInput{
		TopicArn: aws.String(n.topicARN),
		Message:  aws.String(string(payload)),
		Subject:  aws.String(fmt.Sprintf("portwatch: %d port change(s) detected", len(events))),
	})
	if err != nil {
		return fmt.Errorf("sns: publish: %w", err)
	}
	return nil
}
