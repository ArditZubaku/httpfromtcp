// Package request provides HTTP request parsing functionality for building HTTP from TCP.
package request

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/ArditZubaku/httpfromtcp/internal/headers"
)

type parserState string

const (
	StateInit    parserState = "init"
	StateHeaders parserState = "headers"
	StateBody    parserState = "body"
	StateDone    parserState = "done"
	StateError   parserState = "error"
)

type RequestLine struct {
	HTTPVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     *headers.Headers
	Body        string
	state       parserState
}

func getIntHeader(headers *headers.Headers, name string, defaultVal int) int {
	valueStr, exists := headers.Get(name)
	if !exists {
		return defaultVal
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultVal
	}

	return value
}

func newRequest() *Request {
	return &Request{
		state:   StateInit,
		Headers: headers.NewHeaders(),
		Body:    "",
	}
}

var (
	ErrMalformedRequestLine   = fmt.Errorf("malformed request-line")
	ErrUnsupportedHTTPVersion = fmt.Errorf("unsupported HTTP version")
	ErrRequestInErrorState    = fmt.Errorf("request in error state")
)

var SEPARATOR = []byte("\r\n")

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := b[:idx]
	read := idx + len(SEPARATOR)

	parts := bytes.Split(startLine, []byte(" "))
	// Check if we have METHOD PATH HTTP_PROTOCOL
	if len(parts) != 3 {
		return nil, 0, ErrMalformedRequestLine
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	// Check if the HTTP_PROTOCOL part is 'HTTP/1.1'
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, ErrMalformedRequestLine
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HTTPVersion:   string(httpParts[1]),
	}

	return rl, read, nil
}

func (r *Request) hasBody() bool {
	// TODO: When doing chunked encoding, this will need to change
	length := getIntHeader(r.Headers, "content-length", 0)
	return length > 0
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		currentData := data[read:]
		if len(currentData) == 0 {
			break outer
		}
		switch r.state {
		case StateError:
			{
				return 0, ErrRequestInErrorState
			}
		case StateInit:
			{
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
			}
		case StateDone:
			{
				break outer
			}
		case StateHeaders:
			{
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
			}
		case StateBody:
			{
				length := getIntHeader(r.Headers, "content-length", 0)
				if length == 0 {
					panic("chunked encoding not supported yet")
				}

				// We now need to parse the body
				remaining := min(length-len(r.Body), len(currentData))
				r.Body += string(currentData[:remaining])
				read += remaining

				if len(r.Body) == length {
					r.state = StateDone
				}
			}
		default:
			panic("somehow we have programmed poorly")
		}
	}

	return read, nil
}

func (r *Request) done() bool {
	return r.state == StateDone || r.state == StateError
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	// NOTE: Buffer could get overrun...
	// A header that exceeds 1k would do that
	// Or the body
	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		// TODO: Handle errors better
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
