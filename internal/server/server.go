package server

import (
	"fmt"
	"net"
	"sync/atomic"
)

type Server struct {
	Port     int
	Listener net.Listener
	Closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	server := &Server{
		Port:     port,
		Listener: listener,
	}

	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	if !s.Closed.CompareAndSwap(false, true) {
		return fmt.Errorf("server already closed")
	}

	return s.Listener.Close()
}

func (s *Server) listen() {
	for {
		if s.Closed.Load() {
			return
		}

		conn, err := s.Listener.Accept()
		if err != nil {
			if s.Closed.Load() {
				return
			}

			fmt.Println("error:", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Length: 13\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		"Hello, world!"

	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Println("error:", err)
	}
}
