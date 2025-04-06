package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kx0101/httpfromtcp/internal/request"
	"github.com/kx0101/httpfromtcp/internal/response"
	"github.com/kx0101/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, func(w *response.Writer, req *request.Request) *server.HandlerError {
		w.SetHeader("X-Custom-Header", "test")

		var status response.StatusCode
		var body string

		if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
			targetPath := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
			resp, err := http.Get("https://httpbin.org/" + targetPath)

			if err != nil {
				w.WriteStatusLine(response.InternalServerError)
				w.Write([]byte(fmt.Sprintf("Failed to proxy request: %v", err)))

				return nil
			}

			defer resp.Body.Close()

			w.WriteStatusLine(response.StatusCode(resp.StatusCode))

			for k, v := range resp.Header {
				w.SetHeader(k, strings.Join(v, ", "))
			}

			w.Headers.Delete("Content-Length")
			w.SetHeader("Transfer-Encoding", "chunked")

			buf := make([]byte, 1024)
			for {
				n, err := resp.Body.Read(buf)
				if n > 0 {
					fmt.Printf("Read %d bytes from httpbin\n", n)
					w.WriteChunkedBody(buf[:n])
				}

				if err == io.EOF {
					break
				}

				if err != nil {
					break
				}
			}

			w.WriteChunkedBodyDone()

			return nil
		}

		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			status = response.BadRequest
			body = `<html>
		<head><title>400 Bad Request</title></head>
		<body><h1>Bad Request</h1><p>Your request honestly kinda sucked.</p></body>
		</html>`
		case "/myproblem":
			status = response.InternalServerError
			body = `<html>
		<head><title>500 Internal Server Error</title></head>
		<body><h1>Internal Server Error</h1><p>Okay, you know what? This one is on me.</p></body>
		</html>`
		default:
			status = response.OK
			body = `<html>
		<head><title>200 OK</title></head>
		<body><h1>Success!</h1><p>Your request was an absolute banger.</p></body>
		</html>`
		}

		w.WriteStatusLine(status)
		w.Write([]byte(body))

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer server.Close()
	fmt.Printf("server started on port %d\n", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("shutting down server")
}
