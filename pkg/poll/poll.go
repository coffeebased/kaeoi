// Package poll provides timed polling for checkers and emits a signal when true is returned.
package poll

import (
	"context"
	"errors"
	"time"
)

const (
	defaultDelay   = 60 * time.Second
	defaultTimeout = 10 * time.Second
)

type Checker interface {
	Check(ctx context.Context) bool
}

type Options struct {
	Delay   time.Duration
	Timeout time.Duration
}

func (o *Options) Normalize() {
	if o.Delay == 0 {
		o.Delay = defaultDelay
	}

	if o.Timeout == 0 {
		o.Timeout = defaultTimeout
	}
}

func (o Options) Validate() error {
	if o.Delay < 0 {
		return errors.New("delay cannot be negative")
	}

	if o.Timeout < 0 {
		return errors.New("timeout cannot be negative")
	}

	return nil
}

type Poller struct {
	options Options
}

func New(options Options) (Poller, error) {
	options.Normalize()

	if err := options.Validate(); err != nil {
		return Poller{}, err
	}

	return Poller{
		options: options,
	}, nil
}

func (p Poller) Run(ctx context.Context, source Checker, signals chan<- struct{}) {
	for {
		checkCtx, cancel := context.WithTimeout(ctx, p.options.Timeout)
		trigger := source.Check(checkCtx)
		cancel()

		if ctx.Err() != nil {
			return
		}

		if trigger {
			select {
			case signals <- struct{}{}:
			default:
			}
		}

		timer := time.NewTimer(p.options.Delay)

		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return
		case <-timer.C:
		}
	}
}
