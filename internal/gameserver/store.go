package gameserver

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("game server not found")

type Store interface {
	List(ctx context.Context) ([]GameServer, error)
	Get(ctx context.Context, code string) (GameServer, error)
	Create(ctx context.Context, server Create) error
	Update(ctx context.Context, code string, patch Update) error
	Delete(ctx context.Context, code string) error
}
