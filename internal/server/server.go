// Package server provides HTTP server functionality.
package server

import (
	"bytes"
	"fmt"
	"io"
	"net"

	"github.com/ArditZubaku/httpfromtcp/internal/headers"
	"github.com/ArditZubaku/httpfromtcp/internal/request"
	"github.com/ArditZubaku/httpfromtcp/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

type Server struct {
	closed  bool
	handler Handler
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	headers := response.GetDefaultHeaders(0)

	r, err := request.RequestFromReader(conn)
	if err != nil {
		s.writeResponseHead(conn, response.StatusBadRequest, headers)
		return
	}

	// TODO: Maybe give it a buffer of some size?
	writer := bytes.NewBuffer([]byte{})
	handlerErr := s.handler(writer, r)

	var body []byte
	var statusCode response.StatusCode

	if handlerErr != nil {
		statusCode = handlerErr.StatusCode
		body = []byte(handlerErr.Message)
	} else {
		statusCode = response.StatusOK
		body = writer.Bytes()
	}

	headers.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
	s.writeResponseHead(conn, statusCode, headers)
	conn.Write(body)
}

func (s *Server) writeResponseHead(
	conn net.Conn,
	statusCode response.StatusCode,
	headers *headers.Headers,
) {
	response.WriteStatusLine(conn, statusCode)
	response.WriteHeaders(conn, headers)
}

func (s *Server) listen(listener net.Listener) {
	defer listener.Close()
	for {
		if s.closed {
			return
		}

		conn, err := listener.Accept()
		if err != nil {
			return
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) Close() error {
	s.closed = true
	return nil
}

func Serve(port uint16, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		closed:  false,
		handler: handler,
	}

	go s.listen(listener)

	return s, nil
}
