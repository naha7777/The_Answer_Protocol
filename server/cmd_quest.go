package server

import (
	"slices"
	"strconv"
	"strings"
	"the_answer_protocol/game"
)

func (d *Dispatcher) handleQuest(client *Client, args []string) {
	str := strings.ToUpper(strings.Join(args, " "))
	d.mu.Lock()
	defer d.mu.Unlock()
	find := false
	// check if the NPC is in the actual room
	for _, npcs := range client.room.NPCs {
		npc := FindNPC(d, client, npcs)
		if strings.ToUpper(npcs) == str || str == strings.ToUpper(npc.Name) {
			for _, q := range d.world.Quests {
				for _, quest := range npc.Quest {
					if q.Id == quest {
						// check if the NPC has an available quest
						if q.Status == "available" && q.Description != "" {
							// give the quest to the player
							client.quests = append(client.quests, q)
							client.Send(`OK quest={"quest_npc": "%s", "quest_id": "%s", "description": "%s", "reward": "%s", "status": "%s"}`, q.NPC, q.Id, q.Description, q.Reward, q.Status)
							SendAll(d, client, `EVT QUEST took=%s`, q.Id)
							client.log().Info("quest taken", "quest", q.Id)
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
		client.Send("ERR 15 NO_QUEST_AVAILABLE")
		return
	}
	client.Send("ERR 404 NPC_NOT_FOUND")
}

func (d *Dispatcher) handleQuests(client *Client, args []string) {
	// Display player's quests
	res := "OK quests=["
	d.mu.Lock()
	defer d.mu.Unlock()
	for i, v := range client.quests {
		res += "{"
		res += `"quest_npc": "`
		res += v.NPC
		res += `", "quest_id": "`
		res += v.Id
		res += `", "status": "`
		res += v.Status
		if v.Status != "completed" {
			_, nb := CheckObjective(v, client)
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
	// check if the quest is not an empty string
	if len(args) == 0 {
		client.Send("ERR 15 NO_QUEST_AVAILABLE")
		return
	}
	for _, quest := range client.quests {
		// check the player has this quest
		if quest.Id == args[0] {
			find = true
			// check if the quest's NPC is in the actual room
			if slices.Contains(client.room.NPCs, quest.NPC) {
				// check if the quest is still available
				if quest.Status != "active" {
					client.Send("ERR 15 NO_QUEST_AVAILABLE")
					return
				}
				// check if the quest's objectives are done
				completed, _ = CheckObjective(quest, client)
				if completed {
					// Complete the quest
					success := CompleteObjective(quest, client, d)
					if success {
						SendAll(d, client, `EVT QUEST completed=%s`, quest.Id)
						client.log().Info("quest completed", "quest", quest.Id)
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
		client.Send("ERR 15 NO_QUEST_AVAILABLE")
		return
	}
}

func (d *Dispatcher) handleAbandonQuest(client *Client, args []string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	// check if the quest is not an empty string
	if len(args) == 0 {
		client.Send("ERR 15 NO_QUEST_AVAILABLE")
		return
	}

	var final_quests []*game.Quest
	// check if the quest is a player's quest, if it is, abandon it and make it
	// available
	for _, quest := range client.quests {
		if quest.Id == args[0] {
			client.Send("OK abandon_quest=%s", quest.Id)
			SendAll(d, client, `EVT QUEST abandoned=%s`, quest.Id)
			client.log().Info("quest abandoned", "quest", quest.Id)
			quest.Status = "available"
		} else {
			final_quests = append(final_quests, quest)
			client.Send("ERR 15 NO_QUEST_AVAILABLE")
			return
		}
	}
}
