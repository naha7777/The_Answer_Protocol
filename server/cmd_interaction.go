package server

import (
	"fmt"
)

func (d *Dispatcher) handleChat(client *Client, args []string) {
	fmt.Print("a")
}

func (d *Dispatcher) handleTalk(client *Client, args []string) {
	fmt.Print("a")
}

func (d *Dispatcher) handleWho(client *Client, args []string) {
	fmt.Print("a")
}
