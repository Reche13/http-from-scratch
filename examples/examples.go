package examples

import (
	"io"
	"os"
	"time"

	"github.com/reche13/http-from-scratch/internal/request"
	"github.com/reche13/http-from-scratch/internal/response"
)

func Handler(w *response.Writer, r *request.Request) {
	switch r.RequestLine.Path {
	case "/":
		home(w, r)
	case "/bad-request":
		badRequest(w, r)
	case "/server-error":
		ServerError(w, r)
	case "/logs":
		streamLogs(w, r)
	case "/video":
		streamVideo(w, r)
	default:
		notFound(w)
	}
}

func notFound(w *response.Writer) {
	body := []byte(`<h1>404 Not Found</h1>`)
	h := response.GetDefaultHeaders(len(body))
	h.Replace("Content-Type", "text/html")

	w.WriteStatusLine(response.StatusNotFound)
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func home(w *response.Writer, _ *request.Request) {
	body := []byte(`<h1>Welcome to HTTP-from-scratch</h1>`)
	h := response.GetDefaultHeaders(len(body))
	h.Replace("Content-Type", "text/html")

	w.WriteStatusLine(response.StatusOk)
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func badRequest(w *response.Writer, _ *request.Request) {
	body := []byte(`<h1>400 Bad Request</h1>`)
	h := response.GetDefaultHeaders(len(body))
	h.Replace("Content-Type", "text/html")

	w.WriteStatusLine(response.StatusBadRequest)
	w.WriteHeaders(h)
	w.WriteBody(body)
}


func ServerError(w *response.Writer, _ *request.Request) {
	body := []byte(`<h1>500 Internal Server Error</h1>`)
	h := response.GetDefaultHeaders(len(body))
	h.Replace("Content-Type", "text/html")

	w.WriteStatusLine(response.StatusInternalServerError)
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func streamLogs(w *response.Writer, _ *request.Request) {
	h := response.GetDefaultHeadersChunked()
	h.Replace("Content-Type", "text/plain")

	w.EnableChunkedEncoding(h)
	w.WriteStatusLine(response.StatusOk)
	w.WriteHeaders(h)

	file, err := os.Open("./sample-data/server.log")
	if err != nil {
		w.WriteChunk([]byte("Error opening log file\n"))
		w.FinalizeChunkedEncoding()
		return
	}
	defer file.Close()

	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if n > 0 {
			time.Sleep(500 * time.Millisecond) // simulate delay
			w.WriteChunk(buf[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			w.WriteChunk([]byte("Error reading file\n"))
			break
		}
	}

	w.FinalizeChunkedEncoding()
}

func streamVideo(w *response.Writer, _ *request.Request) {
	file, err := os.Open("./sample-data/video.mp4")
	if err != nil {
		h := response.GetDefaultHeaders(0)
		w.WriteStatusLine(response.StatusInternalServerError)
		w.WriteHeaders(h)
		return
	}
	defer file.Close()

	h := response.GetDefaultHeadersChunked()
	h.Replace("Content-Type", "video/mp4")

	w.EnableChunkedEncoding(h)
	w.WriteStatusLine(response.StatusOk)
	w.WriteHeaders(h)

	buf := make([]byte, 1024*64)
	for {
		n, err := file.Read(buf)
		if n > 0 {
			w.WriteChunk(buf[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
	}

	w.FinalizeChunkedEncoding()
}