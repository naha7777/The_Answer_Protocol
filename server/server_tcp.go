package server

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"the_answer_protocol/game"
)

func HandleConnection(connection net.Conn, dispatcher *Dispatcher,
	world *game.World) {

	defer connection.Close()
	client := &Client{conn: connection}
	for _, room := range world.Rooms {
		if room.Id == world.Start {
			// item := &game.Item{Name: "Herbs", Description: "de l'herbe",
			// 	Obtainable: true, Id: "item.herbs"}
			// var inv []*game.Item
			// inv = append(inv, item)
			client = &Client{conn: connection, name: "", group: nil, room: room,
				inventory: nil, hp: 100, hp_max: 100,
				status: "very healthy", strengh: 10, target: "",
				real_strengh: 10, quit: false}
		}
	}

	defer func() {
		dispatcher.removeClient(client)
		dispatcher.displayRoomStats()
		dispatcher.displayStats()
	}()
	client.Send("OK hello proto=1")

	reader := bufio.NewReaderSize(connection, 1024)
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			break
		}
		if isPrefix {
			for isPrefix && err == nil {
				_, isPrefix, err = reader.ReadLine()
			}
			if err != nil {
				break
			}
			client.Send("ERR 16 LINE_TOO_LONG")
			continue
		}

		text := string(line)
		upper := strings.ToUpper(text)
		if client.name == "" && strings.HasPrefix(upper, "CONNECT") {
			dispatcher.Dispatch(client, text)
			dispatcher.displayStats()
			dispatcher.displayRoomStats()
		} else if client.name != "" {
			dispatcher.Dispatch(client, text)
		} else {
			client.Send("ERR 02 YOU_MUST_BE_CONNECTED")
		}
		if client.quit {
			break
		}
	}
}

func Server(listener net.Listener) {

	world, err := game.LoadWorld("data/world.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}

	for _, r := range world.Rooms {
		if r.NPCs == nil {
			r.NPCs = []string{}
		}
		if r.Items == nil {
			r.Items = []string{}
		}
		if r.Players == nil {
			r.Players = []string{}
		}
	}

	dispatcher := NewDispatcher(world)

	for {
		connection, err := listener.Accept()
		if err != nil {
			// verrification vu que listener.Close est appele listenner.Accept
			// retourne une net.ErrClosed donc on a ferme le listenner volontairement
			if errors.Is(err, net.ErrClosed) {
				return
			}
			fmt.Fprintln(os.Stderr, "Error accepting:", err)
			continue
		}
		go HandleConnection(connection, dispatcher, world)
	}
}
