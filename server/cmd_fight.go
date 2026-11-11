package server

import "strings"

func (d *Dispatcher) handleAttack(client *Client, args []string) {
	target := strings.ToUpper(strings.Join(args, " "))
	damage := 10
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, npcs := range client.room.NPCs {
		if strings.ToUpper(npcs) == target {
			if  d.world.NPCs[npcs].Hp > 0 {
				if d.world.NPCs[npcs].Hostile {
				client.status = "combat"
				if d.world.NPCs[npcs].Hp > 0 {
					d.world.NPCs[npcs].Hp -= damage
				}
				if d.world.NPCs[npcs].Hp < 0 {
					d.world.NPCs[npcs].Hp = 0
				}
				client.Send("OK {'attacker_hp': %d, 'target_hp': %v, 'damage': %d, 'status': %s}", client.hp, d.world.NPCs[npcs].Hp, damage, client.status)
				return
				} else {
					client.Send("ERR 405 NPC_NOT_HOSTILE")
					return
				}
			}
		} else if target == strings.ToUpper(d.world.NPCs[npcs].Name) {
			if d.world.NPCs[npcs].Hp > 0 {
				if d.world.NPCs[npcs].Hostile {
				client.status = "combat"
				if d.world.NPCs[npcs].Hp > 0 {
					d.world.NPCs[npcs].Hp -= damage
				}
				if d.world.NPCs[npcs].Hp < 0 {
					d.world.NPCs[npcs].Hp = 0
				}
				client.Send("OK {'attacker_hp': %d, 'target_hp': %v, 'damage': %d, 'status': %s}", client.hp, d.world.NPCs[npcs].Hp, damage, client.status)
				return
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

func ChangeStatus(client *Client) {
	if client.hp <= client.hp_max && client.hp >= (80%*&client.hp_max) {
		client.status = "very healthy"
	} else if client.hp < (80%*&client.hp_max) && client.hp >= (60%*&client.hp_max) {
		client.status = "good"
	} else if client.hp < (60%*&client.hp_max) && client.hp >= (40%*&client.hp_max) {
		client.status = "mid"
	} else if client.hp < (40%*&client.hp_max) && client.hp >= (20%*&client.hp_max) {
		client.status = "suffering"
	} else if client.hp < (20%*&client.hp_max) && client.hp > 0 {
		client.status = "almost dead"
	} else if client.hp == 0 {
		client.status = "dead"
	} else {
		client.status = "far too healthy"
	}
}
