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
	d.mu.Lock()
	defer d.mu.Unlock()

	ChangeStatus(client)
	client.Send("OK {'hp': %d, 'max_hp': %d, 'status': %s}", client.hp,
		client.hp_max, client.status)
}

func ChangeStatus(client *Client) {
	if client.hp <= client.hp_max && client.hp >= (80%*&client.hp_max) {
		client.status = "very healthy"
	} else if client.hp < (80%*&client.hp_max) && client.hp >= (60%*&client.hp_max) {
		client.status = "good"
	} else if client.hp < (60%*&client.hp_max) && client.hp >= (40%*&client.hp_max) {
		client.status = "mid"
	} else if client.hp < (40%*&client.hp_max) && client.hp >= (20%*&client.hp_max) {
		client.status = "suffering"
	} else if client.hp < (20%*&client.hp_max) && client.hp > 0 {
		client.status = "almost dead"
	} else if client.hp == 0 {
		client.status = "dead"
	} else {
		client.status = "far too healthy"
	}
}
