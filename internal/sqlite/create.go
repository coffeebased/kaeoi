package sqlite

import (
	"context"
	"errors"

	"github.com/coffeebased/kaeoi/internal/gameserver"
)

func (s *Store) Create(ctx context.Context, request gameserver.Create) error {
	if ctx == nil {
		return errors.New("ctx is required")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.ExecContext(
		ctx,
		`
		INSERT INTO game_servers (
			code,
			host,
			port,
			query_host,
			query_port,
			kind,
			password,
			displays_ip,
			displays_port
		) VALUES (?,?,?,?,?,?,?,?,?)
		`,
		request.Code,
		request.Host,
		request.Port,
		request.QueryHost,
		request.QueryPort,
		request.Kind,
		request.Password,
		request.DisplaysIP,
		request.DisplaysPort,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		`
		INSERT INTO game_server_defaults (
			game_server_code,
			application,
			version,
			max_players,
			title,
			description
		) VALUES (?,?,?,?,?,?)
		`,
		request.Code,
		request.Application,
		request.Version,
		request.MaxPlayers,
		request.Title,
		request.Description,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}
