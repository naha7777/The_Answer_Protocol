package server

import "slices"

func SendRoom(d *Dispatcher, client *Client, format string, a... any) {
	for _, c := range d.clients {
		if slices.Contains(client.room.Players, c.name) {
			if c.name != client.name {
				c.Send(format, a...)
			}
		}
	}
}

func SendAll(d *Dispatcher, client *Client, format string, a... any) {
	for _, c := range d.clients {
		if c != client {
			c.Send(format, a...)
		}
	}
}

func SendGroup(client *Client, group *Group, format string, a... any) {
	for _, c := range group.members {
		c.Send(format, a...)
	}
}
