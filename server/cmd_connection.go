package server

import (
	"strings"
)

func (d *Dispatcher) handleConnect(client *Client, args []string) {
	if client.name != "" {
		client.Send("ERR 10 YOU_HAVE_A_NAME")
		return
	}
	if len(args) == 0 {
		client.Send("ERR 09 THIS_IS_NOT_A_NAME")
		return
	}
	username := strings.Join(args, " ")

	d.mu.Lock()
	defer d.mu.Unlock()
	for _, c := range d.clients {
		if c.name == username {
			client.Send("ERR 201 NAME_IN_USE")
			return
		}
	}
	client.name = username
	client.room.Players = append(client.room.Players, client.name)
	d.clients[client.name] = client
	client.Send("OK connected")
}

func (d *Dispatcher) handleQuit(client *Client, args []string) {
	defer d.removeClient(client)
}
