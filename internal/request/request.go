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
	Status      Status
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Status int

const (
	Initialized Status = iota
	Done
)

var statusName = map[Status]string{
	Initialized: "initialized",
	Done:        "done",
}

func (s Status) String() string {
	return statusName[s]
}

func parseRequestLine(data string) (RequestLine, int, error) {
	endIndex := strings.Index(data, "\r\n")
	if endIndex == -1 {
		return RequestLine{}, 0, nil
	}

	requestLine := data[:endIndex]
	components := strings.Split(requestLine, " ")

	if len(components) != 3 {
		return RequestLine{}, 0, fmt.Errorf("invalid request line: %s, expected 3 components", requestLine)
	}

	method, requestTarget, httpVersion := components[0], components[1], components[2]

	if !slices.Contains(methods, method) {
		return RequestLine{}, 0, fmt.Errorf("invalid method: %s, expected one of: %v", method, methods)
	}

	if !strings.HasPrefix(requestTarget, "/") {
		return RequestLine{}, 0, fmt.Errorf("invalid request target: %s, expected to start with /", requestTarget)
	}

	if httpVersion != "HTTP/1.1" {
		return RequestLine{}, 0, fmt.Errorf("invalid HTTP version: %s, expected: HTTP/1.1", httpVersion)
	}

	return RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   httpVersion,
	}, endIndex + 2, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	r := string(buf)

	requestLine, bytesRead, err := parseRequestLine(string(r))
	if err != nil {
		return nil, err
	}

	if bytesRead == 0 {
		return nil, fmt.Errorf("invalid request line, give more data dawg")
	}

	return &Request{
		RequestLine: requestLine,
		Status:      Done,
	}, nil
}
