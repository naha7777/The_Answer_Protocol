package game

import (
	"fmt"
	"os"
)

type World struct {
	Start	string
	Rooms	map[string]*Room
	Items	map[string]*Item
	NPCs	map[string]*NPC
	Quests	map[string]*Quest
}

type Room struct {
	Id			string
	Name		string
	Description	string
	Exits		map[string]string
	Players		[]string
	Items		[]string
	NPCs		[]string
}

type Item struct {
	Id			string
	Name		string
	Description	string
	Obtainable	bool
}

type NPC struct {
	Name		string
	Hp			int
	Hostile		bool
	Dialogue	[]string
	Quest		string
}

type Quest struct {
	NPC			string
	Id			string
	Description	string
	Reward		string
	Status		string
}

func LoadWorld(path string) (*World, error) {

	w, err := extractor(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	for id, room := range w.Rooms {
		for dir, target := range room.Exits {
			if _, ok := w.Rooms[target]; !ok {
				return nil, fmt.Errorf("room %s: exit %s leads to unknown room %q", id, dir, target)
			}
		}
	}

	room_lst := make(map[string]*Room, len(w.Rooms))
	for str, room := range w.Rooms {
		r := &Room{Id: room.Id, Name: room.Name, Exits: room.Exits,
				   Description: room.Description,
				   Items: room.Items, NPCs: room.NPCs,
				   Players: []string{}}
		room_lst[str] = r
	}

	item_lst := make(map[string]*Item, len(w.Items))
	for str, item := range w.Items {
		i := &Item{Name: item.Name, Description: item.Description,
				   Obtainable: item.Obtainable, Id: item.Id}
		item_lst[str] = i
	}

	npc_lst := make(map[string]*NPC, len(w.NPCs))
	for str, npc := range w.NPCs {
		n := &NPC{Name: npc.Name, Hp: npc.HP, Hostile: npc.Hostile,
				  Dialogue: npc.Dialogue, Quest: npc.Quest}
		npc_lst[str] = n
	}

	quest_lst := make(map[string]*Quest, len(w.Quests))
	for str, quest := range w.Quests {
		q := &Quest{NPC: quest.NPC, Id: quest.Id,
					Description: quest.Description,
					Reward: quest.Reward, Status: quest.Status}
		quest_lst[str] = q
	}

	real_w := World{Start: w.Start, Rooms: room_lst, Items: item_lst,
					NPCs: npc_lst, Quests: quest_lst,}
	return &real_w, nil
}
