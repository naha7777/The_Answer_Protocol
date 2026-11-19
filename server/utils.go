package server

import (
	"strings"
	"slices"
	"the_answer_protocol/game"
)

func (d *Dispatcher) removeClient(client *Client) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if client.name == "" {
		return
	}
	if client.group != nil {
		delete(client.group.members, client.name)
		client.group = nil
	}
	if client.inventory != nil {
		for _, item := range client.inventory {
			client.room.Items = append(client.room.Items, item.Id)
		}
		client.inventory = nil
	}
	if client.room != nil {
		client.room.Players = remove(client.name, client.room.Players)
		client.room = nil
	}
	delete(d.clients, client.name)
	client.name = ""
	client.Send("OK bye")
}

func remove(name string, place []string) []string {
	var new []string
	for _, item := range place {
		if item != name {
			new = append(new, item)
		}
	}
	return new
}

func (d *Dispatcher) displayStats() {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, c := range d.clients {
		c.Send("EVT STATS players=%d", len(d.clients))
	}
}

func (d *Dispatcher) displayRoomStats() {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, c := range d.clients {
		for _, str := range c.room.Players {
			if c.name == str {
				c.Send("EVT STATS players_in_room=%d", len(c.room.Players))
			}
		}
	}
}

func Move(direction string, client *Client, d *Dispatcher) {
	find := false
	for k, v := range client.room.Exits {
		if strings.ToUpper(k) == direction {
			for _, val := range d.world.Rooms {
				if v == val.Id {
					for _, c := range d.clients {
						if slices.Contains(client.room.Players, c.name) {
							if c.name != client.name {
								c.Send("EVT ROOM PRESENCE LEAVE %s", client.name)
							}
						}
					}
					client.room.Players = remove(client.name, client.room.Players)
					client.room = val
					client.room.Players = append(client.room.Players, client.name)
					client.Send("OK room=%s", v)
					find = true
					for _, c := range d.clients {
						if slices.Contains(client.room.Players, c.name) {
							if c.name != client.name {
								c.Send("EVT ROOM PRESENCE ENTER %s", client.name)
							}
						}
					}
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

func RemoveItemInventory(client *Client, item_id string) {
	final_inventory := []*game.Item{}
	count := 0
	for _, i := range client.inventory {
		if i.Id == item_id && count == 0{
			count += 1
		} else {
			final_inventory = append(final_inventory, i)
		}
	}
	client.inventory = final_inventory
}
