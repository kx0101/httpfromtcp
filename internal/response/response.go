package response

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kx0101/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	OK                  StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

type Writer struct {
	StatusCode StatusCode
	Headers    headers.Headers
	Body       []byte
	State      int
}

func NewWriter() *Writer {
	return &Writer{
		Headers: headers.Headers{},
		Body:    []byte{},
		State:   0,
	}
}

func (w *Writer) SetHeader(key, value string) {
	if w.State >= 2 {
		fmt.Println("Warning: Attempted to modify headers after they were written")
		return
	}

	w.Headers.Set(key, value)
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.State != 0 {
		return fmt.Errorf("error: status line already written")
	}

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

	w.Body = append(w.Body, []byte(statusLine)...)
	w.State = 1

	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.State != 1 {
		return fmt.Errorf("error: headers already written")
	}

	var headerStr strings.Builder
	for k, v := range headers {
		headerStr.WriteString(k + ": " + v + "\r\n")
	}

	headerStr.WriteString("\r\n")

	w.Body = append(w.Body, []byte(headerStr.String())...)
	w.State = 2

	return nil
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.State < 2 {
		w.Headers.Delete("Content-Length")
		w.Headers.Set("Transfer-Encoding", "chunked")

		if err := w.WriteHeaders(w.Headers); err != nil {
			return 0, err
		}
	}

	chunkSize := fmt.Sprintf("%x\r\n", len(p))
	chunk := append([]byte(chunkSize), p...)
	chunkData := append(chunk, []byte("\r\n")...)

	w.Body = append(w.Body, chunkData...)

	return len(chunkData), nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	finalChunk := []byte("0\r\n\r\n")
	w.Body = append(w.Body, finalChunk...)

	return len(w.Body), nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	if w.State < 2 {
		return fmt.Errorf("error: headers not written yet")
	}

	var trailerStr strings.Builder
	for k, v := range h {
		trailerStr.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	trailerStr.WriteString("\r\n")

	w.Body = append(w.Body, []byte(trailerStr.String())...)
	w.State = 3

	return nil
}

func (w *Writer) Write(p []byte) (int, error) {
	if w.State != 1 {
		return 0, fmt.Errorf("error: body already written")

	}

	if w.State < 2 {
		w.WriteHeaders(w.Headers)
	}

	w.Body = append(w.Body, p...)
	w.State = 3

	return len(p), nil
}

func GetDefaultHeaders(contentLen int, contentType string) headers.Headers {
	contentLength := strconv.Itoa(contentLen)

	return headers.Headers{
		"Content-Type":   contentType,
		"Content-Length": contentLength,
		"Connection":     "close",
	}
}
