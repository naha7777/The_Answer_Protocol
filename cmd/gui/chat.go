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
	"fyne.io/fyne/v2/widget"
)

// /===========================================\
// |				   Struct				   |
// \===========================================/

type Chat struct {
	window   fyne.Window
	scope    string
	txtColor color.Color
	inGroup  bool

	messages *fyne.Container
	scroll   *container.Scroll
	entry    *widget.Entry
	root     fyne.CanvasObject
	player   *Player
	conn     *client.Conn
}

// /===========================================\
// |				Constructor 			   |
// \===========================================/

func NewChat(w fyne.Window, p *Player, conn *client.Conn) *Chat {
	c := &Chat{
		window:   w,
		scope:    "GLOBAL",
		txtColor: lime,
		player:   p,
		conn:     conn,
	}

	c.messages = container.NewVBox()
	c.scroll = container.NewVScroll(c.messages)
	c.entry = widget.NewEntry()
	c.entry.SetPlaceHolder("GLOBAL chat...")
	// Send with enter key
	c.entry.OnSubmitted = func(string) { c.send() }

	// Chat Buttons
	cbut_glob := widget.NewButton("Global", func() { c.setChannel("GLOBAL", lime) })
	cbut_room := widget.NewButton("Room", func() { c.setChannel("ROOM", cyan) })
	cbut_grp := widget.NewButton("Group", func() {
		if c.player.group.Id != "" {
			c.setChannel("GROUP", orange)
		} else {
			Popup(w, "You are not in a group")
		}
	})
	cbut_say := widget.NewButton("SAY", c.send)

	// frame around the chat
	chatFrame := canvas.NewRectangle(color.NRGBA{R: 20, G: 20, B: 20, A: 200})
	chatFrame.StrokeColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	chatFrame.StrokeWidth = 2

	// Containers on top of each others
	cbut_box := container.NewBorder(nil, nil, container.NewHBox(cbut_glob, cbut_room, cbut_grp), nil)
	send_box := container.NewBorder(nil, nil, nil, cbut_say, c.entry)
	inside := container.NewBorder(cbut_box, send_box, nil, nil, c.scroll)
	c.root = container.NewStack(chatFrame, container.NewPadded(inside))

	return c
}

// /===========================================\
// |				   Method				   |
// \===========================================/

func (c *Chat) Content() fyne.CanvasObject {
	return c.root
}

func (c *Chat) addLine(channel, author, text string, col color.Color) {
	c.messages.Add(canvas.NewText("["+channel+"] "+get_time()+" | "+author+": "+text, col))
	c.scroll.ScrollToBottom()
}

func (c *Chat) send() {
	text := strings.TrimSpace(c.entry.Text)
	if text == "" {
		return
	}
	if c.scope == "GROUP" && c.player.group.Id == "" {
		Popup(c.window, "You are not in a group")
		return
	}
	if err := c.conn.Send("CHAT " + c.scope + " " + text); err != nil {
		Popup(c.window, "901 SEND_FAILED")
		return
	}
	if c.scope != "GROUP" {
		c.addLine(c.scope, c.player.name, text, c.txtColor)
	}
	c.entry.SetText("")
}

func (c *Chat) setChannel(scope string, col color.Color) {
	c.scope = scope
	c.txtColor = col
	c.entry.SetPlaceHolder(scope + " chat...")
}

func (c *Chat) Receive(channel, author, text string) {
	col := color.Color(lime)
	switch channel {
	case "ROOM":
		col = cyan
	case "GROUP":
		col = orange
	case "NPC":
		col = yellow
	}
	c.addLine(channel, author, text, col)
}
