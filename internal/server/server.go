// Package server provides HTTP server functionality.
package server

import (
	"fmt"
	"net"

	"github.com/ArditZubaku/httpfromtcp/internal/request"
	"github.com/ArditZubaku/httpfromtcp/internal/response"
)

const maxConcurrentConnections = 100

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	closed  bool
	handler Handler
	sem     chan struct{}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	rp := response.NewWriter(conn)

	r, err := request.RequestFromReader(conn)
	if err != nil {
		rp.WriteStatusLine(response.StatusBadRequest)
		rp.WriteHeaders(response.GetDefaultHeaders(0))
		return
	}

	s.handler(rp, r)
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

		s.sem <- struct{}{}
		go func() {
			defer func() { <-s.sem }()
			s.handleConnection(conn)
		}()
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
		sem:     make(chan struct{}, maxConcurrentConnections),
	}

	go s.listen(listener)

	return s, nil
}
