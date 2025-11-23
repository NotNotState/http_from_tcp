package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/NotNotState/httpfromtcp/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	Headers     *headers.Headers
	State       requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

// RECALL iota is Go's kinda sorta enum thingy
const (
	requestStateInitialized requestState = iota // = 0
	requestStateDone                            // = 1
	requestStateParsingHeaders
)

const crlf = "\r\n"
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	req := &Request{
		State:   requestStateInitialized,
		Headers: headers.NewHeaders(),
	}

	for req.State != requestStateDone {
		if readToIndex >= len(buffer) {
			newBuff := make([]byte, len(buffer)*2)
			copy(newBuff, buffer)
			buffer = newBuff
		}

		nBytesRead, err := reader.Read(buffer[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.State != requestStateDone {
					return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", req.State, nBytesRead)
				}
				break
			}
			return nil, err
		}
		readToIndex += nBytesRead

		nBytesParsed, err := req.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buffer, buffer[nBytesParsed:])
		readToIndex -= nBytesParsed

	}
	return req, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", str)
	}

	method := parts[0]
	//By convention, standardized methods are defined in all-uppercase US-ASCII letters
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}
	//ignore for now
	requestTarget := parts[1]

	//validating HTTP version
	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}

	versionPart := versionParts[1]
	if versionPart != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", versionPart)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionPart,
	}, nil

}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf)) // returns first index where byte argument appears within byte slice
	if idx == -1 {
		return nil, 0, nil
	}
	requestLineText := string(data[:idx])

	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}

	return requestLine, idx + 2, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.State {
	case requestStateInitialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			// something is ackkksually wrong
			return 0, err
		}
		if n == 0 {
			// just need more data can't do anything
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.State = requestStateParsingHeaders
		return n, nil
	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.State = requestStateDone
		}
		return n, nil
	case requestStateDone:
		return 0, errors.New("error: attempting to read into done state")
	default:
		return 0, errors.New("error: Unknown State")
	}
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.State != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += n
		if n == 0 {
			break
		}
	}
	return totalBytesParsed, nil
}
