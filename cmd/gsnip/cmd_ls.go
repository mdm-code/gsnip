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
			descr:   "gsnip\tlist\tlist all snippets",
			aliases: []string{"ls"},
		},
	)
}

func cmdLs(c net.Conn, args []string) error {
	fs := flag.NewFlagSet("ls", flag.ContinueOnError)
	err := fs.Parse(args)
	if err != nil {
		return err
	}
	args = fs.Args()
	err = transact(c, "@LIST")
	if err != nil {
		return err
	}
	return nil
}
