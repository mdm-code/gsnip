package main

import (
	"flag"
	"net"
)

func init() {
	addCmd(
		cmd{
			name:    "reload",
			fn:      cmdReload,
			descr:   "gsnip\treload\treload snippet source file",
			aliases: []string{"rld", "rl"},
		},
	)
}

func cmdReload(c net.Conn, args []string) error {
	fs := flag.NewFlagSet("reload", flag.ContinueOnError)
	err := fs.Parse(args)
	if err != nil {
		return err
	}
	args = fs.Args()
	err = transact(c, "@RELOAD")
	if err != nil {
		return err
	}
	return nil
}
