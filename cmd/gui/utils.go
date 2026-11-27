package main

// /===========================================\
// |				  Import				   |
// \===========================================/

import (
	"image/color"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// /===========================================\
// |				 Function				   |
// \===========================================/

// Put popup that disappear after 3 seconds
func Popup(w fyne.Window, txt string) {
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

// Put gameoverpopup that disappear after 3 seconds
func PopupGameOver(w fyne.Window, onDone func()) {
	title := canvas.NewText("GAME OVER", color.NRGBA{R: 220, G: 20, B: 20, A: 255})
	title.TextSize = 72
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	subtitle := canvas.NewText("You are dead.... BYE!", color.NRGBA{R: 230, G: 230, B: 230, A: 255})
	subtitle.TextSize = 20
	subtitle.Alignment = fyne.TextAlignCenter

	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(450, 20))

	content := container.NewPadded(
		container.NewVBox(
			spacer,
			container.NewCenter(title),
			container.NewCenter(subtitle),
			spacer,
		),
	)

	pop := widget.NewModalPopUp(content, w.Canvas())
	pop.Show()

	time.AfterFunc(3*time.Second, func() {
		fyne.Do(func() {
			pop.Hide()
			if onDone != nil {
				onDone()
			}
		})
	})
}

// Put popup that disappear after duration seconds
func PopupWithDuration(w fyne.Window, txt string, duration time.Duration) {
	if duration <= 0 {
		duration = 1
	}

	width := w.Canvas().Size().Width
	height := w.Canvas().Size().Height

	popup := widget.NewPopUp(widget.NewLabel(txt), w.Canvas())
	popup.ShowAtPosition(fyne.Position{X: width / 2, Y: 3 * (height / 4)})

	time.AfterFunc(duration*time.Second, func() {
		fyne.Do(func() {
			popup.Hide()
		})
	})
}

// Return time as 24h:min format
func get_time() string {
	formatted := time.Now().Format("15:04")
	return formatted
}

// Pop a string from a string slice
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

// Generate a container with all the strings in options with appropriate button
func Generate_option(options []string, action string, d *Data) *fyne.Container {
	cont := container.NewVBox()
	for _, option := range options {
		label := option
		if action == "TAKE" || action == "DROP" || action == "TALK" || action == "QUEST" || action == "ATTACK" || action == "USE_ITEM" {
			label = IDToName(option)
		}
		but := widget.NewButton(label, func() { d.conn.Send(action + " " + option) })
		if action == "USE_ITEM" {
			but = widget.NewButton(label, func() { d.conn.Send(action + " " + label); d.g.inventory.remove(option) })
		}

		if res := loadIcon(option); res != nil {
			but.SetIcon(res)
		}
		cont.Add(but)
	}
	return cont
}

// Load the asset from the id
func loadIcon(id string) fyne.Resource {
	var path string
	switch {
	case strings.HasPrefix(id, "item."):
		path = "data/items/" + strings.TrimPrefix(id, "item.") + ".png"
	case strings.HasPrefix(id, "npc."):
		path = "data/npc/" + strings.TrimPrefix(id, "npc.") + ".png"
	default:
		return nil
	}
	res, err := fyne.LoadResourceFromPath(path)
	if err != nil {
		return nil
	}
	return res
}

// Put popup at item position displaying all objects
func Option_popup(item fyne.CanvasObject, w fyne.Window, objects ...fyne.CanvasObject) *widget.PopUp {
	var pop *widget.PopUp

	var hookButtons func(objs []fyne.CanvasObject)
	hookButtons = func(objs []fyne.CanvasObject) {
		for _, o := range objs {
			switch v := o.(type) {
			case *widget.Button:
				action := v.OnTapped
				v.OnTapped = func() {
					pop.Hide()
					if action != nil {
						action()
					}
				}
			case *fyne.Container:
				hookButtons(v.Objects)
			}
		}
	}
	hookButtons(objects)

	menu := container.NewVBox(
		objects...,
	)
	pop = widget.NewPopUp(menu, w.Canvas())
	pop.ShowAtRelativePosition(fyne.NewPos(0, -pop.MinSize().Height), item)
	return pop
}

// Convert id string to its name
func IDToName(id string) string {

	id = strings.TrimPrefix(id, "item.")
	id = strings.TrimPrefix(id, "npc.")
	id = strings.TrimPrefix(id, "loc.")

	id = strings.Title(strings.ToLower(id))
	id = strings.ReplaceAll(id, "_", " ")
	return id
}

// Put confirm box
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
