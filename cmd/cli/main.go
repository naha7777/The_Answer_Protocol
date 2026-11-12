package main

import (
	"the_answer_protocol/client"
	"fmt"
	"os"
	"bufio"
)

func main() {
	c, err := client.Dial("localhost:8080")
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERR 900 CONNECTION_FAILED")
		os.Exit(1)
	}
	defer c.Close()

	go func() {
		for msg := range c.Messages() {
			fmt.Println(msg.Raw)
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
