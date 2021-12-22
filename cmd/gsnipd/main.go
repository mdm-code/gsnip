package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"syscall"

	"github.com/mdm-code/gsnip/server"
	"github.com/mdm-code/gsnip/xdg"
)

var (
	port int
	addr string
	file string
)

func main() {
	flag.IntVar(&port, "port", 7862, "UDP server port")
	flag.StringVar(&addr, "addr", "127.0.0.1", "UDP server IP address")
	flag.StringVar(&file, "file", "", "snippet source file")
	setupFlags(flag.CommandLine)
	flag.Parse()

	if file == "" {
		dirs := xdg.Arrange()
		dir, ok := xdg.Discover(dirs)
		if !ok {
			fmt.Fprintf(os.Stderr, "gsnipd ERROR: could not find any snippet file")
			os.Exit(1)
		}
		file = path.Join(dir.Item(), "snippets")
	}

	s, err := server.NewServer("udp", addr, port, file)

	if err != nil {
		fmt.Fprintf(os.Stderr, "gsnipd ERROR: %s", err)
		os.Exit(1)
	}

	s.Log("INFO", "reading source file: "+file)
	s.Listen()
	defer s.ShutDown()
	s.AwaitSignal(syscall.SIGHUP)
	s.AwaitConn()
}

func setupFlags(f *flag.FlagSet) {
	f.Usage = func() {
		fmt.Fprintf(f.Output(), "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintf(f.Output(), "Start the snippet server.\n\n")
		f.PrintDefaults()
	}
}
