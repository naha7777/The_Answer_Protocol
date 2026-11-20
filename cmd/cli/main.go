package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"the_answer_protocol/client"
	"golang.org/x/term"
)

func main() {
	c, err := client.Dial("localhost:8080")
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERR 900 CONNECTION_FAILED")
		os.Exit(1)
	}
	defer c.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer term.Restore(fd, oldState)

	in := &Input{}
	done := make(chan struct{})

	go func() {
		defer close(done)
		for msg := range c.Messages() {
			in.Println(Colorize(msg))
		}
	}()

	go func() {
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

	<-done
	term.Restore(fd, oldState)
	fmt.Println("\nCLI disconnected")
}
