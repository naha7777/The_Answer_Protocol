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

// LoginScreen builds the initial authentication and character creation screen.
// It first prompts the user to enter a player name and connect to the server,
// then switches to a character gender selection step before invoking onSuccess.
// It returns the root CanvasObject view and a message handler callback for server responses.
func LoginScreen(conn *client.Conn, player *Player, onSuccess func(name string)) (fyne.CanvasObject, func(client.Message)) {
	contentBox := container.NewVBox()

	nameBanner := canvas.NewText("ENTER YOUR NAME", color.White)
	nameBanner.TextSize = 40
	nameBanner.TextStyle = fyne.TextStyle{Bold: true}
	nameBanner.Alignment = fyne.TextAlignCenter

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
			errorLabel.SetText("ERR 901 SEND_FAILED")
		}
		pendingName = name
	}
	nameEntry.OnSubmitted = func(string) { connect() } // touche entree
	butConnect := widget.NewButton("CONNECT", connect)

	nameForm := container.NewVBox(nameBanner, nameEntry, butConnect, errorLabel)

	genderBanner := canvas.NewText("CHOOSE YOUR CHARACTER", color.White)
	genderBanner.TextSize = 40
	genderBanner.TextStyle = fyne.TextStyle{Bold: true}
	genderBanner.Alignment = fyne.TextAlignCenter

	selectGender := func(gender string) {
		player.avatar = "./data/npc/" + gender + ".png"
		player.gender = gender
		onSuccess(pendingName)
	}

	girlChoice := container.NewVBox(widget.NewButton("GIRL", func() { selectGender("girl") }))
	boyChoice := container.NewVBox(widget.NewButton("BOY", func() { selectGender("boy") }))

	genderForm := container.NewVBox(
		genderBanner,
		container.NewCenter(container.NewHBox(girlChoice, boyChoice)),
	)

	contentBox.Add(nameForm)

	bg := canvas.NewImageFromFile("./data/welcome.jpeg")
	bg.FillMode = canvas.ImageFillCover
	overlay := canvas.NewRectangle(color.NRGBA{A: 150})

	view := container.NewStack(bg, overlay, container.NewCenter(contentBox))

	onMessage := func(msg client.Message) {
		switch {
		case msg.Kind == "OK" && msg.Rest == "connected":
			player.name = pendingName
			contentBox.Objects = []fyne.CanvasObject{genderForm}
			contentBox.Refresh()
		case msg.Kind == "ERR":
			e := strings.Split(msg.Raw, " ")
			switch e[1] {
			case "021":
				errorLabel.SetText("This name is too long")
			case "019":
				errorLabel.SetText("The server is full try again later")
			case "009":
				errorLabel.SetText("Invalid name")
			case "010":
				errorLabel.SetText("Already connected")
			case "002":
				errorLabel.SetText("You must be connected")
			case "201":
				errorLabel.SetText("This name is taken")
			default:
				errorLabel.SetText(msg.Rest)
			}
		}
	}
	return view, onMessage
}
