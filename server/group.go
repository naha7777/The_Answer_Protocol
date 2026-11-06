package server

import (
	"fmt"
)

func (d *Dispatcher) handleGroupCreate(client *Client, args []string) {
	fmt.Print("a")
}

func (d *Dispatcher) handleGroupInvite(client *Client, args []string) {
	fmt.Print("a")
}

func (d *Dispatcher) handleGroupJoin(client *Client, args []string) {
	fmt.Print("a")
}

func (d *Dispatcher) handleGroupLeave(client *Client, args []string) {
	fmt.Print("a")
}
