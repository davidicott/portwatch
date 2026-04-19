// Package notify provides alert delivery integrations for portwatch.
//
// # HipChat
//
// The HipChatNotifier delivers port-change alerts to a HipChat room using
// the HipChat v2 REST API.
//
// Configuration fields:
//
//	enabled    – set to true to activate this notifier
//	server_url – base URL of the HipChat server (default: https://api.hipchat.com)
//	room_id    – numeric or string ID of the target room
//	token      – personal or room API token with Send Notification scope
//
// Example YAML:
//
//	hipchat:
//	  enabled: true
//	  server_url: https://api.hipchat.com
//	  room_id: "42"
//	  token: "your-api-token"
package notify
