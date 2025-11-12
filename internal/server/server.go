package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	Addr string
	ln net.Listener
	done chan struct{}
}

func New(Addr string) *Server {
	return &Server{
		Addr: Addr,
		done: make(chan struct{}),
	}
}

func (s *Server) ListenAndServe() error {
	var err error
	s.ln, err = net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("listener error: %w", err)
	}

	log.Printf("listening on %s", s.Addr)

	go s.handleShutdownSignals()

	for {
		conn, err := s.ln.Accept()
		if err != nil {
			select {
			case <-s.done:
				return nil
			default:
				log.Printf("connection accept error: %v", err)
				continue
			}
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)

	_, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("Error reading connection: %v", err)
		return
	}

	resp := "HTTP/1.1 200 OK\r\n" +
	"Content-Length: 23\r\n" +
	"Content-Type: text/plain\r\n" +
	"Connection: close\r\n\r\n" +
	"Hello from http server\n"

	conn.Write([]byte(resp))
}

func (s *Server) Close() {
	close(s.done)
	if s.ln != nil {
		s.ln.Close()
	}
}

func (s *Server) handleShutdownSignals() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Printf("shutting down...")
	s.Close()
}