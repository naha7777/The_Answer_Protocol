package server

import (
	"strings"
)

func (d *Dispatcher) handleConnect(client *Client, args []string) {
	if client.name != "" {
		client.Send("ERR 100 YOU_HAVE_A_NAME")
		return
	}
	if len(args) == 0 {
		client.Send("ERR 101 THIS_IS_NOT_A_NAME")
		return
	}
	username := strings.Join(args, " ")

	d.mu.Lock()
	defer d.mu.Unlock()
	if d.names[username] {
		client.Send("ERR 201 NAME_IN_USE")
		return
	}
	d.names[username] = true
	client.name = username
	client.Send("OK connected")
}

func (d *Dispatcher) handleQuit(client *Client, args []string) {
	defer d.removeClient(client)
}
