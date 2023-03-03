package player

import (
	"time"
)

type Ticker interface {
	// Ticks returns a read-only channel delivering ticks according to a
	// prescribed interval. The value returned does not need to be the same
	// channel, and may be nil.
	//
	// NOTE: Callers should assume that reads from Ticks() are stale after
	// any invocations of Resume, Pause, or Stop.
	Ticks() <-chan time.Time

	// Resume starts or resumes the underlying ticker, such that Ticks()
	// will fire at regular intervals. After calling Resume, Ticks() should
	// minimally send ticks at the prescribed interval.
	//
	// NOTE: It MUST be safe to call Resume at any time, and more than once
	// successively.
	Resume()

	// Pause suspends the underlying ticker, such that Ticks() stops
	// signaling at regular intervals. After calling Pause, the ticker
	// should not send any ticks scheduled with the chosen interval. Forced
	// ticks are still permissible, as in the case of the Force Ticker.
	//
	// NOTE: It MUST be safe to call Pause at any time, and more than once
	// successively.
	Pause()

	// Stop suspends the underlying ticker, such that Ticks() stops
	// signaling at regular intervals, and permanently frees up any
	// remaining resources.
	//
	// NOTE: The behavior of a Ticker is undefined after calling Stop.
	Stop()
}

// T is the production implementation of the resumable Ticker interface. This
// allows various components to toggle their need for tick events, which may
// vary depending on system load.
type T struct {
	// interval is the desired duration between ticks when active.
	interval time.Duration

	// ticker is the ephemeral, underlying time.Ticker. We keep a reference
	// to this ticker so that it can be stopped and cleaned up on Pause or
	// Stop.
	ticker *time.Ticker

	// time start - reflects the exact time of starting the ticker
	startTime time.Time

	// time stop - reflects the exact time of stopping the ticker
	elapsedTime time.Duration
}

// A compile-time constraint to ensure T satisfies the Ticker interface.
var _ Ticker = (*T)(nil)

// New returns a new ticker that signals with the given interval when not
// paused.
// The ticker starts active.
func New(interval time.Duration) *T {
	return &T{
		interval:  interval,
		ticker:    time.NewTicker(interval),
		startTime: time.Now(),
	}
}

// Ticks returns a receive-only channel that delivers times at the ticker's
// prescribed interval. This method returns nil when the ticker is paused.
//
// NOTE: Part of the Ticker interface.
func (t *T) Ticks() <-chan time.Time {
	if t.ticker == nil {
		return nil
	}
	return t.ticker.C
}

// Resume starts underlying time.Ticker and causes the ticker to begin
// delivering scheduled events.
//
// NOTE: Part of the Ticker interface.
func (t *T) Resume() {
	if t.ticker == nil {
		//t.ticker = time.NewTicker(t.interval)
		t.ticker.Stop()
		t.ticker = time.NewTicker(t.elapsedTime)
		t.startTime = time.Now()
	}
}

// Pause suspends the underlying ticker, such that Ticks() stops signaling at
// regular intervals.
//
// NOTE: Part of the Ticker interface.
func (t *T) Pause() {
	if t.ticker != nil {
		//t.ticker.Stop()
		//t.ticker = nil
		t.elapsedTime = time.Since(t.startTime)
	}
}

// Stop suspends the underlying ticker, such that Ticks() stops signaling at
// regular intervals, and permanently frees up any resources. For this
// implementation, this is equivalent to Pause.
//
// NOTE: Part of the Ticker interface.
func (t *T) Stop() {
	if t.ticker != nil {
		t.ticker.Stop()
		t.ticker = nil
	}
}
