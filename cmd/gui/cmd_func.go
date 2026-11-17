package main

// /===========================================\
// |				  Import				   |
// \===========================================/

import (
	"encoding/json"
	"fmt"
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

// QUESTS
func Str_to_quests(str string, w fyne.Window) {
	width := w.Canvas().Size().Width
	height := w.Canvas().Size().Height

	// parse
	raw := `[{"quest_id": "fetch_herbs", "status": "active", "progress": "1/3"}, {"quest_id": "defeat_goblin", "status": "completed"}]`

	var quests []map[string]string
	if err := json.Unmarshal([]byte(raw), &quests); err != nil {
		fmt.Println("erreur de parsing:", err)
		return
	}

	// cont
	quests_cont := container.NewVBox()
	quests_scroll := container.NewScroll(quests_cont)
	q_cont := container.NewBorder(container.NewCenter(widget.NewLabel("Quests")), nil, nil, nil, quests_scroll)

	for _, q := range quests {

		quest := q["quest_id"]
		quest = strings.ReplaceAll(quest, "_", " ")
		quest = strings.Title(strings.ToLower(quest))
		single_cont := container.NewBorder(nil, nil, widget.NewLabel(quest), nil, container.NewCenter(widget.NewLabel(q["progress"])))
		// img
		if q["status"] == "completed" {
			image := canvas.NewImageFromFile("pictures/items/placeholder.png")
			if _, err := os.Stat("pictures/icons/tick.png"); err == nil {
				image = canvas.NewImageFromFile("pictures/icons/tick.png")
			}
			image.FillMode = canvas.ImageFillContain
			image.SetMinSize(fyne.NewSize(30, 30))
			single_cont = container.NewBorder(nil, nil, container.NewCenter(widget.NewLabel(quest)), nil, image)
		}
		quests_cont.Add(single_cont)
		quests_cont.Add(widget.NewSeparator())
	}

	quests_frame := canvas.NewRectangle(color.NRGBA{R: 20, G: 20, B: 20, A: 200})
	quests_frame.StrokeColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	quests_frame.StrokeWidth = 2
	q := container.NewStack(quests_frame, container.NewPadded(q_cont))

	popup := widget.NewPopUp(q, w.Canvas())
	popup.ShowAtPosition(fyne.Position{X: 7*(width/9) - quests_cont.MinSize().Width*2, Y: height - 50*4 - 15})
	popup.Resize(fyne.Size{Width: quests_cont.MinSize().Width * 2, Height: 50 * 4})
}

func Option_popup(item fyne.CanvasObject, w fyne.Window, objects ...fyne.CanvasObject) {
	menu := container.NewVBox(
		objects...,
	)
	pop := widget.NewPopUp(menu, w.Canvas())
	pop.ShowAtRelativePosition(fyne.NewPos(0, -pop.MinSize().Height), item)
}
