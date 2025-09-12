// Package server provides HTTP server functionality.
package server

import (
	"fmt"
	"net"
)

type Server struct {
	closed bool
}

func runConnection(s *Server, conn net.Conn) {
	defer conn.Close()

	out := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!")
	conn.Write(out)
}

func runServer(s *Server, listener net.Listener) {
	defer listener.Close()
	for {
		if s.closed {
			return
		}

		conn, err := listener.Accept()
		if err != nil {
			return
		}

		go runConnection(s, conn)
	}
}

func Serve(port uint16) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{closed: false}
	go runServer(server, listener)

	return server, nil
}

func (s *Server) Close() error {
	s.closed = true
	return nil
}
