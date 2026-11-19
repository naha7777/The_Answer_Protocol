package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"the_answer_protocol/server"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listening : %s", err)
		os.Exit(1)
	}

	var once sync.Once
	shutdown := func() {
		once.Do(func() {
			if err := listener.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Listener closure error: %s\n", err)
			}
		})
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nSignal received, shutting down...")
		shutdown()
	}()

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
