package response

import (
	"io"
	"strconv"

	"github.com/kx0101/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	OK                  StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var statusLine string

	switch {
	case statusCode == OK:
		statusLine = "HTTP/1.1 200 OK\r\n"
	case statusCode == BadRequest:
		statusLine = "HTTP/1.1 400 Bad Request\r\n"
	case statusCode == InternalServerError:
		statusLine = "HTTP/1.1 500 Internal Server Error\r\n"
	default:
		statusLine = "HTTP/1.1 500 Internal Server Error\r\n"
	}

	_, err := w.Write([]byte(statusLine))
	if err != nil {
		return err
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	contentLength := strconv.Itoa(contentLen)

	return headers.Headers{
		"Content-Type":   "text/plain",
		"Content-Length": contentLength,
		"Connection":     "close",
	}
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		_, err := w.Write([]byte(k + ": " + v + "\r\n"))
		if err != nil {
			return err
		}
	}

	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}

	return nil
}
