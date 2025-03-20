package headers

import (
	"errors"
	"strings"
)

var (
	ErrMalformedHeaderWhitespace = errors.New("malformed header: spaces before colon")
	ErrMalformedHeaderNotFound   = errors.New("malformed header: colon not found")
	ErrInvalidHeaderChars        = errors.New("invalid header characters")
)

var validCharsMap = map[rune]bool{
	'a': true, 'b': true, 'c': true, 'd': true, 'e': true, 'f': true, 'g': true, 'h': true, 'i': true, 'j': true,
	'k': true, 'l': true, 'm': true, 'n': true, 'o': true, 'p': true, 'q': true, 'r': true, 's': true, 't': true,
	'u': true, 'v': true, 'w': true, 'x': true, 'y': true, 'z': true, 'A': true, 'B': true, 'C': true, 'D': true,
	'E': true, 'F': true, 'G': true, 'H': true, 'I': true, 'J': true, 'K': true, 'L': true, 'M': true, 'N': true,
	'O': true, 'P': true, 'Q': true, 'R': true, 'S': true, 'T': true, 'U': true, 'V': true, 'W': true, 'X': true,
	'Y': true, 'Z': true, '0': true, '1': true, '2': true, '3': true, '4': true, '5': true, '6': true, '7': true,
	'8': true, '9': true, '!': true, '#': true, '$': true, '%': true, '&': true, '\'': true, '/': true, '*': true, '+': true,
	'-': true, '.': true, '^': true, '_': true, '`': true, '|': true, '~': true, ':': true,
}

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
		if colonIndex == -1 {
			return 0, false, ErrMalformedHeaderNotFound
		}

		if headerLine[colonIndex-1] == ' ' {
			return 0, false, ErrMalformedHeaderWhitespace
		}

		key := strings.TrimSpace(headerLine[:colonIndex])
		value := strings.TrimSpace(headerLine[colonIndex+1:])

		for _, c := range key {
			if !isValidHeaderKey(c) {
				return 0, false, ErrInvalidHeaderChars
			}
		}

		for _, c := range value {
			if !isValidHeaderValue(c) {
				return 0, false, ErrInvalidHeaderChars
			}
		}

		key = strings.ToLower(key)
		h[key] = value

		bytesParsed = crlfIndex + 2

		if bytesParsed >= len(d) {
			break
		}
	}

	return bytesParsed, false, nil
}

func (h Headers) Get(key string) string {
	return h[strings.ToLower(key)]
}

func NewHeaders() *Headers {
	headers := make(Headers)
	return &headers
}

func isValidHeaderValue(c rune) bool {
	_, exists := validCharsMap[c]
	return exists
}

func isValidHeaderKey(c rune) bool {
	return c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' &&
		c <= '9' || c == '-' || c == '_'
}
