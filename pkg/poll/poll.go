// Package poll provides repeated callback execution with per-call timeouts.
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

type Options struct {
	Delay   time.Duration
	Timeout time.Duration
}

func (o *Options) normalize() {
	if o.Delay == 0 {
		o.Delay = defaultDelay
	}

	if o.Timeout == 0 {
		o.Timeout = defaultTimeout
	}
}

func (o Options) validate() error {
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
	options.normalize()

	if err := options.validate(); err != nil {
		return Poller{}, err
	}

	return Poller{
		options: options,
	}, nil
}

func (p Poller) Run(ctx context.Context, callback func(ctx context.Context)) {
	if ctx == nil {
		panic("nil context")
	}

	if callback == nil {
		panic("nil callback function")
	}

	for {
		if ctx.Err() != nil {
			return
		}

		checkCtx, cancel := context.WithTimeout(ctx, p.options.Timeout)
		callback(checkCtx)
		cancel()

		timer := time.NewTimer(p.options.Delay)

		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
	}
}
