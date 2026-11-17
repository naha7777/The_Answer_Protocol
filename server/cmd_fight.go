package server

import (
	"strings"
	"math/rand/v2"
)

func (d *Dispatcher) handleAttack(client *Client, args []string) {
	target := strings.ToUpper(strings.Join(args, " "))
	damage := client.strengh
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, npcs := range client.room.NPCs {
		if strings.ToUpper(npcs) == target || target == strings.ToUpper(d.world.NPCs[npcs].Name) {
			if  d.world.NPCs[npcs].Hp > 0 {
				if d.world.NPCs[npcs].Hostile {
				client.status = "combat"
				if d.world.NPCs[npcs].Hp > 0 {
					d.world.NPCs[npcs].Hp -= damage
					client.Send("OK {'player_hp': %d, 'NPC_hp': %v, 'damage': %d, 'status': %s}", client.hp, d.world.NPCs[npcs].Hp, damage, client.status)
					if d.world.NPCs[npcs].Hp <= 0 {
						d.world.NPCs[npcs].Hp = 0
						ChangeStatus(client)
						client.Send(`OK %s=dead, "status"="%s"`, npcs, client.status)
						return
					}
					return
				}
				} else {
					client.Send("ERR 405 NPC_NOT_HOSTILE")
					return
				}
			}
		}
	}
	client.Send("ERR 404 NPC_NOT_FOUND")
}

func (d *Dispatcher) handleStatus(client *Client, args []string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	ChangeStatus(client)
	client.Send("OK {'hp': %d, 'max_hp': %d, 'status': %s}", client.hp,
		client.hp_max, client.status)
}

func (d *Dispatcher) handleFlee(client *Client, args []string) {
	if client.status == "combat" {
		if rand.Float64() < 0.5 {
			ChangeStatus(client)
			client.Send(`OK succed_flee, "status"="%s"`, client.status)
			return
		} else {
			client.Send(`OK fail_flee, "status"="combat"`)
			return
		}
	} else {
		client.Send("ERR 13 NOT_IN_FIGHT")
		return
	}

}

func ChangeStatus(client *Client) {
	if client.hp <= client.hp_max && (client.hp*100/client.hp_max >= 80) {
		client.status = "very healthy"
	} else if (client.hp*100/client.hp_max < 80) && (client.hp*100/client.hp_max >= 60) {
		client.status = "good"
	} else if (client.hp*100/client.hp_max < 60) && (client.hp*100/client.hp_max >= 40) {
		client.status = "mid"
	} else if (client.hp*100/client.hp_max < 40) && (client.hp*100/client.hp_max >= 20) {
		client.status = "suffering"
	} else if (client.hp*100/client.hp_max < 20) && client.hp > 0 {
		client.status = "almost dead"
	} else if client.hp == 0 {
		client.status = "dead"
	} else {
		client.status = "far too healthy"
	}
}
