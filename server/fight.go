package server

import (
	"fmt"
)

func (d *Dispatcher) handleAttack(client *Client, args []string) {
	fmt.Print("a")
}

func (d *Dispatcher) handleDefend(client *Client, args []string) {
	fmt.Print("a")
}

func (d *Dispatcher) handleFlee(client *Client, args []string) {
	fmt.Print("a")
}

func (d *Dispatcher) handleUseItem(client *Client, args []string) {
	fmt.Print("a")
}

func (d *Dispatcher) handleStatus(client *Client, args []string) {
	fmt.Print("a")
}
