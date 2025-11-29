package server

import (
	"fmt"
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

type Handler func(w *response.Writer, req *request.Request)

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
	rw := response.NewWriter(conn)
	if err != nil {
		rw.WriteStatusLine(response.Bad_Request)
		body := []byte("Error Parsing Request!")
		rw.WriteHeaders(response.GetDefaultHeaders(len(body)))
		rw.WriteBody(body)
		return
	}
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
	//add handler to server
	server := &Server{
		listener: listener,
		handler:  handler,
	}
	go server.listen() // spawns go routine which will spin in the backround until stopped
	return server, nil
}
