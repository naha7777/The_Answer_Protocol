package server

import (
	"fmt"
)

func SendAll(d *Dispatcher, sender *Client, format string, a ...any) {
	var targets []*Client
	for _, c := range d.clients {
		if c != sender {
			targets = append(targets, c)
		}
	}
	broadcast(targets, format, a...)
}

func SendRoom(d *Dispatcher, sender *Client, format string, a ...any) {
	var targets []*Client
	for _, c := range d.clients {
		if c != sender && c.room == sender.room {
			targets = append(targets, c)
		}
	}
	broadcast(targets, format, a...)
}

func SendGroup(d *Dispatcher, sender *Client, format string, a ...any) {
	var targets []*Client
	for _, c := range d.clients {
		if c != sender && sender.group != nil && c.group == sender.group {
			targets = append(targets, c)
		}
	}
	broadcast(targets, format, a...)
}

func broadcast(targets []*Client, format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	recipients := make([]string, 0, len(targets))

	for _, t := range targets {
		if err := t.write(msg); err != nil {
			Srv().Error("send failed", "to", t.label(), "err", err)
			continue
		}
		recipients = append(recipients, t.label())
	}

	if len(recipients) > 0 {
		Srv().Info("event sent", "event", msg, "recipients", recipients,
				   "count", len(recipients))
	}
}
