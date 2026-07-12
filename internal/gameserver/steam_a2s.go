package gameserver

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/coffeebased/kaeoi/pkg/a2s"
)

func querySteamA2S(ctx context.Context, server GameServer, queryTimeout time.Duration) State {
	localCtx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	dialer := net.Dialer{}

	conn, err := dialer.DialContext(localCtx, "udp", server.GetQueryAddress())
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

	if deadline, ok := localCtx.Deadline(); ok {
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
