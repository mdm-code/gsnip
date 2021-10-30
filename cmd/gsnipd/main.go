package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/mdm-code/gsnip/manager"
	"github.com/mdm-code/gsnip/parsing"
	"github.com/mdm-code/gsnip/snippets"
)

const src = "/usr/local/share/gsnip/snippets"

var (
	port int
	file string
)

func main() {
	flag.IntVar(&port, "port", 7862, "UDP server port")
	addr := net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: port,
	}
	conn, err := net.ListenUDP("udp", &addr)
	defer conn.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, msgf("ERROR", err))
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, msgf("INFO", fmt.Sprintf("running %s", &addr)))

	mgr := getSnippetContainer()

	for {
		buf := make([]byte, 512)
		length, rAddr, err := conn.ReadFromUDP(buf)
		fmt.Fprintf(os.Stdout, "gsnipd READ FROM %v: %s\n", rAddr, buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gsnipd ERROR: %s\n", err)
			continue
		}
		go respond(conn, rAddr, &mgr, buf[:length])
	}
}

func msgf(class, msg interface{}) string {
	return fmt.Sprintf("gsnipd %s: %s\n", class, msg)
}

func getSnippetContainer() manager.Manager {
	var snpts snippets.Container
	var f *os.File
	f, err := os.Open(src)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			msgf("ERROR", fmt.Sprintf("unable to open %s", src)),
		)
		os.Exit(1)
		f, err = os.Open(src)
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				msgf("ERROR", fmt.Sprintf("could not find %s", src)),
			)
			os.Exit(1)
		}
	}
	defer f.Close()
	fmt.Fprintf(os.Stdout, msgf("INFO", fmt.Sprintf("reading %s", src)))

	parser := parsing.NewParser()
	snpts, err = parser.Parse(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, msgf("ERROR", err))
		os.Exit(1)
	}

	mgr, ok := manager.NewManager(snpts)
	if !ok {
		fmt.Fprint(os.Stderr, msgf("ERROR", "failed to start snippet manager"))
		os.Exit(1)
	}
	return mgr
}

func respond(conn *net.UDPConn, addr *net.UDPAddr, m *manager.Manager, buf []byte) {
	output, err := m.Execute(string(buf))
	if err != nil {
		fmt.Fprintf(os.Stderr, msgf("ERROR", err))
		_, err = conn.WriteToUDP([]byte("ERROR"), addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, msgf("ERROR", err))
			return
		}
		return
	}
	msg := []byte(output)
	_, err = conn.WriteToUDP(msg, addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, msgf("ERROR", err))
		return
	}
	fmt.Fprintf(os.Stdout, msgf("WRITE", "success"))
}
