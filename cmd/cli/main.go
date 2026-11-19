package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"the_answer_protocol/client"
)

var Reset	= "\033[0m"
var Red		= "\033[31m"
var Green	= "\033[32m"
var Cyan	= "\033[36m"
var Italic	= "\033[3m"
var Magenta	= "\033[35m"
var Yellow  = "\033[33m"

func main() {
	c, err := client.Dial("localhost:8080")
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERR 900 CONNECTION_FAILED")
		os.Exit(1)
	}
	defer c.Close()

	go func() {
		for msg := range c.Messages() {
			switch msg.Kind {
			case "EVT":
				if strings.HasPrefix(msg.Rest, "DEATH") || strings.HasPrefix(msg.Rest, "dead") {
					fmt.Println(Italic + Red + msg.Raw + Reset)
				} else if strings.HasPrefix(msg.Rest, "NPC") {
					fmt.Println(Yellow + msg.Raw + Reset)
				}	else {
					fmt.Println(Green + msg.Raw + Reset)
				}
			case "ERR":
				fmt.Println(Red + msg.Raw + Reset)
			case "OK":
				if strings.HasPrefix(msg.Rest, "attack=") || strings.HasPrefix(msg.Rest, "flee=") || strings.HasPrefix(msg.Rest, "defend") {
					fmt.Println(Italic + Magenta + msg.Raw + Reset)
				} else {
					fmt.Println(Cyan + msg.Raw + Reset)
				}
			default:
				fmt.Println(msg.Raw)
			}
		}
		fmt.Println("CLI disconnected")
		os.Exit(0)
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if err := c.Send(scanner.Text()); err != nil {
			fmt.Fprintln(os.Stderr, "901 SEND_FAILED")
			return
		}
	}
	if err := scanner.Err(); err != nil {
		c.Close()
	}
}
