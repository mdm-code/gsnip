package main

import (
	"bufio"
	"flag"
	"os"
	"strings"

	"github.com/mdm-code/gsnip/internal/editor"
	"github.com/mdm-code/gsnip/internal/stream"
)

func init() {
	addCmd(
		cmd{
			name:    "insert",
			fn:      cmdInsert,
			desc:    "insert a new snippet",
			aliases: []string{"i", "in", "ins"},
		},
	)
}

func cmdInsert(args []string) error {
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
	err = transact(stream.Insert, []byte(data))
	return err
}

func insert() (string, error) {
	if isPiped() {
		var lines []string

		s := bufio.NewScanner(os.Stdin)
		s.Split(bufio.ScanLines)
		for s.Scan() {
			lines = append(lines, s.Text())
		}

		return strings.Join(lines, "\n"), nil
	}
	e, err := editor.NewEditor(nil)
	if err != nil {
		return "", err
	}
	defer e.Exit()

	data, err := e.Run()
	if err != nil {
		return "", err
	}

	return string(data), nil
}
