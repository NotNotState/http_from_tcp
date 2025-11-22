package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

const crlf = "\r\n"

var ERROR_MALFORMED_FIELD_LINE = fmt.Errorf("Malformed Field Line")
var ERROR_MALFORMED_FIELD_NAME = fmt.Errorf("Malformed Field NAME")

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

func (h Headers) Parse(data []byte) (int, bool, error) {

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
		fmt.Println(value)
		h[name] = value
	}

	return read, done, nil
}
