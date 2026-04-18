// Package notify provides notifier implementations for portwatch.
//
// # SignalR Notifier
//
// The SignalRNotifier sends port change alerts to a SignalR-compatible HTTP
// endpoint such as the Azure SignalR Service REST API.
//
// Each event is serialised as a string argument under the "portwatch" target
// and posted to: <endpoint>/api/v1/hubs/<hub>
//
// Configuration fields:
//
//	enabled:   enable the notifier (default: false)
//	endpoint:  base URL of the SignalR service
//	hub:       hub name to broadcast to (default: "portwatch")
//	api_key:   bearer token for Authorization header (optional)
package notify
