package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Item struct {
	name	string
}

type Room struct {
	visitors	map[string]*Client
	items		map[string]*Item
}

type Client struct {
	conn	net.Conn
	name	string
	group	*Group
	room	*Room
}

func handleClient(connection net.Conn, dispatcher *Dispatcher) {

	defer connection.Close()
	client := &Client{conn: connection, name: "", group: nil}

	defer dispatcher.removeClient(client)

	client.Send("OK hello proto=1")
	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		line := strings.ToUpper(scanner.Text())
		if client.name == "" && strings.HasPrefix(line, "CONNECT") {
			dispatcher.Dispatch(client, scanner.Text())
		} else if client.name != "" {
			dispatcher.Dispatch(client, scanner.Text())
		} else {
			client.Send("ERR 100 YOU_MUST_BE_CONNECTED")
		}
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
