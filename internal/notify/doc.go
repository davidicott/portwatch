// Package notify provides notifier implementations for delivering port-change
// alerts to external destinations.
//
// Available notifiers:
//
//   - StdoutNotifier  – writes formatted events to standard output.
//   - WebhookNotifier – HTTP POST of a JSON payload to a configurable URL.
//   - SlackNotifier   – posts a human-readable message to a Slack incoming webhook.
//   - MultiNotifier   – fan-out wrapper that delegates to multiple notifiers.
//
// All notifiers satisfy the alert.Notifier interface:
//
//	type Notifier interface {
//		Notify(events []alert.Event) error
//	}
//
// Compose notifiers with NewMultiNotifier to deliver alerts to several
// destinations simultaneously.
package notify
