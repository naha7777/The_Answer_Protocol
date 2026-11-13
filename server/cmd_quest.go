package server

import (
	"strings"
	"slices"
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
						client.quests = append(client.quests, q)
						client.Send(`OK  {"quest_id": "%s", "description": "%s", "reward": "%s", "status": "%s"}`, q.Id, q.Description, q.Reward, q.Status)
						q.Status = "active"
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
	res := "OK ["
	d.mu.Lock()
	defer d.mu.Unlock()
	for i, v := range client.quests {
		res += "{"
		res += `"quest_id: "`
		res += v.Id
		res += `", "status": "`
		res += v.Status
		if v.Status != "completed" {
			res += `", "progress": "0/1"`
		} else {
			res += `"`
		}
		if i != len(client.quests) - 1 {
			res += "},"
		} else {
			res += "}"
		}
	}
	res += "]"
	client.Send("OK %s", res)
}

func (d *Dispatcher) handleCompleteQuest(client *Client, args []string) {
	find := false
	presence := false
	completed := false
	d.mu.Lock()
	defer d.mu.Unlock()
	if len(args) == 0 {
		client.Send("ERR 406 NO_QUEST_AVAILABLE")
		return
	}
	for i, quest := range client.quests {
		if quest.Id == args[0] {
			find =  true
			if slices.Contains(client.room.NPCs, quest.NPC) {
				presence = true
				parts := strings.Split(strings.ToUpper(quest.Description), " ")
				for _, obj := range client.inventory {
					if slices.Contains(parts, strings.ToUpper(obj.Name)) || slices.Contains(parts, strings.ToUpper(obj.Id)){
						completed = true
					}
				}
				if completed {
					quest.Status = "completed"
					client.Send("OK %s: %s, reward: %s", quest.Id, quest.Status, quest.Reward)
				} else {
					client.Send("ERR 11 INCOMPLETE_QUEST_REQUIREMENTS")
					return
				}
			} else if quest.Id != "active" {
				client.Send("ERR 406 NO_QUEST_AVAILABLE")
				return
			}
		} else if i == len(client.quests) -1 {
			client.Send("ERR 406 NO_QUEST_AVAILABLE")
			return
		}
	}
	if !find {
		client.Send("ERR 406 NO_QUEST_AVAILABLE")
		return
	} else if !presence {
		client.Send("ERR 404 NPC_NOT_FOUND")
		return
	}

}
