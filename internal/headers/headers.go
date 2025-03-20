package headers

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMalformedHeaderWhitespace = errors.New("malformed header: spaces before colon")
	ErrMalformedHeaderNotFound   = errors.New("malformed header: colon not found")
)

const (
	crlf = "\r\n"
)

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	d := string(data)
	bytesParsed := 0

	for {
		crlfIndex := strings.Index(d[bytesParsed:], "\r\n")
		if crlfIndex == -1 {
			return bytesParsed, false, nil
		}

		if crlfIndex == 0 {
			return bytesParsed + 2, true, nil
		}

		crlfIndex += bytesParsed
		headerLine := d[bytesParsed:crlfIndex]

		colonIndex := strings.Index(headerLine, ":")
		fmt.Println("headerLine: ", headerLine)
		fmt.Println("colonIndex: ", colonIndex)

		if colonIndex == -1 {
			return 0, false, ErrMalformedHeaderNotFound
		}

		if headerLine[colonIndex-1] == ' ' {
			return 0, false, ErrMalformedHeaderWhitespace
		}

		key := strings.TrimSpace(headerLine[:colonIndex])
		value := strings.TrimSpace(headerLine[colonIndex+1:])

		h[key] = value

		bytesParsed = crlfIndex + 2

		if bytesParsed >= len(d) {
			break
		}
	}

	return bytesParsed, false, nil
}

func (h Headers) Get(key string) string {
	return h[key]
}

func NewHeaders() *Headers {
	headers := make(Headers)
	return &headers
}
