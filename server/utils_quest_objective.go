package server

import (
	"strings"
	"slices"
	"the_answer_protocol/game"
)

func CheckObjective(quest *game.Quest, client *Client) (bool, int) {
	nb := 0
	for {
		clone_obj := slices.Clone(quest.Objective)
		clone_it := slices.Clone(client.inventory)
		var new_obj []string
		var new_it []*game.Item

		for _, obj := range clone_obj {
			find := false
			for _, it := range clone_it {
				if it.Id != obj {
					new_it = append(new_it, it)
				} else {
					find = true
				}
			}
			if !find {
				new_obj = append(new_obj, obj)
			} else {
				nb += 1
			}
		}

		clone_obj = new_obj
		clone_it = new_it

		if len(clone_obj) <= 0 {
			return true, nb
		} else {
			return false, nb
		}
	}
}

func CompleteObjective(quest *game.Quest, client *Client, d *Dispatcher) {
	quest.Status = "completed"
	var final_inventory []*game.Item
	// remove from the player's inventory the objective's items
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
	// Add the rewards on the player's inventory
	for _, item_id := range quest.Reward {
		for _, real_item := range d.world.Items {
			if item_id == real_item.Id {
				client.inventory = append(client.inventory, real_item)
				client.Send("OK taken=%s", real_item.Id)
			}
		}
	}
}
