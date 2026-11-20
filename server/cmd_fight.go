package server

import (
	"math/rand/v2"
	"strings"
)

func (d *Dispatcher) handleAttack(client *Client, args []string) {
	target := strings.ToUpper(strings.Join(args, " "))

	d.mu.Lock()
	defer d.mu.Unlock()

	// percentage chance of landing a critical hit
	damage_type := "normal"
	if rand.Float64() < 0.1 {
		client.strengh += (client.strengh * 40 / 100)
		damage_type = "critical"
	}

	for _, npcs := range client.room.NPCs {
		if strings.ToUpper(npcs) == target || target == strings.ToUpper(d.world.NPCs[npcs].Name) {
			if d.world.NPCs[npcs].Hp > 0 {
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
						CheckBuffDebuff(client, d)
						// if the hit kill the NPC
						if d.world.NPCs[npcs].Hp <= 0 {
							DeadNPC(client, client.real_strengh, npcs, d)
							return
						}
						// NPC hit back
						d.world.NPCs[npcs].Hp += NPCsAction("ATTACK", client, d.world.NPCs[npcs].Strengh, d.world.NPCs[npcs].Hp, d.world.NPCs[npcs].Max_hp, d.world.NPCs[npcs].Name, d)
						if damage_type == "critical" {
							client.strengh -= (client.strengh * 40 / 100)
						}
						return
					}
				} else {
					client.Send("ERR 405 NPC_NOT_HOSTILE")
					client.strengh = client.real_strengh
					return
				}
			}
		}
	}
	client.Send("ERR 404 NPC_NOT_FOUND")
	client.strengh = client.real_strengh
}

func (d *Dispatcher) handleStatus(client *Client, args []string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	client.status = "exploration"
	client.Send(`OK status={"hp": %d, "max_hp": %d, "status": "%s"}`, client.hp, client.hp_max, client.status)
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
				if strings.EqualFold(strings.ToUpper(npcs), strings.ToUpper(client.target)) || strings.EqualFold(strings.ToUpper(client.target), strings.ToUpper(d.world.NPCs[npcs].Name)) {
					client.Send(`OK flee=fail, "status"="combat"`)
					CheckBuffDebuff(client, d)
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
		CheckBuffDebuff(client, d)
		for _, npcs := range client.room.NPCs {
			if strings.EqualFold(strings.ToUpper(npcs), strings.ToUpper(client.target)) || strings.EqualFold(strings.ToUpper(client.target), strings.ToUpper(d.world.NPCs[npcs].Name)) {
				npc := d.world.NPCs[npcs]
				npc_strengh := npc.Strengh - (npc.Strengh * 40 / 100)
				npc.Hp += NPCsAction("DEFEND", client, npc_strengh, npc.Hp, npc.Max_hp, npc.Name, d)
				return
			}
		}
	} else {
		client.Send("ERR 13 NOT_IN_FIGHT")
		return
	}
}

func (d *Dispatcher) handleUseItem(client *Client, args []string) {
	if client.status == "combat" {
		if len(args) == 0 {
			client.Send("ERR 404 ITEM_NOT_FOUND ")
			return
		}
		query := strings.ToUpper(strings.Join(args, " "))
		d.mu.Lock()
		defer d.mu.Unlock()

		for _, item := range client.inventory {
			if !strings.EqualFold(item.Id, query) && !strings.EqualFold(item.Name, query) {
				continue
			}
			switch query {
			case "BEER":
				client.hp -= UseBeer(client, 3, d)
				npc := d.world.NPCs[client.target]
				npc.Hp += NPCsAction("USE_ITEM", client, npc.Strengh, npc.Hp, npc.Max_hp, npc.Name, d)
				return
			case "SWORD":
				client.strengh += UseSword(client, 10, d)
				npc := d.world.NPCs[client.target]
				npc.Hp += NPCsAction("USE_ITEM", client, npc.Strengh, npc.Hp, npc.Max_hp, npc.Name, d)
				return
			case "POISON POTION":
				d.world.NPCs[client.target].Hp -= UsePoisonPotion(client, 5, d)
				npc := d.world.NPCs[client.target]
				npc.Hp += NPCsAction("USE_ITEM", client, npc.Strengh, npc.Hp, npc.Max_hp, npc.Name, d)
				return
			case "HEALTH POTION":
				client.hp += UseHealthPotion(client, 4, d)
				npc := d.world.NPCs[client.target]
				npc.Hp += NPCsAction("USE_ITEM", client, npc.Strengh, npc.Hp, npc.Max_hp, npc.Name, d)
				return
			case "STRENGH POTION":
				client.strengh += UseStrenghPotion(client, 4, d)
				npc := d.world.NPCs[client.target]
				npc.Hp += NPCsAction("USE_ITEM", client, npc.Strengh, npc.Hp, npc.Max_hp, npc.Name, d)
				return
			default:
				client.Send("ERR 12 IN_FIGHT")
				return
			}
		}
		client.Send("ERR 404 ITEM_NOT_FOUND")
		return
	} else {
		client.Send("ERR 13 NOT_IN_FIGHT")
		return
	}
}
