package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/NotNotState/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	Ok                    StatusCode = 200
	Bad_Request           StatusCode = 400
	Internal_Server_Error StatusCode = 500
)

type WriterStatus int

const (
	StatusLineWrite WriterStatus = iota
	HeadersWrite
	BodyWrite
)

type Writer struct {
	httpStatusCode StatusCode
	responseWriter io.Writer
	writerState    WriterStatus
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		responseWriter: w,
		writerState:    StatusLineWrite,
	}
}

func (w *Writer) WriteStatusLine(statuscode StatusCode) error {

	if w.writerState != StatusLineWrite {
		return fmt.Errorf(
			"Calling Writer Out of Order. Writing in state %d when it should be %d",
			w.writerState,
			StatusLineWrite,
		)
	}

	var err error
	switch statuscode {
	case Ok:
		_, err = w.responseWriter.Write([]byte(fmt.Sprintf("HTTP/1.1 %d OK\r\n", statuscode)))
	case Bad_Request:
		_, err = w.responseWriter.Write([]byte(fmt.Sprintf("HTTP/1.1 %d Bad Request\r\n", statuscode)))
	case Internal_Server_Error:
		_, err = w.responseWriter.Write([]byte(fmt.Sprintf("HTTP/1.1 %d Internal Server Error\r\n", statuscode)))
	default:
		_, err = w.responseWriter.Write([]byte(fmt.Sprintf("HTTP/1.1 %d\r\n", statuscode)))
	}

	if err != nil {
		return err
	}
	w.writerState = HeadersWrite
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != HeadersWrite {
		return fmt.Errorf(
			"Calling Writer Out of Order. Writing in state %d when it should be %d",
			w.writerState,
			HeadersWrite,
		)
	}
	var err error
	headers.ForEach(
		func(key, value string) {
			res := fmt.Sprintf("%s: %s\r\n", key, value)
			_, err1 := w.responseWriter.Write([]byte(res))
			if err1 != nil {
				err = err1
			}
		},
	)
	_, err = w.responseWriter.Write([]byte("\r\n"))
	w.writerState = BodyWrite
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != BodyWrite {
		return 0, fmt.Errorf(
			"Calling Writer Out of Order. Writing in state %d when it should be %d",
			w.writerState,
			BodyWrite,
		)
	}
	w.writerState = StatusLineWrite
	return w.responseWriter.Write(p)
}

func WriteStatusLine(w io.Writer, statuscode StatusCode) error {
	var err error
	switch statuscode {
	case Ok:
		_, err = w.Write([]byte(fmt.Sprintf("HTTP/1.1 %d OK\r\n", statuscode)))
	case Bad_Request:
		_, err = w.Write([]byte(fmt.Sprintf("HTTP/1.1 %d Bad Request\r\n", statuscode)))
	case Internal_Server_Error:
		_, err = w.Write([]byte(fmt.Sprintf("HTTP/1.1 %d Internal Server Error\r\n", statuscode)))
	default:
		_, err = w.Write([]byte(fmt.Sprintf("HTTP/1.1 %d\r\n", statuscode)))
	}

	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	var err error
	n, err := w.responseWriter.Write([]byte(fmt.Sprintf("%X\r\n", len(p))))
	if err != nil {
		return 0, err
	}
	n2, err := w.responseWriter.Write(p)
	if err != nil {
		return n, err
	}
	n3, err := w.responseWriter.Write([]byte("\r\n"))
	if err != nil {
		return n2, err
	}
	return n + n2 + n3, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	return w.responseWriter.Write([]byte("0\r\n\r\n"))
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	heads := headers.NewHeaders()
	heads.Set("Content-Length", strconv.Itoa(contentLen))
	heads.Set("connection", "close")
	heads.Set("Content-Type", "text/plain")
	return *heads
}
