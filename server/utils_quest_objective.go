package server

import (
	"strings"
	"the_answer_protocol/game"
)

func CheckObjective(completed bool, quest *game.Quest, client *Client) bool {
	nb := 0
	for _, obj := range quest.Objective {
		for _, it := range client.inventory {
			if it.Id == obj {
				nb += 1
			}
		}
		for _, kill := range client.kills {
			if kill == obj {
				nb += 1
			}
		}
	}
	if nb >= len(quest.Objective) {
		completed = true
	}
	return completed
}

func CompleteObjective(quest *game.Quest, client *Client, d *Dispatcher) {
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
	client.Send("OK complete_quest= %s: %s", quest.Id, quest.Status)
	for _, item_id := range quest.Reward {
		for _, real_item := range d.world.Items {
			if item_id == real_item.Id {
				client.inventory = append(client.inventory, real_item)
				client.Send("OK taken=%s", real_item.Id)
			}
		}
	}
}
