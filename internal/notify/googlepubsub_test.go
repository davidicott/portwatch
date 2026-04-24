package notify

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPubSubClient struct {
	published [][]byte
	err       error
}

func (m *mockPubSubClient) Publish(ctx context.Context, topic string, data []byte) error {
	if m.err != nil {
		return m.err
	}
	m.published = append(m.published, data)
	return nil
}

func makeGPSEvent(kind, proto string, port int) Event {
	return Event{
		Kind:     kind,
		Proto:    proto,
		Port:     port,
		Hostname: "testhost",
	}
}

func TestGooglePubSubNotifier_SkipsEmptyEvents(t *testing.T) {
	client := &mockPubSubClient{}
	n := newTestPubSubNotifier(client, "projects/test/topics/portwatch")
	err := n.Notify(context.Background(), nil)
	require.NoError(t, err)
	assert.Empty(t, client.published)
}

func TestGooglePubSubNotifier_PostsPayload(t *testing.T) {
	client := &mockPubSubClient{}
	n := newTestPubSubNotifier(client, "projects/test/topics/portwatch")
	events := []Event{
		makeGPSEvent("opened", "tcp", 8080),
		makeGPSEvent("closed", "udp", 53),
	}
	err := n.Notify(context.Background(), events)
	require.NoError(t, err)
	require.Len(t, client.published, 1)

	var payload map[string]interface{}
	require.NoError(t, json.Unmarshal(client.published[0], &payload))
	eventsRaw, ok := payload["events"]
	assert.True(t, ok)
	eventsSlice, ok := eventsRaw.([]interface{})
	assert.True(t, ok)
	assert.Len(t, eventsSlice, 2)
}

func TestGooglePubSubNotifier_PayloadContainsCount(t *testing.T) {
	client := &mockPubSubClient{}
	n := newTestPubSubNotifier(client, "projects/test/topics/portwatch")
	events := []Event{
		makeGPSEvent("opened", "tcp", 443),
	}
	err := n.Notify(context.Background(), events)
	require.NoError(t, err)
	require.Len(t, client.published, 1)

	var payload map[string]interface{}
	require.NoError(t, json.Unmarshal(client.published[0], &payload))
	count, ok := payload["count"]
	assert.True(t, ok)
	assert.Equal(t, float64(1), count)
}
