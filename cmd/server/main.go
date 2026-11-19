package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"the_answer_protocol/server"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listening : %s", err)
		os.Exit(1)
	}

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			cmd := strings.TrimSpace(strings.ToLower(scanner.Text()))
			switch cmd {
			case "stop", "quit", "exit":
				fmt.Println("Server shutdown request")
				if err := listener.Close(); err != nil {
					fmt.Fprintf(os.Stderr, "Listener closure error: %s\n", err)
				}
				return
			case "":
				continue
			default:
				fmt.Printf("Unknown command: %q (use 'quit', 'stop' or 'exit').\n", cmd)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %s\n", err)
		}
	}()
	server.Server(listener)
	fmt.Println("Shut down the server properly")
}
