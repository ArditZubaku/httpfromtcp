# httpfromtcp

An HTTP/1.1 server built from raw TCP sockets in Go — no `net/http` on the server side.

## What this is

A from-scratch implementation of the HTTP/1.1 protocol on top of TCP, including:

- **Request parsing** — a finite state machine that incrementally parses request lines, headers, and bodies from a TCP stream
- **Response writing** — status lines, headers, and body output with buffered I/O
- **Header management** — case-insensitive storage, multi-value support, validation per RFC token rules
- **Chunked transfer encoding** — streaming responses with trailers (SHA256 checksum, content length)
- **Reverse proxy** — proxies `/httpbin/*` routes to httpbin.org using chunked encoding
- **Video streaming** — serves binary files by streaming in 32KB chunks instead of loading into memory
- **Concurrency** — goroutine-per-connection model capped with a semaphore (100 max) and race-free shutdown via `atomic.Bool`

## Project structure

```
cmd/
  httpserver/    HTTP server on port 42069
  tcplistener/   Raw TCP listener that prints parsed requests
  udpsender/     Interactive UDP client
internal/
  server/        TCP accept loop, connection handling, concurrency control
  request/       State machine HTTP request parser
  response/      Buffered HTTP response writer
  headers/       Header parsing, storage, and validation
```

## Running

```
go run cmd/httpserver/main.go
```

Place a video at `assets/video.mp4` for the `/video` endpoint.

## Routes

| Route | Description |
|-------|-------------|
| `/` | 200 OK |
| `/yourproblem` | 400 Bad Request |
| `/myproblem` | 500 Internal Server Error |
| `/video` | Streams `assets/video.mp4` |
| `/httpbin/*` | Proxies to httpbin.org with chunked encoding and SHA256 trailer |
