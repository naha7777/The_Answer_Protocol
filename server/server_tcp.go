package server

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	conn net.Conn
	name string
	// room
	// quest
	// health ...
}


func handleClient(connection net.Conn, dispatcher *Dispatcher) {

	defer connection.Close()
	client := &Client{conn: connection, name: ""}

	defer dispatcher.removeClient(client)

	client.Send("OK hello proto=1")
	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		dispatcher.Dispatch(client, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Print(scanner.Err())
	}
}

func Server() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	dispatcher := NewDispatcher()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		go handleClient(connection, dispatcher)
	}
}
