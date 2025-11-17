package main

import (
	"io"
	"log"

	"github.com/reche13/http-from-scratch/internal/request"
	"github.com/reche13/http-from-scratch/internal/response"
	"github.com/reche13/http-from-scratch/internal/server"
)

func main() {
	srv := server.New(8080, func(w io.Writer, r *request.Request) *server.HandlerError {
		if r.RequestLine.Path == "/bad-request" {
			return &server.HandlerError{
				StatusCode: response.StatusBadRequest,
				Message: "bad request bro\n",
			}
		} else if r.RequestLine.Path == "/server-error" {
			return &server.HandlerError{
				StatusCode: response.StatusInternalServerError,
				Message: "Internal server error bro\n",
			}
		} else {
			w.Write([]byte("all good bro\n"))
		}
		return nil
	})

	if err := srv.Serve(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}