package main

// /===========================================\
// |				  Import				   |
// \===========================================/

import (
	"image/color"
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
	window		fyne.Window
	activeChat	string
	txtColor	color.Color
	inGroup		bool

	messages	*fyne.Container
	scroll		*container.Scroll
	entry		*widget.Entry
	root		fyne.CanvasObject
	player		*Player
	msg			bool
}

// /===========================================\
// |				Constructor 			   |
// \===========================================/

func NewChat(w fyne.Window, p *Player, conn *client.Conn) *Chat {
	c := &Chat{
		window:		w,
		activeChat:	"Global",
		txtColor:	lime,
		player:		p,
		msg:		false,
	}
	
	c.messages = container.NewVBox()
	c.scroll = container.NewScroll(c.messages)
	c.entry = widget.NewEntry()
	c.entry.SetPlaceHolder("Global chat...")
	// Send with enter key
	c.entry.OnSubmitted = func(string) {c.send()}

	// Chat Buttons
	cbut_glob := widget.NewButton("Global", func() {c.setChannel("Global", lime)})
	cbut_room := widget.NewButton("Room", func() {c.setChannel("Room", cyan)})
	cbut_grp := widget.NewButton("Group", func() {
		if c.inGroup {
			c.setChannel(c.player.group, orange)
		} else {
			popup(w, "You are not in a group")
		}
	})
	// cbut_say := widget.NewButton("SAY", c.send)
	cbut_say := widget.NewButton("SAY", func() {
		conn.Send("CHAT" + c.activeChat + c.entry.Text)
		if c.msg {c.send()}
		c.msg = false
	})

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
	c.messages.Add(canvas.NewText("[" + channel + "] " + get_time() + " [" + author +" ]: " + text, col))
	c.scroll.ScrollToBottom()
}

func (c *Chat) send() {
	if c.entry.Text == "" {
		return
	}
	c.addLine(c.activeChat, c.player.name, c.entry.Text, c.txtColor)
}

func (c *Chat) setChannel(name string, col color.Color) {
	c.activeChat = name
	c.txtColor = col
	c.entry.SetPlaceHolder(name + " chat...")
}

func (c *Chat) Receive(channel, author, text string) {
	col := color.Color(lime)
	switch channel {
	case "ROOM":
		col = cyan
	case "GROUP":
		col = orange
	}
	c.addLine(channel, author, text, col)
}