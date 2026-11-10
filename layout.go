package main

// /===========================================\
// |				  Import				   |
// \===========================================/

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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
	// if len(objects) < 4 {
	// 	fmt.Println("Error: Hud needs 4 arguments")
	// }
	// w, h := float32(0), float32(0)
	// for _, o := range objects {
	// 	childSize := o.MinSize()

	// 	w += childSize.Width
	// 	h += childSize.Height
	// }
	// return fyne.NewSize(w, h)
	if len(objects) < 4 {
		return fyne.NewSize(0,0)
	}
	minW := objects[2].MinSize().Width + 200
	minH := objects[0].MinSize().Height + objects[2].MinSize().Height + 30
	return fyne.NewSize(minW, minH)
}

func (ui *Hud) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) < 4 {
		return
	}
	w, h := containerSize.Width, containerSize.Height

	// for _, o := range objects {
	// 	o.Resize(o.MinSize())
	// }

	// objects[0].Move(fyne.NewPos(0, 0))
	// objects[1].Move(fyne.NewPos(w - objects[1].MinSize().Width, 0))
	// objects[2].Move(fyne.NewPos(0, h-objects[2].MinSize().Height))
	// objects[3].Move(fyne.NewPos(w-objects[3].MinSize().Width, h-objects[3].MinSize().Height))
	margin := float32(10)

	exitSize := objects[0].MinSize()
	objects[0].Resize(exitSize)
	objects[0].Move(fyne.NewPos(margin, margin))

	txtSize := objects[1].MinSize()
	objects[1].Resize(txtSize)
	objects[1].Move(fyne.NewPos(w-txtSize.Width-margin, margin))

	btnSize := objects[2].MinSize()
	btnX := margin
	btnY := h - btnSize.Height - margin
	objects[2].Resize(btnSize)
	objects[2].Move(fyne.NewPos(btnX, btnY))

	chatX := btnX + btnSize.Width + margin
	chatWidth := w - chatX - margin
	if chatWidth < 50 {
		chatWidth = 50
	}

	chatHeight := objects[3].MinSize().Height
	chatY := h - chatHeight - margin

	objects[3].Resize(fyne.NewSize(chatWidth, chatHeight))
	objects[3].Move(fyne.NewPos(chatX, chatY))
}

func newButtonContainer(app fyne.App) *fyne.Container {
	
	newBtn := func(label string, action func()) *widget.Button {
		btn := widget.NewButton(label, action)
		btn.Importance = widget.LowImportance
		return btn
	}

	// Actions Buttons
	butAttack := newBtn("ATTACK", func() { fmt.Println("ATTACK") })
	butLook := newBtn("LOOK", func() { fmt.Println("LOOK") })
	butMove := newBtn("MOVE", func() { fmt.Println("MOVE") })
	butTake := newBtn("TAKE", func() { fmt.Println("TAKE") })
	butDrop := newBtn("DROP", func() { fmt.Println("DROP") })
	butTalk := newBtn("TALK", func() { fmt.Println("TALK") })
	butStatus := newBtn("STATUS", func() { fmt.Println("STATUS") })
	butQuest := newBtn("QUEST", func() { fmt.Println("QUEST") })
	butQuests := newBtn("QUESTS", func() { fmt.Println("QUESTS") })
	butWho := newBtn("WHO", func() { fmt.Println("WHO") })
	butGroup := newBtn("GROUP", func() { fmt.Println("GROUP") })
	butQuit := widget.NewButton("QUIT", func() { app.Quit() })
	butQuit.Importance = widget.LowImportance

	// Cardinals Buttons
	butNorth := newBtn("NORTH", func() { fmt.Println("NORTH") })
	butSouth := newBtn("SOUTH", func() { fmt.Println("SOUTH") })
	butEast := newBtn("EAST", func() { fmt.Println("EAST") })
	butWest := newBtn("WEST", func() { fmt.Println("WEST") })

	butsActions := container.NewGridWithColumns(3,
		butLook, butMove, butTake,
		butDrop, butTalk, butAttack,
		butStatus, butQuest, butQuests,
		butWho, butGroup, butQuit,
	)

	butsCardinals := container.NewGridWithColumns(3,
		widget.NewLabel(""), butNorth, widget.NewLabel(""),
		butWest, widget.NewLabel(""), butEast,
		widget.NewLabel(""), butSouth, widget.NewLabel(""),
	)

	btnsFrame := canvas.NewRectangle(color.NRGBA{R: 20, G: 20, B: 20, A: 200})
	btnsFrame.StrokeColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	btnsFrame.StrokeWidth = 2

	insideBtn := container.NewHBox(
		butsActions,
		widget.NewSeparator(),
		butsCardinals,
	)

	return container.NewStack(
		btnsFrame,
		container.NewPadded(insideBtn),
	)
}
