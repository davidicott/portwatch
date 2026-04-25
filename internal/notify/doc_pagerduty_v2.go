// Package notify provides alert notifier implementations for portwatch.
//
// # PagerDuty Events API v2 Notifier
//
// The PagerDutyV2 notifier sends alert events to PagerDuty using the
// Events API v2 (https://developer.pagerduty.com/docs/ZG9jOjExMDI5NTgw).
//
// Configuration fields:
//
//	  pagerduty_v2:
//	    enabled: true
//	    routing_key: "<your integration key>"
//	    endpoint: "https://events.pagerduty.com/v2/enqueue"  # optional
//
// Each scan cycle that produces at least one port change event will trigger
// a single PagerDuty incident with a summary describing the number of changes.
// The severity is set to "warning" and the source is "portwatch".
package notify
