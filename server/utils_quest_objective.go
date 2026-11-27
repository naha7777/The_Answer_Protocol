package server

import (
	"slices"
	"strings"
	"the_answer_protocol/game"
)

func CheckObjective(quest *game.Quest, client *Client) (bool, int) {
	nb := 0
	clone_obj := slices.Clone(quest.Objective)
	clone_it := slices.Clone(client.inventory)
	clone_kill := slices.Clone(client.kills)
	var new_obj []string
	var new_it []*game.Item
	var new_kill []string
	for {
		for _, obj := range clone_obj {
			find := false
			for _, it := range clone_it {
				if it.Id != obj {
					new_it = append(new_it, it)
				} else {
					find = true
				}
			}
			for _, kill := range clone_kill {
				if kill != obj {
					new_kill = append(new_kill, kill)
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

		clone_kill = new_kill
		clone_obj = new_obj
		clone_it = new_it

		if len(clone_obj) <= 0 {
			return true, nb
		} else {
			return false, nb
		}
	}
}

func CompleteObjective(quest *game.Quest, client *Client, d *Dispatcher) bool {
	quest.Status = "completed"

	var obj_it []string
	for _, obj := range quest.Objective {
		if strings.HasPrefix(obj, "item") {
			obj_it = append(obj_it, obj)
		}
	}
	if len(client.inventory)-len(obj_it)+len(quest.Reward) > 100 {
		client.Send("ERR 17 FULL_INVENTORY")
		return false
	}

	for _, obj := range quest.Objective {
		find := false
		for _, it := range client.inventory {
			if strings.EqualFold(it.Id, obj) || strings.EqualFold(it.Name, obj) {
				if it.Id == obj && find == false {
					RemoveItemInventory(client, it.Id, d, false, "quest")
					find = true
				}
			}
		}
	}
	client.Send("OK complete_quest=%s: %s", quest.Id, quest.Status)

	// Add the rewards on the player's inventory

	for _, item_id := range quest.Reward {
		for _, real_item := range d.world.Items {
			if item_id == real_item.Id {
				client.inventory = append(client.inventory, real_item)
				client.Send("OK taken=%s", real_item.Id)
			}
		}
	}
	return true
}
