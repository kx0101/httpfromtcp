package main

import (
	"fmt"
	"net"

	request "github.com/kx0101/httpfromtcp/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:42069")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("server started on port 42069")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("error:", err)
			continue
		}

		fmt.Println("connection accepted")

		go func(conn net.Conn) {
			defer func() {
				conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length: 5\r\n\r\nDone\n"))

				conn.Close()
				fmt.Println("connection closed")
			}()

			req, err := request.RequestFromReader(conn)
			if err != nil {
				fmt.Println("error:", err)
				return
			}

			fmt.Printf("Request Line:\n"+
				"  - Method: %s\n"+
				"  - Target: %s\n"+
				"  - Version: %s\n",
				req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
		}(conn)
	}
}
