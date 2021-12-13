package main

import (
	"flag"
	"net"
)

func init() {
	addCmd(
		cmd{
			name:    "list",
			fn:      cmdLs,
			desc:    "gsnip\tlist\tlist all snippets",
			aliases: []string{"ls"},
		},
	)
}

func cmdLs(c net.Conn, args []string) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	err := fs.Parse(args)
	if err != nil {
		return err
	}
	args = fs.Args()
	err = transact(c, "@LST", "")
	if err != nil {
		return err
	}
	return nil
}
