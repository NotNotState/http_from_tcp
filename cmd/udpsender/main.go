package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal("Unable to find address", err)
	}

	updConn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal("Unable to Establish UDP Dial", err)
	}
	defer updConn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Error in reading user input", err)
		}
		n, err := updConn.Write([]byte(input))

		if err != nil {
			log.Fatal("Error in writing user input", err)
		}

		fmt.Println("Wrote", n, "bytes to upd stream")
	}

}
