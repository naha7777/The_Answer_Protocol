package main

// /===========================================\
// |				  Import				   |
// \===========================================/

import (
	"image/color"
	"strings"
	"the_answer_protocol/client"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// /===========================================\
// |				   Struct				   |
// \===========================================/

// struct for the HUD needed to implement the fyne layout Interface
// /!\ Need at least 4 widget
type Hud struct {
}

type Layout struct {
	backgroundFile string
	container      *fyne.Container
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
	if len(objects) < 3 {
		return
	}
	w, h := containerSize.Width, containerSize.Height

	margin := float32(10)

	// Buttons
	btnSize := objects[0].MinSize()
	btnX := margin
	Y := h - btnSize.Height - margin
	objects[0].Resize(btnSize)
	objects[0].Move(fyne.NewPos(btnX, Y))

	// Chat
	chatX := btnX + btnSize.Width + margin
	chatWidth := (w - chatX - margin) / 2
	if chatWidth < 50 {
		chatWidth = 50
	}

	chatHeight := btnSize.Height

	objects[1].Resize(fyne.NewSize(chatWidth, chatHeight))
	objects[1].Move(fyne.NewPos(chatX, Y))

	// Inventory
	invX := chatX + chatWidth + margin
	invY := Y - 50
	invWidth := (w - invX - margin) / 3
	if invWidth < 50 {
		invWidth = 50
	}

	invHeight := h - invY - margin
	objects[2].Resize(fyne.NewSize(invWidth, invHeight))
	objects[2].Move(fyne.NewPos(invX, invY))
}

func newButtonContainer(w fyne.Window, conn *client.Conn, background *canvas.Image, p *Player) *fyne.Container {

	mapImg := "./data/map.jpeg"
	previous := background.File

	toggleMap := func() {
		if background.File == mapImg {
			background.File = previous
		} else {
			previous = background.File
			background.File = mapImg
		}
		background.Refresh()
	}

	newBtn := func(label string, action func()) *widget.Button {
		btn := widget.NewButton(label, action)
		btn.Importance = widget.LowImportance
		return btn
	}

	send := func(cmd string) func() {
		return func() {
			if strings.HasPrefix(cmd, "MOVE") {
				p.move = true
			}
			if err := conn.Send(cmd); err != nil {
				fyne.Do(func() { popup(w, "901 SEND_FAILED") })
			}
		}
	}

	// Actions Buttons
	butAttack := newBtn("ATTACK", send("ATTACK"))
	butLook := newBtn("LOOK", send("LOOK"))
	// butMove := newBtn("MOVE", send("MOVE") )
	butTake := newBtn("TAKE", send("TAKE"))
	// butDrop := newBtn("DROP", send("DROP") )
	// butTalk := newBtn("TALK", send("TALK") )
	butStatus := newBtn("STATUS", send("STATUS"))
	butQuest := newBtn("QUEST", send("QUEST"))
	butQuests := newBtn("QUESTS", send("QUESTS"))
	butWho := newBtn("WHO", send("WHO"))
	butMap := newBtn("MAP", toggleMap)

	butQuit := widget.NewButton("QUIT", func() {
		d := dialog.NewConfirm(
			"QUIT",
			"YOU REALLY WANT TO QUIT ?",
			func(confirmed bool) {
				if !confirmed {
					return
				}
				send("QUIT")()
			},
			w,
		)
		d.SetConfirmText("YES")
		d.SetDismissText("NO")
		d.Show()
	})
	butQuit.Importance = widget.LowImportance

	var butGroup *widget.Button
	butGroup = newBtn("GROUP", func() {
		Option_popup(butGroup, w,
			newBtn("CREATE", func() { send("GROUP CREATE")() }),
			newBtn("QUIT", func() { send("GROUP LEAVE")() }),
			newBtn("INVIT", func() { send("GROUP INVITE")() }),
		)
	})

	// Cardinals Buttons
	butNorth := newBtn("NORTH", send("MOVE NORTH"))
	butSouth := newBtn("SOUTH", send("MOVE SOUTH"))
	butEast := newBtn("EAST", send("MOVE EAST"))
	butWest := newBtn("WEST", send("MOVE WEST"))

	butsActions := container.NewGridWithColumns(3,
		// widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""),
		butLook, butTake, butAttack,
		widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""),
		// butDrop, butTalk, butMove,
		butStatus, butQuest, butQuests,
		widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""),
		butWho, butGroup, butQuit,
		// widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""),
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
