package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/reche13/http-from-scratch/internal/request"
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

	r, err := request.ReadRequest(conn)
	if err != nil {
		log.Printf("failed to read: %v", err)
	}

	body := fmt.Sprintf("method: %s, http-version: %s, path: %s", r.RequestLine.Method, r.RequestLine.HttpVersion, r.RequestLine.Path )

	resp := "HTTP/1.1 200 OK\r\n" +
	fmt.Sprintf("Content-Length: %d\r\n", len(body) + 1) +
	"Content-Type: text/plain\r\n" +
	"Connection: close\r\n\r\n" +
	fmt.Sprintf("%s\n", body)

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