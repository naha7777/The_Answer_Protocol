package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ServerResponse struct {
	RawRoom	ServerData	`json:"room"`
}

type ServerData struct {
	Id			string	`json:"Id"`
	Name		string	`json:"Name"`
	Description	string	`json:"Description"`
	Exits		map[string]string	`json:"Exits"`
	Players		[]string	`json:"Players"`
	Items		[]string	`json:"Items"`
	NPCs		[]string	`json:"NPCs"`
}

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
	if p.group.Id == "" {
		b.group.SetText("Group: -")
		b.members.SetText("")
		return
	}
	b.group.SetText(fmt.Sprintf("Group: %s (leader %s)", p.group.Id, p.group.Leader))
	b.members.SetText("Members: " + strings.Join(p.group.Members, ", "))
}

func (player *Player) PlayerRoomUpdate(msg string) {
	rawJSON := strings.TrimPrefix(msg, "look=")
	rawJSON = strings.TrimSpace(rawJSON)
	if rawJSON == "" {
		return
	}
	var resp ServerResponse
	if err := json.Unmarshal([]byte(rawJSON), &resp); err != nil {
		fmt.Println("Update playerinfo error")
		return
	}
	if resp.RawRoom.Name != "" {
		player.room.name = strings.ToUpper(resp.RawRoom.Name)
	}
	player.room.exits = map[string]bool {
		"north": resp.RawRoom.Exits["north"] != "",
		"south": resp.RawRoom.Exits["south"] != "",
		"east": resp.RawRoom.Exits["east"] != "",
		"west": resp.RawRoom.Exits["west"] != "",
	}
}