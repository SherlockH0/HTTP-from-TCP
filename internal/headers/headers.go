package headers

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/request"
	"regexp"
)

var INVALID_HEADER = fmt.Errorf("invalid header")

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	idx := bytes.Index(data, request.SEPARATOR)
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 0, true, nil
	}

	data = data[:idx]
	parts := bytes.SplitN(data, []byte(":"), 2)
	if len(parts) != 2 {
		return 0, false, INVALID_HEADER
	}

	fieldName := parts[0]
	if fieldName[len(fieldName)-1] == ' ' {
		return 0, false, INVALID_HEADER
	}
	fieldName = bytes.ToLower(bytes.TrimSpace(fieldName))
	m, err := regexp.Match(`^[a-z0-9!#$%&'*+\-.^_\x60|~]*$`, fieldName)

	if err != nil {
		return 0, false, err
	}
	if !m {
		return 0, false, INVALID_HEADER
	}
	value := string(bytes.TrimSpace(parts[1]))

	header := h[string(fieldName)]

	if header != "" {
		value = header + ", " + value
	}

	h[string(fieldName)] = value
	return idx + len(request.SEPARATOR), false, nil
}
