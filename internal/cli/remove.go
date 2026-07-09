package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
)

func (c *CLI) remove(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("remove", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	fs.Usage = func() {
		_, _ = fmt.Fprintln(c.stderr, `Usage:
	servermonitor remove <code>

Arguments:
  <code>  server identifier code`)
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return errors.New("remove requires exactly one server identifier code")
	}

	code := fs.Arg(0)

	if err := c.store.Delete(ctx, code); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(c.stdout, "removed server: %s\n", code)
	_, err := fmt.Fprintln(c.stdout, "daemon restart needed to apply changes")
	return err
}
