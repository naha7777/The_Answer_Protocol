package main

// /===========================================\
// |				  Import				   |
// \===========================================/

import (
	"fmt"
	"strconv"
	"strings"
	"the_answer_protocol/client"
	"time"

	"fyne.io/fyne/v2"
)

func ServerResponseEVT(msg client.Message, player *Player, data *Data) {
	var err error
	m := strings.Split(msg.Raw, " ")
	if len(m) == 3 && msg.Raw == "EVT PLAYER DEATH" {
		fmt.Println(msg.Rest)
		Popup(data.w, "You are DEAD !!!")
		data.battleUI.EndFight()
		data.battleUI.SetButtonsEnabled(false)
		time.AfterFunc(3*time.Second, func() {
			fyne.Do(func() {
				player.currentHP = player.maxHP
				Popup(data.w, "You had respawned")
				data.w.SetContent(data.gameContent)
				data.conn.Send("LOOK")
				data.conn.Send("INVENTORY")
				data.conn.Send("STATUS")
			})
		})
	}
	if len(m) == 2 && strings.HasPrefix(m[1], "dead=") {
		deathName := strings.TrimPrefix(m[1], "dead=")
		data.chat.Receive("Global", "Game", deathName+" is dead !")
	}
	if len(m) >= 5 && m[2] == "CHAT" {
		parts := strings.SplitN(msg.Rest, " ", 4)
		if len(parts) == 4 {
			scope, username, message := parts[0], parts[2], parts[3]
			data.chat.Receive(scope, username, message)
		}
	}
	if len(m) >= 3 {
		switch m[1] {
		case "GROUP":
			switch m[2] {
			case "INVITE":
				if len(m) >= 4 {
					confirmGroup(data.w, []string{m[3], m[2]}, func() {
						data.conn.Send("GROUP JOIN " + m[3])
					})
				}
			case "JOIN", "LEAVE":
				data.conn.Send("INFO_GROUP")
			}
		case "ROOM":
			switch m[2] {
			case "ENTER":
				player.room.nb_players += 1
			case "LEAVE":
				player.room.nb_players -= 1
			}
		case "NPC":
			switch m[2] {
			case "DEAD":
				data.conn.Send("LOOK")
			case "FLEE":
				npcID := strings.Split(msg.Rest, "=")
				npcMsg := ""
				if len(npcID) >= 2 {
					npcMsg = npcID[1] + " fled"
					Popup(data.w, npcMsg)
				}
				data.battleUI.EndFight()
				data.battleUI.SetButtonsEnabled(false)
				data.g.inventory.Clear()
				time.AfterFunc(3*time.Second, func() {
					fyne.Do(func() {
						data.w.SetContent(data.gameContent)
						data.conn.Send("LOOK")
					})
				})
			case "ATTACK":
				data.battleUI.npcAttack = true
			case "REGENERATION":
				data.battleUI.npcLog.SetText(fmt.Sprintf("%s has regenerated", data.battleUI.targetNPC))
			}
		case "BEER":
			// -5hp pour le player juste pendant le combat ?
		case "SWORD":
			// +5 de point de force ou attack pour le joueur
		case "POISON":
			if m[3] == "POTION" {
				// -5 npc hp
			} else {
				// -5 player hp
			}
		case "STRENGTH":
			if m[3] == "POTION" {
				//  +5 de point de force ou attack pour le joueur
			} else {
				// +5 de point de force ou attack pour le npc
			}
		case "HEALTH":
			if m[3] == "POTION" {
				// +5 hp pour le player
			} else {
				// +5 hp pour le npc
			}
		case "STATS":
			if strings.HasPrefix(m[2], "players=") {
				data.nb_players, err = strconv.Atoi(strings.TrimPrefix(m[2], "players="))
				if err != nil {
					fmt.Println("Failed to get number of players online")
				}
			} else if strings.HasPrefix(m[2], "players_in_room=") {
				data.g.player.room.nb_players, err = strconv.Atoi(strings.TrimPrefix(m[2], "players_in_room="))
				if err != nil {
					fmt.Println("Failed to get number of players in room")
				}
			}
		case "PLAYER":
			data.conn.Send("LOOK")
		}
	}
}
