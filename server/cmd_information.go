package server

import (
	"fmt"
)

func (d *Dispatcher) handleInventory(client *Client, args []string) {
	fmt.Print("a")
}

func (d *Dispatcher) handleLook(client *Client, args []string) {
	fmt.Print("a")
}
