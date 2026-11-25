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
func Str_to_quests(str string, d Data) {
	// width := w.Canvas().Size().Width
	height := d.w.Canvas().Size().Height

	// parse
	var quests []map[string]string
	if err := json.Unmarshal([]byte(str), &quests); err != nil {
		fmt.Println("erreur de parsing:", err)
		return
	}

	// cont
	quests_cont := container.NewVBox()
	quests_scroll := container.NewScroll(quests_cont)
	title := container.NewVBox(container.NewCenter(widget.NewLabel("Quests")), widget.NewSeparator())
	q_cont := container.NewBorder(title, nil, nil, nil, quests_scroll)

	for _, q := range quests {

		quest := q["quest_id"]
		quest = strings.ReplaceAll(quest, "_", " ")
		quest = strings.Title(strings.ToLower(quest))
		completeBut := widget.NewButton("Complete", func() { d.conn.Send("COMPLETE_QUEST " + q["quest_id"]) })
		label_cont := container.NewHBox(widget.NewLabel(quest), widget.NewLabel(IDToName(q["quest_npc"])))
		single_cont := container.NewBorder(nil, nil, label_cont, nil, container.NewHBox(container.NewCenter(widget.NewLabel(q["progress"])), completeBut))
		// img
		if q["status"] == "completed" {
			image := canvas.NewImageFromFile("data/icons/placeholder.png")
			if _, err := os.Stat("data/icons/tick.png"); err == nil {
				image = canvas.NewImageFromFile("data/icons/tick.png")
			}
			image.FillMode = canvas.ImageFillContain
			image.SetMinSize(fyne.NewSize(30, 30))
			single_cont = container.NewBorder(nil, nil, label_cont, nil, image)
		}
		quests_cont.Add(single_cont)
		quests_cont.Add(widget.NewSeparator())
	}

	quests_frame := canvas.NewRectangle(color.NRGBA{R: 20, G: 20, B: 20, A: 200})
	quests_frame.StrokeColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	quests_frame.StrokeWidth = 2
	q := container.NewStack(quests_frame, container.NewPadded(q_cont))

	popup := widget.NewPopUp(q, d.w.Canvas())
	popup.ShowAtPosition(fyne.Position{X: 15, Y: height - 50*8 - 30})
	if quests_cont.MinSize().Width < 200 {
		popup.Resize(fyne.Size{Width: 200, Height: 50 * 4})
	} else {
		popup.Resize(fyne.Size{Width: quests_cont.MinSize().Width + 30, Height: 50 * 4})
	}
}
