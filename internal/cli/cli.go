// Package cli exposes CRUD cli commands to manage game server's stored data
package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/coffeebased/kaeoi/internal/gameserver"
)

type CLI struct {
	store  gameserver.Store
	stdout io.Writer
	stderr io.Writer
}

func New(store gameserver.Store, stdout, stderr io.Writer) *CLI {
	if store == nil {
		panic("store is required")
	}

	if stdout == nil {
		panic("stdout is required")
	}

	if stderr == nil {
		panic("stderr is required")
	}

	return &CLI{
		store:  store,
		stdout: stdout,
		stderr: stderr,
	}
}

func (c *CLI) Run(ctx context.Context, args []string) error {
	if ctx == nil {
		panic("ctx is required")
	}

	if len(args) == 0 {
		return c.help()
	}

	switch args[0] {
	case "list":
		return c.list(ctx, args[1:])
	case "get":
		return c.get(ctx, args[1:])
	case "add":
		return c.add(ctx, args[1:])
	case "update":
		return c.update(ctx, args[1:])
	case "remove":
		return c.remove(ctx, args[1:])
	case "help", "-h", "--help":
		return c.help()
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func (c *CLI) help() error {
	_, err := fmt.Fprintln(c.stdout, `Usage:
	servermonitor list
	servermonitor get <code>
	servermonitor add [flags] <code>
	servermonitor update [flags] <code>
	servermonitor remove <code>`)
	return err
}
