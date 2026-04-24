// Package notify provides notifier implementations for portwatch alerts.
//
// # Google Chat Notifier
//
// The GoogleChatNotifier sends port change alerts to a Google Chat space
// using an incoming webhook URL.
//
// Usage:
//
//	notifier := notify.NewGoogleChatNotifier("https://chat.googleapis.com/v1/spaces/.../messages?key=...")
//	err := notifier.Notify(events)
//
// The message includes a summary count and a bulleted list of each
// port event with its kind (opened/closed) and port identifier.
package notify
