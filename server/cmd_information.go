package server

import (
	"fmt"
)

func (d *Dispatcher) handleWho(client *Client, args []string) {
	d.mu.Lock()
	nb := len(d.clients)
	d.mu.Unlock()
	client.Send("OK players=%d", nb)
}

func (d *Dispatcher) handleInventory(client *Client, args []string) {
	fmt.Print("a")
}

func (d *Dispatcher) handleLook(client *Client, args []string) {
	fmt.Print("a")
}
