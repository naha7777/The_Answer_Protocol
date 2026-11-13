package server

import (
	"strings"
	"sync"
	"the_answer_protocol/game"
)

type CommandFunc func(client *Client, args []string)

type Dispatcher struct {
	mu			sync.Mutex
	commands	map[string]CommandFunc
	clients		map[string]*Client
	groups		map[string]*Group
	world		*game.World
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
	d.commands["ABANDON_QUEST"] = d.handleAbandonQuest
}

func (d *Dispatcher) Dispatch(client *Client, line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}
	parts := strings.Fields(line)
	var cmd string
	var args []string
	first := strings.ToUpper(parts[0])
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
		client.Send("ERR 01 UNKNOWN_COMMAND: %s", first)
		return
	}
}
