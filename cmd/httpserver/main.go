package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/NotNotState/httpfromtcp/internal/request"
	"github.com/NotNotState/httpfromtcp/internal/response"
	"github.com/NotNotState/httpfromtcp/internal/server"
)

const port = 42069

func handler(w io.Writer, req *request.Request) *server.HandleError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandleError{
			StatusCode: response.Bad_Request,
			Message:    "Your problem is not my problem\n",
		}
	case "/myproblem":
		return &server.HandleError{
			StatusCode: response.Internal_Server_Error,
			Message:    "Woopsie, my bad\n",
		}
	default:
		w.Write([]byte("All good, frfr\n"))
		return nil
	}
}

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	// will add a signal to sigChan when any program interrupts occure
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // will block until sigChan can pop/discard an item
	log.Println("Server gracefully stopped")
}
