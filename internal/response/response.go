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

type Writer struct {
	// TODO: I could maybe embed this instead...
	writer io.Writer
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{writer}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
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

	_, err := w.writer.Write(statusLine)
	return err
}

func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()
	h.Set(headers.ContentType, "text/plain")
	h.Set(headers.ContentLength, fmt.Sprintf("%d", contentLen))
	h.Set(headers.Connection, "close")

	return h
}

func (w *Writer) WriteHeaders(headers *headers.Headers) error {
	var b []byte

	headers.ForEach(func(n, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", n, v)
	})
	b = fmt.Append(b, "\r\n")

	_, err := w.writer.Write(b)
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.writer.Write(p)

	return n, err
}
