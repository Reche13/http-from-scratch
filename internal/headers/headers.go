package headers

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

var SEPARATOR = "\r\n"

var ERROR_MALFORMED_FIELD_LINE = fmt.Errorf("malformed field-line")
var ERROR_MALFORMED_FIELD_NAME = fmt.Errorf("malformed field-name")

func (h *Headers) Set(name, value string) {
	name = strings.ToLower(name)
	if v, ok := h.headers[name]; ok {
		h.headers[name] = fmt.Sprintf("%s,%s", v, value)
	} else {
		h.headers[name] = value
	}
}

func (h *Headers) Replace(name, value string) {
	name = strings.ToLower(name)
	h.headers[name] = value
}

func (h *Headers) Get(name string) (string, bool) {
	val, ok := h.headers[strings.ToLower(name)]
	return val, ok
}

func (h *Headers) Remove(name string) {
	name = strings.ToLower(name)
	delete(h.headers, name)
}

func (h *Headers) ForEach(cb func(n, v string)) {
	for k, v := range h.headers {
		cb(k, v)
	}
}

func isValidToken(str []byte) bool {
	for _, ch := range str {
		if (ch >= 'A' && ch <= 'Z') ||
           (ch >= 'a' && ch <= 'z') ||
           (ch >= '0' && ch <= '9') {
            continue
        }
        switch ch {
        case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
            continue
        }
        return false
	}
	return true
}

func parseHeader(data []byte) (string, string, error) {
	parts := bytes.SplitN(data, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", ERROR_MALFORMED_FIELD_LINE
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) || !isValidToken(name) {
		return "", "", ERROR_MALFORMED_FIELD_NAME
	}

	return string(name), string(value), nil
}


func (h *Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false

	for {
		idx := bytes.Index(data[read:], []byte(SEPARATOR))

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

		read += idx + len(SEPARATOR)
		h.Set(name, value)
	}

	return read, done, nil

}