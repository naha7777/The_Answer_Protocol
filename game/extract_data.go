package game

import (
	"fmt"
	"encoding/json"
	"os"
)

type WorldData struct {
	Start	string				`json:"start"`
	Respawn	string				`json:"respawn"`
	Rooms	map[string]RoomData	`json:"rooms"`
	Items	map[string]ItemData	`json:"items"`
	NPCs	map[string]NPCData	`json:"npcs"`
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
	Name		string	`json:"name"`
	Description	string	`json:"description"`
	Obtainable	bool	`json:"obtainable"`
}

type NPCData struct {
	Name		string		`json:"name"`
	HP			int			`json:"hp"`
	Hostile		bool		`json:"hostile"`
	Dialogue	[]string	`json:"dialogue"`
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
