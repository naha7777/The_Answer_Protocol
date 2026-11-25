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
	id          string
	name        string
	description string
	exits       map[string]bool
	npc         []string
	players     []string
	nb_players  int
	items       []string
}

type Data struct {
	w          fyne.Window
	conn       *client.Conn
	cb         *CardinalButtons
	g          *GameState
	chat       *Chat
	nb_players int
}

type Player struct {
	name   string
	health int
	group  GroupInfo
	room   Room
	move   bool
}

type GameState struct {
	inventory *Inventory
	player    *Player
	npc       []string
	popup     *widget.PopUp
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
