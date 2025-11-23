package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const crlf = "\r\n" // This is the break line in HTTP protocol

var ERROR_MALFORMED_FIELD_LINE = fmt.Errorf("Malformed Field Line")
var ERROR_MALFORMED_FIELD_NAME = fmt.Errorf("Malformed Field NAME")
var ERROR_INVALID_FIELD_NAME_CHARS = fmt.Errorf("Invalid Field Name Characters")

func IsValidFieldName(name string) (string, error) {
	if len(name) < 1 {
		return "", ERROR_INVALID_FIELD_NAME_CHARS
	}

	validChars := map[rune]bool{
		'!': true, '#': true, '$': true, '%': true,
		'&': true, '\'': true, '*': true, '+': true,
		'-': true, '.': true, '^': true, '_': true,
		'`': true, '|': true, '~': true,
	}

	validated := make([]byte, 0, len(name))

	for _, c := range name {
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
			validated = append(validated, byte(c))
		} else if c >= '0' && c <= '9' {
			validated = append(validated, byte(c))
		} else if validChars[c] {
			validated = append(validated, byte(c))
		} else {
			return "", ERROR_INVALID_FIELD_NAME_CHARS
		}
	}

	return string(validated), nil
}

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", ERROR_MALFORMED_FIELD_LINE
	}
	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", ERROR_MALFORMED_FIELD_NAME
	}

	name = bytes.TrimSpace(name)

	return string(name), string(value), nil

}

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func (h *Headers) Get(name string) string {
	return h.headers[strings.ToLower(name)]
}

func (h *Headers) Set(name, value string) {
	name = strings.ToLower(name)
	if _, ok := h.headers[name]; ok {
		h.headers[name] = h.headers[name] + "," + value
	} else {
		h.headers[name] = value
	}

}

func (h *Headers) Parse(data []byte) (int, bool, error) {

	read := 0
	done := false

	for {
		idx := bytes.Index(data[read:], []byte(crlf))
		if idx == -1 {
			break
		}
		// hit Empty header (registered nurs)
		if idx == 0 {
			done = true
			read += len(crlf)
			break
		}

		name, value, err := parseHeader(data[read : read+idx])

		if err != nil {
			return 0, false, err
		}

		read += idx + len(crlf)

		validName, err := IsValidFieldName(name)
		if err != nil {
			return 0, false, err
		}

		h.Set(validName, value)
	}

	return read, done, nil
}
