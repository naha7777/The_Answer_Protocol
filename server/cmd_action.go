package server

import (
	"strings"
	"slices"
)

func (d *Dispatcher) handleTake(client *Client, args []string) {
	if len(args) == 0 {
		client.Send("ERR 115 MISSING_ITEM")
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
			client.Send("ERR 116 ITEM_NOT_OBTAINABLE")
			return
		}
		if item.Owner != "" {
			client.Send("ERR 404 ITEM_NOT_FOUND")
			return
		}
		client.room.Items = slices.Delete(client.room.Items, i, i+1)
		client.inventory = append(client.inventory, item)
		item.Owner = client.name
		client.Send("OK taken=%s", id)
		return
	}
	client.Send("ERR 404 ITEM_NOT_FOUND")
}

func (d *Dispatcher) handleDrop(client *Client, args []string) {
	if len(args) == 0 {
		client.Send("ERR 115 MISSING_ITEM")
		return
	}
	query := strings.ToUpper(strings.Join(args, " "))
	d.mu.Lock()
	defer d.mu.Unlock()

	for i, item := range client.inventory {
		if !strings.EqualFold(item.Id, query) && !strings.EqualFold(item.Name, query) {
			continue
		}
		client.room.Items = append(client.room.Items, item.Id)
		client.inventory = slices.Delete(client.inventory, i, i+1)
		item.Owner = ""
		client.Send("OK dropped=%s", item.Id)
		return
	}
	client.Send("ERR 404 ITEM_NOT_FOUND")
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
	case "EST":
		Move("EST", client, d)
	case "WEST":
		Move("WEST", client, d)
	default:
		client.Send("ERR 301 NO_EXIT")
		return
	}
}

func Move(direction string, client *Client, d *Dispatcher) {
	find := false
	for k, v := range(client.room.Exits) {
		if strings.ToUpper(k) == direction {
			for _, val := range(d.world.Rooms) {
				if v == val.Id {
					client.room.Players = remove(client.name, client.room.Players)
					client.room = val
					client.room.Players = append(client.room.Players, client.name)
					client.Send("OK room=%s", v)
					find = true
					return
				}
			}
		}
	}
	if !find {
		client.Send("ERR 301 NO_EXIT")
		return
	}
}
