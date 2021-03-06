package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/mdm-code/gsnip/internal/server"
	"github.com/mdm-code/xdg"
)

var (
	sock string
	file string
)

func main() {
	flag.StringVar(
		&sock,
		"sock",
		"/tmp/gsnip.sock",
		"UDS server socket name",
	)
	flag.StringVar(&file, "file", "", "snippet source file")
	setupFlags(flag.CommandLine)
	flag.Parse()

	if file == "" {
		var ok bool
		file, ok = xdg.Find(xdg.Data, "gsnip/snippets")
		if !ok {
			fmt.Fprintln(os.Stderr, "gsnipd ERROR: could not find snippet file")
			os.Exit(1)
		}
	}

	cleanup()
	defer cleanup()

	s, err := server.NewServer("unix", sock, file)
	defer s.ShutDown()

	if err != nil {
		fmt.Fprintf(os.Stderr, "gsnipd ERROR: %s", err)
		os.Exit(1)
	}

	s.Log("INFO", fmt.Sprintf("reading source file: %s", file))
	err = s.Listen()
	if err != nil {
		s.Log(
			"ERROR",
			fmt.Sprintf("UDS socket file taken: %s", sock),
		)
		os.Exit(2)
	}
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

func cleanup() {
	if _, err := os.Stat(sock); err == nil {
		if err := os.RemoveAll(sock); err != nil {
			log.Fatal(err)
		}
	}
}
