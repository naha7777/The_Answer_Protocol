package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"strings"
	"the_answer_protocol/client"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type GroupInfo struct {
	id      string
	leader  string
	members []string
	invited []string
}

type groupInfosMsg struct {
	Group GroupInfo `json:"group"`
}

type Player struct {
	name   string
	health int
	group  GroupInfo
	room   string
	move   bool
}

// /===========================================\
// |				 Var/Const				   |
// \===========================================/

// Colors
var (
	lime   = color.RGBA{R: 33, G: 222, B: 20, A: 255}
	cyan   = color.RGBA{R: 17, G: 233, B: 237, A: 255}
	orange = color.RGBA{R: 237, G: 123, B: 17, A: 255}
)

// /===========================================\
// |				  Display				   |
// \===========================================/

func main() {
	player := Player{health: 100, move: false}

	conn, err := client.Dial("localhost:8080")
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERR 900 CONNECTION_FAILED")
		os.Exit(1)
	}
	defer conn.Close()

	a := app.New()
	w := a.NewWindow("TAP")
	w.Resize(fyne.NewSize(1920, 1080))

	// Key hook for esc to close the app
	w.Canvas().SetOnTypedKey(func(keyEvent *fyne.KeyEvent) {
		if keyEvent.Name == fyne.KeyEscape {
			a.Quit()
		}
	})

	overlay := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 90})
	background := canvas.NewImageFromFile("./data/square.jpeg")
	background.FillMode = canvas.ImageFillStretch

	logLabel := widget.NewLabel("")
	logLabel.Wrapping = fyne.TextWrapWord
	logView := container.NewVScroll(logLabel)

	// ajouter conn en parametre de la fonction
	chat := NewChat(w, &player, conn)
	inv := NewInv(w, &player, conn)

	// ajouter conn en parametre de la fonction
	btnContainer := newButtonContainer(w, conn, background, &player)
	hudContent := container.New(&Hud{}, btnContainer, chat.Content(), inv.Content())
	gameContent := container.NewStack(
		background,
		overlay,
		hudContent,
	)

	loggedIn := false
	// playerinfo :=
	loginView, onLoginMsg := LoginScreen(conn, func(name string) {
		player.name = name
		loggedIn = true
		w.SetContent(gameContent)
		conn.Send("INVENTORY")
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
					switch {
					case strings.HasPrefix(msg.Rest, "room=") && player.move:
						roomID := strings.TrimPrefix(msg.Rest, "room=loc.")
						player.room = roomID
						background.File = "./data/" + roomID + ".jpeg"
						background.Refresh()
						player.move = false

					case msg.Rest == "bye":
						a.Quit()

					// GROUP
					case strings.HasPrefix(msg.Rest, "group_created="):
						popup(w, strings.ToUpper(player.group.id)+" group was created")
						conn.Send("INFO_GROUP")
					case strings.HasPrefix(msg.Rest, "group_leaved="):
						popup(w, "You leave "+strings.ToUpper(player.group.id)+" group")
						player.group = GroupInfo{}
					case strings.HasPrefix(msg.Rest, "group_joined="):
						popup(w, "You join the "+strings.ToUpper(player.group.id)+"group")
						conn.Send("INFO_GROUP")
					case strings.HasPrefix(msg.Rest, "group_infos="):
						raw := strings.TrimPrefix(msg.Rest, "group_infos=")
						var infos groupInfosMsg
						if err := json.Unmarshal([]byte(raw), &infos); err != nil {
							popup(w, "invalid group infos")
							break
						}
						player.group = infos.Group
					case strings.HasPrefix(msg.Rest, "taken="):
						item := strings.TrimPrefix(msg.Rest, "taken=item.")
						item = strings.Title(strings.ToLower(item))
						item = strings.ReplaceAll(item, "_", " ")
						inv.take(item)
					case strings.HasPrefix(msg.Rest, "dropped="):
						item := strings.TrimPrefix(msg.Rest, "dropped=item.")
						item = strings.Title(strings.ToLower(item))
						item = strings.ReplaceAll(item, "_", " ")
						inv.drop(item)
					case strings.HasPrefix(msg.Rest, "inventory="):
						inv.BuildInv(msg.Rest)
						inv.Update(msg.Rest)
					case strings.HasPrefix(msg.Rest, "look="):
						inv.Update(msg.Rest)
					}

				case "EVT":
					m := strings.Split(msg.Raw, " ")
					if len(m) >= 3 && m[1] == "GROUP" {
						switch m[2] {
						case "INVITE":
							if len(m) >= 5 {
								confirmGroup(w, []string{m[4], m[3]}, func() {
									conn.Send("GROUP JOIN " + m[3])
								})
							}
						case "JOIN", "LEAVE":
							conn.Send("GROUP INFOS")
						}
					}
					parts := strings.SplitN(msg.Rest, " ", 4)
					if len(parts) == 4 && parts[1] == "CHAT" {
						chat.Receive(parts[0], parts[2], parts[3])
					}

				case "ERR":
					popup(w, msg.Rest)
				}
			})
		}

		fyne.Do(func() {
			logLabel.SetText(logLabel.Text + "-- lost connection --\n")
			logView.ScrollToBottom()
		})
		fmt.Fprintln(os.Stderr, "-- lost connection --")
	}()

	w.ShowAndRun()
	fmt.Println("Exited, end of TAP")
}
