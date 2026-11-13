package game

import (
	"fmt"
	"encoding/json"
	"os"
)

type WorldData struct {
	Start	string					`json:"start"`
	Respawn	string					`json:"respawn"`
	Rooms	map[string]RoomData		`json:"rooms"`
	Items	map[string]ItemData		`json:"items"`
	NPCs	map[string]NPCData		`json:"npcs"`
	Quests	map[string]QuestData	`json:"quests"`
}

type RoomData struct {
	Id			string				`json:"id"`
	Name		string				`json:"name"`
	Description	string				`json:"description"`
	Exits		map[string]string	`json:"exits"`
	Items		[]string			`json:"items"`
	NPCs		[]string			`json:"npcs"`
}

type ItemData struct {
	Id			string				`json:"id"`
	Name		string				`json:"name"`
	Description	string				`json:"description"`
	Obtainable	bool				`json:"obtainable"`
}

type NPCData struct {
	Name		string				`json:"name"`
	HP			int					`json:"hp"`
	Hostile		bool				`json:"hostile"`
	Dialogue	[]string			`json:"dialogue"`
	Quest		[]string			`json:"quests"`
}

type QuestData struct {
	NPC			string				`json:"quest_npc"`
	Id			string				`json:"quest_id"`
	Description	string				`json:"description"`
	Reward		[]string			`json:"reward"`
	Status		string				`json:"status"`
	Objective	[]string			`json:"objective"`
}

func extractor(path string) (*WorldData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening world file: %w", err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()

	var w WorldData
	if err := dec.Decode(&w); err != nil {
		return nil, fmt.Errorf("decoding world %s: %w", path, err)
	}

	return &w, nil
}
