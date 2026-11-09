package server

import (
	"fmt"
)

func (d *Dispatcher) handleTake(client *Client, args []string) {
	fmt.Print("a")
}

func (d *Dispatcher) handleDrop(client *Client, args []string) {
	fmt.Print("a")
}

func (d *Dispatcher) handleMove(client *Client, args []string) {
	fmt.Print("a")
}
