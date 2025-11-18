package response

import (
	"bytes"
	"testing"

	"github.com/reche13/http-from-scratch/internal/headers"
)

func TestWriteStatusLine(t *testing.T) {
	tests := []struct {
		name       string
		statusCode StatusCode
		want       string
	}{
		{
			name:       "200 OK",
			statusCode: StatusOk,
			want:       "HTTP/1.1 200 OK\r\n",
		},
		{
			name:       "400 Bad Request",
			statusCode: StatusBadRequest,
			want:       "HTTP/1.1 400 Bad Request\r\n",
		},
		{
			name:       "500 Internal Server Error",
			statusCode: StatusInternalServerError,
			want:       "HTTP/1.1 500 Internal Server Error\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			w := NewWriter(&buf)

			err := w.WriteStatusLine(tt.statusCode)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if buf.String() != tt.want {
				t.Fatalf("got %q, want %q", buf.String(), tt.want)
			}
		})
	}
}

func TestWriteHeaders(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	h := headers.NewHeaders()
	h.Set("Content-Type", "text/html")
	h.Set("Content-Length", "42")

	err := w.WriteHeaders(h)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	expected := "content-type: text/html\r\ncontent-length: 42\r\n\r\n"
	if output != expected {
		t.Fatalf("expected headers %s, got %s", expected, output)
	}
}

func TestWriteBody(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	data := []byte("Hello World")
	n, err := w.WriteBody(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if n != len(data) {
		t.Fatalf("got %d bytes written, want %d", n, len(data))
	}

	if buf.String() != string(data) {
		t.Fatalf("got %q, want %q", buf.String(), string(data))
	}
}

func TestGetDefaultHeaders(t *testing.T) {
	h := GetDefaultHeaders(100)

	contentLen, _ := h.Get("content-length")
	if contentLen != "100" {
		t.Fatalf("content-length: got %q, want %q", contentLen, "100")
	}

	conn, _ := h.Get("connection")
	if conn != "close" {
		t.Fatalf("connection: got %q, want %q", conn, "close")
	}

	contentType, _ := h.Get("content-type")
	if contentType != "text/plain" {
		t.Fatalf("content-type: got %q, want %q", contentType, "text/plain")
	}
}

func TestChunkedResponse(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	h := GetDefaultHeadersChunked()
	w.EnableChunkedEncoding(h)

	// Write status line
	err := w.WriteStatusLine(StatusOk)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Write headers
	err = w.WriteHeaders(h)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Write chunks
	_, err = w.WriteChunk([]byte("Hello"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = w.WriteChunk([]byte(" World"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Finalize
	err = w.FinalizeChunkedEncoding()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	expected := "HTTP/1.1 200 OK\r\ntransfer-encoding: chunked\r\nconnection: close\r\ncontent-type: text/plain\r\n\r\n5\r\nHello\r\n6\r\n World\r\n0\r\n\r\n"
	if output != expected {
		t.Fatalf("got %q, want %q", output, expected)
	}
}

