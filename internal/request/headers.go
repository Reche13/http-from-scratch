package request

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: make(map[string]string),
	}
}

var ERROR_MALFORMED_FIELD_LINE = fmt.Errorf("malformed field-line")
var ERROR_MALFORMED_FIELD_NAME = fmt.Errorf("malformed field-name")

func (h *Headers) Set(name, value string) {
	h.headers[strings.ToLower(name)] = value
}

func (h *Headers) Get(name string) string {
	return h.headers[strings.ToLower(name)]
}

func parseHeader(data []byte) (string, string, error) {
	parts := bytes.SplitN(data, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", ERROR_MALFORMED_FIELD_LINE
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", ERROR_MALFORMED_FIELD_NAME
	}

	return string(name), string(value), nil
}


func (h *Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false

	for {
		idx := bytes.Index(data, []byte(SEPARATOR))

		if idx == -1 {
			break
		}

		if idx == 0 {
			done = true
			read += len(SEPARATOR)
			break
		}

		name, value, err := parseHeader(data[read:read + idx])
		if err != nil {
			return 0, false, err
		}

		idx += read + len(SEPARATOR)
		h.Set(name, value)
	}

	return read, done, nil

}