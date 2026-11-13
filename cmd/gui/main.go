package main

import (
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
	chat := NewChat(w)

	// ajouter conn en parametre de la fonction
	btnContainer := newButtonContainer(w, conn, background)
	hudContent := container.New(&Hud{}, btnContainer, chat.Content())
	gameContent := container.NewStack(
		background,
		overlay,
		hudContent,
	)

	loggedIn := false

	loginView, onLoginMsg := LoginScreen(conn ,func(name string) {
		chat.name = name
		loggedIn = true
		w.SetContent(gameContent)
	})

	w.SetContent(loginView)

	go func() {
		for msg := range conn.Messages() {
			fyne.Do(func() {
				if !loggedIn {
					onLoginMsg(msg)
					return
				}
				if msg.Kind == "OK" && strings.HasPrefix(msg.Rest, "room=") {
					roomID := strings.TrimPrefix(msg.Rest, "room=loc.")
					background.File = "./data/" + roomID + ".jpeg"
					background.Refresh()
				}
				logView.ScrollToBottom()
				logLabel.SetText(logLabel.Text + "-- lost connection --\n")
				fmt.Fprintln(os.Stderr, "-- lost connection --\n")
			})
		}
	}()

	w.ShowAndRun()
	fmt.Println("Exited, end of TAP")
}
