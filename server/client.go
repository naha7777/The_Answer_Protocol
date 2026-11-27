package server

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"the_answer_protocol/game"
	"time"
)

type Client struct {
	conn          net.Conn
	name          string
	group         *Group
	room          *game.Room
	inventory     []*game.Item
	hp            int
	hp_max        int
	status        string
	quests        []*game.Quest
	strength      int
	target        string
	buff          string
	time_buff     int
	debuff        string
	time_debuff   int
	real_strength int
	quit          bool
	kills         []string
	lastChat      time.Time
	wmu           sync.Mutex
	cmdTimes      []time.Time
	floodLogged   bool
	dead          atomic.Bool
	nb_death      int
}

func (c *Client) label() string {
	if c.name != "" {
		return c.name
	}
	return c.conn.RemoteAddr().String()
}

func (c *Client) write(msg string) error {
	if c.dead.Load() {
		return net.ErrClosed
	}
	c.wmu.Lock()
	defer c.wmu.Unlock()
	_, err := fmt.Fprintf(c.conn, "%s\n", msg)
	if err != nil && c.dead.CompareAndSwap(false, true) {
		Srv().Warn("client unreachable", "to", c.label(), "err", err)
	}
	return err
}

func (c *Client) Send(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	if err := c.write(msg); err != nil {
		return
	}

	kind, rest, _ := strings.Cut(msg, " ")
	switch kind {
	case "ERR":
		code, _, _ := strings.Cut(rest, " ")
		if code == "500" {
			Srv().Error("error sent", "to", c.label(), "code", code, "message",
				msg)
		} else {
			Srv().Warn("error sent", "to", c.label(), "code", code, "message",
				msg)
		}
	case "EVT":
		Srv().Info("event sent", "event", msg, "recipients",
			[]string{c.label()}, "count", 1)
	default:
		Srv().Info("response sent", "to", c.label(), "message", msg)
	}
}

func (client *Client) Death(d *Dispatcher) {
	// handle client death : reset his status, inventory, quests, target
	RemoveItemInventory(client, "false_item", d, false, "drop")
	client.nb_death += 1

	client.Send(`EVT PLAYER DEATH`)
	Srv().Info("death", "player", client.name)
	client.status = "exploration"
	npc := FindNPC(d, client, client.target)
	npc.Fighter = ""
	client.target = ""
	// Send an EVT to the players
	SendAll(d, client, `EVT PLAYER DEAD %s`, client.name)

	// reset client HP
	client.hp = client.hp_max - (client.nb_death * 5)
	if client.hp <= 0 {
		client.Send("GAME OVER")
		client.quit = true
		return
	}

	// handle respawn
	for _, val := range d.world.Rooms {
		if d.world.Respawn == val.Id {
			SendRoom(d, client, `EVT ROOM PRESENCE LEAVE %s`, client.name)
			client.room.Players = Remove(client.name, client.room.Players)
			client.room = val
			client.room.Players = append(client.room.Players, client.name)
			client.Send(`OK room=%s, player_hp=%d, player_max_hp=%d`,
				d.world.Respawn, client.hp, client.hp_max)
			Srv().Info("respawn", "player", client.name, "room",
				client.room.Name, "player_hp", client.hp)
			SendRoom(d, client, `EVT ROOM PRESENCE ENTER %s`, client.name)
		}
	}
}

func (c *Client) tooManyCommands(limit int, window time.Duration) bool {
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
