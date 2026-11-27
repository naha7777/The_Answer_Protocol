package server

import (
	"errors"
	"net"
	"os"
	"sync"
	"the_answer_protocol/game"
	"time"
)

const connection_limit = 30

type connTracker struct {
	mu   sync.Mutex
	seen map[string][]time.Time
}

var tracker = &connTracker{seen: map[string][]time.Time{}}

var (
	wg      sync.WaitGroup
	connsMu sync.Mutex
	conns   = map[net.Conn]struct{}{}
)

func trackConn(c net.Conn) {
	// track the connection
	connsMu.Lock()
	conns[c] = struct{}{}
	connsMu.Unlock()
}

func untrackConn(c net.Conn) {
	// untrack the connection
	connsMu.Lock()
	delete(conns, c)
	connsMu.Unlock()
}

func closeAllConns() {
	// close all the connections
	connsMu.Lock()
	defer connsMu.Unlock()
	for c := range conns {
		c.Close()
	}
}

func checkWorld() *game.World {
	path := os.Getenv("WORLD_FILE")
	if path == "" {
		path = "data/world.json"
	}

	// world generation
	world, err := game.LoadWorld(path)
	if err != nil {
		Srv().Error("cannot load world", "file", path, "err", err)
		os.Exit(1)
	}

	// world validation
	errs, warns := world.Validate()
	for _, w := range warns {
		Srv().Warn("world validation", "problem", w)
	}
	if len(errs) > 0 {
		for _, e := range errs {
			Srv().Error("invalid world", "problem", e)
		}
		Srv().Error("refusing to start", "errors", len(errs))
		os.Exit(1)
	}
	Srv().Info("world loaded", "file", path, "rooms", len(world.Rooms),
		"npcs", len(world.NPCs), "quests", len(world.Quests),
		"items", len(world.Items))

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
	return world
}

func Server(listener net.Listener) {
	// generate and check the world
	world := checkWorld()

	// dispatcher creation
	dispatcher := NewDispatcher(world)

	Srv().Info("server started", "addr", listener.Addr().String())

	for {
		connection, err := listener.Accept()
		if err != nil {
			// verrification vu que listener.Close est appele listenner.Accept
			// retourne une net.ErrClosed donc on a ferme le listenner
			// volontairement
			if errors.Is(err, net.ErrClosed) {
				Srv().Info("server shutting down")
				closeAllConns()
				wg.Wait()
				Srv().Info("server stopped")
				return
			}
			Srv().Error("accept failed", "err", err)
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
