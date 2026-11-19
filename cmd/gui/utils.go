package main

// /===========================================\
// |				  Import				   |
// \===========================================/

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"the_answer_protocol/client"
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
	for _, item := range place {
		if item != name && !find {
			new = append(new, item)
		} else if find {
			new = append(new, item)
		} else {
			find = true
		}
	}
	return new
}

func dialWithRetry(addr string) (*client.Conn, error) {
	deadline := time.Now().Add(30*time.Second)
	var lastErr error

	for {
		conn, err := client.Dial(addr)
		if err == nil {
			return conn, nil
		}
		lastErr = err

		if time.Now().After(deadline) {
			break
		}
		fmt.Println("attempt to reconnect")
		time.Sleep(3*time.Second)
	}
	return nil, lastErr
}