package main

import (
	"log"

	"github.com/reche13/http-from-scratch/internal/server"
)

func main() {
	srv := server.New(":8080")

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}