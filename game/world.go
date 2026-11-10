package game

import (
	"fmt"
	"os"
)

type World struct {
	Start	string
	Rooms	map[string]*Room
}

type Room struct {
	Id		string
	Name	string
	Exits	map[string]string
	Items	[]string
	NPCs	[]string
	Players	[]string
}

type Item struct {
	name		string
	description	string
	obtainable	bool
}

type NPC struct {
	name		string
	hp			int
	hostile		bool
	dialogue	string
}

func LoadWorld(path string) (*World, error) {

	w, err := extractor(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, nil
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
				   Items: room.Items, NPCs: room.NPCs,
				   Players: []string{}}
		room_lst[str] = r
	}

	real_w := World{Start: w.Start, Rooms: room_lst}
	return &real_w, nil
}
