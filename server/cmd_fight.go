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
						NPCsAction("ATTACK", client, d.world.NPCs[npcs].Strengh, d.world.NPCs[npcs].Hp)
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
			NPCsAction("FLEE", client, 0, 0)
			return
		}
	} else {
		client.Send("ERR 13 NOT_IN_FIGHT")
		return
	}

}

func NPCsAction(player_action string, client *Client, strengh int, npc_hp int) {
	switch player_action {
	case "ATTACK":
		if rand.Float64() < 0.2 {
			ChangeStatus(client)
			client.Send(`EVT NPC FLEE : "status"="%s"`, client.status)
			return
		} else if rand.Float64() < 0.7 {
			client.Send(`EVT NPC ATTACK`)
			client.hp -= strengh
			// gerer la mort
			client.Send("OK {'player_hp': %d, 'NPC_hp': %v, 'damage': %d, 'status': %s}", client.hp, npc_hp, strengh, client.status)
			return
		} else {
			client.Send(`EVT NPC REGENERATION`)
			npc_hp += 5
			return
		}
	case "FLEE":

	case "DEFEND":

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
