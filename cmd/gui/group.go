package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

)
type GroupInfo struct {
	id		string
	leader	string
	members	[]string
	invited	[]string
}

type groupInfosMsg struct {
	Group	GroupInfo	`json:"group"`
}

func confirmGroup(w fyne.Window, group []string, onYes func()) {
	var groupName, leaderName string
	if len(group) == 2 {
		groupName, leaderName = group[0], group[1]
	}
	d := dialog.NewConfirm(
		"GROUP",
		fmt.Sprintf("%s invit you in the %s group.\nDo you really want to integrate ?", leaderName, groupName),
		func(confirmed bool) {
			if confirmed {
				onYes()
			}
		},
		w,
	)
	d.SetConfirmText("join")
	d.SetDismissText("no")
	d.Show()
}

func askText(w fyne.Window, title, label string, onSubmit func(text string)) {
	entry := widget.NewEntry()
	d := dialog.NewForm(title, "OK", "CANCEL",
		[]*widget.FormItem{widget.NewFormItem(label, entry)},
	func(ok bool) {
		if ok {
			onSubmit(strings.TrimSpace(entry.Text))
		}
	}, w)
	d.Resize(fyne.NewSize(400, d.MinSize().Height))
	d.Show()
	w.Canvas().Focus(entry)
}
