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

	if len(objects) < 2 {
		return fyne.NewSize(0, 0)
	}
	btnMin := objects[0].MinSize()
	chatMin := objects[1].MinSize()

	margin := float32(10)
	minW := btnMin.Width + chatMin.Width + 3*margin
	minH := btnMin.Height + 2*margin
	return fyne.NewSize(minW, minH)
}

func (ui *Hud) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) < 2 {
		return
	}
	w, h := containerSize.Width, containerSize.Height

	margin := float32(10)

	btnSize := objects[0].MinSize()
	btnX := margin
	btnY := h - btnSize.Height - margin
	objects[0].Resize(btnSize)
	objects[0].Move(fyne.NewPos(btnX, btnY))

	chatX := btnX + btnSize.Width + margin
	chatWidth := (w - chatX - margin) / 2
	if chatWidth < 50 {
		chatWidth = 50
	}

	chatHeight := btnSize.Height
	chatY := btnY

	objects[1].Resize(fyne.NewSize(chatWidth, chatHeight))
	objects[1].Move(fyne.NewPos(chatX, chatY))
}

func newButtonContainer(app fyne.App, g *Game) *fyne.Container {
	
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
	butMap := newBtn("MAP", func() { g.ToggleMap() })
	butQuit.Importance = widget.LowImportance

	// Cardinals Buttons
	butNorth := newBtn("NORTH", func() { g.Move("north") })
	butSouth := newBtn("SOUTH", func() { g.Move("south") })
	butEast := newBtn("EAST", func() { g.Move("east") })
	butWest := newBtn("WEST", func() { g.Move("west") })

	butsActions := container.NewGridWithColumns(3,
		butLook, butMove, butTake,
		butDrop, butTalk, butAttack,
		butStatus, butQuest, butQuests,
		butWho, butGroup, butQuit,
	)

	butsCardinals := container.NewGridWithColumns(3,
		widget.NewLabel(""), butNorth, widget.NewLabel(""),
		butWest, butMap, butEast,
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
