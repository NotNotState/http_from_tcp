package server

import (
	"fmt"
	"io"
	"net"
	"sync/atomic"

	"github.com/NotNotState/httpfromtcp/internal/request"
	"github.com/NotNotState/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

type HandleError struct {
	StatusCode response.StatusCode
	Message    string
}

func (hErr *HandleError) WriteHandlerError(w io.Writer) {
	response.WriteStatusLine(w, hErr.StatusCode)
	message := hErr.Message
	body_length := len(message)
	defaultHeaders := response.GetDefaultHeaders(body_length)
	response.WriteHeaders(w, defaultHeaders)
	w.Write([]byte(message))
}

type Handler func(w *response.Writer, req *request.Request)

//type Handler func(w io.Writer, req *request.Request) *HandleError

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		handler_error := HandleError{
			StatusCode: response.Bad_Request,
			Message:    err.Error(),
		}
		handler_error.WriteHandlerError(conn)
		return
	}

	rw := response.NewWriter(conn)
	s.handler(rw, req)
}

func (s *Server) listen() { // modify to take in request.Request
	for {
		conn, err := s.listener.Accept()
		if s.closed.Load() {
			return
		}
		if err != nil {
			return
		}
		go s.handle(conn) // every successful non-server closed connection will be spun off as it's own routine
	}
}

func Serve(port uint16, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	// how and what do i put within the Hanlder here?
	//add handler to server
	server := &Server{
		listener: listener,
		handler:  handler,
	}
	go server.listen() // spawns go routine which will spin in the backround until stopped
	return server, nil
}
