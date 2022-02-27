package server

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/signal"

	"github.com/mdm-code/gsnip/internal/fs"
	"github.com/mdm-code/gsnip/internal/manager"
	"github.com/mdm-code/gsnip/internal/stream"
)

type logger interface {
	log(string, interface{})
}

type loggerAdapter func(string, interface{})

func newLogger() logger {
	// NOTE: Add new loggers here
	switch {
	default:
		return loggerAdapter(toStderr)
	}
}

func (l loggerAdapter) log(level string, msg interface{}) {
	l(level, msg)
}

func toStderr(level string, msg interface{}) {
	fmt.Fprintf(os.Stderr, "%s %s: %s\n", "gsnipd", level, msg)
}

// Server specifies the functional server interface.
type Server interface {
	Listen() error
	ShutDown()
	AwaitSignal(...os.Signal)
	AwaitConn()
	Log(string, interface{})
}

// unixServer represents a server connecting over a Unix Domain Socket.
type unixServer struct {
	socket      string
	listener    net.Listener
	manager     *manager.Manager
	signals     chan os.Signal
	logger      logger
	fileHandler *fs.FileHandler
}

// NewServer creates a server connecting over the specified network. The address
// of the sever could be a file or an address with a port.
func NewServer(ntwrk string, addr string, fname string) (Server, error) {
	switch ntwrk {
	case "unix":
		srv, err := newUnixServer(addr, fname)
		if err != nil {
			return nil, err
		}
		return srv, nil
	default:
		return nil, fmt.Errorf("unimplemented protocol: %s", ntwrk)
	}
}

func newUnixServer(sock string, fname string) (*unixServer, error) {
	fh, err := fs.NewFileHandler(fname, fs.Perm)
	if err != nil {
		return nil, err
	}
	m, err := manager.NewManager(fh)
	if err != nil {
		return nil, err
	}
	return &unixServer{
		socket:      sock,
		manager:     m,
		signals:     make(chan os.Signal, 1),
		logger:      newLogger(),
		fileHandler: fh,
	}, nil
}

// Listen causes the server to start listening on the socket.
func (s *unixServer) Listen() (err error) {
	s.listener, err = net.Listen("unix", s.socket)
	if err != nil {
		return err
	}
	err = rpc.Register(s.manager)
	if err != nil {
		return err
	}
	s.Log("INFO", "listening on "+s.socket)
	return
}

// ShutDown closes the server down.
func (s *unixServer) ShutDown() {
	// NOTE: file handler closes down the moment the server is closed
	s.fileHandler.Close()
	s.listener.Close()
}

// AwaitSignal orders the server to wait signals and call reload when one of
// them is received.
func (s *unixServer) AwaitSignal(sig ...os.Signal) {
	signal.Notify(s.signals, sig...)
	// NOTE: Goroutine runs until the program terminates. There is no reason
	// to call close(s.signals) to explicitly relieve the scheduler.
	go func() {
		for {
			select {
			case <-s.signals:
				rq := stream.Request{Operation: stream.Reload, Body: []byte{}}
				var rp stream.Reply
				err := s.manager.Execute(rq, &rp)
				if err != nil {
					s.Log("ERROR", err)
					continue
				}
				s.Log("INFO", "reloaded snippet source file")
			}
		}
	}()
}

// AwaitConn waits for incoming connections. This is a blocking function.
func (s *unixServer) AwaitConn() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.Log("INFO", err)
			continue
		}
		s.Log(
			"INFO",
			fmt.Sprintf("received connection from %v", conn.RemoteAddr().Network()),
		)
		go jsonrpc.ServeConn(conn)
	}
}

// Log logs the message with a provided severity level.
func (s *unixServer) Log(level string, msg interface{}) {
	s.logger.log(level, msg)
}
