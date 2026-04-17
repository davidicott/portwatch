// Package notify provides notification backends for portwatch alerts.
//
// # Mattermost
//
// The MattermostNotifier posts alert messages to a Mattermost channel via
// an incoming webhook URL. Configure the webhook in your Mattermost instance
// under Integrations → Incoming Webhooks.
//
// Example configuration (config.yaml):
//
//	notifiers:
//	  mattermost:
//	    enabled: true
//	    webhook_url: "https://mattermost.example.com/hooks/<token>"
//	    channel: "#security-alerts"
//
// Each notification includes the number of changes and a list of affected
// ports with their event kind (opened/closed).
package notify
