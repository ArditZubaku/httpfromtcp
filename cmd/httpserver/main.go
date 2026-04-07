package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ArditZubaku/httpfromtcp/internal/headers"
	"github.com/ArditZubaku/httpfromtcp/internal/request"
	"github.com/ArditZubaku/httpfromtcp/internal/response"
	"github.com/ArditZubaku/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	s, err := server.Serve(
		port,
		func(w *response.Writer, req *request.Request) {
			h := response.GetDefaultHeaders(0)
			body := respondWith200()
			statusCode := response.StatusOK

			switch req.RequestLine.RequestTarget {
			case "/yourproblem":
				body = respondWith400()
				statusCode = response.StatusBadRequest
			case "/myproblem":
				body = respondWith500()
				statusCode = response.StatusInternalServerError
			}

			h.Replace(headers.ContentLength, fmt.Sprintf("%d", len(body)))
			h.Replace(headers.ContentType, "text/html")
			w.WriteStatusLine(statusCode)
			w.WriteHeaders(h)
			w.WriteBody(body)
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

func respondWith400() []byte {
	return []byte(`
	<html>
	  <head>
	    <title>400 Bad Request</title>
	  </head>
	  <body>
	    <h1>Bad Request</h1>
	  </body>
	</html>
	`)
}

func respondWith500() []byte {
	return []byte(`
	<html>
	  <head>
	    <title>500 Internal Server Error</title>
	  </head>
	  <body>
	    <h1>Internal Server Error</h1>
	  </body>
	</html>
	`)
}

func respondWith200() []byte {
	return []byte(`
	<html>
	  <head>
	    <title>200 OK</title>
	  </head>
	  <body>
	    <h1>Success!</h1>
	    <p>Your request was an absolute banger.</p>
	  </body>
	</html>
	`)
}
