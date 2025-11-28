package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/NotNotState/httpfromtcp/internal/request"
	"github.com/NotNotState/httpfromtcp/internal/response"
	"github.com/NotNotState/httpfromtcp/internal/server"
)

const port = 42069

func handler(w *response.Writer, req *request.Request) {
	heads := response.GetDefaultHeaders(0)
	heads.SetHard("Content-Type", "text/html")

	switch req.RequestLine.RequestTarget {

	case "/yourproblem":
		respBody := []byte("Your request honestly kinda sucked.\n")
		respLen := len(respBody)
		heads.SetHard("Content-Length", strconv.Itoa(respLen))
		w.WriteStatusLine(response.Bad_Request)
		w.WriteHeaders(heads)
		w.WriteBody(respBody)

	case "/myproblem":
		respBody := []byte("Okay, you know what? This one is on me.\n")
		respLen := len(respBody)
		heads.SetHard("Content-Length", strconv.Itoa(respLen))
		w.WriteStatusLine(response.Internal_Server_Error)
		w.WriteHeaders(heads)
		w.WriteBody(respBody)

	default:
		respBody := []byte("Your request was an absolute banger.\n")
		respLen := len(respBody)
		heads.SetHard("Content-Length", strconv.Itoa(respLen))
		w.WriteStatusLine(response.Ok)
		w.WriteHeaders(heads)
		w.WriteBody(respBody)
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
