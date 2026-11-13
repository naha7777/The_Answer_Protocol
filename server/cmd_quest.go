package server

import (
	"strings"
	"slices"
	"strconv"
	"the_answer_protocol/game"
)

func (d *Dispatcher) handleQuest(client *Client, args []string) {
	npc := strings.ToUpper(strings.Join(args, " "))
	d.mu.Lock()
	defer d.mu.Unlock()
	find := false
	for _, npcs := range client.room.NPCs {
		if strings.ToUpper(npcs) == npc || npc == strings.ToUpper(d.world.NPCs[npcs].Name){
			for _, q := range(d.world.Quests) {
				for _, quest := range d.world.NPCs[npcs].Quest {
					if q.Id == quest {
						if q.Status == "available" && q.Description != "" {
							client.quests = append(client.quests, q)
							client.Send(`OK {"quest_id": "%s", "description": "%s", "reward": "%s", "status": "%s"}`, q.Id, q.Description, q.Reward, q.Status)
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
	res := "OK ["
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
			}
			res += `", "progress": "` + strconv.Itoa(nb) + "/"
			res += strconv.Itoa(len(v.Objective))
		}
		res += `"`
		if i != len(client.quests) - 1 {
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
			find =  true
			if slices.Contains(client.room.NPCs, quest.NPC) {
				if quest.Status != "active" {
					client.Send("ERR 406 NO_QUEST_AVAILABLE")
					return
				}
				nb := 0
				for _, obj := range quest.Objective {
					for _, it := range client.inventory {
						if it.Id == obj {
							nb += 1
						}
					}
				}
				if nb == len(quest.Objective) {
					completed = true
				}
				if completed {
					quest.Status = "completed"
					var final_inventory []*game.Item
					for _, item := range client.inventory {
						isObjective := false
						for _, obj := range quest.Objective {
							if strings.EqualFold(item.Id, obj) || strings.EqualFold(item.Name, obj) {
								isObjective = true
								break
							}
						}
						if isObjective {
							continue
						} else {
							final_inventory = append(final_inventory, item)
						}
					}
					client.inventory = final_inventory
					client.Send("OK %s: %s", quest.Id, quest.Status)
					for _, item_id := range quest.Reward {
						for _, real_item := range d.world.Items {
							if item_id == real_item.Id {
								client.inventory = append(client.inventory, real_item)
								client.Send("OK taken=%s", real_item.Id)
							}
						}
					}
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
			quest.Status = "available"
		} else {
			final_quests = append(final_quests, quest)
		}
	}
}
