package request

import "strings"

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: make(map[string]string),
	}
}

func (h *Headers) Set(name, value string) {
	h.headers[strings.ToLower(name)] = value
}

func (h *Headers) Get(name string) string {
	return h.headers[strings.ToLower(name)]
}