package main

import (
	"bufio"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"the_answer_protocol/server"
)

func main() {
	// the server listens on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		server.Log.Error("listen failed", "err", err)
		os.Exit(1)
	}

	// handle listener error
	var once sync.Once
	shutdown := func(reason string) {
		once.Do(func() {
			server.Log.Info("shutdown requested", "reason", reason)
			if err := listener.Close(); err != nil {
				server.Log.Error("listener closure failed", "err", err)
			}
		})
	}

	// handle forcing closure (CRTL+C for example)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		shutdown("signal")
	}()

	// handle closure with 'quit', 'exit' or 'stop'
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			cmd := strings.TrimSpace(strings.ToLower(scanner.Text()))
			switch cmd {
			case "stop", "quit", "exit":
				shutdown("stdin command")
				return
			case "":
				continue
			default:
				server.Log.Warn("unknown server command", "cmd", cmd,
								"hint", "use quit, stop or exit")
			}
		}
		if err := scanner.Err(); err != nil {
			server.Log.Error("stdin read failed", "err", err)
		}
	}()
	server.Server(listener)
}
