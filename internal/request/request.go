package request

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type ParserState string

const (
	StateInit ParserState = "init"
	StateHeaders ParserState = "headers"
	StateBody ParserState = "body"
	StateDone ParserState = "done"
	StateError ParserState = "error"
)

type Request struct {
	RequestLine RequestLine
	Headers *Headers
	Body string
	state ParserState
}

type RequestLine struct {
	Method string
	Path string
	HttpVersion string
}

var ERROR_INCOMPLETE_START_LINE = fmt.Errorf("incomplete start-line")
var ERROR_MALFORMED_REQUEST_LINE = fmt.Errorf("malformed request-line")
var ERROR_REQUEST_IN_ERROR_STATE = fmt.Errorf("request in error state")

var SEPARATOR = []byte("\r\n")

func newRequest() *Request {
	return &Request{
		state: StateInit,
		Headers: NewHeaders(),
		Body: "",
	}
}

func getIntHeader(h *Headers, name string, defaultvalue int) int {
	valStr, ok := h.Get(name)
	if !ok {
		return defaultvalue
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultvalue
	}

	return val
}

func (r *Request) hasBody() bool {
	contentLength := getIntHeader(r.Headers, "content-length", 0)
	return contentLength > 0
}

func (r *Request) Done() bool {
	return r.state == StateDone || r.state == StateError
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0

	outer:
	for {
		currentData := data[read:]
		switch r.state {
		case StateError:
			return 0, ERROR_REQUEST_IN_ERROR_STATE

		case StateInit:
			rl, n, err := parseRequestLine(currentData)
			if err != nil {
				r.state = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n
			r.state = StateHeaders

		case StateHeaders:
			n, done, err := r.Headers.Parse(currentData)
			if err != nil {
				r.state = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			read += n

			if done {
				if r.hasBody() {
					r.state = StateBody
				} else {
					r.state = StateDone
				}
			}
		
		case StateBody:
			contentLength := getIntHeader(r.Headers, "content-length", 0)
			if contentLength == 0 {
				r.state = StateDone
				break
			}

			remaining := min(contentLength - len(r.Body), len(currentData))
			r.Body += string(currentData[:remaining])
			read += remaining

			if len(r.Body) == contentLength {
				r.state = StateDone
			}
			// body
		case StateDone:
			break outer 

		default:
			panic("Bad code bro")
		}
	}

	return read, nil
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPARATOR)

	if idx == -1 {
		return nil, 0, ERROR_INCOMPLETE_START_LINE
	}

	startLine := b[:idx]
	read := idx + len(SEPARATOR)



	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ERROR_MALFORMED_REQUEST_LINE
	}

	httpVersionParts := bytes.Split(parts[2], []byte("/"))
	if len(httpVersionParts) != 2 || string(httpVersionParts[0]) != "HTTP" || string(httpVersionParts[1]) != "1.1" {
		return nil, 0, ERROR_MALFORMED_REQUEST_LINE
	}


	rl := &RequestLine{
		Method: string(parts[0]),
		Path: string(parts[1]),
		HttpVersion: string(parts[2]),
	}

	return rl, read, nil
}

func ReadRequest(reader io.Reader) (*Request, error) {
	request := newRequest()

	buf := make([]byte, 1024)
	bufLen := 0

	for !request.Done() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}

		bufLen += n
		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	return request, nil
}