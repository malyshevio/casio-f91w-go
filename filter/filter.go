// ticker implements functionality for generating events
// at user specified time instances (hourly, daily, etc).
package ticker_filter

import (
	"time"
)

type RoundingFunc func(time.Time) time.Time

func FloorToSeconds(t time.Time) time.Time {
	y, m, d := t.Date()
	hh, mm, ss := t.Clock()
	return time.Date(
		y, m, d,
		hh, mm, ss,
		0,
		t.Location(),
	)
}

func FloorToTensOfSeconds(t time.Time) time.Time {
	y, m, d := t.Date()
	hh, mm, ss := t.Clock()
	new_ss := (ss - (ss % 10))
	return time.Date(
		y, m, d,
		hh, mm, new_ss,
		0,
		t.Location(),
	)
}

func FloortoMinute(t time.Time) time.Time {
	y, m, d := t.Date()
	hh, mm, _ := t.Clock()
	return time.Date(
		y, m, d,
		hh, mm, 0,
		0,
		t.Location(),
	)
}

func FloorToHour(t time.Time) time.Time {
	y, m, d := t.Date()
	hh := t.Hour()
	return time.Date(
		y, m, d,
		hh, 0, 0,
		0,
		t.Location(),
	)
}

type TickerFilter struct {
	RoundingFunc RoundingFunc
	prev         time.Time
}

func NewTickerFilter(roundingFunc RoundingFunc) *TickerFilter {
	return &TickerFilter{prev: time.Time{}, RoundingFunc: roundingFunc}
}

// Sets the initial "previous" value against which the next HasChanged() argument
// is compared against.
func (tf *TickerFilter) Init(reference time.Time) {
	tf.prev = tf.RoundingFunc(reference)
}

// Clears out any previous value
func (tf *TickerFilter) Reset() {
	tf.prev = time.Time{}
}

// HasChanged return true if rounded version of "curr" results in different (rounded)
// representation than the internal (rounded) "prev" reference value. This implies that
// the first call to this function will always return false since there is no "prev"
// reference to be compared against. However, note that you may initialize the "prev"
// to a specific value with Init() if needed.
//
// NOTE: it is the user's responsibility to pass values to this call at sufficiently
// high rate; for example, if you wish to track changes to the Seconds value, then
// the difference between two consecutive values of "curr" should be less than a
// second (e.g. 100ms). In other words, the sampling rate of "curr" should be sufficiently
// high compared to the change rate implied by TickerFilter.RoundingFunc.
func (tf *TickerFilter) HasChanged(curr time.Time) bool {
	currTrunc := tf.RoundingFunc(curr)

	// NOTE: We treat the "zero" time.Time{} as "no value set" sentinel
	res := !tf.prev.IsZero() && currTrunc.After(tf.prev)
	tf.prev = currTrunc
	return res
}
