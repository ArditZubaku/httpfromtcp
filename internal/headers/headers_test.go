package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\nFooFoo:    barbar\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	header, ok := headers.Get("HOST")
	assert.True(t, ok)
	assert.Equal(t, "localhost:42069", header)
	header, ok = headers.Get("FooFoo")
	assert.True(t, ok)
	assert.Equal(t, "barbar", header)
	header, ok = headers.Get("MissingKey")
	assert.False(t, ok)
	assert.Equal(t, "", header)
	assert.Equal(t, 44, n)
	assert.True(t, done)

	// Test: Invalid header with special character in name
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Multiple headers with same name (should be concatenated)
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nHost: localhost:42069\r\n\r\n")
	_, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)

	header, ok = headers.Get("host")
	assert.True(t, ok)
	assert.Equal(t, "localhost:42069,localhost:42069", header)
	assert.True(t, done)
}
