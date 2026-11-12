package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type Game struct {
	window		fyne.Window
	rooms		map[string]room
	current		string
	background	*canvas.Image
	showMap		bool
}

func(g *Game) ToggleMap() {
	g.showMap = !g.showMap
	g.updateBackgroung()
}

func NewGame(w fyne.Window, bg *canvas.Image) *Game {
	g := &Game {
		window:		w,
		rooms:		rooms(),
		current:	"square",
		background:	bg,
	}
	g.updateBackgroung()
	return g
}

func (g *Game) Move(dir string) {
	if g.showMap {
		return
	}

	r := g.rooms[g.current]

	next := ""
	switch dir {
	case "north":
		next = r.north
	case "south":
		next = r.south
	case "east":
		next = r.east
	case "west":
		next = r.west
	}

	if next == "" {
		popup(g.window, "You can't go that way")
		return
	}

	if _, ok := g.rooms[next]; !ok {
		fmt.Println("Unknown room : ", next)
		return
	}

	g.current = next
	g.updateBackgroung()
}

func (g *Game) updateBackgroung() {
	if g.showMap {
		g.background.File = g.rooms["map"].image
	} else {
		g.background.File = g.rooms[g.current].image
	}
	g.background.Refresh()
}
