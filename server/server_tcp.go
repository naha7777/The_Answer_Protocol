package server

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"io"
	"strings"
	"the_answer_protocol/game"
)

func HandleConnection(connection net.Conn, dispatcher *Dispatcher,
	world *game.World) {

	defer connection.Close()
	client := &Client{conn: connection}
	for _, room := range world.Rooms {
		if room.Id == world.Start {
			client = &Client{conn: connection, name: "", group: nil, room: room,
				inventory: nil, hp: 100, hp_max: 100,
				status: "very healthy", strengh: 10, target: "", real_strengh: 10,}
		}
	}

	defer dispatcher.removeClient(client)
	client.Send("OK hello proto=1")

	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		line := strings.ToUpper(scanner.Text())
		if client.name == "" && strings.HasPrefix(line, "CONNECT") {
			dispatcher.Dispatch(client, scanner.Text())
			dispatcher.displayStats()
			dispatcher.displayRoomStats()
		} else if client.name != "" {
			dispatcher.Dispatch(client, scanner.Text())
		} else {
			client.Send("ERR 02 YOU_MUST_BE_CONNECTED")
		}
	}
	if err := scanner.Err(); err != nil && !errors.Is(err, net.ErrClosed) &&
    !errors.Is(err, io.EOF) {
		fmt.Fprintln(os.Stderr, scanner.Err())
	}
}

func Server(listener net.Listener) {

	world, err := game.LoadWorld("data/world.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}

	for _, r := range world.Rooms {
		if r.NPCs == nil { r.NPCs = []string{} }
		if r.Items == nil { r.Items = []string{} }
		if r.Players == nil { r.Players = []string{} }
	}

	dispatcher := NewDispatcher(world)

	for {
		connection, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			fmt.Fprintln(os.Stderr, "Error accepting:", err)
			continue
		}
		go HandleConnection(connection, dispatcher, world)
	}
}
