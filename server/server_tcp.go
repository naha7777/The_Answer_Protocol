package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"the_answer_protocol/game"
)

func handleClient(connection net.Conn, dispatcher *Dispatcher,
				  world *game.World) {

	defer connection.Close()
	client := &Client{conn: connection}
	for _, room := range world.Rooms {
		if room.Id == world.Start {
			client = &Client{conn: connection, name: "", group: nil, room: room,
							 inventory: nil, hp: 100, hp_max: 100,
							 status: "very healthy"}
		}
	}

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
		fmt.Fprintln(os.Stderr, scanner.Err())
	}
}

func Server() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error listening:", err)
		return
	}
	defer listener.Close()

	world, err := game.LoadWorld("data/world.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}

	dispatcher := NewDispatcher(world)

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error accepting:", err)
			continue
		}
		go handleClient(connection, dispatcher, world)
	}
}
