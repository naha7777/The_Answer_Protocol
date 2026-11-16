package main

// /===========================================\
// |				  Import				   |
// \===========================================/

import (
	"image/color"
	"os"
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

type Inventory struct {
	window fyne.Window

	items  *fyne.Container
	scroll *container.Scroll
	root   fyne.CanvasObject
}

// /===========================================\
// |				Constructor 			   |
// \===========================================/

func NewInv(w fyne.Window, p *Player, conn *client.Conn) *Inventory {
	i := &Inventory{
		window: w,
	}

	i.items = container.NewVBox()
	i.scroll = container.NewScroll(i.items)
	container.NewBorder(container.NewCenter(widget.NewLabel("Inventory")), nil, nil, nil, i.scroll)

	// frame around the inventory
	invFrame := canvas.NewRectangle(color.NRGBA{R: 20, G: 20, B: 20, A: 200})
	invFrame.StrokeColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	invFrame.StrokeWidth = 2

	// Buttons
	take_but := widget.NewButton("TAKE", func() { conn.Send("TAKE Herbs") })
	drop_but := widget.NewButton("DROP", func() { conn.Send("DROP Herbs") })

	// Containers on top of each others
	inv_title := container.NewCenter(widget.NewLabel("Inventory"))
	but_cont := container.NewCenter(container.NewHBox(take_but, drop_but))
	inv_cont := container.NewBorder(container.NewVBox(inv_title, widget.NewSeparator()), but_cont, nil, nil, i.scroll)
	inside := container.NewBorder(nil, nil, nil, nil, inv_cont)
	i.root = container.NewStack(invFrame, container.NewPadded(inside))

	return i
}

// /===========================================\
// |				   Method				   |
// \===========================================/

func (i *Inventory) Content() fyne.CanvasObject {
	return i.root
}

func (i *Inventory) BuildInv(str string) {
	// parse
	str = strings.Trim(str, "[]")
	items := strings.Split(str, ", ")

	for _, item := range items {
		item = strings.Trim(item, `"`)
		item, _ = strings.CutPrefix(item, "item.")

		// img
		image := canvas.NewImageFromFile("data/icons/placeholder.png")
		if _, err := os.Stat("data/items/" + item + ".png"); err == nil {
			image = canvas.NewImageFromFile("data/items/" + item + ".png")
		}

		image.FillMode = canvas.ImageFillContain
		image.SetMinSize(fyne.NewSize(50, 50))

		item = strings.ReplaceAll(item, "_", " ")
		item = strings.Title(strings.ToLower(item))
		single_cont := container.NewBorder(nil, nil, image, nil, container.NewCenter(widget.NewLabel(item)))
		i.items.Add(single_cont)
		i.items.Add(widget.NewSeparator())
	}
}

func (i *Inventory) drop(item string) {
	for n, obj := range i.items.Objects {
		if cont, ok := obj.(*fyne.Container); ok {
			text := cont.Objects[0].(*fyne.Container).Objects[0].(*widget.Label).Text
			if text == item {
				sep := i.items.Objects[n+1]
				i.items.Remove(obj)
				i.items.Remove(sep)
				return
			}
		}
	}
}

func (i *Inventory) take(item string) {
	// img
	item_file := strings.ReplaceAll(item, " ", "_")
	item_file = strings.ToLower(item_file)
	image := canvas.NewImageFromFile("data/icons/placeholder.png")
	if _, err := os.Stat("data/items/" + item_file + ".png"); err == nil {
		image = canvas.NewImageFromFile("data/items/" + item_file + ".png")
	}

	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(50, 50))

	// cont
	single_cont := container.NewBorder(nil, nil, image, nil, container.NewCenter(widget.NewLabel(item)))
	i.items.Add(single_cont)
	i.items.Add(widget.NewSeparator())
}
