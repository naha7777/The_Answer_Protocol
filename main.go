package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/canvas"
)

func get_time() string {
	formatted := time.Now().Format("02/01/06 03:04:05")
	return formatted
}

func makeUI() (*widget.Label, *widget.Button, *widget.Entry) {
	return widget.NewLabel(get_time()),
		widget.NewButton("ATTACK", func() {
			fmt.Println("yo")
		}),
		widget.NewEntry()
}

// /===========================================\
// |				   Layout				   |
// \===========================================/

// struct for HUD  
// /!\ Need at least 4 widget 
type Hud struct {

}

// Need to implement the Layout Interface with MinSize and Layout
func (ui *Hud) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) < 4 { fmt.Println("Error: Hud needs 4 arguments") }
	w, h := float32(0), float32(0)
	for _, o := range objects {
		childSize := o.MinSize()

		w += childSize.Width
		h += childSize.Height
	}
	return fyne.NewSize(w, h)
}

func (ui *Hud) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) < 4 { return }
	w, h := containerSize.Width, containerSize.Height

	for _, o := range objects {
		o.Resize(o.MinSize())
	}

	chat_size := fyne.Size{Width: containerSize.Width / 3, Height: containerSize.Height / 5}
	objects[3].Resize(chat_size)

	objects[0].Move(fyne.NewPos(0, 0))
	objects[1].Move(fyne.NewPos(w - objects[1].MinSize().Width, 0))
	objects[2].Move(fyne.NewPos(0, h - objects[2].MinSize().Height))
	objects[3].Move(fyne.NewPos(w - chat_size.Width - 10, h - chat_size.Height))
}


// /===========================================\
// |				   Display				   |
// \===========================================/

func main() {
	a := app.New()
	w := a.NewWindow("TAP")

	// /======================================\
	// |	             Chat	  			  |
	// \======================================/
	
	boxMessages := container.NewVBox()
	scrollMessages := container.NewScroll(boxMessages)
	
	// Key hook for esc to close the app
	w.Canvas().SetOnTypedKey(func(keyEvent *fyne.KeyEvent) {
		if keyEvent.Name == fyne.KeyEscape {
			a.Quit()
		}
	})

	// all widget
	txt1 := widget.NewLabel("topl")
	chat := widget.NewEntry()
	chat.SetPlaceHolder("Global chat...")

	send_msg := func() {
		if chat.Text != "" {
			boxMessages.Add(widget.NewLabel("Me: " + chat.Text))
			chat.SetText("")
			scrollMessages.ScrollToBottom()
		}
	}

	but_say := widget.NewButton("SAY", send_msg)
	send_box := container.NewBorder(nil, nil, nil, but_say, chat)
	chat_box := container.NewBorder(nil, send_box, nil, nil, scrollMessages)
	chat.OnSubmitted = func(text string) {
		send_msg()
	}
	
	
	img := canvas.NewImageFromFile("tavern.png")
	img.Resize(fyne.Size{Width: 200, Height: 200})
	img.FillMode = canvas.ImageFillOriginal
	but1 := widget.NewButton("ATTACK", func() {fmt.Println("boom")})
	but2 := widget.NewButton("MOVE", func() {fmt.Println("boom")})
	
	buttonsContainer := container.NewHBox(but1, but2)

	w.SetContent(container.New(&Hud{}, txt1, img, buttonsContainer, chat_box))
	
	w.Resize(fyne.NewSize(1920, 1080))
	w.ShowAndRun()
	fmt.Println("Exited, end of TAP")
}