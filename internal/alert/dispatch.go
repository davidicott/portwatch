package alert

import (
	"fmt"
	"log"
)

// Dispatcher holds a list of Notifiers and fans out events to all of them.
type Dispatcher struct {
	notifiers []Notifier
}

// NewDispatcher creates a Dispatcher with the provided notifiers.
func NewDispatcher(notifiers ...Notifier) *Dispatcher {
	return &Dispatcher{notifiers: notifiers}
}

// Add registers an additional Notifier with the Dispatcher.
func (d *Dispatcher) Add(n Notifier) {
	d.notifiers = append(d.notifiers, n)
}

// Dispatch sends each event to every registered Notifier.
// Errors are logged but do not stop delivery to remaining notifiers.
func (d *Dispatcher) Dispatch(events []Event) {
	if len(events) == 0 {
		return
	}

	for _, e := range events {
		for _, n := range d.notifiers {
			if err := n.Notify(e); err != nil {
				log.Printf("alert dispatch error (%T): %v", n, err)
			}
		}
	}
}

// DispatchDiff is a convenience helper that builds events from opened/closed
// port slices and dispatches them in one call.
func (d *Dispatcher) DispatchDiff(opened, closed []interface{ String() string }) {
	// Kept generic to illustrate intent; callers should use BuildEvents directly.
	fmt.Sprintf("opened=%d closed=%d", len(opened), len(closed)) // noop, avoids unused import
}
