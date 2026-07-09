// Package sqlite provides concrete implementations of stores under sqlite
package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDB(path string) (*sql.DB, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("db path is required")
	}

	dsn := fmt.Sprintf(
		"file:%s?_pragma=foreign_keys(on)&_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)",
		path,
	)

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	if err := initSchema(db); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func initSchema(db *sql.DB) error {
	const schema = `
CREATE TABLE IF NOT EXISTS game_servers (
	code TEXT NOT NULL PRIMARY KEY CHECK (trim(code) <> '' AND code = trim(code) AND code = upper(code)),
	host TEXT NOT NULL CHECK (trim(host) <> ''),
	port INTEGER NOT NULL CHECK (port BETWEEN 1 AND 65535),
	query_host TEXT NOT NULL DEFAULT ''CHECK (query_host = trim(query_host)),
	query_port INTEGER NOT NULL DEFAULT 0 CHECK (query_port BETWEEN 0 AND 65535),
	kind TEXT NOT NULL CHECK (kind IN ('none', 'generic_tcp', 'generic_udp', 'steam_a2s', 'minecraft_java')) DEFAULT 'none',
	password TEXT NOT NULL DEFAULT '',
	displays_ip INTEGER NOT NULL DEFAULT 0 CHECK (displays_ip IN (0, 1)),
	displays_port INTEGER NOT NULL DEFAULT 0 CHECK (displays_port IN (0, 1)),

	CHECK (query_host = '' OR query_port BETWEEN 1 AND 65535)
);

CREATE TABLE IF NOT EXISTS game_server_defaults (
	game_server_code TEXT NOT NULL PRIMARY KEY,
	application TEXT NOT NULL DEFAULT '',
	version TEXT NOT NULL DEFAULT '',
	max_players INTEGER NOT NULL DEFAULT 0 CHECK (max_players >= 0),
	title TEXT NOT NULL DEFAULT '',
	description TEXT NOT NULL DEFAULT '',

	FOREIGN KEY (game_server_code) REFERENCES game_servers(code) ON DELETE CASCADE
);
`

	_, err := db.Exec(schema)
	return err
}
