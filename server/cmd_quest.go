package server

import (
	"fmt"
	"strings"
)

func (d *Dispatcher) handleQuest(client *Client, args []string) {
	npc := strings.ToUpper(strings.Join(args, " "))
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, npcs := range client.room.NPCs {
		if strings.ToUpper(npcs) == npc || npc == strings.ToUpper(d.world.NPCs[npcs].Name){
			for _, q := range(d.world.Quests) {
				if q.Id == d.world.NPCs[npcs].Quest {
					if q.Status == "available" && q.Description != "" {
						client.Send(`OK  {"quest_id": "%s", "description": "%s", "reward": "%s", "status": "%s"}`, q.Id, q.Description, q.Reward, q.Status)
						return
					} else {
						client.Send("ERR 406 NO_QUEST_AVAILABLE")
						return
					}
				} else {
					client.Send("ERR 406 NO_QUEST_AVAILABLE")
					return
				}
			}
		}
	}
	client.Send("ERR 404 NPC_NOT_FOUND")
}

func (d *Dispatcher) handleQuests(client *Client, args []string) {
	fmt.Print("a")
}

func (d *Dispatcher) handleCompleteQuest(client *Client, args []string) {
	fmt.Print("a")
}
