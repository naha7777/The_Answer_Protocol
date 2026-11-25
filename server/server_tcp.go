package server

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"unicode"
	"unicode/utf8"
	"strings"
	"runtime/debug"
	"sync/atomic"
	"the_answer_protocol/game"
)

var conn_nb atomic.Int32
const connection_limit = 30

var (
	wg      sync.WaitGroup
	connsMu sync.Mutex
	conns   = map[net.Conn]struct{}{}
)

func trackConn(c net.Conn)   { connsMu.Lock(); conns[c] = struct{}{}; connsMu.Unlock() }
func untrackConn(c net.Conn) { connsMu.Lock(); delete(conns, c); connsMu.Unlock() }

func closeAllConns() {
	connsMu.Lock()
	defer connsMu.Unlock()
	for c := range conns {
		c.Close()
	}
}

func hasControlChars(s string) bool {
	for _, r := range s {
		if unicode.IsControl(r) {
			return true
		}
	}
	return false
}

func HandleConnection(connection net.Conn, dispatcher *Dispatcher,
					  world *game.World) {

	defer connection.Close()
	ip := connection.RemoteAddr().String()

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
				status: "very healthy", strengh: 10, target: "",
				real_strengh: 10, quit: false}
		}
	}

	// limit the number of connections
	if conn_nb.Add(1) > connection_limit {
		Log.Warn("connection refused: server full", "ip", ip)
		client.Send("ERR 19 SERVER_FULL")
		conn_nb.Add(-1)
		return
	}
	defer conn_nb.Add(-1)

	// Server log for client connection
	Log.Info("client connected", "ip", ip)
	defer func() {
		Log.Info("client disconnected", "ip", ip)
		dispatcher.removeClient(client)
		dispatcher.displayRoomStats()
		dispatcher.displayStats()
	}()

	// avoid panic
	defer func() {
		if r := recover(); r != nil {
			Log.Error("panic recovered", "player", client.name, "ip", ip,
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
			client.Send("ERR 16 LINE_TOO_LONG")
			continue
		}

		// manipulation attacks prevention
		text := string(line)
		if !utf8.ValidString(text) || hasControlChars(text) {
			client.Send("ERR 20 INVALID_INPUT")
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
			client.Send("ERR 02 YOU_MUST_BE_CONNECTED")
		}
		if client.quit {
			break
		}
	}
}

func Server(listener net.Listener) {
	// world generation
	world, err := game.LoadWorld("data/world.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}

	for _, r := range world.Rooms {
		if r.NPCs == nil {
			r.NPCs = []string{}
		}
		if r.Items == nil {
			r.Items = []string{}
		}
		if r.Players == nil {
			r.Players = []string{}
		}
	}

	// dispatcher creation
	dispatcher := NewDispatcher(world)

	Log.Info("server started", "addr", listener.Addr().String())

	for {
		connection, err := listener.Accept()
		if err != nil {
			// verrification vu que listener.Close est appele listenner.Accept
			// retourne une net.ErrClosed donc on a ferme le listenner volontairement
			if errors.Is(err, net.ErrClosed) {
				Log.Info("server shutting down")
				closeAllConns()
				wg.Wait()
				Log.Info("server stopped")
				return
			}
			Log.Error("accept failed", "err", err)
			continue
		}
		// one goroutine for one connection
		trackConn(connection)
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer untrackConn(connection)
			HandleConnection(connection, dispatcher, world)
		}()
	}
}
