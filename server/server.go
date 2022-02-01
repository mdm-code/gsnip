package server

import (
	"fmt"
	"io"
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

type LoggerAdapter func(string, interface{})

func NewLogger() Logger {
	// NOTE: Add new loggers here
	switch {
	default:
		return LoggerAdapter(toStderr)
	}
}

func (l LoggerAdapter) Log(level string, msg interface{}) {
	l(level, msg)
}

func toStderr(level string, msg interface{}) {
	fmt.Fprintf(os.Stderr, "%s %s: %s\n", "gsnipd", level, msg)
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

func (s *UnixServer) Listen() (err error) {
	s.listener, err = net.Listen("unix", s.socket)
	if err != nil {
		return err
	}
	s.Log("INFO", "listening on "+s.socket)
	return
}

func (s *UnixServer) ShutDown() {
	// NOTE: file handler closes down the moment the server is closed
	s.fileHandler.Close()
	s.listener.Close()
}

func (s *UnixServer) AwaitSignal(sig ...os.Signal) {
	signal.Notify(s.signals, sig...)
	go func() {
		for {
			select {
			case <-s.signals:
				err := s.manager.Reload()
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
func (s *UnixServer) AwaitConn() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.Log("INFO", err)
			continue
		}

		buff := make([]byte, 2048)

		length, err := conn.Read(buff)
		if err != nil {
			if err == io.EOF {
				s.Log("INFO", "client sent EOF")
				continue
			}
			s.Log("INFO", err)
			continue
		}

		s.Log(
			"INFO",
			fmt.Sprintf("read %s from %v", buff, conn.RemoteAddr().Network()),
		)
		go s.respond(conn, buff[:length])
	}
}

func (s *UnixServer) respond(conn net.Conn, buff []byte) {
	defer conn.Close()

	msg := s.interpreter.Eval(buff)

	switch msg.T() {
	case stream.Rld:
		s.signals <- syscall.SIGHUP
		_, err := conn.Write([]byte(""))
		if err != nil {
			s.Log("ERROR", err)
			return
		}
		return
	default:
		resp, err := s.manager.Execute(msg)
		if err != nil {
			s.Log("ERROR", err)
			_, err = conn.Write([]byte("ERROR"))
			if err != nil {
				s.Log("ERROR", err)
				return
			}
			return
		}
		outMsg := []byte(resp)
		_, err = conn.Write(outMsg)
		if err != nil {
			s.Log("ERROR", err)
			return
		}
		s.Log("INFO", "write successful")
	}
}

func (s *UnixServer) Log(level string, msg interface{}) {
	s.logger.Log(level, msg)
}
