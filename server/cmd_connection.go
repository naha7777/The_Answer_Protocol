package server

import (
	"strings"
)

func (d *Dispatcher) handleConnect(client *Client, args []string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if client.name != "" {
		client.Send("ERR 10 YOU_HAVE_A_NAME")
		return
	}
	if len(args) == 0 || len(args) > 1 {
		client.Send("ERR 09 THIS_IS_NOT_A_NAME")
		return
	}
	username := strings.Join(args, " ")

	for _, c := range d.clients {
		if strings.ToUpper(c.name) == strings.ToUpper(username) {
			client.Send("ERR 201 NAME_IN_USE")
			return
		}
	}

	for _, room := range d.world.Rooms {
		if room.Id == d.world.Start {
			client.room = room
		}
	}

	client.name = username
	client.room.Players = append(client.room.Players, client.name)
	d.clients[client.name] = client
	client.Send("OK connected")
}

func (d *Dispatcher) handleQuit(client *Client, args []string) {
	d.removeClient(client)
	defer d.displayRoomStats()
	defer d.displayStats()
}
