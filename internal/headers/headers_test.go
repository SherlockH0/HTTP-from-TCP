package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// TEST: valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("         Host: localhost:42069         \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 41, n)
	assert.False(t, done)

	// TEST:  Valid 2 headers with existing headers
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nAuthorization: Bearer token\r\n\r\n")
	n, done, err = headers.Parse(data)
	n, done, err = headers.Parse(data[n:])
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "Bearer token", headers["authorization"])

	// TEST: Valid done
	headers = NewHeaders()
	data = []byte("\r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.True(t, done)

	// TEST: Invalid spacing header
	headers = NewHeaders()
	data = []byte("Host : localhost :42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)

	// TEST: Invalid character in header key
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)

	// TEST: Multiple values
	headers = NewHeaders()
	data = []byte("Set-Person: lane-loves-go\r\nSet-Person: prime-loves-zig\r\nSet-Person: tj-loves-ocaml\r\n\r\n")
	n = 0
	for range 3 {
		i, _, _ := headers.Parse(data[n:])
		n += i
	}
	assert.Equal(t, "lane-loves-go, prime-loves-zig, tj-loves-ocaml", headers["set-person"])
}
