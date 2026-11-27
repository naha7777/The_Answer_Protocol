package server

import (
	"bufio"
	"fmt"
	"net"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"the_answer_protocol/game"
	"time"
	"unicode"
	"unicode/utf8"
)

var conn_nb atomic.Int32

func hasControlChars(s string) bool {
	// identify control chars to protect the program against attacks
	for _, r := range s {
		if unicode.IsControl(r) {
			return true
		}
	}
	return false
}

func (t *connTracker) count(host string, window time.Duration) int {
	// quick connections protection
	t.mu.Lock()
	defer t.mu.Unlock()
	now := time.Now()

	kept := t.seen[host][:0]
	for _, ts := range t.seen[host] {
		if now.Sub(ts) < window {
			kept = append(kept, ts)
		}
	}
	kept = append(kept, now)
	t.seen[host] = kept

	if len(t.seen) > 1000 {
		for h, l := range t.seen {
			if now.Sub(l[len(l)-1]) > window {
				delete(t.seen, h)
			}
		}
	}
	return len(kept)
}

func HandleConnection(connection net.Conn, dispatcher *Dispatcher,
	world *game.World) {

	defer connection.Close()
	ip := connection.RemoteAddr().String()

	host, _, _ := net.SplitHostPort(ip)
	if n := tracker.count(host, 10*time.Second); n == 6 {
		Srv().Warn("abuse: rapid connections", "ip", host, "limit", "5 per 10s")
	}

	// client creation
	client := &Client{conn: connection}
	for _, room := range world.Rooms {
		if room.Id == world.Start {
			// item := &game.Item{Name: "Herbs", Description: "de l'herbe",
			// 	Obtainable: true, Id: "item.herbs"}
			// var inv []*game.Item
			// inv = append(inv, item)
			client = &Client{conn: connection, name: "", group: nil, room: room,
				inventory: nil, hp: 100, hp_max: 100,
				status: "exploration", strength: 10, target: "",
				real_strength: 10, quit: false}
		}
	}

	// limit the number of connections
	if conn_nb.Add(1) > connection_limit {
		Srv().Warn("connection refused: server full", "ip", ip)
		client.Send("ERR 019 SERVER_FULL")
		conn_nb.Add(-1)
		return
	}
	defer conn_nb.Add(-1)

	// Server log for client connection
	Srv().Info("client connected", "ip", ip)
	defer func() {
		Srv().Info("client disconnected", "ip", ip, "player", client.name)
		client.dead.Store(true)
		dispatcher.removeClient(client)
		dispatcher.displayRoomStats()
		dispatcher.displayStats()
	}()

	// avoid panic
	defer func() {
		if r := recover(); r != nil {
			Srv().Error("panic recovered", "player", client.name, "ip", ip,
				"panic", fmt.Sprint(r), "stack", string(debug.Stack()))
			client.Send("ERR 500 INTERNAL")
		}
	}()
	client.Send("OK hello proto=1")

	reader := bufio.NewReaderSize(connection, 1024)
	// read input
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			break
		}
		if isPrefix {
			for isPrefix && err == nil {
				_, isPrefix, err = reader.ReadLine()
			}
			if err != nil {
				break
			}
			client.log().Warn("abuse: line too long")
			client.Send("ERR 016 LINE_TOO_LONG")
			continue
		}

		// manipulation attacks prevention
		text := string(line)
		if !utf8.ValidString(text) || hasControlChars(text) {
			client.Send("ERR 020 INVALID_INPUT")
			continue
		}
		upper := strings.ToUpper(text)

		// the first command must be CONNECT with a valid name
		if client.name == "" && strings.HasPrefix(upper, "CONNECT") {
			dispatcher.Dispatch(client, text)
			if client.name != "" {
				dispatcher.displayStats()
				dispatcher.displayRoomStats()
			}
			// then, all commands are available
		} else if client.name != "" {
			dispatcher.Dispatch(client, text)
		} else {
			client.Send("ERR 002 YOU_MUST_BE_CONNECTED")
		}
		if client.quit {
			break
		}
	}
}
