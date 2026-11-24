package server

func UseBeer(client *Client, time int, d *Dispatcher) int {
	// Use the beer if it's on the player's inventory
	found := false
	for _, i := range client.inventory {
		if i.Id == "item.beer"{
			found = true
		}
	}
	if found {
		client.Send(`OK use=item.beer`)
		RemoveItemInventory(client, "item.beer", d, false, "fight")
	}
	client.debuff = "poison"
	client.time_debuff = time - 1
	client.Send(`EVT BEER : -5 player hp during %d turns`, client.time_debuff)
	return 5
}

func UseSword(client *Client, time int, d *Dispatcher) int {
	// Use the sword
	UseBuff(client, "item.sword", "strengh", time, d)
	client.Send(`EVT SWORD : +5 player strengh during %d turns`, client.time_buff)
	return 5
}

func UsePoisonPotion(client *Client, time int, d *Dispatcher) int {
	// Use the poison potion
	UseBuff(client, "item.poison_potion", "poison", time, d)
	client.Send(`EVT POISON POTION : -5 npc hp during %d turns`, client.time_buff)
	return 5
}

func UseStrenghPotion(client *Client, time int, d *Dispatcher) int {
	// Use the strengh potion
	UseBuff(client, "item.strengh_potion", "strengh", time, d)
	client.Send(`EVT STRENGH POTION : +5 player strengh during %d turns`, client.time_buff)
	return 5
}

func UseHealthPotion(client *Client, time int, d *Dispatcher) int {
	// Use the health potion
	UseBuff(client, "item.health_potion", "health", time, d)
	client.Send(`EVT HEALTH POTION : +5 player hp during %d turns`, client.time_buff)
	return 5
}
