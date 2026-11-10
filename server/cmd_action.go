package server

import (
	"fmt"
	"strings"
)

func (d *Dispatcher) handleTake(client *Client, args []string) {
	fmt.Print("a")
}

func (d *Dispatcher) handleDrop(client *Client, args []string) {
	fmt.Print("a")
}

func (d *Dispatcher) handleMove(client *Client, args []string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	switch direction := strings.ToUpper(args[0]); direction {
	case "NORTH":
		find := false
		for k, v := range(client.room.Exits) {
			if strings.ToUpper(k) == "NORTH" {
				for _, val := range(d.world.Rooms) {
					if v == val.Id {
						client.room.Players = remove(client.name, client.room.Players)
						client.room = val
						client.room.Players = append(client.room.Players, client.name)
						client.Send("OK room=%s", v)
						find = true
					}
				}
			}
		}
		if !find {
			client.Send("ERR 112 NO_ROOM_IN_THIS_DIRECTION")
			return
		}
	case "SOUTH":
		find := false
		for k, v := range(client.room.Exits) {
			if strings.ToUpper(k) == "SOUTH" {
				for _, val := range(d.world.Rooms) {
					if v == val.Id {
						client.room.Players = remove(client.name, client.room.Players)
						client.room = val
						client.room.Players = append(client.room.Players, client.name)
						client.Send("OK room=%s", v)
						find = true
					}
				}
			}
		}
		if !find {
			client.Send("ERR 112 NO_ROOM_IN_THIS_DIRECTION")
			return
		}
	case "EST":
		find := false
		for k, v := range(client.room.Exits) {
			if strings.ToUpper(k) == "EST" {
				for _, val := range(d.world.Rooms) {
					if v == val.Id {
						client.room.Players = remove(client.name, client.room.Players)
						client.room = val
						client.room.Players = append(client.room.Players, client.name)
						client.Send("OK room=%s", v)
						find = true
					}
				}
			}
		}
		if !find {
			client.Send("ERR 112 NO_ROOM_IN_THIS_DIRECTION")
			return
		}
	case "WEST":
		find := false
		for k, v := range(client.room.Exits) {
			if strings.ToUpper(k) == "WEST" {
				for _, val := range(d.world.Rooms) {
					if v == val.Id {
						client.room.Players = remove(client.name, client.room.Players)
						client.room = val
						client.room.Players = append(client.room.Players, client.name)
						client.Send("OK room=%s", v)
						find = true
					}
				}
			}
		}
		if !find {
			client.Send("ERR 112 NO_ROOM_IN_THIS_DIRECTION")
			return
		}
	default:
		client.Send("ERR 111 INVALID_DIRECTION")
		return
	}
}
