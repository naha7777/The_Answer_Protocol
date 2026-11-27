package main

// /===========================================\
// |				  Import				   |
// \===========================================/

import (
	"fmt"
	"image/color"
	"os"
	"strings"

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
		a:    a,
		w:    w,
		conn: conn,
		cb:   &CardinalButtons{},
		g:    &GameState{inventory: NewInv(w, &player, conn), player: &player},
		chat: NewChat(w, &player, conn),
		bg:   background,
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
		conn.Send("WHO")
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
				fmt.Println(msg.Raw)
				if msg.Raw == "GAME OVER" {
					PopupGameOver(data.w)
					a.Quit()
				}
				switch msg.Kind {
				case "OK":
					ServerResponseOK(msg, &player, &data)

				case "EVT":
					ServerResponseEVT(msg, &player, &data)

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
