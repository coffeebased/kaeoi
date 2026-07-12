package cli

import (
	"context"
	"flag"
	"fmt"

	"github.com/coffeebased/kaeoi/internal/gameserver"
)

func (c *CLI) add(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	var request gameserver.Create
	var kind string

	fs.StringVar(&request.Host, "host", "", "server join host")
	fs.IntVar(&request.Port, "port", 0, "server join port")
	fs.StringVar(&request.QueryHost, "query-host", "", "server query host")
	fs.IntVar(&request.QueryPort, "query-port", 0, "server query port")
	fs.StringVar(&kind, "kind", "", "server kind")
	fs.StringVar(&request.Password, "password", "", "server join password")
	fs.BoolVar(&request.DisplaysIP, "displays-ip", false, "display IP instead of host on server address")
	fs.BoolVar(&request.DisplaysPort, "displays-port", false, "display port on server address")
	fs.BoolVar(&request.Ignore, "ignore", false, "exclude from monitoring process")
	fs.BoolVar(&request.OverrideDefaults, "override-defaults", false, "do not override populated defaults with state")
	fs.StringVar(&request.Application, "application", "", "application name")
	fs.StringVar(&request.Version, "version", "", "application version")
	fs.IntVar(&request.MaxPlayers, "max-players", 0, "maximum number of players")
	fs.StringVar(&request.Title, "title", "", "server title")
	fs.StringVar(&request.Description, "description", "", "server description text")

	fs.Usage = func() {
		_, _ = fmt.Fprintln(c.stderr, `Usage:
  servermonitor add [flags] <code>

Arguments:
  <code>  server identifier code

Flags:`)
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		fs.Usage()
		return fmt.Errorf("add expects exactly one server identifier code argument")
	}

	request.Code = fs.Arg(0)
	request.Kind = gameserver.Kind(kind)

	request.Normalize()

	if err := request.Validate(); err != nil {
		return err
	}

	if err := c.store.Create(ctx, request); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(c.stdout, "added server: %s\n", request.Code)
	_, err := fmt.Fprintln(c.stdout, "daemon restart needed to apply changes")
	return err
}
