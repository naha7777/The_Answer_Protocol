package main

// /===========================================\
// |				  Import				   |
// \===========================================/

import (
	"strings"
	"the_answer_protocol/client"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// /===========================================\
// |				 Function				   |
// \===========================================/

func popup(w fyne.Window, txt string) {
	width := w.Canvas().Size().Width
	height := w.Canvas().Size().Height

	popup := widget.NewPopUp(widget.NewLabel(txt), w.Canvas())
	popup.ShowAtPosition(fyne.Position{X: width / 2, Y: 3 * (height / 4)})

	time.AfterFunc(3*time.Second, func() {
		fyne.Do(func() {
			popup.Hide()
		})
	})
}

func get_time() string {
	formatted := time.Now().Format("15:04")
	return formatted
}

func Pop(name string, place []string) []string {
	var new []string
	find := false
	for _, index := range place {
		if index != name && !find {
			new = append(new, index)
		} else if find {
			new = append(new, index)
		} else {
			find = true
		}
	}
	return new
}

func Generate_option(options []string, action string, conn *client.Conn) *fyne.Container {
	cont := container.NewVBox()
	for _, option := range options {
		var label string
		if action == "TAKE" || action == "DROP" || action == "TALK" || action == "QUEST" || action == "ATTACK" {
			label = IDToName(option)
		} else {
			label = option
		}
		but := widget.NewButton(label, func() { conn.Send(action + " " + option) })
		cont.Add(but)
	}
	return cont
}

func IDToName(id string) string {

	id = strings.TrimPrefix(id, "item.")
	id = strings.TrimPrefix(id, "npc.")
	id = strings.TrimPrefix(id, "loc.")

	id = strings.Title(strings.ToLower(id))
	id = strings.ReplaceAll(id, "_", " ")
	return id
}

func confirmAction(w fyne.Window, title, message string, onYes func()) {
	d := dialog.NewConfirm(title, message, func(confirmed bool) {
		if confirmed {
			onYes()
		}
	}, w)
	d.SetConfirmText("YES")
	d.SetDismissText("NO")
	d.Show()
}
