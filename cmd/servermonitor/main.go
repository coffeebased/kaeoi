package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/coffeebased/kaeoi/internal/cli"
	"github.com/coffeebased/kaeoi/internal/sqlite"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	db, err := sqlite.OpenDB("servermonitordb.sqlite")
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("open sqlite db: %w", err))
		os.Exit(1)
	}

	store, err := sqlite.NewStore(db)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("server config store: %w", err))
		os.Exit(1)
	}

	cli := cli.New(store, os.Stdout, os.Stderr)

	if err := cli.Run(ctx, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("run command: %w", err))
		os.Exit(1)
	}
}
