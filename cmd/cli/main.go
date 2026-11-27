package main

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"os"
	"strings"
	"the_answer_protocol/client"
)

func main() {
	// connect to the server
	c, err := client.Dial()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERR 900 CONNECTION_FAILED")
		os.Exit(1)
	}
	// close connection when we are out of the main
	defer c.Close()

	fd := int(os.Stdin.Fd())
	// Put the terminal in mode raw. oldState reput the terminal in a normal
	// mode.
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer term.Restore(fd, oldState)

	in := &Input{}
	// tell if the session is over
	done := make(chan struct{})

	// First goroutine : colorize server messages
	go func() {
		// when this goroutine is finished, main is unlock
		defer close(done)
		// c.Messages is a channel : the loop stop when it's close
		for msg := range c.Messages() {
			// delete input, print the message
			in.Println(Colorize(msg))
		}
	}()

	// Second goroutine : read the input
	go func() {
		// if ReadLine return an error : close connection
		defer c.Close()
		r := bufio.NewReader(os.Stdin)
		in.mu.Lock()
		in.redraw()
		in.mu.Unlock()
		for {
			line, err := in.ReadLine(r)
			if err != nil {
				return
			}
			if strings.TrimSpace(line) == "" {
				continue
			}
			if err := c.Send(line); err != nil {
				in.Println("901 SEND_FAILED")
				return
			}
		}
	}()

	// main is waiting for the session end
	<-done
	if err := c.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "connection error:", err)
	}
	// restore the terminal before writing
	term.Restore(fd, oldState)
	fmt.Println("\nCLI disconnected")
}
