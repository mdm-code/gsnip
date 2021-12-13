package main

import (
	"bufio"
	"flag"
	"net"
	"os"
)

func init() {
	addCmd(
		cmd{
			name:    "delete",
			fn:      cmdDel,
			desc:    "gsnip\tdelete\tdelete a snippet",
			aliases: []string{"d"},
		},
	)
}

func cmdDel(c net.Conn, args []string) error {
	fs := flag.NewFlagSet("delete", flag.ContinueOnError)
	err := fs.Parse(args)
	if err != nil {
		return err
	}
	args = fs.Args()

	var names []string
	if isPiped() {
		s := bufio.NewScanner(os.Stdin)
		s.Split(bufio.ScanWords)
		for s.Scan() {
			names = append(names, s.Text())
		}
	}
	for _, n := range names {
		err := transact(c, "@DEL", n)
		if err != nil {
			return err
		}
	}
	return nil
}
