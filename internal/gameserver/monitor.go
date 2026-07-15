package gameserver

import (
	"context"
	"errors"
	"sync"

	"github.com/coffeebased/kaeoi/pkg/poll"
)

var ErrMonitorRunning = errors.New("monitor is already running")

type Target interface {
	Latest() GameServer
	Refresh(ctx context.Context) (server GameServer, changed bool)
}

type Monitor struct {
	poller      poll.Poller
	targets     []Target
	subscribers map[chan GameServer]struct{}
	running     bool
	mu          sync.RWMutex
}

func NewMonitor(targets []Target, poller poll.Poller) *Monitor {
	snapshot := append([]Target(nil), targets...)

	return &Monitor{
		poller:      poller,
		targets:     snapshot,
		subscribers: make(map[chan GameServer]struct{}),
	}
}

func (m *Monitor) Subscribe(ctx context.Context) <-chan GameServer {
	if ctx == nil {
		panic("nil context")
	}

	ch := make(chan GameServer, len(m.targets))

	m.mu.Lock()
	if ctx.Err() != nil {
		close(ch)
		m.mu.Unlock()
		return ch
	}

	m.subscribers[ch] = struct{}{}
	for _, target := range m.targets {
		ch <- target.Latest()
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

func (m *Monitor) Run(ctx context.Context) error {
	if ctx == nil {
		panic("nil context")
	}

	m.mu.Lock()

	if m.running {
		m.mu.Unlock()
		return ErrMonitorRunning
	}

	m.running = true
	m.mu.Unlock()

	defer func() {
		m.mu.Lock()
		m.running = false
		m.mu.Unlock()
	}()

	var wg sync.WaitGroup

	for _, target := range m.targets {
		wg.Add(1)

		go func() {
			defer wg.Done()

			m.poller.Run(ctx, func(pollCtx context.Context) {
				server, changed := target.Refresh(pollCtx)
				if changed {
					m.publish(server)
				}
			})
		}()
	}

	<-ctx.Done()
	wg.Wait()

	return nil
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
