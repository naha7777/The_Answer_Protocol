package server

import (
	"math/rand/v2"
)

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
			return Flee(client, d)
		} else if rand.Float64() < 0.8 {
			return Attack(client, strengh, d, npc_name, npc_hp, NPC_max_hp, damage_type)
		} else {
			return Regen(client, strengh, d, npc_name, npc_hp, NPC_max_hp, damage_type)
		}
	case "FLEE":
		 if rand.Float64() < 0.95 {
			return Attack(client, strengh, d, npc_name, npc_hp, NPC_max_hp, damage_type)
		} else {
			return Regen(client, strengh, d, npc_name, npc_hp, NPC_max_hp, damage_type)
		}
	case "DEFEND":
		if rand.Float64() < 0.05 {
			return Flee(client, d)
		} else if rand.Float64() < 0.75 {
			return Attack(client, strengh, d, npc_name, npc_hp, NPC_max_hp, damage_type)
		} else {
			return Regen(client, strengh, d, npc_name, npc_hp, NPC_max_hp, damage_type)
		}
	case "USE_ITEM":
		if rand.Float64() < 0.05 {
			return Flee(client, d)
		} else if rand.Float64() < 0.55 {
			return Attack(client, strengh, d, npc_name, npc_hp, NPC_max_hp, damage_type)
		} else {
			return Regen(client, strengh, d, npc_name, npc_hp, NPC_max_hp, damage_type)
		}
	}
	strengh = initial_strengh
	return 0
}

func Attack(client *Client, strengh int, d *Dispatcher, npc_name string, npc_hp int, NPC_max_hp int, damage_type string) int {
	client.Send(`EVT NPC ATTACK`)
	client.hp -= strengh
	if client.hp <= 0 {
		client.status = "exploration"
		client.Send(`OK attack={"player_hp": %d, "player_max_hp": %d, "NPC_name": %s, "NPC_hp": %d, "NPC_max_hp": %d, "damage": %d, "damage_type": %s, "status": %s}`, 0, client.hp_max, npc_name, npc_hp, NPC_max_hp, strengh, damage_type, client.status)
		client.Death(d)
		return 0
	}
	client.Send(`OK attack={"player_hp": %d, "player_max_hp": %d, "NPC_name": %s, "NPC_hp": %d, "NPC_max_hp": %d, "damage": %d, "damage_type": %s, "status": %s}`, client.hp, client.hp_max, npc_name, npc_hp, NPC_max_hp, strengh, damage_type, client.status)
	return 0
}

func Regen(client *Client, strengh int, d *Dispatcher, npc_name string, npc_hp int, NPC_max_hp int, damage_type string) int {
	client.Send(`EVT NPC REGENERATION : +5 hp`)
	new_hp := npc_hp + 5
	if new_hp > NPC_max_hp {
		new_hp = NPC_max_hp
	}
	client.Send(`OK attack={"player_hp": %d, "player_max_hp": %d, "NPC_name": %s, "NPC_hp": %d, "NPC_max_hp": %d, "damage": %d, "damage_type": %s, "status": %s}`, client.hp, client.hp_max, npc_name, new_hp, NPC_max_hp, strengh, damage_type, client.status)
	return 5
}

func Flee(client *Client, d *Dispatcher) int {
	client.status = "exploration"
	client.Send(`EVT NPC FLEE, "status"="%s"`, client.status)
	d.world.NPCs[client.target].Fighter = ""
	client.target = ""
	return 0
}

func DeadNPC(client *Client, initial_strengh int, npcs string, d *Dispatcher) {
	client.status = "exploration"
	client.Send(`OK attack=%s=dead, "status"="%s"`, npcs, client.status)
	client.kills = append(client.kills, npcs)
	room_npc := []string{}
	for _, np := range client.room.NPCs {
		if np != npcs {
			room_npc = append(room_npc, np)
		}
	}
	client.room.NPCs = room_npc
	client.target = ""
	client.strengh = initial_strengh
	SendRoom(d, client, "EVT NPC %s DEAD", npcs)
}
