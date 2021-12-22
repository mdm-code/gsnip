package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/mdm-code/gsnip/fs"
	"github.com/mdm-code/gsnip/manager"
	"github.com/mdm-code/gsnip/stream"
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
	Log(string, interface{})
}

type UDPServer struct {
	addr net.UDPAddr
	conn *net.UDPConn
	mngr *manager.Manager
	sigs chan os.Signal
	logr Logger
	itrp stream.Interpreter
	fh   *fs.FileHandler
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
	fh, err := fs.NewFileHandler(fname, fs.Perm)
	if err != nil {
		return nil, err
	}
	m, err := manager.NewManager(fh)
	if err != nil {
		return nil, err
	}
	return &UDPServer{
		addr: net.UDPAddr{
			IP:   net.ParseIP(addr),
			Port: port,
		},
		mngr: m,
		sigs: make(chan os.Signal, 1),
		logr: NewLogger(),
		itrp: stream.NewInterpreter(),
		fh:   fh,
	}, nil
}

func (s *UDPServer) Listen() (err error) {
	s.conn, err = net.ListenUDP("udp", &s.addr)
	s.Log("INFO", "listening on "+s.addr.String())
	return
}

func (s *UDPServer) ShutDown() {
	// NOTE: file handler closes down the moment the server is closed
	s.fh.Close()
	s.conn.Close()
}

func (s *UDPServer) AwaitSignal(sig ...os.Signal) {
	signal.Notify(s.sigs, sig...)
	go func() {
		for {
			select {
			case <-s.sigs:
				err := s.mngr.Reload()
				if err != nil {
					s.Log("ERROR", err)
					continue
				}
				s.Log("INFO", "reloaded snippet source file")
			}
		}
	}()
}

// Await for incoming connections. This is a blocking function.
func (s *UDPServer) AwaitConn() {
	for {
		buff := make([]byte, 2048)
		length, respAddr, err := s.conn.ReadFromUDP(buff)
		if err != nil {
			s.Log("INFO", err)
			continue
		}
		s.Log("INFO", fmt.Sprintf("read %s from %v", buff, respAddr))
		go s.respond(respAddr, buff[:length])
	}
}

func (s *UDPServer) respond(addr *net.UDPAddr, buff []byte) {
	msg := s.itrp.Eval(buff)
	switch msg.T() {
	case stream.Rld:
		s.sigs <- syscall.SIGHUP
		_, err := s.conn.WriteToUDP([]byte(""), addr)
		if err != nil {
			s.Log("ERROR", err)
			return
		}
		return
	default:
		resp, err := s.mngr.Execute(msg)
		if err != nil {
			s.Log("ERROR", err)
			_, err = s.conn.WriteToUDP([]byte("ERROR"), addr)
			if err != nil {
				s.Log("ERROR", err)
				return
			}
			return
		}
		outMsg := []byte(resp)
		_, err = s.conn.WriteToUDP(outMsg, addr)
		if err != nil {
			s.Log("ERROR", err)
			return
		}
		s.Log("INFO", "write successful")
	}
}

func (s *UDPServer) Log(level string, msg interface{}) {
	s.logr.Log(level, msg)
}
