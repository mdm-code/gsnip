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
			desc:    "reload snippet source file",
			aliases: []string{"r", "rld"},
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
