package request

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

var (
	methods = []string{"GET", "POST", "PUT", "DELETE"}
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	rBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	r := string(rBytes)

	lines := strings.Split(r, "\r\n")
	requestLine := lines[0]

	components := strings.Split(requestLine, " ")
	if len(components) != 3 {
		return nil, fmt.Errorf("there's need to be always 3 parts to the request line. Found: %s", requestLine)
	}

	method := components[0]
	requestTarget := components[1]
	httpVersion := components[2]

	if !slices.Contains(methods, method) {
		return nil, fmt.Errorf("invalid method: %s, expected one of: %v", method, methods)
	}

	if !strings.HasPrefix(requestTarget, "/") {
		return nil, fmt.Errorf("invalid request target: %s, expected to start with /", requestTarget)
	}

	if httpVersion != "HTTP/1.1" {
		return nil, fmt.Errorf("invalid HTTP version: %s, expected: HTTP/1.1", httpVersion)
	}

	return &Request{
		RequestLine: RequestLine{
			Method:        method,
			RequestTarget: requestTarget,
			HttpVersion:   httpVersion,
		},
	}, nil
}
