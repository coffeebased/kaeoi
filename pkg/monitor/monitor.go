// Package monitor provides an agnostic polling coordinator for detecting status changes.
package monitor

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrMonitorAlreadyRunning = errors.New("monitor is already running")

const (
	defaultPollInterval   = 60 * time.Second
	defaultCheckTimeout   = 10 * time.Second
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

type Options struct {
	PollInterval   time.Duration
	CheckTimeout   time.Duration
	SubscriberSize int
}

func (o *Options) normalize() {
	if o.PollInterval == 0 {
		o.PollInterval = defaultPollInterval
	}

	if o.CheckTimeout == 0 {
		o.CheckTimeout = defaultCheckTimeout
	}

	if o.SubscriberSize == 0 {
		o.SubscriberSize = defaultSubscriberSize
	}
}

func (o Options) validate() error {
	if o.PollInterval < 0 {
		return errors.New("poll interval cannot be negative")
	}

	if o.CheckTimeout < 0 {
		return errors.New("check timeout cannot be negative")
	}

	if o.SubscriberSize < 0 {
		return errors.New("subscriber size cannot be negative")
	}

	return nil
}

type Monitor struct {
	items       map[ID]Monitorable
	nextID      ID
	options     Options
	subscribers map[chan Update]struct{}

	cancel context.CancelFunc
	done   chan struct{}

	mu sync.RWMutex
}

func New(options Options) (*Monitor, error) {
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
	if monitorable == nil {
		return 0, errors.New("monitorable is required")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.nextID++
	id := m.nextID

	m.items[id] = monitorable

	return id, nil
}

func (m *Monitor) Remove(id ID) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.items[id]; !exists {
		return false
	}

	delete(m.items, id)
	return true
}

func (m *Monitor) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	clear(m.items)
}

func (m *Monitor) Subscribe(ctx context.Context) <-chan Update {
	if ctx == nil {
		panic("nil context")
	}

	ch := make(chan Update, m.options.SubscriberSize)

	m.mu.Lock()
	if ctx.Err() != nil {
		close(ch)
		m.mu.Unlock()
		return ch
	}

	m.subscribers[ch] = struct{}{}
	m.mu.Unlock()

	go func() {
		<-ctx.Done()

		m.mu.Lock()
		if _, exists := m.subscribers[ch]; exists {
			delete(m.subscribers, ch)
			close(ch)
		}
		m.mu.Unlock()
	}()

	return ch
}

func (m *Monitor) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cancel != nil {
		return ErrMonitorAlreadyRunning
	}

	ctx, cancel := context.WithCancel(context.Background())

	m.cancel = cancel
	m.done = make(chan struct{})

	go m.run(ctx, m.done)

	return nil
}

func (m *Monitor) Stop() {
	m.mu.Lock()

	cancel := m.cancel
	done := m.done

	if cancel == nil {
		m.mu.Unlock()
		return
	}

	m.mu.Unlock()

	cancel()
	<-done

	m.mu.Lock()
	if m.done == done {
		m.cancel = nil
		m.done = nil
	}
	m.mu.Unlock()
}

type item struct {
	id          ID
	monitorable Monitorable
}

func (m *Monitor) snapshot() []item {
	m.mu.RLock()
	defer m.mu.RUnlock()

	items := make([]item, 0, len(m.items))

	for id, monitorable := range m.items {
		items = append(items, item{
			id:          id,
			monitorable: monitorable,
		})
	}

	return items
}

func (m *Monitor) run(ctx context.Context, done chan struct{}) {
	defer close(done)

	for {
		m.poll(ctx)

		timer := time.NewTimer(m.options.PollInterval)

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

func (m *Monitor) poll(ctx context.Context) {
	for _, item := range m.snapshot() {
		checkCtx, cancel := context.WithTimeout(ctx, m.options.CheckTimeout)
		changed, err := item.monitorable.Check(checkCtx)
		cancel()

		if ctx.Err() != nil {
			return
		}

		if changed || err != nil {
			m.publish(Update{
				ID:      item.id,
				Changed: changed,
				Err:     err,
			})
		}
	}
}

func (m *Monitor) publish(update Update) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for subscriber := range m.subscribers {
		select {
		case subscriber <- update:
		default:
		}
	}
}
