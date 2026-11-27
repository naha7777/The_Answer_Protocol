package server

import (
	"encoding/json"
	"strings"
)

type RoomQuest struct {
	AllQuests []*NPCQuest `json:"quests"`
}

type NPCQuest struct {
	Npc    string   `json:"npc"`
	Quests []string `json:"quests"`
}

func (d *Dispatcher) handleWho(client *Client, args []string) {
	// Display the number of connected players
	d.mu.Lock()
	nb := len(d.clients)
	d.mu.Unlock()
	client.Send("OK players=%d", nb)
}

func (d *Dispatcher) handleLook(client *Client, args []string) {
	// Display all informations about the actual room
	d.mu.Lock()
	defer d.mu.Unlock()
	jsonData, err := json.Marshal(client.room)
	if err != nil {
		Srv().Error("json marshal failed", "player", client.name, "cmd",
			"LOOK", "err", err)
		client.Send("ERR 500 INTERNAL")
	}
	client.Send(`OK look={"room": %s}`, jsonData)
}

func (d *Dispatcher) handleInventory(client *Client, args []string) {
	// Display the player's inventory
	d.mu.Lock()
	defer d.mu.Unlock()
	var item_lst []string
	for _, v := range client.inventory {
		item_lst = append(item_lst, v.Id)
	}
	lst := strings.Join(item_lst, " ")
	client.Send("OK inventory=[%s]", lst)
}

func (d *Dispatcher) handleFindQuests(client *Client, args []string) {
	// Display available quests in the actual room
	d.mu.Lock()
	defer d.mu.Unlock()
	lst := []*NPCQuest{}
	for _, npcID := range client.room.NPCs {
		npc, ok := d.world.NPCs[npcID]
		if !ok {
			continue
		}

		available := []string{}
		for _, questID := range npc.Quest {
			q := d.world.Quests[questID]
			if q != nil && q.Status == "available" {
				available = append(available, questID)
			}
		}

		if len(available) > 0 {
			lst = append(lst, &NPCQuest{Npc: npc.Name, Quests: available})
		}
	}
	if len(lst) == 0 {
		client.Send("ERR 015 NO_QUEST_AVAILABLE")
	}
	jsonData, err := json.Marshal(RoomQuest{AllQuests: lst})
	if err != nil {
		client.log().Error("json marshal failed", "cmd", "LOOK", "err", err)
		client.Send("ERR 500 INTERNAL")
		return
	}
	client.Send(`OK find_quests=%s`, jsonData)
}
