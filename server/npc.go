package server

import (
	"math/rand/v2"
	"the_answer_protocol/game"
)

func NPCsAction(player_action string, client *Client, strength int, npc_hp int,
	NPC_max_hp int, npc_name string, d *Dispatcher) int {
	initial_strength := strength

	// percentage chance that the NPC will attack, flee, or regenerate depending
	// on the player's action
	switch player_action {
	case "ATTACK":
		if rand.Float64() < 0.1 {
			return Flee(client, d)
		} else if rand.Float64() < 0.8 {
			return Attack(client, strength, d, npc_name, npc_hp, NPC_max_hp,
				"normal")
		} else {
			return Regen(client, strength, d, npc_name, npc_hp, NPC_max_hp)
		}
	case "FLEE":
		if rand.Float64() < 0.95 {
			return Attack(client, strength, d, npc_name, npc_hp, NPC_max_hp,
				"normal")
		} else {
			return Regen(client, strength, d, npc_name, npc_hp, NPC_max_hp)
		}
	case "DEFEND":
		if rand.Float64() < 0.05 {
			return Flee(client, d)
		} else if rand.Float64() < 0.75 {
			return Attack(client, strength, d, npc_name, npc_hp, NPC_max_hp,
				"normal")
		} else {
			return Regen(client, strength, d, npc_name, npc_hp, NPC_max_hp)
		}
	case "USE_ITEM":
		if rand.Float64() < 0.05 {
			return Flee(client, d)
		} else if rand.Float64() < 0.55 {
			return Attack(client, strength, d, npc_name, npc_hp, NPC_max_hp,
				"normal")
		} else {
			return Regen(client, strength, d, npc_name, npc_hp, NPC_max_hp)
		}
	}
	strength = initial_strength
	return 0
}

func Attack(client *Client, strength int, d *Dispatcher, npc_name string,
	npc_hp int, NPC_max_hp int, damage_type string) int {
	// percentage chance of landing a critical hit
	if rand.Float64() < 0.1 {
		strength += (strength * 40 / 100)
		damage_type = "critical"
	}
	// NPC hit back
	client.Send(`EVT NPC ATTACK`)
	// player loose HP
	client.hp -= strength
	// handle player's death
	if client.hp <= 0 {
		client.Send(`OK attack={"player_hp": %d, "player_max_hp": %d, `+
			`"NPC_name": "%s", "NPC_hp": %d, "NPC_max_hp": %d, "damage": %d, `+
			`"damage_type": "%s", "status": "%s"}`, 0, client.hp_max, npc_name,
			npc_hp, NPC_max_hp, strength, damage_type, client.status)
		Log.Info("npc action", "source", npc_name, "action", "attack",
			"target", client.name, "damage", strength)
		client.Death(d)
		return 0
	}
	// print attack's datas
	client.Send(`OK attack={"player_hp": %d, "player_max_hp": %d, "NPC_name": `+
		`"%s", "NPC_hp": %d, "NPC_max_hp": %d, "damage": %d, "damage_type": "%s", `+
		`"status": "%s"}`, client.hp, client.hp_max, npc_name, npc_hp, NPC_max_hp,
		strength, damage_type, client.status)
	Log.Info("npc action", "source", npc_name, "action", "attack", "target",
		client.name, "damage", strength)
	return 0
}

func Regen(client *Client, strength int, d *Dispatcher, npc_name string,
	npc_hp int, NPC_max_hp int) int {
	client.Send(`EVT NPC REGENERATION : +5 hp`)
	// add 5 to NPC's HP
	new_hp := npc_hp + 5
	if new_hp > NPC_max_hp {
		new_hp = NPC_max_hp
	}
	client.Send(`OK attack={"player_hp": %d, "player_max_hp": %d, "NPC_name": `+
		`"%s", "NPC_hp": %d, "NPC_max_hp": %d, "damage": %d, "damage_type": "%s", `+
		`"status": "%s"}`, client.hp, client.hp_max, npc_name, new_hp, NPC_max_hp,
		strength, "normal", client.status)
	Log.Info("npc action", "source", npc_name, "action", "regenerate",
		"regeneration", 5)
	return new_hp - npc_hp
}

func Flee(client *Client, d *Dispatcher) int {
	// Change the player's status
	client.status = "exploration"
	client.Send(`EVT NPC FLEE, "status"="%s"`, client.status)
	Log.Info("npc action", "source", client.target, "action", "flee", "status",
		client.status)
	// reinitialise the NPC's fight and the player's target
	d.world.NPCs[client.target].Fighter = ""
	client.target = ""
	return 0
}

func DeadNPC(client *Client, initial_strength int, npcs string, d *Dispatcher) {
	// reinitialise player's status, player's target and player's strength
	client.status = "exploration"
	client.Send(`OK attack="%s"=dead, "status"="%s"`, npcs, client.status)
	// Add the NPC on the player's kill list
	client.kills = append(client.kills, npcs)
	room_npc := []string{}
	// Delete the NPC of the room
	for _, np := range client.room.NPCs {
		if np != npcs {
			room_npc = append(room_npc, np)
		}
	}
	client.room.NPCs = room_npc
	client.target = ""
	client.strength = initial_strength
	// Send an EVT to tell players that an NPC is dead
	SendRoom(d, client, `EVT NPC DEAD "%s"`, npcs)
	Srv().Info("death", "dead npc", npcs)
}

func FindNPC(d *Dispatcher, client *Client, npc_id string) *game.NPC {
	npc, ok := d.world.NPCs[npc_id]
	if !ok {
		Srv().Error("world inconsistency", "player", client.name, "npc",
			client.target)
		client.Send("ERR 500 INTERNAL")
		return nil
	}
	return npc
}
