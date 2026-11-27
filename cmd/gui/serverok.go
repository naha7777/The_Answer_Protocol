package main

// /===========================================\
// |				  Import				   |
// \===========================================/

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"the_answer_protocol/client"
	"time"

	"fyne.io/fyne/v2"
)

func ServerResponseOK(msg client.Message, player *Player, data *Data) {
	switch {
	// MOVE
	case strings.HasPrefix(msg.Rest, "room=") && player.move:
		roomID := strings.TrimPrefix(msg.Rest, "room=loc.")
		player.room.name = strings.ToUpper(roomID)
		data.bg.File = "./data/" + roomID + ".jpeg"
		data.bg.Refresh()
		player.move = false
		data.playerInfo.Update(*data)
		data.cb.Update(player.room.exits)
		data.conn.Send("LOOK")

	case msg.Rest == "bye":
		data.a.Quit()

	// GROUP
	case strings.HasPrefix(msg.Rest, "group_created="):
		data.conn.Send("INFO_GROUP")
		Popup(data.w, strings.TrimPrefix(msg.Rest, "group_created=")+" group was created")
	case strings.HasPrefix(msg.Rest, "group_leaved="):
		Popup(data.w, "You leave "+strings.ToUpper(player.group.Id)+" group")
		player.group = GroupInfo{}
		data.playerInfo.Update(*data)
	case strings.HasPrefix(msg.Rest, "group_joined="):
		data.conn.Send("INFO_GROUP")
		Popup(data.w, "You join the "+strings.ToUpper(strings.TrimPrefix(msg.Rest, "group_joined="))+"group")
	case strings.HasPrefix(msg.Rest, "group_infos="):
		raw := strings.TrimPrefix(msg.Rest, "group_infos=")
		var infos groupInfosMsg
		if err := json.Unmarshal([]byte(raw), &infos); err != nil {
			Popup(data.w, "invalid group infos")
			break
		}
		player.group = infos.Group
		data.playerInfo.Update(*data)

	// TAKE
	case strings.HasPrefix(msg.Rest, "taken="):
		item := strings.TrimPrefix(msg.Rest, "taken=")
		data.g.inventory.take(item)
		data.conn.Send("LOOK")

	// DROP
	case strings.HasPrefix(msg.Rest, "dropped="):
		item := strings.TrimPrefix(msg.Rest, "dropped=")
		data.g.inventory.drop(item)
		data.conn.Send("LOOK")

	// INVENTORY
	case strings.HasPrefix(msg.Rest, "inventory="):
		data.g.inventory.BuildInv(msg.Rest)
		data.g.Update(msg.Rest)

	// LOOK
	case strings.HasPrefix(msg.Rest, "look="):
		data.g.Update(msg.Rest)
		player.PlayerRoomUpdate(msg.Rest)
		data.playerInfo.Update(*data)
		data.cb.Update(player.room.exits)

	// QUESTS
	case strings.HasPrefix(msg.Rest, "quests="):
		msg.Rest = strings.TrimPrefix(msg.Rest, "quests=")
		Str_to_quests(msg.Rest, *data)

	// TALK
	case strings.HasPrefix(msg.Rest, "talk="):
		msg.Rest = strings.TrimPrefix(msg.Rest, "talk=")
		author, message, _ := strings.Cut(msg.Rest, ": ")
		author = IDToName(author)
		data.chat.Receive("NPC", author, message)

	// ATTACK
	case strings.HasPrefix(msg.Rest, "flee="):
		if strings.Contains(msg.Rest, "=succed") {
			data.battleUI.EndFight()
			data.playerInfo.Update(*data)
			data.w.SetContent(data.gameContent)
			data.conn.Send("LOOK")
		}
	case strings.HasPrefix(msg.Rest, "attack="):
		raw := strings.TrimPrefix(msg.Rest, "attack=")
		if npcID, _, ok := strings.Cut(raw, "=dead"); ok {
			npcID = strings.Trim(npcID, `"`)
			Popup(data.w, IDToName(npcID)+" is dead!")
			data.battleUI.EndFight()
			data.battleUI.SetButtonsEnabled(false)
			data.playerInfo.Update(*data)
			time.AfterFunc(3*time.Second, func() {
				fyne.Do(func() {
					player.currentHP = player.maxHP
					data.w.SetContent(data.gameContent)
					data.conn.Send("LOOK")
					data.conn.Send("INVENTORY")
					data.conn.Send("STATUS")
				})
			})
			break
		}
		var atk AttackData
		if err := atk.parse(strings.TrimPrefix(msg.Rest, "attack=")); err != nil {
			break
		}
		if !data.battleUI.inFight {
			data.battleUI.InitFight(atk)
			data.w.SetContent(data.battleContent)
		} else if !data.battleUI.npcAttack {
			data.battleUI.Update(atk)
		} else {
			data.battleUI.Update(atk)
		}
	case strings.HasPrefix(msg.Rest, "status="):
		status := strings.TrimPrefix(msg.Rest, "status=")
		if status != "combat" && data.battleUI.inFight {
			data.battleUI.EndFight()
			data.playerInfo.Update(*data)
			data.w.SetContent(data.gameContent)
			data.conn.Send("LOOK")
		}
		HandleStatus(status, player)
		data.playerInfo.Update(*data)

	// QUEST
	case strings.HasPrefix(msg.Rest, "quest="):
		Popup(data.w, "Quest accepted!")
		HandleQuest(msg.Rest, *data)

	// ITEM REMOVAL
	case strings.HasPrefix(msg.Rest, "gave="):
		item := strings.TrimPrefix(msg.Rest, "gave=")
		data.g.inventory.remove(item)
		data.conn.Send("LOOK")

	// WHO
	case strings.HasPrefix(msg.Rest, "players="):
		var err error
		data.nb_players, err = strconv.Atoi(strings.TrimPrefix(msg.Rest, "players="))
		if err != nil {
			fmt.Println("Failed to get number of players online")
		}
	}
}
