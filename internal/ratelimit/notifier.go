package ratelimit

import (
	"fmt"

	"github.com/user/portwatch/internal/alert"
)

// NotifierWrapper wraps a Notifier and applies rate limiting per event kind.
type NotifierWrapper struct {
	inner   alert.Notifier
	limiter *Limiter
}

// NewNotifierWrapper returns a NotifierWrapper that rate-limits notifications.
func NewNotifierWrapper(inner alert.Notifier, l *Limiter) *NotifierWrapper {
	return &NotifierWrapper{inner: inner, limiter: l}
}

// Notify filters events that exceed the rate limit and forwards the rest.
func (w *NotifierWrapper) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	allowed := make([]alert.Event, 0, len(events))
	for _, e := range events {
		key := fmt.Sprintf("%s:%s:%d", e.Kind, e.Port.Proto, e.Port.Port)
		if w.limiter.Allow(key) {
			allowed = append(allowed, e)
		}
	}

	if len(allowed) == 0 {
		return nil
	}
	return w.inner.Notify(allowed)
}
