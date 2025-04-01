package main

import (
	"fmt"
	"os"
	"os/signal"
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
