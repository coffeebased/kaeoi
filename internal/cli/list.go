package cli

import (
	"context"
	"flag"
	"fmt"
	"text/tabwriter"
)

func (c *CLI) list(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	fs.Usage = func() {
		_, _ = fmt.Fprintln(c.stderr, `Usage:
	servermonitor list`)
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() != 0 {
		return fmt.Errorf("unexpected argument: %s", fs.Arg(0))
	}

	servers, err := c.store.List(ctx)
	if err != nil {
		return err
	}

	if len(servers) == 0 {
		_, err := fmt.Fprintln(c.stdout, "there are no stored servers")
		return err
	}

	w := tabwriter.NewWriter(c.stdout, 0, 0, 2, ' ', 0)

	_, _ = fmt.Fprintln(w, "CODE\tJOIN\tQUERY\tKIND")
	for _, server := range servers {
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			server.Code,
			formatAddress(server.Host, server.Port),
			formatAddress(server.QueryHost, server.QueryPort),
			server.Kind,
		)
	}

	return w.Flush()
}

func formatAddress(address string, port int) string {
	if address == "" && port == 0 {
		return "-"
	}

	return fmt.Sprintf("%s:%d", address, port)
}
