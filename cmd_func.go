package main

// /===========================================\
// |				  Import				   |
// \===========================================/

import (
	"image/color"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// /===========================================\
// |				 Function				   |
// \===========================================/

// INVENTORY
func Str_to_inv(str string, w fyne.Window) {
	width := w.Canvas().Size().Width
	height := w.Canvas().Size().Height

	// parse
	str = strings.Trim(str, "[]")
	items := strings.Split(str, ", ")

	// cont
	items_cont := container.NewVBox()
	items_scroll := container.NewScroll(items_cont)
	inv_cont := container.NewBorder(container.NewCenter(widget.NewLabel("Inventory")), nil, nil, nil, items_scroll)

	for _, item := range items {
		item, _ = strings.CutPrefix(item, "item.")
		
		// img
		image := canvas.NewImageFromFile("pictures/items/placeholder.png")
		if _, err := os.Stat("pictures/items/" + item + ".png"); err == nil {
			image = canvas.NewImageFromFile("pictures/items/" + item + ".png")
		}

		image.FillMode = canvas.ImageFillContain
		image.SetMinSize(fyne.NewSize(50, 50))

		item = strings.ReplaceAll(item, "_", " ")
		item = strings.Title(strings.ToLower(item))
		single_cont := container.NewBorder(nil, nil, image, nil, container.NewCenter(widget.NewLabel(item)))
		items_cont.Add(single_cont)
		items_cont.Add(widget.NewSeparator())
	}

	inv_frame := canvas.NewRectangle(color.NRGBA{R: 20, G: 20, B: 20, A: 200})
	inv_frame.StrokeColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	inv_frame.StrokeWidth = 2
	inv := container.NewStack(inv_frame, container.NewPadded(inv_cont))

	popup := widget.NewPopUp(inv, w.Canvas())
	popup.ShowAtPosition(fyne.Position{X: 7*(width/9) - items_cont.MinSize().Width*2, Y: height - 50*4 - 15})
	popup.Resize(fyne.Size{Width: items_cont.MinSize().Width * 2, Height: 50 * 4})
}
