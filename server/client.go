package server

import (
	"fmt"
	"net"
	"slices"
	"the_answer_protocol/game"
)

type Client struct {
	conn      net.Conn
	name      string
	group     *Group
	room      *game.Room
	inventory []*game.Item
	hp        int
	hp_max    int
	status    string
	quests    []*game.Quest
	strengh   int
	target    string
}

func (c *Client) Send(format string, a ...any) {
	fmt.Fprintf(c.conn, format+"\n", a...)
}

func (c *Client) Name() string {
	return c.name
}

func (client *Client) Death(d *Dispatcher) {
	var final_inventory []*game.Item
	for _, item := range client.inventory {
		client.room.Items = append(client.room.Items, item.Id)
	}
	client.inventory = final_inventory

	var final_quests []*game.Quest
	for _, quest := range client.quests {
		quest.Status = "available"
	}
	client.quests = final_quests

	client.Send(`EVT DEATH`)
	client.status = "exploration"
	client.target = ""
	for _, c := range d.clients {
		if c != client {
			c.Send(`EVT dead=%s`, client.name)
		}
	}

	for _, val := range d.world.Rooms {
		if d.world.Respawn == val.Id {
			for _, c := range d.clients {
				if slices.Contains(client.room.Players, c.name) {
					if c.name != client.name {
						c.Send(`EVT ROOM PRESENCE LEAVE %s`, client.name)
					}
				}
			}
			client.room.Players = remove(client.name, client.room.Players)
			client.room = val
			client.room.Players = append(client.room.Players, client.name)
			client.Send(`OK room=%s`, d.world.Respawn)
			for _, c := range d.clients {
				if slices.Contains(client.room.Players, c.name) {
					if c.name != client.name {
						c.Send(`EVT ROOM PRESENCE ENTER %s`, client.name)
					}
				}
			}
		}
	}
	client.hp = client.hp_max
}
