package server

import (
	"strings"
	"math/rand/v2"
)

func (d *Dispatcher) handleAttack(client *Client, args []string) {
	target := strings.ToUpper(strings.Join(args, " "))

	d.mu.Lock()
	defer d.mu.Unlock()

	// percentage chance of landing a critical hit
	initial_strengh := client.strengh
	damage_type := "normal"
	if rand.Float64() < 0.1 {
		client.strengh += (client.strengh*40/100)
		damage_type = "critical"
	}

	for _, npcs := range client.room.NPCs {
		if strings.ToUpper(npcs) == target || target == strings.ToUpper(d.world.NPCs[npcs].Name) {
			if  d.world.NPCs[npcs].Hp > 0 {
				if d.world.NPCs[npcs].Hostile {
					// Hit the NPC
					client.status = "combat"
					client.target = npcs
					if d.world.NPCs[npcs].Hp > 0 {
						d.world.NPCs[npcs].Hp -= client.strengh
						if d.world.NPCs[npcs].Hp < 0 {
							d.world.NPCs[npcs].Hp = 0
						}
						client.Send(`OK attack={"player_hp": %d, "player_max_hp": %d, "NPC_name": %s, "NPC_hp": %d, "NPC_max_hp": %d, "damage": %d, "damage_type": %s, "status": %s}`, client.hp, client.hp_max, d.world.NPCs[npcs].Name, d.world.NPCs[npcs].Hp, d.world.NPCs[npcs].Max_hp, client.strengh, damage_type, client.status)
						// if the hit kill the NPC
						if d.world.NPCs[npcs].Hp <= 0 {
							DeadNPC(client, initial_strengh, npcs, d)
							return
						}
						// NPC hit back
						d.world.NPCs[npcs].Hp += NPCsAction("ATTACK", client, d.world.NPCs[npcs].Strengh, d.world.NPCs[npcs].Hp, d.world.NPCs[npcs].Max_hp, d.world.NPCs[npcs].Name, d)
						client.strengh = initial_strengh
						return
					}
				} else {
					client.Send("ERR 405 NPC_NOT_HOSTILE")
					client.strengh = initial_strengh
					return
				}
			}
		}
	}
	client.Send("ERR 404 NPC_NOT_FOUND")
	client.strengh = initial_strengh
}

func (d *Dispatcher) handleStatus(client *Client, args []string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	client.status = "exploration"
	client.Send(`OK {"hp": %d, "max_hp": %d, "status": "%s"}`, client.hp, client.hp_max, client.status)
}

func (d *Dispatcher) handleFlee(client *Client, args []string) {
	if client.status == "combat" {
		if rand.Float64() < 0.4 {
			client.status = "exploration"
			client.Send(`OK flee=succed, "status"="%s"`, client.status)
			client.target = ""
			return
		} else {
			for _, npcs := range client.room.NPCs {
				if strings.EqualFold(strings.ToUpper(npcs), strings.ToUpper(client.target))  || strings.EqualFold(strings.ToUpper(client.target), strings.ToUpper(d.world.NPCs[npcs].Name)) {
					client.Send(`OK flee=fail, "status"="combat"`)
					npc := d.world.NPCs[npcs]
					npc.Hp += NPCsAction("FLEE", client, npc.Strengh, npc.Hp, npc.Max_hp, npc.Name, d)
					return
				}
			}
		}
	} else {
		client.Send("ERR 13 NOT_IN_FIGHT")
		return
	}
}

func (d *Dispatcher) handleDefend(client *Client, args []string) {
	if client.status == "combat" {
		client.Send(`OK defend, "status"="combat"`)
		for _, npcs := range client.room.NPCs {
			if strings.EqualFold(strings.ToUpper(npcs), strings.ToUpper(client.target))  || strings.EqualFold(strings.ToUpper(client.target), strings.ToUpper(d.world.NPCs[npcs].Name)) {
				npc := d.world.NPCs[npcs]
				npc_strengh := npc.Strengh - (npc.Strengh*40/100)
				npc.Hp += NPCsAction("DEFEND", client, npc_strengh, npc.Hp, npc.Max_hp, npc.Name, d)
				return
			}
		}
	} else {
		client.Send("ERR 13 NOT_IN_FIGHT")
		return
	}
}
