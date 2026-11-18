package server

import (
	"strings"
	"math/rand/v2"
	"slices"
	"the_answer_protocol/game"
)

func (d *Dispatcher) handleAttack(client *Client, args []string) {
	target := strings.ToUpper(strings.Join(args, " "))

	d.mu.Lock()
	defer d.mu.Unlock()

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
					client.status = "combat"
					client.target = npcs
					if d.world.NPCs[npcs].Hp > 0 {
						d.world.NPCs[npcs].Hp -= client.strengh
						if d.world.NPCs[npcs].Hp < 0 {
							d.world.NPCs[npcs].Hp = 0
						}
						client.Send(`OK attack={"player_hp": %d, "player_max_hp": %d, "NPC_name": %s, "NPC_hp": %d, "NPC_max_hp": %d, "damage": %d, "damage_type": %s, "status": %s}`, client.hp, client.hp_max, d.world.NPCs[npcs].Name, d.world.NPCs[npcs].Hp, d.world.NPCs[npcs].Max_hp, client.strengh, damage_type, client.status)
						if d.world.NPCs[npcs].Hp <= 0 {
							client.status = "exploration"
							client.Send(`OK attack=%s=dead, "status"="%s"`, npcs, client.status)
							room_npc := []string{}
							for _, np := range client.room.NPCs {
								if np != npcs {
									room_npc = append(room_npc, np)
								}
							}
							client.room.NPCs = room_npc
							client.target = ""
							client.strengh = initial_strengh
							for _, c := range d.clients {
								if slices.Contains(client.room.Players, c.name) {
									if c.name != client.name {
										c.Send("EVT NPC %s DEAD", npcs)
									}
								}
							}
							return
						}
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
		} else if rand.Float64() < 0.6 {
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

func NPCsAction(player_action string, client *Client, strengh int, npc_hp int, NPC_max_hp int, npc_name string, d *Dispatcher) int {
	initial_strengh := strengh
	damage_type := "normal"
	if rand.Float64() < 0.1 {
		strengh += (strengh*40/100)
		damage_type = "critical"
	}
	switch player_action {
	case "ATTACK":
		if rand.Float64() < 0.1 {
			client.status = "exploration"
			client.Send(`EVT NPC FLEE, "status"="%s"`, client.status)
			client.target = ""
			return 0
		} else if rand.Float64() < 0.8 {
			client.Send(`EVT NPC ATTACK`)
			client.hp -= strengh
			if client.hp <= 0 {
				client.status = "exploration"
				client.Send(`OK attack={"player_hp": %d, "player_max_hp": %d, "NPC_name": %s, "NPC_hp": %d, "NPC_max_hp": %d, "damage": %d, "damage_type": %s, "status": %s}`, 0, client.hp_max, npc_name, npc_hp, NPC_max_hp, strengh, damage_type, client.status)
				DeadPlayer(client, d)
				return 0
			}
			client.Send(`OK attack={"player_hp": %d, "player_max_hp": %d, "NPC_name": %s, "NPC_hp": %d, "NPC_max_hp": %d, "damage": %d, "damage_type": %s, "status": %s}`, client.hp, client.hp_max, npc_name, npc_hp, NPC_max_hp, strengh, damage_type, client.status)
			return 0
		} else {
			client.Send(`EVT NPC REGENERATION : +5 hp`)
			new_hp := npc_hp + 5
			if new_hp > NPC_max_hp {
				new_hp = NPC_max_hp
			}
			client.Send(`OK attack={"player_hp": %d, "player_max_hp": %d, "NPC_name": %s, "NPC_hp": %d, "NPC_max_hp": %d, "damage": %d, "damage_type": %s, "status": %s}`, client.hp, client.hp_max, npc_name, new_hp, NPC_max_hp, strengh, damage_type, client.status)
			return 5
		}
	case "FLEE":
		 if rand.Float64() < 0.95 {
			client.Send(`EVT NPC ATTACK`)
			client.hp -= strengh
			if client.hp <= 0 {
				client.status = "exploration"
				client.Send(`OK attack={"player_hp": %d, "player_max_hp": %d, "NPC_name": %s, "NPC_hp": %d, "NPC_max_hp": %d, "damage": %d, "damage_type": %s, "status": %s}`, 0, client.hp_max, npc_name, npc_hp, NPC_max_hp, strengh, damage_type, client.status)
				DeadPlayer(client, d)
				return 0
			}
			client.Send(`OK attack={"player_hp": %d, "player_max_hp": %d, "NPC_name": %s, "NPC_hp": %d, "NPC_max_hp": %d, "damage": %d, "damage_type": %s, "status": %s}`, client.hp, client.hp_max, npc_name, npc_hp, NPC_max_hp, strengh, damage_type, client.status)
			return 0
		} else {
			client.Send(`EVT NPC REGENERATION : +5 hp`)
			new_hp := npc_hp + 5
			if new_hp > NPC_max_hp {
				new_hp = NPC_max_hp
			}
			client.Send(`OK attack={"player_hp": %d, "player_max_hp": %d, "NPC_name": %s, "NPC_hp": %d, "NPC_max_hp": %d, "damage": %d, "damage_type": %s, "status": %s}`, client.hp, client.hp_max, npc_name, new_hp, NPC_max_hp, strengh, damage_type, client.status)
			return 5
		}
	case "DEFEND":
		if rand.Float64() < 0.05 {
			client.status = "exploration"
			client.Send(`EVT NPC FLEE, "status"="%s"`, client.status)
			client.target = ""
			return 0
		} else if rand.Float64() < 0.75 {
			client.Send(`EVT NPC ATTACK`)
			client.hp -= (strengh - (strengh*40/100))
			if client.hp <= 0 {
				client.status = "exploration"
				client.Send(`OK attack={"player_hp": %d, "player_max_hp": %d, "NPC_name": %s, "NPC_hp": %d, "NPC_max_hp": %d, "damage": %d, "damage_type": %s, "status": %s}`, 0, client.hp_max, npc_name, npc_hp, NPC_max_hp, strengh, damage_type, client.status)
				DeadPlayer(client, d)
				return 0
			}
			client.Send(`OK attack={"player_hp": %d, "player_max_hp": %d, "NPC_name": %s, "NPC_hp": %d, "NPC_max_hp": %d, "damage": %d, "damage_type": %s, "status": %s}`, client.hp, client.hp_max, npc_name, npc_hp, NPC_max_hp, strengh, damage_type, client.status)
			return 0
		} else {
			client.Send(`EVT NPC REGENERATION : +5 hp`)
			new_hp := npc_hp + 5
			if new_hp > NPC_max_hp {
				new_hp = NPC_max_hp
			}
			client.Send(`OK attack={"player_hp": %d, "player_max_hp": %d, "NPC_name": %s, "NPC_hp": %d, "NPC_max_hp": %d, "damage": %d, "damage_type": %s, "status": %s}`, client.hp, client.hp_max, npc_name, new_hp, NPC_max_hp, strengh, damage_type, client.status)
			return 5
		}
	}
	strengh = initial_strengh
	return 0
}

func DeadPlayer(client *Client, d *Dispatcher) {
	var final_inventory []*game.Item
	for _, item := range client.inventory {
		client.room.Items = append(client.room.Items, item.Id)
	}
	client.inventory = final_inventory

	var final_quests []*game.Quest
	for _, quest := range client.quests {
		quest.Status = "available"
	}
	client.quests = final_quests

	client.Send(`EVT DEATH`)
	client.status = "exploration"
	client.target = ""
	for _, c := range d.clients {
		if c != client {
			c.Send(`EVT dead=%s`, client.name)
		}
	}

	for _, val := range d.world.Rooms {
		if d.world.Respawn == val.Id {
			for _, c := range d.clients {
				if slices.Contains(client.room.Players, c.name) {
					if c.name != client.name {
						c.Send(`EVT ROOM PRESENCE LEAVE %s`, client.name)
					}
				}
			}
			client.room.Players = remove(client.name, client.room.Players)
			client.room = val
			client.room.Players = append(client.room.Players, client.name)
			client.Send(`OK room=%s`, d.world.Respawn)
			for _, c := range d.clients {
				if slices.Contains(client.room.Players, c.name) {
					if c.name != client.name {
						c.Send(`EVT ROOM PRESENCE ENTER %s`, client.name)
					}
				}
			}
		}
	}
	client.hp = client.hp_max
}
