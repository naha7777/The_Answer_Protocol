package server

import (
	"strings"
)

func (d *Dispatcher) handleConnect(client *Client, args []string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	// return an error for invalid names
	if client.name != "" {
		client.Send("ERR 010 YOU_HAVE_A_NAME")
		return
	}
	if len(args) == 0 || len(args) > 1 {
		client.Send("ERR 009 THIS_IS_NOT_A_NAME")
		return
	}
	username := strings.Join(args, " ")

	for _, c := range d.clients {
		if strings.EqualFold(strings.ToUpper(c.name), strings.ToUpper(username)) {
			client.Send("ERR 201 NAME_IN_USE")
			return
		}
	}

	// username length limit
	if len(username) > 20 {
		client.Send("ERR 021 NAME_TOO_LONG")
		return
	}

	// put the actual room on the client struct
	for _, room := range d.world.Rooms {
		if room.Id == d.world.Start {
			client.room = room
		}
	}
	// save the client everywhere : dispatcher, room... and Send a success msg
	client.name = username
	client.room.Players = append(client.room.Players, client.name)
	d.clients[client.name] = client
	client.Send("OK connected")
	SendRoom(d, client, "EVT ROOM PRESENCE ENTER %s", client.name)
}

func (d *Dispatcher) handleQuit(client *Client, args []string) {
	// Send a success msg and save that the client quit
	client.Send("OK bye")
	client.quit = true
}
