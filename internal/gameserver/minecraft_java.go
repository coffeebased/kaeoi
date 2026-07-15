package gameserver

import (
	"context"
	"sync"

	"github.com/mcstatus-io/mcutil/v4/status"
)

type MinecraftJavaServer struct {
	server GameServer
	mu     sync.RWMutex
}

func (s *MinecraftJavaServer) Latest() GameServer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.server
}

func (s *MinecraftJavaServer) Refresh(ctx context.Context) (GameServer, bool) {
	s.mu.Lock()
	host := s.server.GetQueryHost()
	port := s.server.GetQueryPort()
	lastChange := s.server.Metadata.LastChange
	s.mu.Unlock()

	state := queryMinecraftJava(ctx, host, port)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.server.SetState(state, "")

	if lastChange == s.server.Metadata.LastChange {
		return s.server, false
	}

	return s.server, true
}

func queryMinecraftJava(ctx context.Context, host string, port int) State {
	data, err := status.Modern(ctx, host, uint16(port))
	if err != nil {
		return State{
			Status: StatusOffline,
		}
	}

	maxPlayers := 0
	if data.Players.Max != nil {
		maxPlayers = int(*data.Players.Max)
	}

	playersCount := 0
	if data.Players.Online != nil {
		playersCount = int(*data.Players.Online)
	}

	return State{
		Version:     data.Version.Name.Clean,
		MaxPlayers:  maxPlayers,
		PlayerCount: playersCount,
		Description: data.MOTD.Clean,
		Status:      StatusOnline,
	}
}
