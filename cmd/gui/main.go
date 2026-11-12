package main

import (
	"fmt"
	"image/color"
	"net"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/canvas"
)

// /===========================================\
// |				 Var/Const				   |
// \===========================================/

// Colors
var (
	lime = color.RGBA{R: 33, G: 222, B: 20, A: 255}
	cyan = color.RGBA{R: 17, G: 233, B: 237, A: 255}
	orange = color.RGBA{R: 237, G: 123, B: 17, A: 255}
)

// /===========================================\
// |				  Display				   |
// \===========================================/

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error listening:", err)
		return
	}

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
	background := canvas.NewImageFromFile("./data/VillageMap.png")
	background.FillMode = canvas.ImageFillStretch

	btnExit := widget.NewButton("EXIT", func() { fmt.Fprintln(conn, "QUIT") })
	btnExit.Importance = widget.HighImportance

	txt1 := widget.NewLabel("topl")
	chat := NewChat(w)

	btnContainer := newButtonContainer(a)
	hudContent := container.New(&Hud{}, btnExit, txt1, btnContainer, chat.Content())
	FullContent := container.NewStack(
		background,
		overlay,
		hudContent,
	)

	w.SetContent(FullContent)
	w.ShowAndRun()
	fmt.Println("Exited, end of TAP")
}
