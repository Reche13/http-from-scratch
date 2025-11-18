package server

import (
	"testing"

	"github.com/reche13/http-from-scratch/internal/request"
	"github.com/reche13/http-from-scratch/internal/response"
)

func TestNew(t *testing.T) {
	handler := func(w *response.Writer, r *request.Request) {}
	srv := New(8080, handler)

	if srv.Addr != ":8080" {
		t.Fatalf("got addr %q, want %q", srv.Addr, ":8080")
	}

	if srv.handler == nil {
		t.Fatalf("handler should not be nil")
	}

	if srv.done == nil {
		t.Fatalf("done channel should not be nil")
	}
}

func TestClose(t *testing.T) {
	handler := func(w *response.Writer, r *request.Request) {}
	srv := New(8080, handler)

	srv.Close()

	select {
	case <-srv.done:
	default:
		t.Fatalf("done channel should be closed")
	}
}
