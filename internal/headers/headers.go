package headers

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/request"
)

var INVALID_HEADER = fmt.Errorf("invalid header")

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	bufIndex := bytes.Index(data, request.SEPARATOR)
	if bufIndex == -1 {
		return 0, false, nil
	}
	if bufIndex == 0 {
		return 0, true, nil
	}

	line := data[:bufIndex]

	index := bytes.Index(line, []byte(":"))
	if index == -1 {
		return 0, false, INVALID_HEADER
	}
	fieldName := line[:index]
	if fieldName[len(fieldName)-1] == ' ' {
		return 0, false, INVALID_HEADER
	}
	fieldName = bytes.Trim(fieldName, " ")
	value := bytes.Trim(line[index+len(":"):], " ")

	h[string(fieldName)] = string(value)
	return bufIndex + len(request.SEPARATOR), false, nil
}
