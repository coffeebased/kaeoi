package gameserver

import (
	"context"
	"sync"

	"github.com/mcstatus-io/mcutil/v4/status"
)

type MinecraftJavaServer struct {
	server GameServer
	mu     sync.Mutex
}

func (s *MinecraftJavaServer) GameServer() GameServer {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.server
}

func (s *MinecraftJavaServer) Check(ctx context.Context) bool {
	s.mu.Lock()
	host := s.GameServer().GetQueryHost()
	port := s.GameServer().GetQueryPort()
	s.mu.Unlock()

	state := queryMinecraftJava(ctx, host, port)

	s.mu.Lock()
	defer s.mu.Unlock()

	if state != s.server.State {
		s.server.State = state

		return true
	}

	return false
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
		Message:     data.MOTD.Clean,
		Status:      StatusOnline,
	}
}
