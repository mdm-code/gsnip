package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	var addr, port string
	flag.StringVar(&addr, "addr", "127.0.0.1", "server address")
	flag.StringVar(&port, "port", "7862", "server port")
	flag.Parse()
	conn, err := net.Dial("udp", addr+":"+port)
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

func isPiped() bool {
	fi, _ := os.Stdin.Stat()
	return (fi.Mode() & os.ModeCharDevice) == 0

}
