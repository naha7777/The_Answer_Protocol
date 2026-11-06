package server

import (
	"fmt"
)

func (c *Client) Send(format string, a ...any) {
	fmt.Fprintf(c.conn, format+"\n", a...)
}

func (d *Dispatcher) removeClient(c *Client) {
	if c.name == "" {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	c.Send("OK bye")
	d.names[c.name] = false
	c.name = ""
}
