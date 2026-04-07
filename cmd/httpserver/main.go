package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

			route := req.RequestLine.RequestTarget
			if route == "/yourproblem" {
				body = respondWith400()
				statusCode = response.StatusBadRequest
			} else if route == "/myproblem" {
				body = respondWith500()
				statusCode = response.StatusInternalServerError
			} else if route == "/video" {
				// NOTE: In order for this to work you should have a video in this path `assets/video.mp4`
				f, err := os.ReadFile("assets/video.mp4")
				if err != nil {
					body = respondWith500()
					statusCode = response.StatusInternalServerError
				} else {
					h.Replace(headers.ContentType, "video/mp4")
					h.Replace(headers.ContentLength, fmt.Sprintf("%d", len(f)))

					w.WriteStatusLine(response.StatusOK)
					w.WriteHeaders(h)
					w.WriteBody(f)
				}
			} else if strings.HasPrefix(route, "/httpbin") {
				res, err := http.Get("https://httpbin.org/" + route[len("/httpbin/"):])
				if err != nil {
					body = respondWith500()
					statusCode = response.StatusInternalServerError
				} else {
					w.WriteStatusLine(response.StatusOK)

					h.Delete(headers.ContentLength)
					h.Set(headers.TransferEncoding, "chunked")
					h.Replace(headers.ContentType, "text/plain")
					h.Set(headers.Trailer, "X-Content-SHA256")
					h.Set(headers.Trailer, "X-Content-Length")

					w.WriteHeaders(h)

					var fullBody []byte
					for {
						data := make([]byte, 32)
						n, err := res.Body.Read(data)
						if err != nil {
							break
						}

						fullBody = append(fullBody, data[:n]...)
						w.WriteBody(fmt.Appendf(nil, "%x\r\n", n))
						w.WriteBody(data[:n])
						w.WriteBody([]byte("\r\n"))
					}
					w.WriteBody([]byte("0\r\n"))

					out := sha256.Sum256(fullBody)
					trailer := headers.NewHeaders()
					trailer.Set("X-Content-SHA256", hex.EncodeToString(out[:]))
					trailer.Set("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))

					w.WriteHeaders(trailer)

					return
				}
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
