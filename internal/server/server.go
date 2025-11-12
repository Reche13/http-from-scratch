package server

import (
	"fmt"
	"log"
	"net"
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

	for {
		conn, err := s.ln.Accept()
		if err != nil {
			return fmt.Errorf("connection accepting error: %w", err)
		}

		conn.Write([]byte("hello from http server"))
	}
}