package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/reche13/http-from-scratch/internal/request"
	"github.com/reche13/http-from-scratch/internal/response"
)

type Server struct {
	Addr string
	ln net.Listener
	handler Handler
	done chan struct{}
}

func New(port uint16, handler Handler ) *Server {
	return &Server{
		Addr: fmt.Sprintf(":%d", port),
		handler: handler,
		done: make(chan struct{}),
	}
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message string
}

type Handler func(w io.Writer, r *request.Request) *HandlerError


func (s *Server) Serve() error {
	var err error
	s.ln, err = net.Listen("tcp", s.Addr)
	if err != nil {
		return  err
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

	headers := response.GetDefaultHeaders(0)
	r, err := request.ReadRequest(conn)
	if err != nil {
		response.WriteStatusLine(conn, response.StatusBadRequest)
		response.WriteHeaders(conn, headers)
		return
	}

	writer := bytes.NewBuffer([]byte{})
	handlerError := s.handler(writer, r)
	
	var body []byte = nil
	var status response.StatusCode = response.StatusOk

	if handlerError != nil {
		status = handlerError.StatusCode
		body = []byte(handlerError.Message)
	} else {
		body = writer.Bytes()
	}

	headers.Replace("Content-Length", fmt.Sprintf("%d", len(body)))

	response.WriteStatusLine(conn, status)
	response.WriteHeaders(conn, headers)
	conn.Write(body)
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