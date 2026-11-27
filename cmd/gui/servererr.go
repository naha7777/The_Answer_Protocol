package main

// /===========================================\
// |				  Import				   |
// \===========================================/

import (
	"strings"
	"the_answer_protocol/client"
)

func ServerResponseERR(msg client.Message, player *Player, data *Data) {
	e := strings.Split(msg.Raw, " ")
	switch e[1] {
	case "201":
		Popup(data.w, "This nickname is already taken")
	case "301":
		Popup(data.w, "You can't go that way")
	case "401":
		Popup(data.w, "You are not in a group")
	case "402":
		Popup(data.w, "The player is already in group")
	case "404":
		Popup(data.w, "You can't do that")
	case "405":
		Popup(data.w, "Let's be friendly")
	case "900":
		Popup(data.w, "Failed to connect to the server")
	case "901":
		Popup(data.w, "Failed to reach the server")
	case "001":
		Popup(data.w, "Command didn't exist")
	case "002":
		Popup(data.w, "You must be connected")
	case "003":
		Popup(data.w, "Player not found")
	case "004":
		Popup(data.w, "You are not the leader")
	case "005":
		Popup(data.w, "You need an invitation")
	case "006":
		Popup(data.w, "To who are you talking?")
	case "007":
		Popup(data.w, "You can't get that")
	case "008":
		Popup(data.w, "Try typing something")
	case "009":
		Popup(data.w, "Name is invalid")
	case "010":
		Popup(data.w, "You are already connected")
	case "011":
		Popup(data.w, "You need to be near the npc with all the requirement")
	case "012":
		Popup(data.w, "Let's focus on the fight")
	case "013":
		Popup(data.w, "Why? It's chill right here")
	case "014":
		Popup(data.w, "I respect the 1v1")
	case "015":
		Popup(data.w, "Don't have no quest for you")
	case "016":
		Popup(data.w, "Line too long")
	case "017":
		Popup(data.w, "Inventory is full")
	case "018":
		Popup(data.w, "Stop the spam")
	case "019":
		Popup(data.w, "Server is full try again later")
	case "020":
		Popup(data.w, "Invalid input")
	case "021":
		Popup(data.w, "Name is too long")
	case "500":
		Popup(data.w, "Internal error")
	default:
		Popup(data.w, msg.Rest)
	}

}
