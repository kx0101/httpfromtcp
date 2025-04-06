package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/kx0101/httpfromtcp/internal/headers"
	"github.com/kx0101/httpfromtcp/internal/request"
	"github.com/kx0101/httpfromtcp/internal/response"
	"github.com/kx0101/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, func(w *response.Writer, req *request.Request) *server.HandlerError {
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

			w.SetHeader("Transfer-Encoding", "chunked")
			w.SetHeader("Trailer", "X-Content-SHA256, X-Content-Length")
			w.Headers.Delete("Content-Length")

			var fullBody = make([]byte, 0)

			buf := make([]byte, 1024)
			for {
				n, err := resp.Body.Read(buf)
				if n > 0 {
					fmt.Printf("Read %d bytes from httpbin\n", n)

					fullBody = append(fullBody, buf[:n]...)
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

			trailers := headers.Headers{
				"X-Content-SHA256": fmt.Sprintf("%x", sha256.Sum256(fullBody)),
				"X-Content-Length": strconv.Itoa(len(fullBody)),
			}
			w.WriteTrailers(trailers)

			return nil
		}

		if req.RequestLine.RequestTarget == "/video" {
			file, err := os.Open("./assets/vim.mp4")
			if err != nil {
				fmt.Println("Error opening video file:", err)

				w.WriteStatusLine(response.InternalServerError)
				w.Write([]byte(fmt.Sprintf("Failed to open video: %v", err)))

				return nil
			}

			w.WriteStatusLine(response.OK)

			w.SetHeader("Content-Type", "video/mp4")
			w.Headers.Delete("Content-Length")
			w.SetHeader("Transfer-Encoding", "chunked")

			buf := make([]byte, 1024)
			for {
				n, err := file.Read(buf)
				if n > 0 {
					_, writeErr := w.WriteChunkedBody(buf[:n])
					if writeErr != nil {
						fmt.Println("Error writing chunked body:", writeErr)
						break
					}
				}

				if err == io.EOF {
					break
				}

				if err != nil {
					break
				}
			}

			_, err = w.WriteChunkedBodyDone()
			if err != nil {
				fmt.Println("Error writing chunked body done:", err)
			}

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
