package main

import (
	"fmt"
	"log"

	"github.com/reche13/http-from-scratch/internal/request"
	"github.com/reche13/http-from-scratch/internal/response"
	"github.com/reche13/http-from-scratch/internal/server"
)

func respond400() []byte {
	return []byte(`
	<html>
		<head>
			<title>400 Bad Request</title>
		</head>
		<body>
			<h1>400 Bad Request</h1>
			<p>bad request bro, try again</p>
		</body>
	</html>
	`)
}

func respond500() []byte {
	return []byte(`
	<html>
		<head>
			<title>500 Internal Server Error</title>
		</head>
		<body>
			<h1>500 Internal Server Error</h1>
			<p>Server error bro, my bad</p>
		</body>
	</html>
	`)
}

func respond200() []byte {
	return []byte(`
	<html>
		<head>
			<title>200 OK</title>
		</head>
		<body>
			<h1>200 OK</h1>
			<p>success message bro!</p>
		</body>
	</html>
	`)
}

func main() {
	srv := server.New(8080, func(w *response.Writer, r *request.Request) {
		h := response.GetDefaultHeaders(0)
		body := respond200()
		status := response.StatusOk

		if r.RequestLine.Path == "/bad-request" {
			status = response.StatusBadRequest
			body = respond400()
		} else if r.RequestLine.Path == "/server-error" {
			status = response.StatusInternalServerError
			body = respond500()
		}
		h.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
		h.Replace("Content-Type", "text/html")
		w.WriteStatusLine(status)
		w.WriteHeaders(h)
		w.WriteBody(body)
	})

	if err := srv.Serve(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}