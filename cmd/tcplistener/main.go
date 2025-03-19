package main

import (
	"fmt"
	"io"
	"net"
	"strings"
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
			defer conn.Close()
			defer fmt.Println("connection closed")

			lines := getLinesChannel(conn)

			for line := range lines {
				fmt.Println(line)
			}
		}(conn)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)

		currentLine := strings.Builder{}
		buf := make([]byte, 8)

		for {
			n, err := f.Read(buf)

			if err != nil {
				if err != io.EOF {
					fmt.Println("Read error:", err)
				}

				if currentLine.Len() > 0 {
					out <- currentLine.String()
				}

				break
			}

			parts := strings.Split(string(buf[:n]), "\n")
			for i, part := range parts {
				if i == len(parts)-1 {
					currentLine.WriteString(part)
					continue
				}

				out <- currentLine.String() + part
				currentLine.Reset()
			}
		}
	}()

	return out
}
