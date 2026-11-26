package server

func UseBuff(client *Client, item string, buff string, time_buff int,
	d *Dispatcher) {
	// check if the item is in the client inventory
	found := false
	for _, i := range client.inventory {
		if i.Id == item {
			found = true
		}
	}
	// use the item during the fight
	if found {
		client.Send(`OK use=%s`, item)
		RemoveItemInventory(client, item, d, false, "fight")
	}
	client.buff = buff
	client.time_buff = time_buff
}

func CheckBuffDebuff(client *Client, d *Dispatcher) {
	// check if the client has a buff or debuff during a fight
	// if he has, active the effects for each function call
	switch client.buff {
	case "poison":
		if client.time_buff > 0 {
			client.time_buff -= 1
			client.Send(`EVT POISON : -5 npc hp during %d turns`,
						client.time_buff)
			npc := FindNPC(d, client, client.target)
			npc.Hp -= 5
		} else {
			client.buff = ""
		}
	case "strength":
		if client.time_buff > 0 {
			client.time_buff -= 1
			client.Send(`EVT Strength : +5 player strength during %d turns`,
						client.time_buff)
		} else {
			client.buff = ""
			client.strength = client.real_strength
		}

	case "health":
		if client.time_buff > 0 {
			client.time_buff -= 1
			client.Send(`EVT HEALTH : +5 player hp during %d turns`,
						client.time_buff)
			client.hp += 5
		} else {
			client.buff = ""
		}
	}

	if client.debuff == "poison" {
		if client.time_debuff > 0 {
			client.time_debuff -= 1
			client.Send(`EVT POISON : -5 player hp during %d turns`,
						client.time_debuff)
			client.hp -= 5
		} else {
			client.debuff = ""
		}
	}
}
