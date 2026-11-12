package main

import (
	"the_answer_protocol/client"
	"fmt"
	"os"
	"bufio"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Cyan = "\033[36m"

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
				fmt.Println(Green + msg.Raw + Reset)
			case "ERR":
				fmt.Println(Red + msg.Raw + Reset)
			case "OK":
				fmt.Println(Cyan + msg.Raw + Reset)
			default:
				fmt.Println(msg.Raw)
			}
		}
		fmt.Println("Server disconnected")
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
