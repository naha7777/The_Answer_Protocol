package main

// /===========================================\
// |				  Import				   |
// \===========================================/

import (
	"fmt"

	"fyne.io/fyne/v2"
)

// /===========================================\
// |				   Struct				   |
// \===========================================/

// struct for the HUD needed to implement the fyne layout Interface
// /!\ Need at least 4 widget
type Hud struct {

}

// /===========================================\
// |				 Function				   |
// \===========================================/

// Need to implement the Layout Interface with MinSize and Layout func

func (ui *Hud) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) < 4 { fmt.Println("Error: Hud needs 4 arguments") }
	w, h := float32(0), float32(0)
	for _, o := range objects {
		childSize := o.MinSize()

		w += childSize.Width
		h += childSize.Height
	}
	return fyne.NewSize(w, h)
}

func (ui *Hud) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) < 4 { return }
	w, h := containerSize.Width, containerSize.Height

	for _, o := range objects {
		o.Resize(o.MinSize())
	}

	chat_size := fyne.Size{Width: containerSize.Width / 3, Height: containerSize.Height / 4}
	objects[3].Resize(chat_size)

	objects[0].Move(fyne.NewPos(0, 0))
	objects[1].Move(fyne.NewPos(w - objects[1].MinSize().Width, 0))
	objects[2].Move(fyne.NewPos(0, h - objects[2].MinSize().Height))
	objects[3].Move(fyne.NewPos(w - chat_size.Width - 10, h - chat_size.Height))
}