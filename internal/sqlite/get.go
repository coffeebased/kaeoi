package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/coffeebased/kaeoi/internal/gameserver"
)

func (s *Store) Get(ctx context.Context, code string) (gameserver.GameServer, error) {
	if ctx == nil {
		panic("nil context")
	}

	if strings.TrimSpace(code) == "" {
		return gameserver.GameServer{}, errors.New("code cannot be empty or white space")
	}

	code = strings.ToUpper(code)

	row := s.db.QueryRowContext(ctx, `
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
	WHERE gs.code = ?
	`, code)

	var server gameserver.GameServer

	err := row.Scan(
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
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gameserver.GameServer{}, gameserver.ErrNotFound
		}
		return gameserver.GameServer{}, err
	}

	return server, nil
}
