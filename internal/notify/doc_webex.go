// Package notify provides notifier implementations for portwatch alerts.
//
// # Webex Notifier
//
// The WebexNotifier sends port change alerts to a Cisco Webex space using
// the Webex REST API (https://webexapis.com/v1/messages).
//
// # Configuration
//
//	 notifiers:
//	   webex:
//	     enabled: true
//	     token: "<bot-access-token>"
//	     room_id: "<webex-room-id>"
//
// The bot must be a member of the target Webex space. Obtain a bot access
// token from https://developer.webex.com.
package notify
