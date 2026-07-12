package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/coffeebased/kaeoi/internal/gameserver"
)

func (c *CLI) update(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("update", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	var request gameserver.Update

	var host string
	var port int
	var queryHost string
	var queryPort int
	var kind string
	var password string
	var displaysIP bool
	var displaysPort bool
	var overrideDefaults bool
	var ignore bool
	var application string
	var version string
	var maxPlayers int
	var title string
	var description string

	fs.StringVar(&host, "host", "", "server join host")
	fs.IntVar(&port, "port", 0, "server join port")
	fs.StringVar(&queryHost, "query-host", "", "server query host")
	fs.IntVar(&queryPort, "query-port", 0, "server query port")
	fs.StringVar(&kind, "kind", "", "server kind")
	fs.StringVar(&password, "password", "", "server join password")
	fs.BoolVar(&displaysIP, "displays-ip", false, "display IP instead of host on server address")
	fs.BoolVar(&displaysPort, "displays-port", false, "display port on server address")
	fs.BoolVar(&ignore, "ignore", false, "exclude from monitoring process")
	fs.BoolVar(&overrideDefaults, "override-defaults", false, "do not override populated defaults with state")
	fs.StringVar(&application, "application", "", "application name")
	fs.StringVar(&version, "version", "", "application version")
	fs.IntVar(&maxPlayers, "max-players", 0, "maximum number of players")
	fs.StringVar(&title, "title", "", "server title")
	fs.StringVar(&description, "description", "", "server description text")

	fs.Usage = func() {
		_, _ = fmt.Fprintln(c.stderr, `Usage:
	servermonitor update [flags] <code>

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
		return fmt.Errorf("update expects exactly one server identifier code argument")
	}

	code := fs.Arg(0)

	if fs.NFlag() == 0 {
		return errors.New("update requires at least one flag")
	}

	fs.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "host":
			request.Host = &host
		case "port":
			request.Port = &port
		case "query-host":
			request.QueryHost = &queryHost
		case "query-port":
			request.QueryPort = &queryPort
		case "kind":
			k := gameserver.Kind(kind)
			request.Kind = &k
		case "password":
			request.Password = &password
		case "displays-ip":
			request.DisplaysIP = &displaysIP
		case "displays-port":
			request.DisplaysPort = &displaysPort
		case "ignore":
			request.Ignore = &ignore
		case "override-defaults":
			request.OverrideDefaults = &overrideDefaults
		case "application":
			request.Application = &application
		case "version":
			request.Version = &version
		case "max-players":
			request.MaxPlayers = &maxPlayers
		case "title":
			request.Title = &title
		case "description":
			request.Description = &description
		}
	})

	request.Normalize()

	if err := request.Validate(); err != nil {
		return err
	}

	if err := c.store.Update(ctx, code, request); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(c.stdout, "updated server: %s\n", code)
	_, err := fmt.Fprintln(c.stdout, "daemon restart needed to apply changes")
	return err
}
