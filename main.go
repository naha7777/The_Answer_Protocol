package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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
	a := app.New()
	w := a.NewWindow("TAP")
	w.Resize(fyne.NewSize(1920, 1080))

	Str_to_quests(`[{"quest_id": "princess_necklace", "status": "active", "progress": "0/1"}, {"quest_id: "example", "status": "completed"}]`, w)

	w.Canvas().SetOnTypedKey(func(keyEvent *fyne.KeyEvent) {
		if keyEvent.Name == fyne.KeyEscape {
			a.Quit()
		}
	})
	overlay := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 90})
	background := canvas.NewImageFromFile("pictures/Village_Map.png")
	background.FillMode = canvas.ImageFillStretch

	game := NewGame(w, background)
	chat := NewChat(w)

	btnContainer := newButtonContainer(a, game)
	hudContent := container.New(&Hud{}, btnContainer, chat.Content())
	FullContent := container.NewStack(
		background,
		overlay,
		hudContent,
	)

	w.SetContent(FullContent)
	w.ShowAndRun()
	fmt.Println("Exited, end of TAP")
}
