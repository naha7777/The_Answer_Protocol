package server

import (
	"strings"
)

func (d *Dispatcher) handleChat(client *Client, args []string) {
	if len(args) <= 1 {
		client.Send("ERR 08 MISSING_MESSAGE")
		return
	}
	message := strings.Join(args[1:], " ")
	d.mu.Lock()
	defer d.mu.Unlock()
	switch receiver := strings.ToUpper(args[0]); receiver {
	case "GLOBAL":
		client.Send("OK")
		SendAll(d, client, "EVT GLOBAL CHAT %s %s", client.name, message)
	case "GROUP":
		if client.group != nil {
			client.Send("OK")
			SendGroup(client, client.group, "EVT GROUP CHAT %s %s", client.name, message)
		} else {
			client.Send("ERR 401 NOT_IN_GROUP")
		}
	case "ROOM":
		if client.room != nil {
			client.Send("OK")
			SendRoom(d, client, "EVT ROOM CHAT %s %s", client.name, message)
		} else {
			client.Send("ERR 401 NOT_IN_GROUP")
		}
	default:
		client.Send("ERR 06 NO_CHAT_SPECIFIED")
		return
	}
}

func (d *Dispatcher) handleTalk(client *Client, args []string) {
	npc := strings.ToUpper(strings.Join(args, " "))
	for _, npcs := range client.room.NPCs {
		if strings.ToUpper(npcs) == npc || npc == strings.ToUpper(d.world.NPCs[npcs].Name) {
			client.Send("OK talk=%s: %v", npcs, d.world.NPCs[npcs].Dialogue[0])
			return
		}
	}
	client.Send("ERR 404 NPC_NOT_FOUND")
}
