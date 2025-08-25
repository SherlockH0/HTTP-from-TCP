package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")

	if err != nil {
		log.Fatal("error", "error", err)
	}
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		}

		r, err := request.RequestFromReader(connection)

		if err != nil {
			log.Fatal("error", "error", err)
		}

		fmt.Printf(
			"Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n",
			r.RequestLine.Method,
			r.RequestLine.RequestTarget,
			r.RequestLine.HttpVersion,
		)
	}

}
