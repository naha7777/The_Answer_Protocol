package main

// /===========================================\
// |				  Import				   |
// \===========================================/

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"slices"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// /===========================================\
// |				   Struct				   |
// \===========================================/

// objects = items ids in the inventory
// accessible = items in the current room
type Inventory struct {
	objects    []string
	accessible []string
	window     fyne.Window

	pop    *widget.PopUp
	items  *fyne.Container
	scroll *container.Scroll
	root   fyne.CanvasObject
}

// struct for json parsing purposes
type Response struct {
	Items []string `json:"items"`
}

// /===========================================\
// |				Constructor 			   |
// \===========================================/

func NewInv(w fyne.Window, p *Player, d *Data) *Inventory {
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
	var take_but, drop_but *widget.Button
	take_but = widget.NewButton("TAKE", func() { i.pop = Option_popup(take_but, w, Generate_option(i.accessible, "TAKE", d)) })
	drop_but = widget.NewButton("DROP", func() { i.pop = Option_popup(drop_but, w, Generate_option(i.objects, "DROP", d)) })

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
	str = strings.TrimPrefix(str, "inventory=")
	str = strings.Trim(str, "[]")
	if str == "" {
		return
	}
	items := strings.Split(str, ", ")

	for _, item := range items {
		item = strings.Trim(item, `"`)
		item = strings.TrimPrefix(item, "item.")

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
			if text == IDToName(item) {
				sep := i.items.Objects[n+1]
				i.items.Remove(obj)
				i.items.Remove(sep)
				i.objects = Pop(item, i.objects)
				i.accessible = append(i.accessible, item)
				fyne.Do(func() {
					if i.pop != nil {
						i.pop.Hide()
					}
				})
				return
			}
		}
	}
}

func (i *Inventory) remove(item string) {
	for n, obj := range i.items.Objects {
		if cont, ok := obj.(*fyne.Container); ok {
			text := cont.Objects[0].(*fyne.Container).Objects[0].(*widget.Label).Text
			if text == IDToName(item) {
				sep := i.items.Objects[n+1]
				i.items.Remove(obj)
				i.items.Remove(sep)
				i.objects = Pop(item, i.objects)
				return
			}
		}
	}
}

func (i *Inventory) take(item string) {
	// img
	item_file := strings.TrimPrefix(item, "item.")
	image := canvas.NewImageFromFile("data/icons/placeholder.png")
	if _, err := os.Stat("data/items/" + item_file + ".png"); err == nil {
		image = canvas.NewImageFromFile("data/items/" + item_file + ".png")
	}

	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(50, 50))

	// cont
	single_cont := container.NewBorder(nil, nil, image, nil, container.NewCenter(widget.NewLabel(IDToName(item))))
	i.items.Add(single_cont)
	i.items.Add(widget.NewSeparator())
	i.objects = append(i.objects, item)
	if slices.Contains(i.accessible, item) {
		i.accessible = Pop(item, i.accessible)
	}
	fyne.Do(func() {
		if i.pop != nil {
			i.pop.Hide()
		}
	})
}

func (i *Inventory) Update(str string) {
	obj := []string{}
	if strings.HasPrefix(str, "inventory=") {
		str = strings.TrimPrefix(str, "inventory=")
		str = strings.Trim(str, "[]")
		items := strings.Split(str, ", ")

		if items[0] != "" {
			i.objects = items
		} else if items[0] == "" {
			i.objects = obj
		}
	} else if strings.HasPrefix(str, "look=") {
		str = strings.TrimPrefix(str, "look=")
		var data map[string]any
		if err := json.Unmarshal([]byte(str), &data); err != nil {
			fmt.Println("erreur de parsing:", err)
			return
		}
		var itemsID []string
		if room, ok := data["room"].(map[string]any); ok {
			if rawItems, ok := room["Items"].([]any); ok {
				for _, item := range rawItems {
					if str, ok := item.(string); ok {
						itemsID = append(itemsID, str)
					}
				}
			}
		}
		i.accessible = itemsID
	}
}

// Clear empties the inventory (used on death). The accessible items are
// refreshed by the next LOOK, once the player has respawned
func (i *Inventory) Clear() {
	i.objects = nil
	i.accessible = nil
	i.items.RemoveAll()
	if i.pop != nil {
		i.pop.Hide()
	}
}
