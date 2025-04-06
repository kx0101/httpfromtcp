package request

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const (
	crlf = "\r\n"
)

var (
	methods                             = []string{"GET", "POST", "PUT", "DELETE"}
	ErrTryingToParseDoneRequest         = errors.New("error: trying to read data in a done request")
	ErrInvalidRequestLine               = errors.New("error: invalid request line")
	ErrInvalidMethod                    = errors.New("error: invalid method")
	ErrInvalidTarget                    = errors.New("error: invalid request target")
	ErrInvalidHTTPVersion               = errors.New("error: invalid HTTP version")
	ErrUnknownState                     = errors.New("error: unknown state")
	ErrInvalidContentLength             = errors.New("error: invalid Content-Length")
	ErrInvalidContentLengthExpectedMore = errors.New("error: invalid request content length not equal to body")
)

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0

	for r.Status != RequestStateDone {
		bytesParsed, err := r.parseSingle(data[totalBytesParsed:])

		if err != nil {
			return totalBytesParsed, err
		}

		if bytesParsed == 0 {
			break
		}

		totalBytesParsed += bytesParsed
	}

	if r.Status == RequestStateDone && totalBytesParsed != len(data) {
		if r.Headers.Get("Content-Length") != "" {
			return totalBytesParsed, ErrInvalidContentLengthExpectedMore
		}
	}

	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.Status {
	case Initialized:
		return r.parseRequestLine(data)
	case RequestStateParsingHeaders:
		return r.parseHeaders(data)
	case RequestStateParsingBody:
		return r.parseBody(data)
	default:
		return 0, fmt.Errorf("%w: %d", ErrUnknownState, r.Status)
	}
}

func (r *Request) parseRequestLine(data []byte) (int, error) {
	requestLine, bytesParsed, err := parseRequestLine(string(data))

	if err != nil {
		return 0, err
	}

	if bytesParsed == 0 {
		return 0, nil
	}

	r.RequestLine = requestLine
	r.Status = RequestStateParsingHeaders

	return bytesParsed, nil
}

func (r *Request) parseHeaders(data []byte) (int, error) {
	bytesParsed, done, err := r.Headers.Parse(data)
	if err != nil {
		return 0, err
	}

	if bytesParsed == 0 {
		return 0, nil
	}

	if done {
		r.Status = RequestStateParsingBody
	}

	if done && r.Headers.Get("Content-Length") == "" {
		r.Status = RequestStateDone
	}

	return bytesParsed, nil
}

func (r *Request) parseBody(data []byte) (int, error) {
	contentLength := r.Headers.Get("Content-Length")

	length, err := strconv.Atoi(contentLength)
	if err != nil || length <= 0 {
		return 0, ErrInvalidContentLength
	}

	if len(r.Body) >= length {
		return 0, ErrInvalidContentLength
	}

	remainingBytes := length - len(r.Body)
	bytesToRead := min(len(data), remainingBytes)

	r.Body = append(r.Body, data[:bytesToRead]...)

	if len(r.Body) == length {
		r.Status = RequestStateDone
	}

	return bytesToRead, nil
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
