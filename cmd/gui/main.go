package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"strconv"
	"strings"
	"the_answer_protocol/client"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

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
		// "NPCs":["npc.village_guard"]
	}
}

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

	conn, err := client.Dial("localhost:8080")
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERR 900 CONNECTION_FAILED")
		os.Exit(1)
	}
	defer conn.Close()

	a := app.New()
	w := a.NewWindow("TAP")
	player := Player{health: 100, move: false}
	w.Resize(fyne.NewSize(1920, 1080))

	// Key hook for esc to close the app
	w.Canvas().SetOnTypedKey(func(keyEvent *fyne.KeyEvent) {
		if keyEvent.Name == fyne.KeyEscape {
			a.Quit()
		}
	})

	overlay := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 90})
	// Struct info
	chat := NewChat(w, &player, conn)
	inv := NewInv(w, &player, conn)
	game := GameState{inventory: inv, player: &player}
	data := Data{w: w, conn: conn, cb: &CardinalButtons{}, g: &game}
	background := canvas.NewImageFromFile("./data/tavern.jpeg")
	background.FillMode = canvas.ImageFillStretch

	logLabel := widget.NewLabel("")
	logLabel.Wrapping = fyne.TextWrapWord
	logView := container.NewVScroll(logLabel)

	// ajouter conn en parametre de la fonction
	btnContainer := newButtonContainer(background, &player, &data)
	playerinfo := showPlayerInfo(&player, data)
	hudContent := container.New(&Hud{}, btnContainer, chat.Content(), inv.Content())
	gameContent := container.NewStack(
		background,
		overlay,
		hudContent,
		container.NewVBox(container.NewHBox(playerinfo.Content())),
	)

	loggedIn := false
	// playerinfo.Update()
	loginView, onLoginMsg := LoginScreen(conn, &player, func(name string) {
		loggedIn = true
		playerinfo.Update(data)
		data.cb.Update(player.room.exits)
		w.SetContent(gameContent)
		conn.Send("INVENTORY")
		conn.Send("LOOK")
	})

	w.SetContent(loginView)

	go func() {
		for msg := range conn.Messages() {
			fyne.Do(func() {
				if !loggedIn {
					onLoginMsg(msg)
					return
				}
				switch msg.Kind {
				case "OK":
					// fmt.Println(msg.Rest)
					switch {

					// MOVE
					case strings.HasPrefix(msg.Rest, "room=") && player.move:
						roomID := strings.TrimPrefix(msg.Rest, "room=loc.")
						player.room.name = strings.ToUpper(roomID)
						background.File = "./data/" + roomID + ".jpeg"
						background.Refresh()
						player.move = false
						playerinfo.Update(data)
						data.cb.Update(player.room.exits)
						data.conn.Send("LOOK")

					case msg.Rest == "bye":
						a.Quit()

					// GROUP
					case strings.HasPrefix(msg.Rest, "group_created="):
						conn.Send("INFO_GROUP")
						popup(w, strings.TrimPrefix(msg.Rest, "group_created=")+" group was created")
					case strings.HasPrefix(msg.Rest, "group_leaved="):
						popup(w, "You leave "+strings.ToUpper(player.group.Id)+" group")
						player.group = GroupInfo{}
						playerinfo.Update(data)
					case strings.HasPrefix(msg.Rest, "group_joined="):
						conn.Send("INFO_GROUP")
						popup(w, "You join the "+strings.ToUpper(strings.TrimPrefix(msg.Rest, "group_joined="))+"group")
					case strings.HasPrefix(msg.Rest, "group_infos="):
						raw := strings.TrimPrefix(msg.Rest, "group_infos=")
						var infos groupInfosMsg
						if err := json.Unmarshal([]byte(raw), &infos); err != nil {
							popup(w, "invalid group infos")
							break
						}
						player.group = infos.Group
						playerinfo.Update(data)

					// TAKE
					case strings.HasPrefix(msg.Rest, "taken="):
						item := strings.TrimPrefix(msg.Rest, "taken=")
						game.inventory.take(item)
						conn.Send("LOOK")

					// DROP
					case strings.HasPrefix(msg.Rest, "dropped="):
						item := strings.TrimPrefix(msg.Rest, "dropped=")
						game.inventory.drop(item)
						conn.Send("LOOK")

					// INVENTORY
					case strings.HasPrefix(msg.Rest, "inventory="):
						game.inventory.BuildInv(msg.Rest)
						game.Update(msg.Rest)

					// LOOK
					case strings.HasPrefix(msg.Rest, "look="):
						game.Update(msg.Rest)
						player.PlayerRoomUpdate(msg.Rest)
						playerinfo.Update(data)
						data.cb.Update(player.room.exits)

					// QUESTS
					case strings.HasPrefix(msg.Rest, "quests="):
						msg.Rest = strings.TrimPrefix(msg.Rest, "quests=")
						Str_to_quests(msg.Rest, data)

					// TALK
					case strings.HasPrefix(msg.Rest, "talk="):
						msg.Rest = strings.TrimPrefix(msg.Rest, "talk=")
						author, message, _ := strings.Cut(msg.Rest, ": ")
						author = IDToName(author)
						chat.Receive("NPC", author, message)

					// ATTACK
					case strings.HasPrefix(msg.Rest, "attack="):
						background.File = "./data/arena.jpeg"
						background.Refresh()
						fmt.Println(background.File)

					case strings.HasPrefix(msg.Rest, "quest="):
						popup(w, "Quest accepted!")
					}

				case "EVT":
					m := strings.Split(msg.Raw, " ")
					if len(m) >= 3 {
						switch m[1] {
						case "GROUP":
							switch m[2] {
							case "INVITE":
								if len(m) >= 5 {
									confirmGroup(w, []string{m[4], m[3]}, func() {
										conn.Send("GROUP JOIN " + m[3])
									})
								}
							case "JOIN", "LEAVE":
								conn.Send("INFO_GROUP")
							}
						case "ROOM":
							switch m[2] {
							case "ENTER":
								player.room.nb_players += 1
							case "LEAVE":
								player.room.nb_players -= 1
							}

						case "CHAT":
							parts := strings.SplitN(msg.Rest, " ", 4)
							scope, username, message := parts[0], parts[2], parts[3]
							chat.Receive(scope, username, message)

						case "NPC":
							fmt.Println(m)
							switch m[2] {
							case "DEAD":
								conn.Send("LOOK")
							case "FLEE":
								// npc flee
							case "ATTACK":
								// npc attack
							case "REGENERATION":
								// npc + 5hp
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
							conn.Send("LOOK")
						}
					}
					// EVT QUEST took=<quest_id>
				case "ERR":
					popup(w, msg.Rest)
				}
			})
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
