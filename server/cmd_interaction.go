package server

import (
	"fmt"
	"strings"
)

func (d *Dispatcher) handleChat(client *Client, args []string) {
	message := strings.Join(args[1:], " ")
	d.mu.Lock()
	defer d.mu.Unlock()
	switch receiver := strings.ToUpper(args[0]); receiver {
	case "GLOBAL":
		for _, c := range(d.clients) {
			if client == c {
				c.Send("OK")
			} else {
				c.Send("%s: %s", client.name, message)
			}
		}
	case "GROUP":
		if client.group != nil {
			for _, c := range(client.group.members) {
				if client == c {
					c.Send("OK")
				} else {
					c.Send("%s: %s", client.name, message)
				}
			}
		} else {
			client.Send("ERR 401 NOT_IN_GROUP")
		}
	case "ROOM":
		if client.room != nil {
			for _, c := range(client.room.visitors) {
				if client == c {
					c.Send("OK")
				} else {
					c.Send("%s: %s", client.name, message)
				}
			}
		} else {
			client.Send("ERR 401 NOT_IN_GROUP")
		}
	default:
		client.Send("ERR 110 NO_CHAT_SPECIFIED")
		return
	}
}

func (d *Dispatcher) handleTalk(client *Client, args []string) {
	fmt.Print("a")
}
