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

type UnixServer struct {
	socket      string
	listener    net.Listener
	manager     *manager.Manager
	signals     chan os.Signal
	logger      Logger
	interpreter stream.Interpreter
	fileHandler *fs.FileHandler
}

func NewServer(ntwrk string, addr string, fname string) (Server, error) {
	switch ntwrk {
	case "unix":
		srv, err := NewUnixServer(addr, fname)
		if err != nil {
			return nil, err
		}
		return srv, nil
	default:
		return nil, fmt.Errorf("unimplemented protocol: %s", ntwrk)
	}
}

func NewUnixServer(sock string, fname string) (*UnixServer, error) {
	fh, err := fs.NewFileHandler(fname, fs.Perm)
	if err != nil {
		return nil, err
	}
	m, err := manager.NewManager(fh)
	if err != nil {
		return nil, err
	}
	return &UnixServer{
		socket:      sock,
		manager:     m,
		signals:     make(chan os.Signal, 1),
		logger:      NewLogger(),
		interpreter: stream.NewInterpreter(),
		fileHandler: fh,
	}, nil
}

func (s *UDPServer) Listen() (err error) {
	s.conn, err = net.ListenUDP("udp", &s.addr)
	if err != nil {
		return err
	}
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
