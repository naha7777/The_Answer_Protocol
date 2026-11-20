package main

// /===========================================\
// |				  Import				   |
// \===========================================/

import (
	"image/color"
	"strings"

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

// type Layout struct {
// 	backgroundFile string
// 	container      *fyne.Container
// }

type CardinalButtons struct {
	North *widget.Button
	South *widget.Button
	East  *widget.Button
	West  *widget.Button
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

func newButtonContainer(background *canvas.Image, p *Player, d *Data) *fyne.Container {

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
			if err := d.conn.Send(cmd); err != nil {
				fyne.Do(func() { popup(d.w, "901 SEND_FAILED") })
			}
		}
	}

	// Actions Buttons
	var butAttack *widget.Button
	butAttack = newBtn("ATTACK", func() { d.g.popup = Option_popup(butAttack, d.w, Generate_option(d.g.npc, "ATTACK", d.conn)) })
	var butTalk *widget.Button
	butTalk = newBtn("TALK", func() { d.g.popup = Option_popup(butTalk, d.w, Generate_option(d.g.npc, "TALK", d.conn)) })
	var butQuest *widget.Button
	butQuest = newBtn("QUEST", func() { d.g.popup = Option_popup(butQuest, d.w, Generate_option(d.g.npc, "QUEST", d.conn)) })

	butLook := newBtn("LOOK", send("LOOK"))
	butStatus := newBtn("STATUS", send("STATUS"))
	butQuests := newBtn("QUESTS", send("QUESTS"))
	butWho := newBtn("WHO", send("WHO"))
	butMap := newBtn("MAP", toggleMap)

	butQuit := widget.NewButton("QUIT", func() {
		confirmAction(d.w, "QUIT", "YOU REALLY WANT TO QUIT ?", send("QUIT"))
	})

	butQuit.Importance = widget.LowImportance

	var butGroup *widget.Button

	askAndSend := func(title, label, cmd string) func() {
		return func() {
			askText(d.w, title, label, func(text string) {
				send(strings.TrimSpace(cmd + " " + text))()
			})
		}
	}
	butGroup = newBtn("GROUP", func() {
		Option_popup(butGroup, d.w,
			newBtn("CREATE", askAndSend("GROUP", "Group name", "GROUP CREATE")),
			newBtn("LEAVE", func() { confirmAction(d.w, "LEAVE", "YOU REALLY WANT TO LEAVE YOUR GROUP ?", send("GROUP LEAVE")) }),
			newBtn("INVIT", askAndSend("GROUP", "Player to invite", "GROUP INVITE")),
		)
	})

	// Cardinals Buttons
	d.cb.North = newBtn("NORTH", send("MOVE NORTH"))
	d.cb.South = newBtn("SOUTH", send("MOVE SOUTH"))
	d.cb.East = newBtn("EAST", send("MOVE EAST"))
	d.cb.West = newBtn("WEST", send("MOVE WEST"))

	butsActions := container.NewGridWithColumns(3,
		butLook, butTalk, butAttack,
		widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""),
		butStatus, butQuest, butQuests,
		widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""),
		butWho, butGroup, butQuit,
	)

	butsCardinals := container.NewGridWithColumns(3,
		widget.NewLabel(""), d.cb.North, widget.NewLabel(""),
		d.cb.West, butMap, d.cb.East,
		widget.NewLabel(""), d.cb.South, widget.NewLabel(""),
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

// Method for showing/hiding buttons based on the “exits” map
func (cb *CardinalButtons) Update(exits map[string]bool) {
	if exits == nil {
		exits = make(map[string]bool)
	}
	updateBtn := func(btn *widget.Button, isOpen bool) {
		if btn == nil {
			return
		}
		if isOpen {
			btn.Enable()
		} else {
			btn.Disable()
		}
	}

	updateBtn(cb.North, exits["north"])
	updateBtn(cb.South, exits["south"])
	updateBtn(cb.East, exits["east"])
	updateBtn(cb.West, exits["west"])
}
