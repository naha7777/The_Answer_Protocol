package main

// /===========================================\
// |				  Import				   |
// \===========================================/

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// /===========================================\
// |				   Struct				   |
// \===========================================/

type Chat struct {
	window	fyne.Window
	activeChat string
	txtColor color.Color
	inGroup bool

	messages *fyne.Container
	scroll *container.Scroll
	entry *widget.Entry
	root fyne.CanvasObject
}

// /===========================================\
// |				Constructor 			   |
// \===========================================/

func NewChat(w fyne.Window) *Chat {
	c := &Chat{
		window: w,
		activeChat: "Global",
		txtColor: lime,
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
			c.setChannel("Group", orange)
		} else {
			popup(w, "You are not in a group")
		}
	})
	cbut_say := widget.NewButton("SAY", c.send)

	// Containers on top of each others
	cbut_box := container.NewBorder(nil, nil, container.NewHBox(cbut_glob, cbut_room, cbut_grp), nil)
	send_box := container.NewBorder(nil, nil, nil, cbut_say, c.entry)
	c.root = container.NewBorder(cbut_box, send_box, nil, nil, c.scroll)

	return c
}

// /===========================================\
// |				   Method				   |
// \===========================================/

func (c *Chat) Content() fyne.CanvasObject {
	return c.root
}

func (c *Chat) send() {
	if c.entry.Text != "" {
		c.messages.Add(canvas.NewText("[" + c.activeChat + "] " + get_time() +" Me: " + c.entry.Text, c.txtColor))
		c.entry.SetText("")
		c.scroll.ScrollToBottom()
	}
}

func (c *Chat) setChannel(name string, col color.Color) {
	c.activeChat = name
	c.txtColor = col
	c.entry.SetPlaceHolder(name + " chat...")
}
