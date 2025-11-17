package response

import (
	"fmt"
	"io"

	"github.com/reche13/http-from-scratch/internal/request"
)

type Response struct {

}

type StatusCode int
const (
	StatusOk StatusCode = 200
	statusBadRequest StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := []byte{}
	switch statusCode {
	case StatusOk:
		statusLine = []byte("HTTP/1.1 200 OK\r\n")
	case statusBadRequest:
		statusLine = []byte("HTTP/1.1 400 Bad Request\r\n")
	case StatusInternalServerError:
		statusLine = []byte("HTTP/1.1 500 Internal Server Error\r\n")
	default:
		return fmt.Errorf("unrecognized error code")
	}

	_, err := w.Write(statusLine)
	return err
}


func GetDefaultHeaders(contentLen int) *request.Headers {
	h := request.NewHeaders()
	h.Set("Content-length", fmt.Sprintf("%d",contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func WriteHeaders(w io.Writer, headers *request.Headers) error {
	b := []byte{}
	headers.ForEach(func(n, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", n, v)
	})
	b = fmt.Appendf(b, "\r\n")
	_, err := w.Write(b)
	return err
}