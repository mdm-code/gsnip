package main

import (
	"flag"

	"github.com/mdm-code/gsnip/internal/stream"
)

func init() {
	addCmd(
		cmd{
			name:    "list",
			fn:      cmdLs,
			desc:    "list all snippets",
			aliases: []string{"l", "ls"},
		},
	)
}

func cmdLs(args []string) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	err := fs.Parse(args)
	if err != nil {
		return err
	}
	args = fs.Args()
	err = transact(stream.List, []byte{})
	if err != nil {
		return err
	}
	return nil
}
