// ticker implements functionality for generating events
// at user specified time instances (hourly, daily, etc).
package ticker

import (
	"context"
	"time"
)

// ChangePred takes two timestampts and returns whether they are
// equal or not by some metric.
type ChangePred func(time.Time, time.Time) bool

func tensCoeff(val int) (int, int) {
	ones := (val % 10)
	return val - ones, ones
}

func HourChangePred(prev, curr time.Time) bool {
	return curr.Minute() == 0 && (prev.Hour() != curr.Hour())
}

func TenSecondChangePred(prev time.Time, curr time.Time) bool {
	lhs, _ := tensCoeff(prev.Second())
	rhs, rhs1 := tensCoeff(curr.Second())
	return rhs1 == 0 && (lhs != rhs)
}

func OneSecondChangePred(prev, curr time.Time) bool {
	return curr.After(prev) && ((prev.Second() % 10) != (curr.Second() % 10))
}

type Ticker struct {
	c            chan time.Time
	pred         ChangePred
	pollInterval time.Duration
}

func NewTicker(pred ChangePred, pollInterval time.Duration) *Ticker {
	ch := make(chan time.Time)
	return &Ticker{c: ch, pred: pred, pollInterval: pollInterval}
}

func (t *Ticker) NotifyChannel() <-chan time.Time {
	return t.c
}

func (t *Ticker) RunWithContext(ctx context.Context) {
	prev := time.Now()
	// internal ticker which we translate
	it := time.NewTicker(t.pollInterval)
loop:
	for {
		select {
		case curr := <-it.C:
			if !t.pred(prev, curr) {
				continue
			}
			t.c <- curr
			prev = curr
		case <-ctx.Done():
			break loop
		}
	}
}

func (t *Ticker) Run() {
	t.RunWithContext(context.Background())
}
