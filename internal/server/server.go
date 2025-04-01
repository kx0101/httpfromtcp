package server

import (
	"fmt"
	"io"
	"net"
	"sync/atomic"

	"github.com/kx0101/httpfromtcp/internal/request"
	"github.com/kx0101/httpfromtcp/internal/response"
)

type Server struct {
	Port     int
	Listener net.Listener
	Closed   atomic.Bool
	Handler  Handler
}

type Handler func(w *response.Writer, req *request.Request) *HandlerError
type HandlerError struct {
	Message string
	Status  response.StatusCode
}

func Serve(port int, handler Handler) (*Server, error) {
	if handler == nil {
		return nil, fmt.Errorf("handler cannot be nil")
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	server := &Server{
		Port:     port,
		Listener: listener,
		Handler:  handler,
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

	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Println("error:", err)

		WriteHandlerError(conn, &HandlerError{
			Message: err.Error(),
			Status:  response.BadRequest,
		})

		return
	}

	writer := response.NewWriter()
	handlerErr := s.Handler(writer, req)
	if handlerErr != nil {
		fmt.Println("error during in handler:", handlerErr.Message)
		WriteHandlerError(conn, handlerErr)
		return
	}

	if writer.State < 2 {
		writer.WriteHeaders(response.GetDefaultHeaders(len(writer.Body), "text/html"))
	}

	_, err = conn.Write(writer.Body)
	if err != nil {
		fmt.Println("error:", err)
	}
}

func WriteHandlerError(w io.Writer, handlerErr *HandlerError) {
	status := handlerErr.Status
	if status < 100 || status > 999 {
		status = response.InternalServerError
	}

	body := []byte(handlerErr.Message)
	headers := response.GetDefaultHeaders(len(body), "text/html")

	writer := response.NewWriter()

	writer.WriteStatusLine(status)
	writer.WriteHeaders(headers)
	writer.Write(body)

	if conn, ok := w.(net.Conn); ok {
		conn.Write(writer.Body)
	}
}
