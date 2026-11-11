package main

import (
	"net"
	"the_answer_protocol/server")

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		Gui()
	}
	defer listener.Close()

	server.Server(listener)
}
