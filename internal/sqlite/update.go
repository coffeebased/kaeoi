package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/coffeebased/kaeoi/internal/gameserver"
)

func (s *Store) Update(ctx context.Context, code string, request gameserver.Update) error {
	if ctx == nil {
		panic("nil context")
	}

	if strings.TrimSpace(code) == "" {
		return errors.New("code cannot be empty or white space")
	}

	code = strings.ToUpper(code)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var exists int
	row := tx.QueryRowContext(ctx, "SELECT 1 FROM game_servers WHERE code = ?", code)
	if err := row.Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gameserver.ErrNotFound
		}
		return err
	}

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

	if request.Ignore != nil {
		if _, err := tx.ExecContext(ctx, `
		UPDATE game_servers
		SET ignore = ?
		WHERE code = ?
		`, *request.Ignore, code); err != nil {
			return err
		}
	}

	if request.OverrideDefaults != nil {
		if _, err := tx.ExecContext(ctx, `
		UPDATE game_servers
		SET override_defaults = ?
		WHERE code = ?
		`, *request.OverrideDefaults, code); err != nil {
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
