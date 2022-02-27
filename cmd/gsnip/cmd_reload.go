package main

import (
	"flag"

	"github.com/mdm-code/gsnip/internal/stream"
)

func init() {
	addCmd(
		cmd{
			name:    "reload",
			fn:      cmdReload,
			desc:    "gsnip\treload\treload snippet source file",
			aliases: []string{"rld", "rl"},
		},
	)
}

func cmdReload(args []string) error {
	fs := flag.NewFlagSet("reload", flag.ContinueOnError)
	err := fs.Parse(args)
	if err != nil {
		return err
	}
	args = fs.Args()
	err = transact(stream.Reload, []byte{})
	if err != nil {
		return err
	}
	return nil
}
