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
	RawRoom ServerData `json:"room"`
}

type ServerData struct {
	Id          string            `json:"Id"`
	Name        string            `json:"Name"`
	Description string            `json:"Description"`
	Exits       map[string]string `json:"Exits"`
	Players     []string          `json:"Players"`
	Items       []string          `json:"Items"`
	NPCs        []string          `json:"NPCs"`
}

type InfoBox struct {
	player		*Player
	name		*widget.Label
	room		*widget.Label
	group		*widget.Label
	members		*widget.Label
	npcList		*fyne.Container
	itemList	*fyne.Container
	content		fyne.CanvasObject
	bgImg		*canvas.Image
}

func showPlayerInfo(p *Player) *InfoBox {

	b := &InfoBox{
		player:		p,
		name:		widget.NewLabel(""),
		room:		widget.NewLabel(""),
		group:		widget.NewLabel(""),
		members:	widget.NewLabel(""),
		npcList:	container.NewVBox(),
		itemList:	container.NewVBox(),
		content:	widget.NewLabel(""),
		bgImg:		canvas.NewImageFromFile(""),
	}

	for _, lbl := range []*widget.Label{b.name, b.room, b.group, b.members} {
		lbl.Wrapping = fyne.TextWrapWord
	}

	frame := canvas.NewRectangle(color.NRGBA{R: 15, G: 15, B: 15, A: 110})
	frame.StrokeColor = color.White
	frame.StrokeWidth = 2

	sizer := canvas.NewRectangle(color.Transparent)
	sizer.SetMinSize(fyne.NewSize(280, 250))


	b.bgImg.FillMode = canvas.ImageFillStretch

	header := container.NewVBox(
		widget.NewSeparator(),
	)

	scrollContent := container.NewVScroll(
		container.NewVBox(b.name, b.room, b.group, b.members, b.npcList, b.itemList),
	)

	body := container.NewBorder(header, nil, nil, nil, scrollContent)

	b.content = container.NewStack(
		sizer,
		b.bgImg,
		frame,
		container.NewPadded(body),
	)
	b.Update()
	return b
}

func (b *InfoBox) Content() fyne.CanvasObject { return b.content }

func (b *InfoBox) Update() {
	p := b.player
	if p.avatar != "" && b.bgImg.File != p.avatar {
		b.bgImg.File = p.avatar
		b.bgImg.Refresh()
	}

	b.name.SetText("NAME: " + p.name)
	b.room.SetText("ROOM: " + p.room.name)
	if p.group.Id == "" {
		b.group.SetText("GROUP: -")
		b.members.SetText("")
	} else {
		b.group.SetText(fmt.Sprintf("Group: %s (leader %s)", p.group.Id, p.group.Leader))
		b.members.SetText("MEMBERS: " + strings.Join(p.group.Members, ", "))
	}
	b.npcList.RemoveAll()
	if len(p.room.npc) == 0 {
		b.npcList.Add(widget.NewLabel("NPC: -"))
	} else {
		b.npcList.Add(widget.NewLabel("NPC:"))
		for _, npcName := range p.room.npc {
			lbl := widget.NewLabel("	" + npcName)
			lbl.Wrapping = fyne.TextWrapWord
			b.npcList.Add(lbl)
		}
	}
	b.npcList.Refresh()

	b.itemList.RemoveAll()
	if len(p.room.items) == 0 {
		b.itemList.Add(widget.NewLabel("ITEMS: -"))
	} else {
		b.itemList.Add(widget.NewLabel("ITEMS:"))
		for _, itemName := range p.room.items {
			lbl := widget.NewLabel("	" + itemName)
			lbl.Wrapping = fyne.TextWrapWord
			b.itemList.Add(lbl)
		}
	}
	b.itemList.Refresh()
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
	r := resp.RawRoom
	exits := make(map[string]bool, len(r.Exits))
	for dir, target := range r.Exits {
		exits[dir] = target != ""
	}
	npcs := make([]string, 0, len(r.NPCs))
	for _, id := range r.NPCs {
		npcs = append(npcs, IDToName(id))
	}
	items := make([]string, 0, len(r.Items))
	for _, id := range r.Items {
		items = append(items, IDToName(id))
	}
	player.room = Room{
		id:          r.Id,
		name:        IDToName(r.Name),
		description: r.Description,
		exits:       exits,
		npc:         npcs,
		players:     r.Players,
		items:       items,
	}
}
