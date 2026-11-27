package main

// /===========================================\
// |				  Import				   |
// \===========================================/

import (
	"fmt"
	"image/color"
	"os"

	"the_answer_protocol/client"

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
	lime   = color.RGBA{R: 33, G: 222, B: 20, A: 255}  // Associated to GLOBAL Chat
	cyan   = color.RGBA{R: 17, G: 233, B: 237, A: 255} // Associated to ROOM Chat
	orange = color.RGBA{R: 237, G: 123, B: 17, A: 255} // Associated to GROUP Chat
	yellow = color.RGBA{R: 255, G: 213, B: 28, A: 255} // Associated to NPC Chat
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
	player := Player{maxHP: 0, currentHP: 0, move: false}
	w.Resize(fyne.NewSize(1920, 1080))

	// Key hook for esc to close the app
	w.Canvas().SetOnTypedKey(func(keyEvent *fyne.KeyEvent) {
		if keyEvent.Name == fyne.KeyEscape {
			a.Quit()
		}
	})

	overlay := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 90})

	// Struct info
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
	hudContent := container.New(&Hud{}, btnContainer, data.chat.Content(), data.g.inventory.Content())
	gameContent := container.NewStack(
		background,
		overlay,
		hudContent,
		container.NewVBox(container.NewHBox(data.playerInfo.Content())),
	)

	// Widgets init for combat layout
	arenaBg := canvas.NewImageFromFile("./data/arena.jpeg")
	arenaBg.FillMode = canvas.ImageFillCover

	data.battleUI = NewBattleUI(&player, &data)
	battleContent := container.NewStack(
		arenaBg,
		overlay,
		data.battleUI.Content(),
	)

	loggedIn := false
	gameOver := false
	loginView, onLoginMsg := LoginScreen(conn, &player, func(name string) {
		loggedIn = true
		data.cb.Update(player.room.exits)
		conn.Send("WHO")
		conn.Send("STATUS")
		conn.Send("INVENTORY")
		conn.Send("LOOK")
		data.playerInfo.Update(data)
		w.SetContent(gameContent)
	})

	data.gameContent = gameContent
	data.battleContent = battleContent
	w.SetContent(loginView)

	go func() {
		for msg := range conn.Messages() {
			fyne.Do(func() {
				if !loggedIn {
					onLoginMsg(msg)
					return
				}
				fmt.Println(msg.Raw)
				if msg.Raw == "GAME OVER" {
					gameOver = true
					PopupGameOver(data.w, func() {
						a.Quit()
					})
					return
				}
				switch msg.Kind {
				case "OK":
					ServerResponseOK(msg, &player, &data)
				case "EVT":
					ServerResponseEVT(msg, &player, &data)
				case "ERR":
					ServerResponseERR(msg, &player, &data)
				}
			},
			)
		}

		fyne.Do(func() {
			if !gameOver {
				logLabel.SetText(logLabel.Text + "-- lost connection --\n")
				logView.ScrollToBottom()
				fmt.Fprintln(os.Stderr, "-- lost connection --")
				a.Quit()
			}
		})
	}()

	w.ShowAndRun()
	fmt.Println("Exited, end of TAP")
}
