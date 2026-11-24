package server

import (
	"fmt"
	"net"
	"time"
	"the_answer_protocol/game"
)

type Client struct {
	conn      		net.Conn
	name      		string
	group     		*Group
	room      		*game.Room
	inventory 		[]*game.Item
	hp        		int
	hp_max    		int
	status    		string
	quests    		[]*game.Quest
	strength   		int
	target    		string
	buff    		string
	time_buff		int
	debuff			string
	time_debuff		int
	real_strength	int
	quit			bool
	kills			[]string
	lastChat		time.Time
}

func (c *Client) Send(format string, a ...any) {
	fmt.Fprintf(c.conn, format+"\n", a...)
}

func (client *Client) Death(d *Dispatcher) {
	// handle client death : reset his status, inventory, quests, target
	RemoveItemInventory(client, "false_item", d, false, "drop")

	client.Send(`EVT DEATH`)
	client.status = "exploration"
	d.world.NPCs[client.target].Fighter = ""
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
