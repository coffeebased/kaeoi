package gameserver

import (
	"context"
	"time"

	"github.com/mcstatus-io/mcutil/v4/status"
)

func queryMinecraftJava(ctx context.Context, server GameServer, queryTimeout time.Duration) State {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	host := server.Host
	if server.QueryHost != "" {
		host = server.QueryHost
	}

	port := server.Port
	if server.QueryPort != 0 {
		port = server.QueryPort
	}

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
