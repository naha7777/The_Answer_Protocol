package main

import (
	"strings"
	"the_answer_protocol/client"
	"unicode"
)

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	RED     = "\033[41m"
	Green   = "\033[32m"
	Cyan    = "\033[36m"
	Italic  = "\033[3m"
	Magenta = "\033[35m"
	Yellow  = "\033[33m"
)

func sanitize(s string) string {
	// manipulation attacks prevention
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, s)
}

func Colorize(msg client.Message) string {
	msg.Raw = sanitize(msg.Raw)
	if msg.Raw == "GAME OVER" {
		return RED + msg.Raw + Reset
	}
	switch msg.Kind {
	case "EVT":
		if strings.HasPrefix(msg.Rest, "DEATH") ||
			strings.HasPrefix(msg.Rest, "dead") {
			return Italic + Red + msg.Raw + Reset
		} else if strings.HasPrefix(msg.Rest, "NPC") {
			return Yellow + msg.Raw + Reset
		}
		return Green + msg.Raw + Reset
	case "ERR":
		return Red + msg.Raw + Reset
	case "OK":
		if strings.HasPrefix(msg.Rest, "attack=") ||
			strings.HasPrefix(msg.Rest, "flee=") ||
			strings.HasPrefix(msg.Rest, "defend") {
			return Italic + Magenta + msg.Raw + Reset
		}
		return Cyan + msg.Raw + Reset
	}
	return msg.Raw
}
