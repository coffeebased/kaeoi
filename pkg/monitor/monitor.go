// Package monitor provides a standard and agnostic poll cordinator to detect status changes
package monitor

import (
	"context"
	"errors"
	"sync"
	"time"
)

const (
	defaultPollInterval   = 60 * time.Second
	defaultQueryTimeout   = 10 * time.Second
	defaultSubscriberSize = 16
)

type ID uint64

type Monitorable interface {
	Check(ctx context.Context) (bool, error)
}

type Update struct {
	ID      ID
	Changed bool
	Err     error
}

type MonitorOptions struct {
	PollInterval   time.Duration
	QueryTimeout   time.Duration
	SubscriberSize int
}

func (o *MonitorOptions) normalize() {
	if o.PollInterval == 0 {
		o.PollInterval = defaultPollInterval
	}

	if o.QueryTimeout == 0 {
		o.QueryTimeout = defaultQueryTimeout
	}

	if o.SubscriberSize == 0 {
		o.SubscriberSize = defaultSubscriberSize
	}
}

func (o MonitorOptions) validate() error {
	if o.PollInterval < 0 {
		return errors.New("poll interval cannot be negative")
	}

	if o.QueryTimeout < 0 {
		return errors.New("query timeout cannot be negative")
	}

	if o.SubscriberSize < 0 {
		return errors.New("subscriber size cannot be negative")
	}

	return nil
}

type Monitor struct {
	items       map[ID]Monitorable
	nextID      ID
	options     MonitorOptions
	subscribers map[chan Update]struct{}

	cancel context.CancelFunc
	done   chan struct{}

	mu sync.RWMutex
}

func New(options MonitorOptions) (*Monitor, error) {
	options.normalize()

	if err := options.validate(); err != nil {
		return nil, err
	}

	return &Monitor{
		items:       make(map[ID]Monitorable),
		options:     options,
		subscribers: make(map[chan Update]struct{}),
	}, nil
}

func (m *Monitor) Add(monitorable Monitorable) (ID, error) {
	panic("not implemented")
}

func (m *Monitor) Remove(id ID) bool {
	panic("not implemented")
}

func (m *Monitor) Subscribe(ctx context.Context) <-chan Update {
	panic("not implemented")
}

func (m *Monitor) Start() error {
	panic("not implemented")
}

func (m *Monitor) Stop() {
	panic("not implemented")
}
