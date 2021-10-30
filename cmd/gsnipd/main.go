package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/mdm-code/gsnip/manager"
	"github.com/mdm-code/gsnip/parsing"
	"github.com/mdm-code/gsnip/snippets"
)

const src = "/usr/local/share/gsnip/snippets"

var (
	port int
	addr string
	file string
)

func main() {
	flag.IntVar(&port, "port", 7862, "UDP server port")
	flag.StringVar(&addr, "addr", "127.0.0.1", "UDP server IP address")
	flag.StringVar(&file, "file", src, "snippet source file")
	flag.Parse()

	addr := net.UDPAddr{
		IP:   net.ParseIP(addr),
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

	// Hot reload config file on SIGHUP.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP)
	go func() {
		for {
			select {
			case <-sigs:
				msgf("INFO", fmt.Sprintf("reload snippet source file"))
				mgr = getSnippetContainer()
			}
		}
	}()

	for {
		buf := make([]byte, 512)
		length, rAddr, err := conn.ReadFromUDP(buf)
		fmt.Fprintf(
			os.Stdout,
			msgf("READ", fmt.Sprintf("%s FROM %v", buf, rAddr)),
		)
		if err != nil {
			msgf("ERROR", err)
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
	f, err := os.Open(file)
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
