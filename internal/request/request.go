package request

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type parserState string

const (
	StateInit parserState = "init"
	StateDone parserState = "done"
)

type Request struct {
	RequestLine RequestLine
	state       parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var ERROR_BAD_START_LINE = fmt.Errorf("invalid start line")
var ERROR_DONE_STATE_READ = fmt.Errorf("trying to read data in a done state")
var ERROR_UNKNOWN_STATE = fmt.Errorf("unknown state")
var SEPARATOR = []byte("\r\n")

func newRequest() *Request {
	return &Request{state: StateInit}
}

func (r *Request) isDone() bool {
	return r.state == StateDone
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.state {
		case StateInit:
			requestLine, n, err := parseRequestLine(data[read:])
			if err != nil {
				return n, err
			}
			if n == 0 {
				return 0, nil
			}
			r.RequestLine = *requestLine

			read += n
			r.state = StateDone
		case StateDone:
			break outer
		default:
			return 0, ERROR_UNKNOWN_STATE
		}
	}
	return read, ERROR_DONE_STATE_READ
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	index := bytes.Index(b, SEPARATOR)

	if index == -1 {
		return nil, 0, nil
	}

	startLine := b[:index]
	messageLength := index + len(SEPARATOR)

	parts := bytes.Split(startLine, []byte(" "))

	if len(parts) != 3 {
		return nil, messageLength, ERROR_BAD_START_LINE
	}

	method := string(parts[0])

	m, err := regexp.Match("[A-Z]", []byte(method))

	if err != nil {
		return nil, messageLength, err
	}
	if !m {
		return nil, messageLength, ERROR_BAD_START_LINE

	}

	target := string(parts[1])

	version := string(parts[2])

	m, err = regexp.MatchString(`HTTP/\d\.\d`, version)

	if err != nil {
		return nil, messageLength, err
	}
	if !m {
		return nil, messageLength, ERROR_BAD_START_LINE
	}

	version = strings.Split(version, "/")[1]

	return &RequestLine{
		RequestTarget: target,
		Method:        method,
		HttpVersion:   version,
	}, messageLength, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	buf := make([]byte, 1024)
	bufLen := 0
	for !request.isDone() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}
		bufLen += n

		readLen, err := request.parse(buf[:bufLen])

		if n == 0 {
			break
		}
		copy(buf, buf[readLen:bufLen])
		bufLen -= readLen
	}
	return request, nil
}
