package gameserver

import (
	"context"
	"errors"
	"sync"

	"github.com/coffeebased/kaeoi/pkg/poll"
)

var ErrMonitorRunning = errors.New("monitor is already running")

const defaultSubscriberSize = 16

type MonitorOptions struct {
	Pollers        poll.Options
	SubscriberSize int
}

func (o *MonitorOptions) normalize() {
	o.Pollers.Normalize()

	if o.SubscriberSize == 0 {
		o.SubscriberSize = defaultSubscriberSize
	}
}

func (o MonitorOptions) validate() error {
	if err := o.Pollers.Validate(); err != nil {
		return err
	}

	if o.SubscriberSize <= 0 {
		return errors.New("subscriber size cannot be less than zero")
	}

	return nil
}

type Monitor struct {
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

func (m *Monitor) Run(ctx context.Context, servers []GameServer) error {
	pollers := make([]*poll.Poller, 0, len(servers))
	channels := make([]chan poll.Signal, 0, len(servers))
	items := make([]GameServer, 0, len(servers))

	for _, server := range servers {
		poller, err := poll.New(m.options.Pollers)
		if err != nil {
			return err
		}

		channel := make(chan poll.Signal, 1)

		pollers = append(pollers, poller)
		channels = append(channels, channel)
		items = append(items, server)
	}

	for i := range pollers {
		go pollers[i].Run(ctx, &items[i], channels[i])
	}

	for {
		if ctx.Err() != nil {
			return nil
		}

		for i := range items {
			select {
			case <-ctx.Done():
				return nil
			case <-channels[i]:
				m.publish(items[i])
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
