package main

import (
	"fmt"
	"image/color"

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
	a := app.New()
	w := a.NewWindow("TAP")
	w.Resize(fyne.NewSize(1920, 1080))
	
	// Key hook for esc to close the app
	w.Canvas().SetOnTypedKey(func(keyEvent *fyne.KeyEvent) {
		if keyEvent.Name == fyne.KeyEscape {
			a.Quit()
		}
	})

	txt1 := widget.NewLabel("topl")
	chat := NewChat(w)

	img := canvas.NewImageFromFile("tavern.png")
	img.Resize(fyne.Size{Width: 200, Height: 200})
	img.FillMode = canvas.ImageFillOriginal

	but1 := widget.NewButton("ATTACK", func() {fmt.Println("boom")})
	but2 := widget.NewButton("GROUP", func() {chat.inGroup = true})
	buttonsContainer := container.NewHBox(but1, but2)


	w.SetContent(container.New(&Hud{}, txt1, img, buttonsContainer, chat.Content()))
	w.ShowAndRun()
	fmt.Println("Exited, end of TAP")
}