package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/mdm-code/gsnip/editor"
)

func main() {
	var addr, port string
	flag.StringVar(&addr, "addr", "127.0.0.1", "server address")
	flag.StringVar(&port, "port", "7862", "server port")
	setupFlags(flag.CommandLine)
	flag.Parse()

	conn, err := net.Dial("udp", addr+":"+port)
	defer conn.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "gsnip ERROR: %s\n", err)
		os.Exit(1)
	}

	var params []string
	if isPiped() {
		s := bufio.NewScanner(os.Stdin)
		s.Split(bufio.ScanWords)
		for s.Scan() {
			params = append(params, s.Text())
		}
	}

	buf := make([]byte, 2048)
	for _, p := range params {
		// TODO: Integrate insert into server cmds
		if p == "@INSERT" {
			err := insert()
			if err != nil {
				fmt.Fprintf(os.Stderr, "gsnip ERROR: %s\n", err)
			}
			continue
		}
		_, err = fmt.Fprintf(conn, p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gsnip ERROR: %s\n", err)
		}
		n, err := bufio.NewReader(conn).Read(buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gsnip ERROR: %s\n", err)
			continue

		}
		fmt.Fprintf(os.Stdout, "%s\n", buf[:n])
	}
}

func setupFlags(f *flag.FlagSet) {
	f.Usage = func() {
		fmt.Fprintf(f.Output(), "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintf(f.Output(), "Use |, < or named pipe to send input.\n\n")
		f.PrintDefaults()
	}
}

func isPiped() bool {
	fi, _ := os.Stdin.Stat()
	return (fi.Mode() & os.ModeCharDevice) == 0
}

func insert() error {
	e, err := editor.NewEditor("nvim", nil)
	defer e.Exit()
	if err != nil {
		return err
	}
	data, err := e.Run()
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
