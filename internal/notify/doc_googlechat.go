// Package notify provides notifier implementations for portwatch alerts.
//
// # Google Chat Notifier
//
// The GoogleChatNotifier delivers port change events to a Google Chat space
// via an incoming webhook URL. Each event is formatted as a single line
// containing the event kind, host, port number, and protocol.
//
// Configuration:
//
//	notifiers:
//	  googlechat:
//	    enabled: true
//	    webhook_url: "https://chat.googleapis.com/v1/spaces/.../messages?key=..."
//
// Obtain a webhook URL by adding an "Incoming Webhook" app to a Google Chat
// space and copying the generated URL.
package notify
