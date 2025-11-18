package main

import (
	"log"

	"github.com/reche13/http-from-scratch/examples"
	"github.com/reche13/http-from-scratch/internal/server"
)

func main() {
	srv := server.New(8080, examples.Handler)

	if err := srv.Serve(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}