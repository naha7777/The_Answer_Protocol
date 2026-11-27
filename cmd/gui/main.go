package main

// /===========================================\
// |				  Import				   |
// \===========================================/

import (
	"fmt"
	"image/color"
	"os"
	// "strconv"
	"strings"
	"the_answer_protocol/client"
	// "time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// /===========================================\
// |				 Var/Const				   |
// \===========================================/

// Colors
var (
	lime   = color.RGBA{R: 33, G: 222, B: 20, A: 255}
	cyan   = color.RGBA{R: 17, G: 233, B: 237, A: 255}
	orange = color.RGBA{R: 237, G: 123, B: 17, A: 255}
	yellow = color.RGBA{R: 255, G: 213, B: 28, A: 255}
)

// /===========================================\
// |				  Display				   |
// \===========================================/

func main() {

	// CLient connection
	conn, err := client.Dial()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERR 900 CONNECTION_FAILED")
		os.Exit(1)
	}
	defer conn.Close()

	// Application and Windows init
	a := app.New()
	w := a.NewWindow("TAP")
	player := Player{maxHP: 100, currentHP: 100, move: false}
	w.Resize(fyne.NewSize(1920, 1080))

	// Key hook for esc to close the app
	w.Canvas().SetOnTypedKey(func(keyEvent *fyne.KeyEvent) {
		if keyEvent.Name == fyne.KeyEscape {
			a.Quit()
		}
	})

	overlay := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 90})

	// Struct info
	// chat := NewChat(w, &player, conn)
	inv := NewInv(w, &player, conn)
	// game := &GameState{inventory: inv, player: &player}
	background := canvas.NewImageFromFile("./data/tavern.jpeg")
	background.FillMode = canvas.ImageFillCover
	data := Data{
		a: a,
		w: w,
		conn: conn,
		cb: &CardinalButtons{},
		g: &GameState{inventory: NewInv(w, &player, conn), player: &player},
		chat: NewChat(w, &player, conn),
		bg: background,
	}

	logLabel := widget.NewLabel("")
	logLabel.Wrapping = fyne.TextWrapWord
	logView := container.NewVScroll(logLabel)

	// Widgets init for base layout
	btnContainer := newButtonContainer(background, &player, &data)
	data.playerInfo = showPlayerInfo(&player, data)
	hudContent := container.New(&Hud{}, btnContainer, data.chat.Content(), inv.Content())
	gameContent := container.NewStack(
		background,
		overlay,
		hudContent,
		container.NewVBox(container.NewHBox(data.playerInfo.Content())),
	)

	// Widgets init for combat layout
	arenaBg := canvas.NewImageFromFile("./data/arena.jpeg")
	arenaBg.FillMode = canvas.ImageFillCover

	// battleUI := NewBattleUI(&player, &data)
	data.battleUI = NewBattleUI(&player, &data)
	battleContent := container.NewStack(
		arenaBg,
		overlay,
		data.battleUI.Content(),
	)

	loggedIn := false
	loginView, onLoginMsg := LoginScreen(conn, &player, func(name string) {
		loggedIn = true
		data.playerInfo.Update(data)
		data.cb.Update(player.room.exits)
		w.SetContent(gameContent)
		conn.Send("STATUS")
		conn.Send("INVENTORY")
		conn.Send("LOOK")
	})

	data.gameContent = gameContent
	data.battleContent = battleContent
	w.SetContent(loginView)
	// npc := false
	// var potion *Potion

	go func() {
		for msg := range conn.Messages() {
			fyne.Do(func() {
				if !loggedIn {
					onLoginMsg(msg)
					return
				}
				switch msg.Kind {
				case "OK":
					ServerResponseOK(msg, &player, &data)
					// switch {

					// // MOVE
					// case strings.HasPrefix(msg.Rest, "room=") && player.move:
					// 	roomID := strings.TrimPrefix(msg.Rest, "room=loc.")
					// 	player.room.name = strings.ToUpper(roomID)
					// 	background.File = "./data/" + roomID + ".jpeg"
					// 	background.Refresh()
					// 	player.move = false
					// 	playerinfo.Update(data)
					// 	data.cb.Update(player.room.exits)
					// 	data.conn.Send("LOOK")

					// case msg.Rest == "bye":
					// 	a.Quit()

					// // GROUP
					// case strings.HasPrefix(msg.Rest, "group_created="):
					// 	conn.Send("INFO_GROUP")
					// 	Popup(w, strings.TrimPrefix(msg.Rest, "group_created=")+" group was created")
					// case strings.HasPrefix(msg.Rest, "group_leaved="):
					// 	Popup(w, "You leave "+strings.ToUpper(player.group.Id)+" group")
					// 	player.group = GroupInfo{}
					// 	playerinfo.Update(data)
					// case strings.HasPrefix(msg.Rest, "group_joined="):
					// 	conn.Send("INFO_GROUP")
					// 	Popup(w, "You join the "+strings.ToUpper(strings.TrimPrefix(msg.Rest, "group_joined="))+"group")
					// case strings.HasPrefix(msg.Rest, "group_infos="):
					// 	raw := strings.TrimPrefix(msg.Rest, "group_infos=")
					// 	var infos groupInfosMsg
					// 	if err := json.Unmarshal([]byte(raw), &infos); err != nil {
					// 		Popup(w, "invalid group infos")
					// 		break
					// 	}
					// 	player.group = infos.Group
					// 	playerinfo.Update(data)

					// // TAKE
					// case strings.HasPrefix(msg.Rest, "taken="):
					// 	item := strings.TrimPrefix(msg.Rest, "taken=")
					// 	game.inventory.take(item)
					// 	conn.Send("LOOK")

					// // DROP
					// case strings.HasPrefix(msg.Rest, "dropped="):
					// 	item := strings.TrimPrefix(msg.Rest, "dropped=")
					// 	game.inventory.drop(item)
					// 	conn.Send("LOOK")

					// // INVENTORY
					// case strings.HasPrefix(msg.Rest, "inventory="):
					// 	game.inventory.BuildInv(msg.Rest)
					// 	game.Update(msg.Rest)

					// // LOOK
					// case strings.HasPrefix(msg.Rest, "look="):
					// 	game.Update(msg.Rest)
					// 	player.PlayerRoomUpdate(msg.Rest)
					// 	playerinfo.Update(data)
					// 	data.cb.Update(player.room.exits)

					// // QUESTS
					// case strings.HasPrefix(msg.Rest, "quests="):
					// 	msg.Rest = strings.TrimPrefix(msg.Rest, "quests=")
					// 	Str_to_quests(msg.Rest, data)

					// // TALK
					// case strings.HasPrefix(msg.Rest, "talk="):
					// 	msg.Rest = strings.TrimPrefix(msg.Rest, "talk=")
					// 	author, message, _ := strings.Cut(msg.Rest, ": ")
					// 	author = IDToName(author)
					// 	chat.Receive("NPC", author, message)

					// // ATTACK
					// case strings.HasPrefix(msg.Rest, "flee="):
					// 	if strings.Contains(msg.Rest, "=succed") {
					// 		battleUI.EndFight()
					// 		playerinfo.Update(data)
					// 		w.SetContent(gameContent)
					// 		conn.Send("LOOK")
					// 	}
					// case strings.HasPrefix(msg.Rest, "attack="):
					// 	raw := strings.TrimPrefix(msg.Rest, "attack=")
					// 	if npcID, _, ok := strings.Cut(raw, "=dead"); ok {
					// 		npcID = strings.Trim(npcID, `"`)
					// 		Popup(w, IDToName(npcID) + " is dead!")
					// 		battleUI.EndFight()
					// 		battleUI.SetButtonsEnabled(false)
					// 		playerinfo.Update(data)
					// 		time.AfterFunc(3*time.Second, func() {
					// 			fyne.Do(func() {
					// 				player.currentHP = player.maxHP
					// 				w.SetContent(gameContent)
					// 				conn.Send("LOOK")
					// 				conn.Send("INVENTORY")
					// 				conn.Send("STATUS")
					// 			})
					// 		})
					// 		break
					// 	}
					// 	var atk AttackData
					// 	if err := atk.parse(strings.TrimPrefix(msg.Rest, "attack=")); err != nil {
					// 		break
					// 	}
					// 	if !battleUI.inFight {
					// 		battleUI.InitFight(atk)
					// 		w.SetContent(battleContent)
					// 	} else if !npc {
					// 		battleUI.npcAttack = false
					// 		battleUI.Update(atk)
					// 	} else {
					// 		battleUI.Update(atk)
					// 		battleUI.npcAttack = false
					// 	}
					// 	case strings.HasPrefix(msg.Rest, `"status"=`), strings.HasPrefix(msg.Rest, "status="):
					// 	status := strings.TrimPrefix(msg.Rest, `"status"=`)
					// 	status = strings.TrimPrefix(status, "status=")
					// 	// status = strings.Trim(status, `"`)

					// 	if status != "combat" && battleUI.inFight {
					// 		battleUI.EndFight()
					// 		playerinfo.Update(data)
					// 		w.SetContent(gameContent)
					// 		conn.Send("LOOK")
					// 	}
					// 	HandleStatus(status, &player)
					// 	playerinfo.Update(data)

					// // QUEST
					// case strings.HasPrefix(msg.Rest, "quest="):
					// 	Popup(w, "Quest accepted!")
					// 	HandleQuest(msg.Rest, data)

					// // ITEM REMOVAL
					// case strings.HasPrefix(msg.Rest, "gave="):
					// 	item := strings.TrimPrefix(msg.Rest, "gave=")
					// 	game.inventory.remove(item)
					// 	conn.Send("LOOK")
					// }

				case "EVT":
					ServerResponseEVT(msg, &player, &data)
					// m := strings.Split(msg.Raw, " ")
					// if len(m) == 3 && msg.Raw == "EVT PLAYER DEATH" {
					// 	fmt.Println(msg.Rest)
					// 	Popup(data.w, "You are DEAD !!!")
					// 	battleUI.EndFight()
					// 	battleUI.SetButtonsEnabled(false)
					// 	time.AfterFunc(3*time.Second, func() {
					// 		fyne.Do(func() {
					// 			player.currentHP = player.maxHP
					// 			Popup(data.w, "You had respawned")
					// 			w.SetContent(gameContent)
					// 			conn.Send("LOOK")
					// 			conn.Send("INVENTORY")
					// 			conn.Send("STATUS")
					// 		})
					// 	})
					// }
					// if len(m) == 2 && strings.HasPrefix(m[1], "dead=") {
					// 	deathName := strings.TrimPrefix(m[1], "dead=")
					// 	data.chat.Receive("Global", "Game", deathName+" is dead !")
					// }
					// if len(m) >= 5 && m[2] == "CHAT" {
					// 	parts := strings.SplitN(msg.Rest, " ", 4)
					// 	if len(parts) == 4 {
					// 		scope, username, message := parts[0], parts[2], parts[3]
					// 		data.chat.Receive(scope, username, message)
					// 	}
					// }
					// if len(m) >= 3 {
					// 	switch m[1] {
					// 	case "GROUP":
					// 		switch m[2] {
					// 		case "INVITE":
					// 			if len(m) >= 4 {
					// 				confirmGroup(w, []string{m[3], m[2]}, func() {
					// 					conn.Send("GROUP JOIN " + m[3])
					// 				})
					// 			}
					// 		case "JOIN", "LEAVE":
					// 			conn.Send("INFO_GROUP")
					// 		}
					// 	case "ROOM":
					// 		switch m[2] {
					// 		case "ENTER":
					// 			player.room.nb_players += 1
					// 		case "LEAVE":
					// 			player.room.nb_players -= 1
					// 		}
					// 	case "NPC":
					// 		switch m[2] {
					// 		case "DEAD":
					// 			conn.Send("LOOK")
					// 		case "FLEE":
					// 			npcID := strings.Split(msg.Rest, "=")
					// 			npcMsg := ""
					// 			if len(npcID) >= 2 {
					// 				npcMsg = npcID[1] + " fled"
					// 				Popup(data.w, npcMsg)
					// 			}
					// 			battleUI.EndFight()
					// 			battleUI.SetButtonsEnabled(false)
					// 			data.g.inventory.Clear()
					// 			time.AfterFunc(3*time.Second, func() {
					// 				fyne.Do(func() {
					// 					w.SetContent(gameContent)
					// 					conn.Send("LOOK")
					// 				})
					// 			})
					// 		case "ATTACK":
					// 			battleUI.npcAttack = true
					// 		case "REGENERATION":
					// 			battleUI.npcLog.SetText(fmt.Sprintf("%s has regenerated", battleUI.targetNPC))
					// 		}
					// 	case "BEER":
					// 		// -5hp pour le player juste pendant le combat ?
					// 	case "SWORD":
					// 		// +5 de point de force ou attack pour le joueur
					// 	case "POISON":
					// 		if m[3] == "POTION" {
					// 			// -5 npc hp
					// 		} else {
					// 			// -5 player hp
					// 		}
					// 	case "STRENGTH":
					// 		if m[3] == "POTION" {
					// 			//  +5 de point de force ou attack pour le joueur
					// 		} else {
					// 			// +5 de point de force ou attack pour le npc
					// 		}
					// 	case "HEALTH":
					// 		if m[3] == "POTION" {
					// 			// +5 hp pour le player
					// 		} else {
					// 			// +5 hp pour le npc
					// 		}
					// 	case "STATS":
					// 		if strings.HasPrefix(m[2], "players=") {
					// 			data.nb_players, err = strconv.Atoi(strings.TrimPrefix(m[2], "players="))
					// 			if err != nil {
					// 				fmt.Println("Failed to get number of players online")
					// 			}
					// 		} else if strings.HasPrefix(m[2], "players_in_room=") {
					// 			data.g.player.room.nb_players, err = strconv.Atoi(strings.TrimPrefix(m[2], "players_in_room="))
					// 			if err != nil {
					// 				fmt.Println("Failed to get number of players in room")
					// 			}
					// 		}
					// 	case "PLAYER":
					// 		conn.Send("LOOK")
					// 	}
					// }
				case "ERR":
					e := strings.Split(msg.Raw, " ")
					switch e[1] {
					case "201":
						Popup(w, "This nickname is already taken")
					case "301":
						Popup(w, "You can't go that way")
					case "401":
						Popup(w, "You are not in a group")
					case "402":
						Popup(w, "The player is already in group")
					case "404":
						Popup(w, "You can't do that")
					case "405":
						Popup(w, "Let's be friendly")
					case "900":
						Popup(w, "Failed to connect to the server")
					case "901":
						Popup(w, "Failed to reach the server")
					case "001":
						Popup(w, "Command didn't exist")
					case "002":
						Popup(w, "You must be connected")
					case "003":
						Popup(w, "Player not found")
					case "004":
						Popup(w, "You are not the leader")
					case "005":
						Popup(w, "You need an invitation")
					case "006":
						Popup(w, "To who are you talking?")
					case "007":
						Popup(w, "You can't get that")
					case "008":
						Popup(w, "Try typing something")
					case "009":
						Popup(w, "Name is invalid")
					case "010":
						Popup(w, "You are already connected")
					case "011":
						Popup(w, "You need to be near the npc with all the requirement")
					case "012":
						Popup(w, "Let's focus on the fight")
					case "013":
						Popup(w, "Why? It's chill right here")
					case "014":
						Popup(w, "I respect the 1v1")
					case "015":
						Popup(w, "Don't have no quest for you")
					case "016":
						Popup(w, "Line too long")
					case "017":
						Popup(w, "Inventory is full")
					case "018":
						Popup(w, "Stop the spam")
					case "019":
						Popup(w, "Server is full try again later")
					case "020":
						Popup(w, "Invalid input")
					case "021":
						Popup(w, "Name is too long")
					case "500":
						Popup(w, "Internal error")
					default:
						Popup(w, msg.Rest)
					}
				}
			},
			)
		}

		fyne.Do(func() {
			logLabel.SetText(logLabel.Text + "-- lost connection --\n")
			logView.ScrollToBottom()
			a.Quit()
		})
		fmt.Fprintln(os.Stderr, "-- lost connection --")
		fyne.Do(func() {
			a.Quit()
		})
	}()

	w.ShowAndRun()
	fmt.Println("Exited, end of TAP")
}
