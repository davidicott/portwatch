package notify

import (
	"context"
	"errors"

	"github.com/user/portwatch/internal/alert"
)

// MultiNotifier fans out events to multiple Notifiers.
// All notifiers are called; errors are combined and returned together.
type MultiNotifier struct {
	notifiers []Notifier
}

// NewMultiNotifier returns a MultiNotifier that dispatches to each n.
func NewMultiNotifier(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{notifiers: notifiers}
}

// Notify sends events to every registered notifier.
// If one or more notifiers fail, all errors are joined and returned.
func (m *MultiNotifier) Notify(ctx context.Context, events []alert.Event) error {
	var errs []error
	for _, n := range m.notifiers {
		if err := n.Notify(ctx, events); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
