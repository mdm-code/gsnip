package main

import (
	"flag"
	"fmt"
	"io"
	"net/rpc/jsonrpc"
	"os"
	"strings"

	"github.com/mdm-code/gsnip/internal/stream"
)

var sock string

var cmdList []cmd

var cmdMap = make(map[string]cmd)

type cmd struct {
	name    string
	fn      func([]string) error
	desc    string
	aliases []string
}

func (c cmd) String() string {
	result := fmt.Sprintf(
		"%-10s %-10s %-10s %-10s",
		"gsnip",
		c.name,
		"["+strings.Join(c.aliases, "|")+"]",
		c.desc,
	)
	return result
}

func addCmd(c cmd) {
	cmdList = append(cmdList, c)
	cmdMap[c.name] = c
	for _, a := range c.aliases {
		cmdMap[a] = c
	}
}

func main() {
	args, err := parseArgs()
	if err != nil {
		os.Exit(1)
	}

	err = dispatchCmd(args)
	if err != nil && err != io.EOF {
		fmt.Fprintf(os.Stderr, "gsnip ERROR: %s\n", err)
		os.Exit(1)
	}
}

func parseArgs() ([]string, error) {
	fs := flag.NewFlagSet("gsnip", flag.ContinueOnError)
	fs.StringVar(&sock, "sock", "/tmp/gsnip.sock", "UDS server socket name")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Global options:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nCommands:\n")
		for _, c := range cmdList {
			fmt.Fprintf(os.Stderr, "%s\n", c)
		}
	}

	err := fs.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}
	args := fs.Args()
	return args, nil
}

func dispatchCmd(args []string) error {
	if len(args) < 1 {
		return nil
	}
	if cmd, ok := cmdMap[args[0]]; ok {
		err := cmd.fn(args[1:])
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("command not found: %s", args[0])
}

func transact(op stream.Opcode, data []byte) error {
	conn, err := jsonrpc.Dial("unix", sock)
	if err != nil {
		return err
	}
	defer conn.Close()

	request := stream.Request{Operation: op, Body: data}
	var reply stream.Reply

	err = conn.Call("Manager.Execute", request, &reply)
	if err != nil {
		return err
	}

	if reply.Result == stream.Failure {
		fmt.Fprintf(os.Stderr, "%s\n", "ERROR")
		return nil
	}

	fmt.Fprintf(os.Stdout, "%s\n", reply.Body)
	return nil
}

func isPiped() bool {
	fi, _ := os.Stdin.Stat()
	return (fi.Mode() & os.ModeCharDevice) == 0
}
