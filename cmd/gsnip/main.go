package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var (
	sock string
	conn net.Conn
)

var errPrefix string = "ERROR"

type cmd struct {
	name    string
	fn      func(net.Conn, []string) error
	desc    string
	aliases []string
}

var cmdList []cmd

var cmdMap = make(map[string]cmd)

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

	conn, err = net.Dial("unix", sock)
	defer conn.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "gsnip ERROR: %s\n", err)
		os.Exit(1)
	}

	err = dispatchCmd(conn, args)
	if err != nil {
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
			fmt.Fprintf(os.Stderr, c.desc+"\n")
		}
	}

	err := fs.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}
	args := fs.Args()
	return args, nil
}

func dispatchCmd(c net.Conn, args []string) error {
	if len(args) < 1 {
		return nil
	}
	if cmd, ok := cmdMap[args[0]]; ok {
		err := cmd.fn(c, args[1:])
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("command not found: %s", args[0])
}

func transact(c net.Conn, kind string, data string) error {
	buf := make([]byte, 4096)
	_, err := fmt.Fprintf(conn, kind+" "+data)
	if err != nil {
		return err
	}
	n, err := bufio.NewReader(conn).Read(buf)
	if err != nil {
		return err
	}
	if strings.HasPrefix(string(buf[:n]), errPrefix) {
		fmt.Fprintf(os.Stderr, "%s\n", buf[:n])
		return nil
	}
	fmt.Fprintf(os.Stdout, "%s\n", buf[:n])
	return nil
}

func isPiped() bool {
	fi, _ := os.Stdin.Stat()
	return (fi.Mode() & os.ModeCharDevice) == 0
}
