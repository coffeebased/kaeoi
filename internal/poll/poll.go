// Package poll provides timed polling for checkers that emit signals.
package poll

import (
	"context"
	"errors"
	"time"
)

const (
	defaultDelay       = 60 * time.Second
	defaultTimeout     = 10 * time.Second
	defaultChannelSize = 1
)

type Checker interface {
	Check(ctx context.Context) (bool, error)
}

type Signal struct {
	Err error
}

type Options struct {
	Delay       time.Duration
	Timeout     time.Duration
	ChannelSize int
}

func (o *Options) normalize() {
	if o.Delay == 0 {
		o.Delay = defaultDelay
	}

	if o.Timeout == 0 {
		o.Timeout = defaultTimeout
	}

	if o.ChannelSize == 0 {
		o.ChannelSize = defaultChannelSize
	}
}

func (o Options) validate() error {
	if o.Delay < 0 {
		return errors.New("delay cannot be negative")
	}

	if o.Timeout < 0 {
		return errors.New("timeout cannot be negative")
	}

	if o.ChannelSize < 0 {
		return errors.New("channel size cannot be negative")
	}

	return nil
}

type Poller struct {
	source  Checker
	options Options
}

func New(source Checker, options Options) (*Poller, error) {
	if source == nil {
		return nil, errors.New("source is required")
	}

	options.normalize()

	if err := options.validate(); err != nil {
		return nil, err
	}

	return &Poller{
		source:  source,
		options: options,
	}, nil
}

func (p *Poller) Run(ctx context.Context) <-chan Signal {
	channel := make(chan Signal, p.options.ChannelSize)

	go p.run(ctx, channel)

	return channel
}

func (p *Poller) run(ctx context.Context, channel chan Signal) {
	defer close(channel)

	for {
		checkCtx, cancel := context.WithTimeout(ctx, p.options.Timeout)
		trigger, err := p.source.Check(checkCtx)
		cancel()

		if ctx.Err() != nil {
			return
		}

		if trigger || err != nil {
			select {
			case channel <- Signal{
				Err: err,
			}:
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
