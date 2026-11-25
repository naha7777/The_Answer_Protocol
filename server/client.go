package server

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"the_answer_protocol/game"
	"time"
)

type Client struct {
	conn         net.Conn
	name         string
	group        *Group
	room         *game.Room
	inventory    []*game.Item
	hp           int
	hp_max       int
	status       string
	quests       []*game.Quest
	strengh      int
	target       string
	buff         string
	time_buff    int
	debuff       string
	time_debuff  int
	real_strengh int
	quit         bool
	kills        []string
	lastChat     time.Time
	wmu          sync.Mutex
}

func (c *Client) label() string {
	if c.name != "" {
		return c.name
	}
	return c.conn.RemoteAddr().String()
}

func (c *Client) write(msg string) error {
	c.wmu.Lock()
	defer c.wmu.Unlock()
	_, err := fmt.Fprintf(c.conn, "%s\n", msg)
	return err
}

func (c *Client) Send(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	if err := c.write(msg); err != nil {
		Log.Error("send failed", "to", c.label(), "err", err)
		return
	}

	kind, rest, _ := strings.Cut(msg, " ")
	switch kind {
	case "ERR":
		code, _, _ := strings.Cut(rest, " ")
		if code == "500" {
			Log.Error("error sent", "to", c.label(), "code", code, "message", msg)
		} else {
			Log.Warn("error sent", "to", c.label(), "code", code, "message", msg)
		}
	case "EVT":
		Log.Info("event sent", "event", msg, "recipients", []string{c.label()}, "count", 1)
	default:
		Log.Info("response sent", "to", c.label(), "message", msg)
	}
}

func (client *Client) Death(d *Dispatcher) {
	// handle client death : reset his status, inventory, quests, target
	RemoveItemInventory(client, "false_item", d, false, "drop")

	client.Send(`EVT DEATH`)
	client.status = "exploration"
	npc := FindNPC(d, client, client.target)
	npc.Fighter = ""
	client.target = ""
	// Send an EVT to the players
	SendAll(d, client, `EVT dead=%s`, client.name)

	// handle respawn
	for _, val := range d.world.Rooms {
		if d.world.Respawn == val.Id {
			SendRoom(d, client, `EVT ROOM PRESENCE LEAVE %s`, client.name)
			client.room.Players = Remove(client.name, client.room.Players)
			client.room = val
			client.room.Players = append(client.room.Players, client.name)
			client.Send(`OK room=%s`, d.world.Respawn)
			SendRoom(d, client, `EVT ROOM PRESENCE ENTER %s`, client.name)
		}
	}
	// reset client HP
	client.hp = client.hp_max
}
