package gameserver

import (
	"context"
	"errors"
	"net"
	"sync"

	"github.com/coffeebased/kaeoi/pkg/a2s"
)

type SteamA2SServer struct {
	server GameServer
	mu     sync.RWMutex
}

func newSteamA2SServer(server GameServer) *SteamA2SServer {
	return &SteamA2SServer{
		server: server,
	}
}

func (s *SteamA2SServer) Latest() GameServer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.server
}

func (s *SteamA2SServer) Refresh(ctx context.Context) (GameServer, bool) {
	s.mu.Lock()
	address := s.server.GetQueryAddress()
	s.mu.Unlock()

	state := querySteamA2S(ctx, address)

	s.mu.Lock()
	defer s.mu.Unlock()

	changed := s.server.UpdateState(state)
	return s.server, changed
}

func querySteamA2S(ctx context.Context, address string) State {
	dialer := net.Dialer{}
	conn, err := dialer.DialContext(ctx, "udp", address)
	if err != nil {
		if ctx.Err() != nil {
			return State{
				Status: StatusUnknown,
			}
		}

		return State{
			Status: StatusOffline,
		}
	}
	defer func() {
		_ = conn.Close()
	}()

	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
	}

	info, err := a2s.QueryInfo(conn)
	if err != nil {
		if ctx.Err() != nil {
			return State{
				Status: StatusUnknown,
			}
		}

		if errors.Is(err, a2s.ErrA2SSplitPacket) ||
			errors.Is(err, a2s.ErrA2SInvalidResponse) ||
			errors.Is(err, a2s.ErrA2SUnsupportedPacket) ||
			errors.Is(err, a2s.ErrA2SChallengeMalformed) {
			return State{
				Status: StatusUnknown,
			}
		}

		return State{
			Status: StatusOffline,
		}
	}

	return State{
		Application: info.Game,
		Version:     info.Version,
		MaxPlayers:  info.MaxPlayers,
		PlayerCount: info.Players,
		Title:       info.Name,
		Description: info.Map,
		Status:      StatusOnline,
	}
}
