package main

// /===========================================\
// |				  Import				   |
// \===========================================/

import (
	"encoding/json"
	"fmt"
	"strings"
	"the_answer_protocol/client"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// /===========================================\
// |				   Struct				   |
// \===========================================/

type Room struct {
	nb_players	int
	id			string
	name		string
	description	string
	npc			[]string
	players		[]string
	items		[]string
	exits		map[string]bool
}

type Data struct {
	nb_players	int
	chat		*Chat
	g			*GameState
	w			fyne.Window
	conn		*client.Conn
	cb			*CardinalButtons
}

type Player struct {
	maxHP		int
	currentHP	int
	room		Room
	move		bool
	avatar		string
	gender		string
	name		string
	group		GroupInfo
}

type GameState struct {
	player		*Player
	npc			[]string
	inventory	*Inventory
	popup		*widget.PopUp
}

// /===========================================\
// |				   Method				   |
// \===========================================/

func (g *GameState) Update(str string) {
	if strings.HasPrefix(str, "inventory=") {
		g.inventory.Update(str)
	} else if strings.HasPrefix(str, "look=") {
		g.inventory.Update(str)
		str = strings.TrimPrefix(str, "look=")
		var data map[string]any
		if err := json.Unmarshal([]byte(str), &data); err != nil {
			fmt.Println("erreur de parsing:", err)
			return
		}
		var NpcsID []string
		if room, ok := data["room"].(map[string]any); ok {
			if rawNpcs, ok := room["NPCs"].([]any); ok {
				for _, npc := range rawNpcs {
					if str, ok := npc.(string); ok {
						NpcsID = append(NpcsID, str)
					}
				}
			}
		}
		g.npc = NpcsID
	}
}
