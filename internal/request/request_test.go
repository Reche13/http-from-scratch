package request

import (
	"strings"
	"testing"
)

func TestReadRequest(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantErr   bool
		wantState ParserState
		wantLine  *RequestLine
		wantHdrs  map[string]string
	}{
		{
			name: "valid simple request",
			input: "GET /hello HTTP/1.1\r\n" +
				"Host: example.com\r\n" +
				"User-Agent: curl\r\n" +
				"\r\n",
			wantErr:   false,
			wantState: StateDone,
			wantLine: &RequestLine{
				Method:      "GET",
				Path:        "/hello",
				HttpVersion: "HTTP/1.1",
			},
			wantHdrs: map[string]string{
				"host":        "example.com",
				"user-agent":  "curl",
			},
		},
		{
			name: "missing http version",
			input: "GET /hello\r\n" +
				"Host: x\r\n\r\n",
			wantErr:   true,
			wantState: StateError,
		},
		{
			name: "invalid header token",
			input: "GET / HTTP/1.1\r\n" +
				"H©st: bad\r\n\r\n",
			wantErr:   true,
			wantState: StateError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)

			req, err := ReadRequest(reader)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				return
			} else{
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}


			if req.state != tt.wantState {
				t.Fatalf("wrong final state: got %v want %v", req.state, tt.wantState)
			}

			if tt.wantLine != nil {
				if req.RequestLine.Method != tt.wantLine.Method ||
					req.RequestLine.Path != tt.wantLine.Path ||
					req.RequestLine.HttpVersion != tt.wantLine.HttpVersion {
					t.Fatalf("request line mismatch: got %+v want %+v",
						req.RequestLine, tt.wantLine)
				}
			}

			for k, v := range tt.wantHdrs {
				got, _ := req.Headers.Get(k)
				if got != v {
					t.Fatalf("header %s mismatch: got %q want %q", k, got, v)
				}
			}
		})
	}
}



func TestBody(t *testing.T) {
	raw := "" +
		"POST /submit HTTP/1.1\r\n" +
		"Host: example.com\r\n" +
		"Content-Length: 11\r\n" +
		"\r\n" +
		"Hello World"

	r, err := ReadRequest(strings.NewReader(raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedBody := "Hello World"
	if r.Body != expectedBody {
		t.Fatalf("expected body %s, got %s", expectedBody, r.Body)
	}

	if r.state != StateDone {
		t.Fatalf("request should be done, state = %v", r.state)
	}
}