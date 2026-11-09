package server

import (
	"strings"
)

func (d *Dispatcher) handleConnect(client *Client, args []string) {
	if client.name != "" {
		client.Send("ERR 101 YOU_HAVE_A_NAME")
		return
	}
	if len(args) == 0 {
		client.Send("ERR 102 THIS_IS_NOT_A_NAME")
		return
	}
	username := strings.Join(args, " ")

	d.mu.Lock()
	defer d.mu.Unlock()
	for _, c := range(d.clients) {
		if c.name == username {
			client.Send("ERR 201 NAME_IN_USE")
			return
		}
	}
	client.name = username
	d.clients[client.name] = client
	client.Send("OK connected")
}

func (d *Dispatcher) handleQuit(client *Client, args []string) {
	defer d.removeClient(client)
}
