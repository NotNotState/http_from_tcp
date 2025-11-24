package main

import (
	"fmt"
	"log"
	"net"

	"github.com/NotNotState/httpfromtcp/internal/request"
)

func main() {

	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println("Fatal Error", err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			//log.Fatal(err)
			fmt.Println("Fatal Error", err)
			continue
		}
		fmt.Println("Connection Accepted Successfully")
		recievedRequest, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf(
			"Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n",
			recievedRequest.RequestLine.Method, recievedRequest.RequestLine.RequestTarget, recievedRequest.RequestLine.HttpVersion)

		fmt.Println("Headers:")
		recievedRequest.Headers.ForEach(
			func(key, value string) {
				fmt.Printf("- %s: %s\n", key, value)
			},
		)

		fmt.Println("Body:")
		fmt.Println(string(recievedRequest.Body))

		fmt.Println("Connection to ", conn.RemoteAddr(), " has been closed")
	}

}
