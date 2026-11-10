package game

import (
	"fmt"
	"encoding/json"
	"os"
)

type Recipient interface {
	Send(format string, a ...any)
	Name()	string
}

type World struct {
    Start string               `json:"start"`
    Rooms map[string]RoomData  `json:"rooms"`
    Items map[string]ItemData  `json:"items"`
    NPCs  map[string]NPCData   `json:"npcs"`
}

type RoomData struct {
    Name        string            `json:"name"`
    Description string            `json:"description"`
    Exits       map[string]string `json:"exits"`
    Items       []string          `json:"items"`
    NPCs        []string          `json:"npcs"`
}

type ItemData struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Obtainable  bool   `json:"obtainable"`
}

type NPCData struct {
    Name     string   `json:"name"`
    HP       int      `json:"hp"`
    Hostile  bool     `json:"hostile"`
    Dialogue []string `json:"dialogue"`
}

func LoadWorld(path string) (*World, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening world file: %w", err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()

	var w World
	if err := dec.Decode(&w); err != nil {
		return nil, fmt.Errorf("decoding world %s: %w", path, err)
	}
	return &w, nil
}

// type Item struct {
// 	name		string
// 	description	string
// 	obtainable	bool
// }

// type NPC struct {
// 	name		string
// 	hp			int
// 	hostile		bool
// 	dialogue	string
// }

// type Room struct {
// 	name		string
// 	description	string
// 	visitors	map[string]Recipient
// 	items		map[string]*Item
// 	exits		map[string]string
// 	npc			map[string]*NPC
// }

