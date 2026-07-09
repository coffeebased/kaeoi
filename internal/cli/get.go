package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
)

func (c *CLI) get(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("get", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	fs.Usage = func() {
		_, _ = fmt.Fprintln(c.stderr, `Usage:
	servermonitor get <code>

Arguments:
  <code>  server identifier code`)
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return errors.New("get requires exactly one server identifier code")
	}

	code := fs.Arg(0)

	server, err := c.store.Get(ctx, code)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(c.stdout, "Code:         %s\n", server.Code)
	_, _ = fmt.Fprintf(c.stdout, "Join Host:    %s\n", server.Host)
	_, _ = fmt.Fprintf(c.stdout, "Join Port:    %d\n", server.Port)
	if server.QueryHost != "" {
		_, _ = fmt.Fprintf(c.stdout, "Query host:   %s\n", server.QueryHost)
	}
	if server.QueryPort != 0 {
		_, _ = fmt.Fprintf(c.stdout, "Query Port:   %d\n", server.QueryPort)
	}
	_, _ = fmt.Fprintf(c.stdout, "Kind:         %s\n", server.Kind)
	if server.Password != "" {
		_, _ = fmt.Fprintf(c.stdout, "Set password: %s\n", server.Password)
	}
	if server.DisplaysIP {
		_, _ = fmt.Fprintln(c.stdout, "-displays ip-")
	}
	if server.DisplaysPort {
		_, _ = fmt.Fprintln(c.stdout, "-displays port-")
	}
	if server.Ignore {
		_, _ = fmt.Fprintln(c.stdout, "-ignore-")
	}
	if server.OverrideDefaults {
		_, _ = fmt.Fprintln(c.stdout, "-override defaults-")
	}
	if server.Defaults.Application != "" ||
		server.Defaults.Version != "" ||
		server.Defaults.MaxPlayers != 0 ||
		server.Defaults.Title != "" ||
		server.Defaults.Description != "" {
		_, _ = fmt.Fprintln(c.stdout, "== STATUS DEFAULTS ==")
	}
	if server.Defaults.Application != "" {
		_, _ = fmt.Fprintf(c.stdout, "Application: %s\n", server.Defaults.Application)
	}
	if server.Defaults.Version != "" {
		_, _ = fmt.Fprintf(c.stdout, "Version:     %s\n", server.Defaults.Version)
	}
	if server.Defaults.MaxPlayers != 0 {
		_, _ = fmt.Fprintf(c.stdout, "Max Players: %d\n", server.Defaults.MaxPlayers)
	}
	if server.Defaults.Title != "" {
		_, _ = fmt.Fprintf(c.stdout, "Title:       %s\n", server.Defaults.Title)
	}
	if server.Defaults.Description != "" {
		_, _ = fmt.Fprintf(c.stdout, "Description: %s\n", server.Defaults.Description)
	}

	return nil
}
