package server

import (
	"fmt"
)

func (c *Client) Send(format string, a ...any) {
	fmt.Fprintf(c.conn, format+"\n", a...)
}

func (d *Dispatcher) removeClient(client *Client) {
	if client.name == "" {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	client.Send("OK bye")
	for _, c := range(d.clients) {
		if c.name == client.name {
			delete(d.clients, c.name)
		}
	}
	client.name = ""
}
