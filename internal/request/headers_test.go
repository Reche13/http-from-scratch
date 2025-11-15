package request_test

import (
	"testing"

	"github.com/reche13/http-from-scratch/internal/request"
)

func TestHeadersParse(t *testing.T) {
	tests := []struct {
		name string
		input string
		wantHeaders map[string]string
		wantRead int
		wantDone bool
		wantErr bool
	}{
		{
			name: "simple headers",
            input: "Host: localhost:6969\r\nUser-Agent: curl\r\n\r\n",
            wantHeaders: map[string]string{
                "host":        "localhost:6969",
                "user-agent":  "curl",
            },
            wantRead:  len("Host: localhost:6969\r\nUser-Agent: curl\r\n\r\n"),
            wantDone: true,
            wantErr:  false,
		},
		{
            name: "malformed field name",
            input: "User : bad\r\n\r\n",
            wantErr: true,
        },
		{
            name: "malformed missing colon",
            input: "Host example.com\r\n\r\n",
            wantErr: true,
        },
		{
            name: "incomplete header line",
            input: "Host: example.com",
            wantDone: false,
            wantErr: false,
            wantRead: 0,
            wantHeaders: map[string]string{},
        },
		{
			name: "invalid token chars",
			input: "H©st: example.com\r\n\r\n",
			wantDone: false,
			wantRead: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := request.NewHeaders()

			read, done, err := h.Parse([]byte(tt.input))

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			if read != tt.wantRead {
				t.Fatalf("read mismatch: got %d, want %d", read, tt.wantRead)
			}

			if done != tt.wantDone {
                t.Fatalf("done mismatch: got %v, want %v", done, tt.wantDone)
            }

			if len(tt.wantHeaders) != len(h.Headers) {
                t.Fatalf("header count mismatch, got %d want %d", len(h.Headers), len(tt.wantHeaders))
            }

            for k, v := range tt.wantHeaders {
                got := h.Get(k)
                if got != v {
                    t.Fatalf("header %q mismatch: got %q want %q", k, got, v)
                }
            }
		})
	}
}