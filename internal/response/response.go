package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/reche13/http-from-scratch/internal/headers"
)

type Response struct {

}

type StatusCode int
const (
	StatusOk StatusCode = 200
	StatusBadRequest StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

type Writer struct {
	writer io.Writer
	chunked bool
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
		chunked: false,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	statusLine := []byte{}
	switch statusCode {
	case StatusOk:
		statusLine = []byte("HTTP/1.1 200 OK\r\n")
	case StatusBadRequest:
		statusLine = []byte("HTTP/1.1 400 Bad Request\r\n")
	case StatusInternalServerError:
		statusLine = []byte("HTTP/1.1 500 Internal Server Error\r\n")
	default:
		return fmt.Errorf("unrecognized error code")
	}

	_, err := w.writer.Write(statusLine)
	return err
}

func (w *Writer) WriteHeaders(h *headers.Headers) error {
	b := []byte{}
	h.ForEach(func(n, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", n, v)
	})
	b = fmt.Appendf(b, "\r\n")
	_, err := w.writer.Write(b)
	return err
}

func (w *Writer) WriteBody(data []byte) (int, error) {
	if w.chunked {
		return w.WriteChunk(data)
	}
	n, err := w.writer.Write(data)
	return n, err
}


func (w *Writer) EnableChunkedEncoding(h *headers.Headers) {
	w.chunked = true
	h.Replace("Transfer-Encoding", "chunked")
	h.Remove("Content-Length")
}


func (w *Writer) WriteChunk(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}

	sizeHex := strconv.FormatInt(int64(len(data)), 16)
	chunk := fmt.Sprintf("%s\r\n", sizeHex)
	
	_, err := w.writer.Write([]byte(chunk))
	if err != nil {
		return 0, err
	}

	n, err := w.writer.Write(data)
	if err != nil {
		return n, err
	}

	_, err = w.writer.Write([]byte("\r\n"))
	if err != nil {
		return n, err
	}

	return n, nil
}

func (w *Writer) FinalizeChunkedEncoding() error {
	if !w.chunked {
		return fmt.Errorf("chunked encoding not enabled")
	}
	_, err := w.writer.Write([]byte("0\r\n\r\n"))
	return err
}

func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-length", fmt.Sprintf("%d",contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func GetDefaultHeadersChunked() *headers.Headers {
	h := headers.NewHeaders()
	h.Set("Transfer-Encoding", "chunked")
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}