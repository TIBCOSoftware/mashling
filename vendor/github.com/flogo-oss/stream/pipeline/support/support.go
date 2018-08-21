package support

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"time"
)

type TimerCallback func(ctx activity.Context) (resume bool)

// TimerSupport is used to create a timer on behalf of the activity
type TimerSupport interface {

	// HasTimer indicates if a timer already exists
	HasTimer(repeating bool) bool

	// CancelTimer cancels the existing timer
	CancelTimer(repeating bool)

	UpdateTimer(repeating bool)

	// CreateTimer creates a timer, note: can only have one active timer at a time for an activity
	CreateTimer(interval time.Duration, callback TimerCallback, repeating bool) error
}

// GetTimerSupport for the activity
func GetTimerSupport(ctx activity.Context) (TimerSupport, bool) {

	ts, ok :=  ctx.(TimerSupport)
	return ts, ok
}
