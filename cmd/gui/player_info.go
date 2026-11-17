package main

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type InfoBox struct {
	player	*Player
	name	*widget.Label
	room	*widget.Label
	group	*widget.Label
	members	*widget.Label
	content	fyne.CanvasObject
}

func showPlayerInfo(p *Player) *InfoBox {

	b := &InfoBox{
		player:		p,
		name:		widget.NewLabel(""),
		room:		widget.NewLabel(""),
		group:		widget.NewLabel(""),
		members:	widget.NewLabel(""),
	}
	b.members.Wrapping = fyne.TextWrapWord

	frame := canvas.NewRectangle(color.NRGBA{R: 20, G: 20, B: 20, A: 200})
	frame.StrokeColor = color.White
	frame.StrokeWidth = 2

	b.content = container.NewStack(
		frame,
		container.NewPadded(container.NewVBox(
			widget.NewLabelWithStyle("PLAYER INFO", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),
			b.name, b.room, b.group, b.members,
		)),
	)
	b.Update()
	return b
}

func (b *InfoBox) Content() fyne.CanvasObject {return b.content}

func (b *InfoBox) Update() {
	p := b.player
	b.name.SetText("Name: " + p.name)
	b.room.SetText("Room: " + p.room.name)
	if p.group.id == "" {
		b.group.SetText("Group: -")
		b.members.SetText("")
		return
	}
	b.group.SetText(fmt.Sprintf("Group: %s (leader %s)", p.group.id, p.group.leader))
	b.members.SetText("Members: " + strings.Join(p.group.members, ", "))
}
