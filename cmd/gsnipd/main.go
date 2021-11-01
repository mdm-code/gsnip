package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"

	"github.com/mdm-code/gsnip/server"
)

const source = "/usr/local/share/gsnip/snippets"

var (
	port int
	addr string
	file string
	prot string
	sigs chan os.Signal
)

func main() {
	flag.IntVar(&port, "port", 7862, "UDP server port")
	flag.StringVar(&addr, "addr", "127.0.0.1", "UDP server IP address")
	flag.StringVar(&file, "file", source, "snippet source file")
	flag.StringVar(&prot, "protocol", "udp", "internet protocol")
	flag.Parse()

	s, err := server.NewServer("udp", addr, port, file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gsnipd ERROR: %s", err)
		os.Exit(1)
	}
	s.Listen()
	defer s.ShutDown()
	s.AwaitSignal(syscall.SIGHUP)
	s.AwaitConn()
}
