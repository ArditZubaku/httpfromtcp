package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ArditZubaku/httpfromtcp/internal/request"
	"github.com/ArditZubaku/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	s, err := server.Serve(
		port,
		func(w io.Writer, req *request.Request) *server.HandlerError {
			if req.RequestLine.RequestTarget == "/yourproblem" {
				
			}
			return nil
		},
	)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer s.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
