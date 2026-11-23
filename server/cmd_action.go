package server

import (
	"slices"
	"strings"
)

func (d *Dispatcher) handleTake(client *Client, args []string) {
	if len(args) == 0 {
		client.Send("ERR 404 ITEM_NOT_FOUND")
		return
	}
	query := strings.ToUpper(strings.Join(args, " "))
	d.mu.Lock()
	defer d.mu.Unlock()

	for i, id := range client.room.Items {
		item, ok := d.world.Items[id]
		if !ok {
			continue
		}
		if !strings.EqualFold(id, query) && !strings.EqualFold(item.Name, query) {
			continue
		}
		if !item.Obtainable {
			client.Send("ERR 07 ITEM_NOT_OBTAINABLE")
			return
		}
		client.room.Items = slices.Delete(client.room.Items, i, i+1)
		client.inventory = append(client.inventory, item)
		client.Send("OK taken=%s", id)
		SendRoom(d, client, "EVT PLAYER TOOK %s", id)
		return
	}
	client.Send("ERR 404 ITEM_NOT_FOUND")
}

func (d *Dispatcher) handleDrop(client *Client, args []string) {
	if len(args) == 0 {
		client.Send("ERR 404 ITEM_NOT_IN_INVENTORY")
		return
	}
	query := strings.ToLower(strings.Join(args, " "))
	d.mu.Lock()
	defer d.mu.Unlock()

	RemoveItemInventory(client, query, d, true)
}

func (d *Dispatcher) handleMove(client *Client, args []string) {
	if len(args) == 0 {
		client.Send("ERR 301 NO_EXIT")
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	switch direction := strings.ToUpper(args[0]); direction {
	case "NORTH":
		Move("NORTH", client, d)
	case "SOUTH":
		Move("SOUTH", client, d)
	case "EAST":
		Move("EAST", client, d)
	case "WEST":
		Move("WEST", client, d)
	default:
		client.Send("ERR 301 NO_EXIT")
		return
	}
}

