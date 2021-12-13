package main

import (
	"flag"
	"net"

	"github.com/mdm-code/gsnip/editor"
)

func init() {
	addCmd(
		cmd{
			name:    "insert",
			fn:      cmdInsert,
			desc:    "gsnip\tinsert\tinsert new snippet",
			aliases: []string{"ins"},
		},
	)
}

func cmdInsert(c net.Conn, args []string) error {
	fs := flag.NewFlagSet("insert", flag.ContinueOnError)
	err := fs.Parse(args)
	if err != nil {
		return err
	}
	args = fs.Args()
	data, err := insert()
	if err != nil {
		return err
	}
	err = transact(c, "@INS", data)
	return err
}

func insert() (string, error) {
	e, err := editor.NewEditor("nvim", nil)
	defer e.Exit()
	if err != nil {
		return "", err
	}
	data, err := e.Run()
	if err != nil {
		return "", err
	}
	return string(data), nil
}
