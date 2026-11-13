package main

import (
	"image/color"
	"strings"
	"the_answer_protocol/client"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func LoginScreen(conn *client.Conn, onSuccess func(name string)) (fyne.CanvasObject, func(client.Message)) {
	banner := canvas.NewText("ENTER YOUR NAME", color.White)
	banner.TextSize = 40
	banner.TextStyle = fyne.TextStyle{Bold: true}
	banner.Alignment = fyne.TextAlignCenter

	nameEntry := widget.NewEntry()

	errorLabel := widget.NewLabel("")
	errorLabel.TextStyle = fyne.TextStyle{Bold: true}

	var pendingName string

	connect := func() {
		name := strings.TrimSpace(nameEntry.Text)
		if name == "" {
			errorLabel.SetText("Please enter a name")
			return
		}
		errorLabel.SetText("")
		if err := conn.Send("CONNECT " + name); err != nil {
			errorLabel.SetText("901 SEND_FAILED")
		}
		pendingName = name
	}
	nameEntry.OnSubmitted = func(string) { connect() } // touche entree
	butConnect := widget.NewButton("CONNECT", connect)

	form := container.NewVBox(banner, nameEntry, butConnect, errorLabel)

	bg := canvas.NewImageFromFile("./data/welcome.jpeg")
	bg.FillMode = canvas.ImageFillStretch
	overlay := canvas.NewRectangle(color.NRGBA{A: 150})

	view := container.NewStack(bg, overlay, container.NewCenter(form))

	onMessage := func(msg client.Message) {
		switch {
		case msg.Kind == "OK" && msg.Rest == "connected":
			onSuccess(pendingName)
		case msg.Kind == "ERR":
			errorLabel.SetText(msg.Rest)
		}
	}
	return view, onMessage
}
