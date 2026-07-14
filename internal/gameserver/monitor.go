package gameserver

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/coffeebased/kaeoi/pkg/poll"
)

var ErrMonitorRunning = errors.New("monitor is already running")

const defaultSubscriberSize = 16

type Target interface {
	GameServer() GameServer
	poll.Checker
}

type MonitorOptions struct {
	SubscriberSize int
	Pollers        poll.Options
}

func (o *MonitorOptions) normalize() {
	if o.SubscriberSize == 0 {
		o.SubscriberSize = defaultSubscriberSize
	}

	o.Pollers.Normalize()
}

func (o MonitorOptions) validate() error {
	if o.SubscriberSize <= 0 {
		return errors.New("subscriber size cannot be less than zero")
	}

	if err := o.Pollers.Validate(); err != nil {
		return fmt.Errorf("poller options: %w", err)
	}

	return nil
}

type Monitor struct {
	targets     []Target
	options     MonitorOptions
	subscribers map[chan GameServer]struct{}
	mu          sync.RWMutex
}

func newMonitor(options MonitorOptions) (*Monitor, error) {
	options.normalize()

	if err := options.validate(); err != nil {
		return nil, err
	}

	return &Monitor{
		options:     options,
		subscribers: make(map[chan GameServer]struct{}),
	}, nil
}

func (m *Monitor) Subscribe(ctx context.Context) <-chan GameServer {
	if ctx == nil {
		panic("nil context")
	}

	ch := make(chan GameServer, m.options.SubscriberSize)

	m.mu.Lock()
	if ctx.Err() != nil {
		close(ch)
		m.mu.Unlock()
		return ch
	}

	m.subscribers[ch] = struct{}{}
	for _, target := range m.targets {
		ch <- target.GameServer()
	}
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

func (m *Monitor) Run(ctx context.Context, targets []Target) error {
	pollers := make([]poll.Poller, 0, len(targets))
	channels := make([]chan struct{}, 0, len(targets))
	m.targets = make([]Target, 0, len(targets))

	m.mu.Lock()
	for _, target := range targets {
		poller, err := poll.New(m.options.Pollers)
		if err != nil {
			m.mu.Unlock()
			return err
		}

		channel := make(chan struct{}, 1)

		pollers = append(pollers, poller)
		channels = append(channels, channel)
		m.targets = append(m.targets, target)
	}
	m.mu.Unlock()

	for i := range pollers {
		go pollers[i].Run(ctx, m.targets[i], channels[i])
	}

	for {
		if ctx.Err() != nil {
			return nil
		}

		for i := range m.targets {
			select {
			case <-ctx.Done():
				return nil
			case <-channels[i]:
				m.publish(m.targets[i].GameServer())
			default:
			}
		}
	}
}

func (m *Monitor) publish(server GameServer) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for ch := range m.subscribers {
		select {
		case ch <- server:
		default:
		}
	}
}
