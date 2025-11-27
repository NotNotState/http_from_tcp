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

// type Writer struct {
// 	statusCode StatusCode
// 	writer     io.Writer
// }

// func (w *Writer) WriteStatusLine(statuscode StatusCode) error {
// 	var err error
// 	switch statuscode {
// 	case Ok:
// 		_, err = w.writer.Write([]byte(fmt.Sprintf("HTTP/1.1 %d OK\r\n", statuscode)))
// 	case Bad_Request:
// 		_, err = w.writer.Write([]byte(fmt.Sprintf("HTTP/1.1 %d Bad Request\r\n", statuscode)))
// 	case Internal_Server_Error:
// 		_, err = w.writer.Write([]byte(fmt.Sprintf("HTTP/1.1 %d Internal Server Error\r\n", statuscode)))
// 	default:
// 		_, err = w.writer.Write([]byte(fmt.Sprintf("HTTP/1.1 %d\r\n", statuscode)))
// 	}

// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (w *Writer) WriteHeaders(headers headers.Headers) error {
// 	var err error
// 	headers.ForEach(
// 		func(key, value string) {
// 			res := fmt.Sprintf("%s: %s\r\n", key, value)
// 			_, err1 := w.writer.Write([]byte(res))
// 			if err1 != nil {
// 				err = err1
// 			}
// 		},
// 	)
// 	_, err = w.writer.Write([]byte("\r\n"))
// 	return err
// }

// func (w *Writer) WriteBody(p []byte) (int, error) {
// 	return w.writer.Write(p)
// }

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

func GetDefaultHeaders(contentLen int) headers.Headers {
	heads := headers.NewHeaders()
	heads.Set("Content-Length", strconv.Itoa(contentLen))
	heads.Set("connection", "close")
	heads.Set("Content-Type", "text/plain")
	return *heads
}

func WriteHeaders(w io.Writer, heads headers.Headers) error {
	var err error
	heads.ForEach(
		func(key, value string) {
			res := fmt.Sprintf("%s: %s\r\n", key, value)
			_, err1 := w.Write([]byte(res))
			if err1 != nil {
				err = err1
			}
		},
	)
	_, err = w.Write([]byte("\r\n"))
	return err
}
