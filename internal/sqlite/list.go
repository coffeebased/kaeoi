package sqlite

import (
	"context"

	"github.com/coffeebased/kaeoi/internal/gameserver"
)

func (s *Store) List(ctx context.Context) ([]gameserver.GameServer, error) {
	if ctx == nil {
		panic("nil context")
	}

	rows, err := s.db.QueryContext(ctx, `
	SELECT
		gs.code,
		gs.host,
		gs.port,
		gs.query_host,
		gs.query_port,
		gs.kind,
		gs.password,
		gs.displays_ip,
		gs.displays_port,
		gsd.application,
		gsd.version,
		gsd.max_players,
		gsd.title,
		gsd.description
	FROM game_servers AS gs
	INNER JOIN game_server_defaults AS gsd
	ON gs.code = gsd.game_server_code
	ORDER BY code;
	`)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	servers := make([]gameserver.GameServer, 0)

	for rows.Next() {
		var server gameserver.GameServer

		if err := rows.Scan(
			&server.Code,
			&server.Host,
			&server.Port,
			&server.QueryHost,
			&server.QueryPort,
			&server.Kind,
			&server.Password,
			&server.DisplaysIP,
			&server.DisplaysPort,
			&server.Defaults.Application,
			&server.Defaults.Version,
			&server.Defaults.MaxPlayers,
			&server.Defaults.Title,
			&server.Defaults.Description,
		); err != nil {
			return nil, err
		}

		servers = append(servers, server)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return servers, nil
}
