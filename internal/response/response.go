package response

import (
	"fmt"
	"io"

	"github.com/ArditZubaku/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusNotFound            StatusCode = 404
	StatusInternalServerError StatusCode = 500
)

type Response struct {
	Status  StatusCode
	Headers map[string]string
	Body    []byte
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var statusLine []byte

	switch statusCode {
	case StatusOK:
		statusLine = []byte("HTTP/1.1 200 OK\r\n")
	case StatusBadRequest:
		statusLine = []byte("HTTP/1.1 400 Bad Request\r\n")
	case StatusNotFound:
		statusLine = []byte("HTTP/1.1 404 Not Found\r\n")
	case StatusInternalServerError:
		statusLine = []byte("HTTP/1.1 500 Internal Server Error\r\n")
	default:
		return fmt.Errorf("unknown status code: %d", statusCode)
	}

	_, err := w.Write(statusLine)
	return err
}

func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Type", "text/plain")
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")

	return h
}

func WriteHeaders(w io.Writer, headers *headers.Headers) error {
	var err error
	headers.ForEach(func(n, v string) {
		if err != nil {
			return
		}
		_, err = fmt.Fprintf(w, "%s: %s\r\n", n, v)
	})
	if err != nil {
		return err
	}

	_, err = w.Write([]byte("\r\n"))
	return err
}
