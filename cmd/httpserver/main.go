package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/NotNotState/httpfromtcp/internal/headers"
	"github.com/NotNotState/httpfromtcp/internal/request"
	"github.com/NotNotState/httpfromtcp/internal/response"
	"github.com/NotNotState/httpfromtcp/internal/server"
)

const port = 42069

func ProxyHandler(w *response.Writer, heads headers.Headers, target string) error {
	requestTarget := "https://httpbin.org/" + target
	buffer := make([]byte, 1024)
	resp, err := http.Get(requestTarget)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	heads.Delete("content-length")
	heads.Delete("connection")
	heads.SetHard("Content-Type", "text/plain")
	heads.Set("Transfer-Encoding", "chunked")
	w.WriteStatusLine(response.Ok)
	w.WriteHeaders(heads)

	var n int
	for {
		n, err = resp.Body.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {

				break
			}
			return err
		}
		//fmt.Printf("Reading %d bytes in proxy\n", n)
		n, err = w.WriteChunkedBody(buffer[:n])
		if err != nil {
			return err
		}
		//fmt.Printf("Writing %d bytes to Client\n", n)
	}

	n, err = w.WriteChunkedBodyDone()
	if err != nil {
		return err
	}
	//fmt.Printf("Writing %d bytes to close response to Client\n", n)
	return nil
}

func handler(w *response.Writer, req *request.Request) {
	heads := response.GetDefaultHeaders(0)
	if req.RequestLine.RequestTarget == "/yourproblem" {
		respBody := []byte("Your request honestly kinda sucked.\n")
		respLen := len(respBody)
		heads.SetHard("Content-Type", "text/html")
		heads.SetHard("Content-Length", strconv.Itoa(respLen))
		w.WriteStatusLine(response.Bad_Request)
		w.WriteHeaders(heads)
		w.WriteBody(respBody)
	} else if req.RequestLine.RequestTarget == "/myproblem" {
		respBody := []byte("Okay, you know what? This one is on me.\n")
		respLen := len(respBody)
		heads.SetHard("Content-Type", "text/html")
		heads.SetHard("Content-Length", strconv.Itoa(respLen))
		w.WriteStatusLine(response.Internal_Server_Error)
		w.WriteHeaders(heads)
		w.WriteBody(respBody)
	} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
		ProxyHandler(w, heads, target)
	} else {
		respBody := []byte("Your request was an absolute banger.\n")
		respLen := len(respBody)
		heads.SetHard("Content-Type", "text/html")
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
