package sqlite

import (
	"context"
	"errors"
	"strings"

	"github.com/coffeebased/kaeoi/internal/gameserver"
)

func (s *Store) Delete(ctx context.Context, code string) error {
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

	result, err := tx.ExecContext(
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

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return gameserver.ErrNotFound
	}

	return tx.Commit()
}
