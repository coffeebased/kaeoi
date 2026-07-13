package sqlite

import (
	"database/sql"
	"errors"
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
