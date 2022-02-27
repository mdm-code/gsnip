package main

import (
	"bufio"
	"flag"
	"io"
	"os"

	"github.com/mdm-code/gsnip/internal/stream"
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

func cmdDel(args []string) error {
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
	for _, name := range names {
		err := transact(stream.Delete, []byte(name))
		if err != nil && err != io.EOF {
			return err
		}
	}
	return nil
}
