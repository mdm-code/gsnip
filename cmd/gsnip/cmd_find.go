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
			name:    "find",
			fn:      cmdFind,
			desc:    "find a snippet",
			aliases: []string{"f", "fnd"},
		},
	)
}

func cmdFind(args []string) error {
	fs := flag.NewFlagSet("find", flag.ContinueOnError)
	err := fs.Parse(args)
	if err != nil {
		return err
	}
	args = fs.Args()

	var params []string
	if isPiped() {
		s := bufio.NewScanner(os.Stdin)
		s.Split(bufio.ScanWords)
		for s.Scan() {
			params = append(params, s.Text())
		}
	}
	for _, p := range params {
		err := transact(stream.Find, []byte(p))
		if err != nil && err != io.EOF {
			return err
		}
	}
	return nil
}
