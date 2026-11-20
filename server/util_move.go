package server

import "strings"

func Move(direction string, client *Client, d *Dispatcher) {
	find := false
	for k, v := range client.room.Exits {
		if strings.ToUpper(k) == direction {
			for _, val := range d.world.Rooms {
				if v == val.Id {
					SendRoom(d, client, "EVT ROOM PRESENCE LEAVE %s", client.name)
					client.room.Players = Remove(client.name, client.room.Players)
					client.room = val
					client.room.Players = append(client.room.Players, client.name)
					client.Send("OK room=%s", v)
					find = true
					SendRoom(d, client, "EVT ROOM PRESENCE ENTER %s", client.name)
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
