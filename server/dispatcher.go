package server

import (
	"strings"
	"sync"
	"the_answer_protocol/game"
	"time"
)

type CommandFunc func(client *Client, args []string)

type Dispatcher struct {
	mu       sync.Mutex
	commands map[string]CommandFunc
	clients  map[string]*Client
	groups   map[string]*Group
	world    *game.World
}

func NewDispatcher(w *game.World) *Dispatcher {
	d := &Dispatcher{
		commands: make(map[string]CommandFunc),
		clients:  make(map[string]*Client),
		groups:   make(map[string]*Group),
		world:    w,
	}
	d.registerCommands()
	return d
}

func (d *Dispatcher) registerCommands() {
	d.commands["CONNECT"] = d.handleConnect
	d.commands["QUIT"] = d.handleQuit
	d.commands["GROUP CREATE"] = d.handleGroupCreate
	d.commands["GROUP INVITE"] = d.handleGroupInvite
	d.commands["GROUP JOIN"] = d.handleGroupJoin
	d.commands["GROUP LEAVE"] = d.handleGroupLeave
	d.commands["MOVE"] = d.handleMove
	d.commands["DROP"] = d.handleDrop
	d.commands["TAKE"] = d.handleTake
	d.commands["TALK"] = d.handleTalk
	d.commands["CHAT"] = d.handleChat
	d.commands["LOOK"] = d.handleLook
	d.commands["WHO"] = d.handleWho
	d.commands["INVENTORY"] = d.handleInventory
	d.commands["ATTACK"] = d.handleAttack
	d.commands["STATUS"] = d.handleStatus
	d.commands["QUEST"] = d.handleQuest
	d.commands["QUESTS"] = d.handleQuests
	d.commands["COMPLETE_QUEST"] = d.handleCompleteQuest
	d.commands["INFO_GROUP"] = d.handleInfoGroup
	d.commands["ABANDON_QUEST"] = d.handleAbandonQuest
	d.commands["FLEE"] = d.handleFlee
	d.commands["DEFEND"] = d.handleDefend
	d.commands["USE_ITEM"] = d.handleUseItem
	d.commands["FIND_QUESTS"] = d.handleFindQuests
	// d.commands["CRASH"] = func(c *Client, args []string) {
	// 	var m map[string]int
	// 	m["x"] = 1
	// }
}

func (c *Client) tooManyCommands(limit int, window time.Duration) bool {
	// count commands frequency to avoid spam
	now := time.Now()
	kept := c.cmdTimes[:0]
	for _, t := range c.cmdTimes {
		if now.Sub(t) < window {
			kept = append(kept, t)
		}
	}
	c.cmdTimes = append(kept, now)
	return len(c.cmdTimes) > limit
}

func (d *Dispatcher) Dispatch(client *Client, line string) {
	if client.tooManyCommands(20, time.Second) {
		if !client.floodLogged {
			client.log().Warn("abuse: command flooding", "limit", "20/s")
			client.floodLogged = true
		}
		client.write("ERR 018 SPAM")
		return
	}
	client.floodLogged = false

	line = strings.TrimSpace(line)
	if line == "" {
		return
	}
	parts := strings.Fields(line)
	var cmd string
	var args []string
	first := strings.ToUpper(parts[0])

	client.log().Info("command", "cmd", first, "args", parts[1:])

	// spam protection
	if first == "CHAT" {
		now := time.Now()
		if !client.lastChat.IsZero() && now.Sub(client.lastChat) < time.Second*2 {
			client.log().Warn("abuse: chat spam")
			client.Send("ERR 018 SPAM")
			return
		}
		client.lastChat = now
	}

	// call appropriate function
	// block some commands during a fight
	if client.status != "combat" {
		if first == "GROUP" {
			if len(parts) >= 2 {
				cmd = strings.ToUpper(parts[0] + " " + parts[1])
				args = parts[2:]
			}
		} else {
			cmd = strings.ToUpper(parts[0])
			args = parts[1:]
		}
		if function, exists := d.commands[cmd]; exists {
			function(client, args)
		} else {
			client.Send("ERR 001 UNKNOWN_COMMAND: %s", first)
			return
		}
	} else {
		if first == "ATTACK" || first == "DEFEND" || first == "FLEE" ||
			first == "STATUS" || first == "USE_ITEM" || first == "INVENTORY" {
			cmd = strings.ToUpper(parts[0])
			args = parts[1:]
		} else {
			client.Send("ERR 012 IN_FIGHT")
			return
		}
		if function, exists := d.commands[cmd]; exists {
			function(client, args)
		} else {
			client.Send("ERR 012 IN_FIGHT")
			return
		}
	}
}
