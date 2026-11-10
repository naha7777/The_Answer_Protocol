package server

import (
	"encoding/json"
	"fmt"
	"os"
)

func (d *Dispatcher) handleWho(client *Client, args []string) {
	d.mu.Lock()
	nb := len(d.clients)
	d.mu.Unlock()
	client.Send("OK players=%d", nb)
}

func (d *Dispatcher) handleLook(client *Client, args []string) {
	jsonData, err := json.Marshal(client.room)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
		return
	}
	client.Send("OK %s", jsonData)
}

func (d *Dispatcher) handleInventory(client *Client, args []string) {
	fmt.Print("a")
}
