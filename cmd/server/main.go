package main

import (
	"fmt"
	"net"
	"os"
	"the_answer_protocol/server"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listening : %s", err)
		os.Exit(1)
	}
	defer listener.Close()

	server.Server(listener)
}
