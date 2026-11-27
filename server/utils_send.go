package server

import (
	"fmt"
)

func SendAll(d *Dispatcher, sender *Client, format string, a ...any) {
	// Send an EVT to everyone except the client
	var targets []*Client
	for _, c := range d.clients {
		if c != sender {
			targets = append(targets, c)
		}
	}
	broadcast(targets, format, a...)
}

func SendRoom(d *Dispatcher, sender *Client, format string, a ...any) {
	// Send an EVT to everyone in the client's room except the client
	var targets []*Client
	for _, c := range d.clients {
		if c != sender && c.room == sender.room {
			targets = append(targets, c)
		}
	}
	broadcast(targets, format, a...)
}

func SendGroup(d *Dispatcher, sender *Client, format string, a ...any) {
	// Send an EVT to everyone in the client's group except the client
	var targets []*Client
	for _, c := range d.clients {
		if c != sender && sender.group != nil && c.group == sender.group {
			targets = append(targets, c)
		}
	}
	broadcast(targets, format, a...)
}

func broadcast(targets []*Client, format string, a ...any) {
	// Display the message to the targets
	msg := fmt.Sprintf(format, a...)
	recipients := make([]string, 0, len(targets))

	for _, t := range targets {
		if err := t.write(msg); err != nil {
			continue
		}
		recipients = append(recipients, t.label())
	}

	if len(recipients) > 0 {
		Srv().Info("event sent", "event", msg, "recipients", recipients,
			"count", len(recipients))
	}
}
