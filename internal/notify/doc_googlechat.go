// Package notify provides notifier implementations for portwatch alerts.
//
// # Google Chat Notifier
//
// The GoogleChatNotifier sends port change alerts to a Google Chat space
// via an incoming webhook URL.
//
// Usage:
//
//	notifier := notify.NewGoogleChatNotifier("https://chat.googleapis.com/...")
//	err := notifier.Notify(events)
//
// The webhook URL can be obtained from a Google Chat space by configuring
// an incoming webhook integration under "Manage webhooks".
package notify
