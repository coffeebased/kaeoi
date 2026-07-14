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
	mu     sync.Mutex
}

func (s *SteamA2SServer) GameServer() GameServer {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.server
}

func (s *SteamA2SServer) Check(ctx context.Context) bool {
	s.mu.Lock()
	address := s.server.GetQueryAddress()
	s.mu.Unlock()

	state := querySteamA2S(ctx, address)

	s.mu.Lock()
	defer s.mu.Unlock()

	if state != s.server.State {
		s.server.State = state

		return true
	}

	return false
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
		Message:     info.Map,
		Status:      StatusOnline,
	}
}
