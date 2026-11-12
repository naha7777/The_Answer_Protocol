package server

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func (d *Dispatcher) handleWho(client *Client, args []string) {
	d.mu.Lock()
	nb := len(d.clients)
	d.mu.Unlock()
	client.Send("OK players=%d", nb)
}

func (d *Dispatcher) handleLook(client *Client, args []string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	jsonData, err := json.Marshal(client.room)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
		return
	}
	client.Send(`OK {"room": %s}`, jsonData)
}

func (d *Dispatcher) handleInventory(client *Client, args []string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	var item_lst []string
	for _, v := range client.inventory {
		item_lst = append(item_lst, v.Id)
	}
	lst := strings.Join(item_lst, " ")
	client.Send("OK [%s]", lst)
}
