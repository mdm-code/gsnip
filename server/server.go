package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/mdm-code/gsnip/manager"
	"github.com/mdm-code/gsnip/signals"
)

type Logger interface {
	Log(string, interface{})
}

type LogMixin struct {
	pName string
	FD    *os.File
}

func NewLogger() Logger {
	// NOTE: Add new loggers here
	switch {
	default:
		return &LogMixin{pName: "gsnipd", FD: os.Stderr}
	}
}

func (l *LogMixin) Log(level string, msg interface{}) {
	fmt.Fprintf(l.FD, "%s %s: %s\n", l.pName, level, msg)
}

type Server interface {
	Listen() error
	ShutDown()
	AwaitSignal(...os.Signal)
	AwaitConn()
}

type UDPServer struct {
	addr   net.UDPAddr
	conn   *net.UDPConn
	mngr   *manager.Manager
	sigs   chan os.Signal
	logr   Logger
	interp signals.Interpreter
	fname  string
}

func NewServer(ntwrk string, addr string, port int, fname string) (Server, error) {
	switch ntwrk {
	case "udp":
		srv, err := NewUDPServer(addr, port, fname)
		if err != nil {
			return nil, err
		}
		return srv, nil
	default:
		return nil, fmt.Errorf("unimplemented protocol: %s", ntwrk)
	}
}

func NewUDPServer(addr string, port int, fname string) (*UDPServer, error) {
	m, err := manager.NewManager(fname)
	if err != nil {
		return nil, err
	}
	return &UDPServer{
		addr: net.UDPAddr{
			IP:   net.ParseIP(addr),
			Port: port,
		},
		mngr:   m,
		sigs:   make(chan os.Signal, 1),
		logr:   NewLogger(),
		interp: signals.NewInterpreter(),
		fname:  fname,
	}, nil
}

func (s *UDPServer) Listen() (err error) {
	s.conn, err = net.ListenUDP("udp", &s.addr)
	s.logr.Log("INFO", "listening on "+s.addr.String())
	return
}

func (s *UDPServer) ShutDown() {
	s.conn.Close()
}

func (s *UDPServer) AwaitSignal(sig ...os.Signal) {
	signal.Notify(s.sigs, sig...)
	go func() {
		for {
			select {
			case <-s.sigs:
				m, err := manager.NewManager(s.fname)
				if err != nil {
					s.logr.Log("ERROR", "failed to reload source file")
					continue
				}
				s.mngr = m
				s.logr.Log("INFO", "reloaded snippet source file")
			}
		}
	}()
}

// Await for incoming connections. This is a blocking function.
func (s *UDPServer) AwaitConn() {
	for {
		buff := make([]byte, 512)
		length, respAddr, err := s.conn.ReadFromUDP(buff)
		if err != nil {
			s.logr.Log("INFO", err)
			continue
		}
		s.logr.Log("INFO", fmt.Sprintf("read %s from %v", buff, respAddr))
		go s.respond(respAddr, buff[:length])
	}
}

func (s *UDPServer) respond(addr *net.UDPAddr, buff []byte) {
	token := s.interp.Eval(string(buff))
	switch token.IsReload() {
	case true:
		s.sigs <- syscall.SIGHUP
		_, err := s.conn.WriteToUDP([]byte(""), addr)
		if err != nil {
			s.logr.Log("ERROR", err)
			return
		}
		return
	default:
		resp, err := s.mngr.Execute(token)
		if err != nil {
			s.logr.Log("ERROR", err)
			_, err = s.conn.WriteToUDP([]byte("ERROR"), addr)
			if err != nil {
				s.logr.Log("ERROR", err)
				return
			}
			return
		}
		outMsg := []byte(resp)
		_, err = s.conn.WriteToUDP(outMsg, addr)
		if err != nil {
			s.logr.Log("ERROR", err)
			return
		}
		s.logr.Log("INFO", "write successful")
	}
}
