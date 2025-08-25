package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
}

var ErrMalformedRequestLine = fmt.Errorf("malformed request-line")
var ErrUnsupportedHttpVersion = fmt.Errorf("unsupported HTTP version")

const SEPARATOR = "\r\n"

func parseRequestLine(b string) (*RequestLine, string, error) {
	idx := strings.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, b, nil
	}

	startLine := b[:idx]
	restOfMsg := b[idx+len(SEPARATOR):]

	parts := strings.Split(startLine, " ")
	// Check if we have METHOD PATH HTTP_PROTOCOL
	if len(parts) != 3 {
		return nil, restOfMsg, ErrMalformedRequestLine
	}

	httpParts := strings.Split(parts[2], "/")
	// Check if the HTTP_PROTOCOL part is 'HTTP/1.1'
	if len(parts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return nil, restOfMsg, ErrMalformedRequestLine
	}

	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   httpParts[1],
	}

	return rl, restOfMsg, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("unable to io.ReadAll"), err)
	}

	str := string(data)
	rl, _, err := parseRequestLine(str)

	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *rl,
	}, err
}
