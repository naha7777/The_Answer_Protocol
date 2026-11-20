package server

import (
	"strings"
	"the_answer_protocol/game"
)

func (d *Dispatcher) removeClient(client *Client) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if client.name == "" {
		return
	}
	if client.group != nil {
		delete(client.group.members, client.name)
		client.group = nil
	}
	if client.inventory != nil {
		RemoveItemInventory(client, "false_item", d, false)
	}
	if client.room != nil {
		client.room.Players = Remove(client.name, client.room.Players)
		client.room = nil
	}
	delete(d.clients, client.name)
	client.name = ""
	client.conn.Close()
}

func RemoveItemInventory(client *Client, item_id string, d *Dispatcher, verbose bool) {
	final_inventory := []*game.Item{}
	found := false
	count := 0
	if item_id == "false_item" {
		for _, i := range client.inventory {
			client.room.Items = append(client.room.Items, i.Id)
		}
		client.inventory = final_inventory
		return
	}
	for _, i := range client.inventory {
		if strings.EqualFold(i.Id, item_id) && count == 0 || strings.EqualFold(i.Name, item_id) && count == 0{
			found = true
			client.room.Items = append(client.room.Items, i.Id)
			count += 1
		} else {
			final_inventory = append(final_inventory, i)
		}
	}
	client.inventory = final_inventory
	if verbose {
		client.Send("OK dropped=%s", item_id)
		SendRoom(d, client, "EVT PLAYER DROPPED %s", item_id)
	}
	if !found {
		client.Send("ERR 404 ITEM_NOT_IN_INVENTORY")
	}
}

func RemoveQuest(client *Client) {
	var final_quests []*game.Quest
	for _, quest := range client.quests {
		if quest.Status == "completed" {
			final_quests = append(final_quests, quest)
		} else {
			quest.Status = "available"
		}
	}
	client.quests = final_quests
}

func Remove(name string, place []string) []string {
	var new []string
	for _, item := range place {
		if item != name {
			new = append(new, item)
		}
	}
	return new
}
