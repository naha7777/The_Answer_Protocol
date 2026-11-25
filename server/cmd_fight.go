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
		client.strength += (client.strength * 40 / 100)
		damage_type = "critical"
	}

	// a player can't fight a NPC in another room
	for _, npcs := range client.room.NPCs {
		npc := FindNPC(d, client, npcs)
		if strings.ToUpper(npcs) == target || target == strings.ToUpper(
			npc.Name) {
			if npc.Hp > 0 {
				// a player can't fight a NPC if he is not hostile
				if npc.Hostile {
					// another player can't attack a NPC if he is in a fight
					if npc.Fighter == "" || npc.Fighter == client.name {
						// a player can't fight two NPC at the same time
						if npcs == client.target || client.target == "" {
							// Hit the NPC
							client.status = "combat"
							client.target = npcs
							npc.Fighter = client.name
							if npc.Hp > 0 {
								npc.Hp -= client.strength
								if npc.Hp < 0 {
									npc.Hp = 0
								}
								client.Send(`OK attack={"player_hp": %d, "player_max_hp": %d, "NPC_name": "%s", "NPC_hp": %d, "NPC_max_hp": %d, "damage": %d, "damage_type": "%s", "status": "%s"}`, client.hp, client.hp_max, npc.Name, npc.Hp, npc.Max_hp, client.strength, damage_type, client.status)
								CheckBuffDebuff(client, d)
								// if the hit kill the NPC
								if npc.Hp <= 0 {
									DeadNPC(client, client.real_strength,
										npcs, d)
									return
								}
								// NPC hit back
								npc.Hp += NPCsAction(
									"ATTACK",
									client,
									npc.Strength,
									npc.Hp,
									npc.Max_hp,
									npc.Name, d)
								if damage_type == "critical" {
									client.strength -= (client.strength * 40 / 100)
								}
								return
							}
						} else {
							client.Send("ERR 12 IN_FIGHT")
							client.strength = client.real_strength
							return
						}
					} else {
						client.Send("ERR 14 NPC_IN_FIGHT")
						client.strength = client.real_strength
						return
					}
				} else {
					client.Send("ERR 405 NPC_NOT_HOSTILE")
					client.strength = client.real_strength
					return
				}
			}
		}
	}
	client.Send("ERR 404 NPC_NOT_FOUND")
	client.strength = client.real_strength
}

func (d *Dispatcher) handleStatus(client *Client, args []string) {
	// displqy client HP, client max HP and client status
	d.mu.Lock()
	defer d.mu.Unlock()

	client.status = "exploration"
	client.Send(`OK status={"hp": %d, "max_hp": %d, "status": "%s"}`, client.hp,
		client.hp_max, client.status)
}

func (d *Dispatcher) handleFlee(client *Client, args []string) {
	// a player has 40% to success his flee during a fight
	if client.status == "combat" {
		if rand.Float64() < 0.4 {
			client.status = "exploration"
			client.Send(`OK flee=succed, "status"="%s"`, client.status)
			d.world.NPCs[client.target].Fighter = ""
			client.target = ""
			return
		} else {
			npc := FindNPC(d, client, client.target)
			client.Send(`OK flee=fail, "status"="combat"`)
			CheckBuffDebuff(client, d)
			// if the player fail his flee, the NPC has an action
			npc.Hp += NPCsAction("FLEE", client, npc.Strength, npc.Hp,
				npc.Max_hp, npc.Name, d)
			return
		}
	} else {
		client.Send("ERR 13 NOT_IN_FIGHT")
		return
	}
}

func (d *Dispatcher) handleDefend(client *Client, args []string) {
	// during a fight, the player can reduce the NPC damages by using
	// Defend command
	if client.status == "combat" {
		client.Send(`OK defend, "status"="combat"`)
		CheckBuffDebuff(client, d)
		npc := FindNPC(d, client, client.target)
		npc_strength := npc.Strength - (npc.Strength * 40 / 100)
		npc.Hp += NPCsAction("DEFEND", client, npc_strength, npc.Hp, npc.Max_hp,
			npc.Name, d)
		return
	} else {
		client.Send("ERR 13 NOT_IN_FIGHT")
		return
	}
}

func (d *Dispatcher) handleUseItem(client *Client, args []string) {
	// during a fight, the player can use items
	if client.status == "combat" {
		if len(args) == 0 {
			client.Send("ERR 404 ITEM_NOT_FOUND ")
			return
		}
		query := strings.ToUpper(strings.Join(args, " "))
		d.mu.Lock()
		defer d.mu.Unlock()

		// check if the item is on his inventory
		for _, item := range client.inventory {
			if !strings.EqualFold(item.Id, query) && !strings.EqualFold(
				item.Name, query) {
				continue
			}
			// call the item function, then the NPC has an action
			switch query {
			case "BEER":
				client.hp -= UseBeer(client, 3, d)
				npc := FindNPC(d, client, client.target)
				if client.hp <= 0 {
					client.Death(d)
					return
				}
				npc.Hp += NPCsAction("USE_ITEM", client, npc.Strength, npc.Hp,
					npc.Max_hp, npc.Name, d)
				return
			case "SWORD":
				client.strength += UseSword(client, 10, d)
				npc := FindNPC(d, client, client.target)
				npc.Hp += NPCsAction("USE_ITEM", client, npc.Strength, npc.Hp,
					npc.Max_hp, npc.Name, d)
				return
			case "POISON POTION":
				npc := FindNPC(d, client, client.target)
				npc.Hp -= UsePoisonPotion(client, 5, d)
				if npc.Hp <= 0 {
					DeadNPC(client, client.real_strength, npc.Name, d)
					return
				}
				npc.Hp += NPCsAction("USE_ITEM", client, npc.Strength, npc.Hp,
					npc.Max_hp, npc.Name, d)
				return
			case "HEALTH POTION":
				client.hp += UseHealthPotion(client, 4, d)
				npc := FindNPC(d, client, client.target)
				npc.Hp += NPCsAction("USE_ITEM", client, npc.Strength, npc.Hp,
					npc.Max_hp, npc.Name, d)
				return
			case "Strength POTION":
				client.strength += UseStrengthPotion(client, 4, d)
				npc := FindNPC(d, client, client.target)
				npc.Hp += NPCsAction("USE_ITEM", client, npc.Strength, npc.Hp,
					npc.Max_hp, npc.Name, d)
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
