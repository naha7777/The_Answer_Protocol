package main

// /===========================================\
// |				  Import				   |
// \===========================================/

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// /===========================================\
// |				 Function				   |
// \===========================================/

func popup(w fyne.Window, txt string) {
	width := w.Canvas().Size().Width
	height := w.Canvas().Size().Height

	popup := widget.NewPopUp(widget.NewLabel(txt), w.Canvas())
	popup.ShowAtPosition(fyne.Position{X: width/2, Y: 3*(height/4)})

	time.AfterFunc(3*time.Second, func() {
			fyne.Do(func() {
				popup.Hide()
		})
	})
}

func get_time() string {
	formatted := time.Now().Format("15:04")
	return formatted
}
