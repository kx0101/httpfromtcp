package request

import (
	"errors"
	"fmt"
	// h "github.com/kx0101/httpfromtcp/internal/headers"
	"slices"
	"strings"
)

const (
	crlf = "\r\n"
)

var (
	methods                     = []string{"GET", "POST", "PUT", "DELETE"}
	ErrTryingToParseDoneRequest = errors.New("error: trying to read data in a done request")
	ErrInvalidRequestLine       = errors.New("invalid request line")
	ErrInvalidMethod            = errors.New("invalid method")
	ErrInvalidTarget            = errors.New("invalid request target")
	ErrInvalidHTTPVersion       = errors.New("invalid HTTP version")
	ErrUnknownState             = errors.New("error: unknown state")
)

func (r *Request) parse(data []byte) (int, error) {
	if r.Status == RequestStateDone {
		return 0, ErrTryingToParseDoneRequest
	}

	if r.Status != Initialized {
		return 0, ErrUnknownState
	}

	requestLine, requestLineBytesParsed, err := parseRequestLine(string(data))
	if err != nil {
		return 0, err
	}

	if requestLineBytesParsed == 0 {
		return 0, nil
	}

	r.RequestLine = requestLine
	r.Status = RequestStateDone

	return requestLineBytesParsed, err
}

func parseRequestLine(data string) (RequestLine, int, error) {
	endIndex := strings.Index(data, crlf)
	if endIndex == -1 {
		return RequestLine{}, 0, nil
	}

	requestLine := data[:endIndex]
	components := strings.Split(requestLine, " ")

	if len(components) != 3 {
		return RequestLine{}, 0, fmt.Errorf("%w: expected 3 components, got: %s", ErrInvalidRequestLine, requestLine)
	}

	method, requestTarget, httpVersion := components[0], components[1], components[2]

	if !slices.Contains(methods, method) {
		return RequestLine{}, 0, fmt.Errorf("%w: %s, expected one of: %v", ErrInvalidMethod, method, methods)
	}

	if !strings.HasPrefix(requestTarget, "/") {
		return RequestLine{}, 0, fmt.Errorf("%w: %s, expected to start with '/'", ErrInvalidTarget, requestTarget)
	}

	if httpVersion != "HTTP/1.1" {
		return RequestLine{}, 0, fmt.Errorf("%w: %s, expected 'HTTP/1.1'", ErrInvalidHTTPVersion, httpVersion)
	}

	return RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   httpVersion,
	}, endIndex + 2, nil
}
