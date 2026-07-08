package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/coffeebased/kaeoi/internal/gameserver"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) (*Store, error) {
	if db == nil {
		return nil, errors.New("db is required")
	}

	return &Store{
		db: db,
	}, nil
}

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
		gsd.description,
	FROM game_servers AS gs
	INNER JOIN game_server_defaults AS gsd
	ON gs.code = gsd.code
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

func (s *Store) Get(ctx context.Context, code string) (gameserver.GameServer, error) {
	if ctx == nil {
		panic("nil context")
	}

	if strings.TrimSpace(code) == "" {
		return gameserver.GameServer{}, errors.New("code cannot be empty or white space")
	}

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

func (s *Store) Create(ctx context.Context, request gameserver.CreateRequest) error {
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

func (s *Store) Update(ctx context.Context, code string, request gameserver.UpdateRequest) error {
	if ctx == nil {
		panic("nil context")
	}

	if strings.TrimSpace(code) == "" {
		return errors.New("code cannot be empty or white space")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if request.Host != nil {
		if _, err := tx.ExecContext(ctx, `
		UPDATE game_servers
		SET host = ?
		WHERE code = ?
		`, *request.Host, code); err != nil {
			return err
		}
	}

	if request.Port != nil {
		if _, err := tx.ExecContext(ctx, `
		UPDATE game_servers
		SET port = ?
		WHERE code = ?
		`, *request.Port, code); err != nil {
			return err
		}
	}

	if request.QueryHost != nil {
		if _, err := tx.ExecContext(ctx, `
		UPDATE game_servers
		SET query_host = ?
		WHERE code = ?
		`, *request.QueryHost, code); err != nil {
			return err
		}
	}

	if request.QueryPort != nil {
		if _, err := tx.ExecContext(ctx, `
		UPDATE game_servers
		SET query_port = ?
		WHERE code = ?
		`, *request.QueryPort, code); err != nil {
			return err
		}
	}

	if request.Kind != nil {
		if _, err := tx.ExecContext(ctx, `
		UPDATE game_servers
		SET kind = ?
		WHERE code = ?
		`, *request.Kind, code); err != nil {
			return err
		}
	}

	if request.Password != nil {
		if _, err := tx.ExecContext(ctx, `
		UPDATE game_servers
		SET password = ?
		WHERE code = ?
		`, *request.Password, code); err != nil {
			return err
		}
	}

	if request.DisplaysIP != nil {
		if _, err := tx.ExecContext(ctx, `
		UPDATE game_servers
		SET displays_ip = ?
		WHERE code = ?
		`, *request.DisplaysIP, code); err != nil {
			return err
		}
	}

	if request.DisplaysPort != nil {
		if _, err := tx.ExecContext(ctx, `
		UPDATE game_servers
		SET displays_port = ?
		WHERE code = ?
		`, *request.DisplaysPort, code); err != nil {
			return err
		}
	}

	if request.Application != nil {
		if _, err := tx.ExecContext(ctx, `
		UPDATE game_server_defaults
		SET application = ?
		WHERE game_server_code = ?
		`, *request.Application, code); err != nil {
			return err
		}
	}

	if request.Version != nil {
		if _, err := tx.ExecContext(ctx, `
		UPDATE game_server_defaults
		SET version = ?
		WHERE game_server_code = ?
		`, *request.Version, code); err != nil {
			return err
		}
	}

	if request.MaxPlayers != nil {
		if _, err := tx.ExecContext(ctx, `
		UPDATE game_server_defaults
		SET max_players = ?
		WHERE game_server_code = ?
		`, *request.MaxPlayers, code); err != nil {
			return err
		}
	}

	if request.Title != nil {
		if _, err := tx.ExecContext(ctx, `
		UPDATE game_server_defaults
		SET title = ?
		WHERE game_server_code = ?
		`, *request.Title, code); err != nil {
			return err
		}
	}

	if request.Description != nil {
		if _, err := tx.ExecContext(ctx, `
		UPDATE game_server_defaults
		SET description = ?
		WHERE game_server_code = ?
		`, *request.Description, code); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Store) Delete(ctx context.Context, code string) error {
	if ctx == nil {
		panic("nil context")
	}

	if strings.TrimSpace(code) == "" {
		return errors.New("code cannot be empty or white space")
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
		DELETE FROM game_servers
		WHERE code = ?
		`,
		code,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}
