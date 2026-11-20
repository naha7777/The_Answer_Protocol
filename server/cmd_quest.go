package server

import (
	"slices"
	"strconv"
	"strings"
	"the_answer_protocol/game"
)

func (d *Dispatcher) handleQuest(client *Client, args []string) {
	npc := strings.ToUpper(strings.Join(args, " "))
	d.mu.Lock()
	defer d.mu.Unlock()
	find := false
	for _, npcs := range client.room.NPCs {
		if strings.ToUpper(npcs) == npc || npc == strings.ToUpper(d.world.NPCs[npcs].Name) {
			for _, q := range d.world.Quests {
				for _, quest := range d.world.NPCs[npcs].Quest {
					if q.Id == quest {
						if q.Status == "available" && q.Description != "" {
							client.quests = append(client.quests, q)
							client.Send(`OK quest={"quest_id": "%s", "description": "%s", "reward": "%s", "status": "%s"}`, q.Id, q.Description, q.Reward, q.Status)
							q.Status = "active"
							find = true
							return
						}
					}
				}
			}
		}
	}
	if !find {
		client.Send("ERR 406 NO_QUEST_AVAILABLE")
		return
	}
	client.Send("ERR 404 NPC_NOT_FOUND")
}

func (d *Dispatcher) handleQuests(client *Client, args []string) {
	res := "OK quests=["
	d.mu.Lock()
	defer d.mu.Unlock()
	for i, v := range client.quests {
		res += "{"
		res += `"quest_id": "`
		res += v.Id
		res += `", "status": "`
		res += v.Status
		if v.Status != "completed" {
			nb := 0
			for _, obj := range v.Objective {
				for _, v := range client.inventory {
					if v.Id == obj {
						nb += 1
					}
				}
				for _, kill := range client.kills {
					if kill == obj {
						nb += 1
					}
				}
			}
			res += `", "progress": "` + strconv.Itoa(nb) + "/"
			res += strconv.Itoa(len(v.Objective))
		}
		res += `"`
		if i != len(client.quests)-1 {
			res += "},"
		} else {
			res += "}"
		}
	}
	res += "]"
	client.Send("%s", res)
}

func (d *Dispatcher) handleCompleteQuest(client *Client, args []string) {
	find := false
	completed := false
	d.mu.Lock()
	defer d.mu.Unlock()
	if len(args) == 0 {
		client.Send("ERR 406 NO_QUEST_AVAILABLE")
		return
	}
	for _, quest := range client.quests {
		if quest.Id == args[0] {
			find = true
			if slices.Contains(client.room.NPCs, quest.NPC) {
				if quest.Status != "active" {
					client.Send("ERR 406 NO_QUEST_AVAILABLE")
					return
				}
				completed = CheckObjective(completed, quest, client)
				if completed {
					CompleteObjective(quest, client, d)
					return
				} else {
					client.Send("ERR 11 INCOMPLETE_QUEST_REQUIREMENTS")
					return
				}
			} else {
				client.Send("ERR 404 NPC_NOT_FOUND")
				return
			}
		}
	}
	if !find {
		client.Send("ERR 406 NO_QUEST_AVAILABLE")
		return
	}
}

func (d *Dispatcher) handleAbandonQuest(client *Client, args []string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if len(args) == 0 {
		client.Send("ERR 406 NO_QUEST_AVAILABLE")
		return
	}

	var final_quests []*game.Quest
	for _, quest := range client.quests {
		if quest.Id == args[0] {
			client.Send("OK abandon_quest=%s", quest.Id)
			quest.Status = "available"
		} else {
			final_quests = append(final_quests, quest)
		}
	}
}
