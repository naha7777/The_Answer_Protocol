package server

func UseBeer(client *Client, time int) int {
	found := false
	for _, i := range client.inventory {
		if i.Id == "item.beer"{
			found = true
		}
	}
	if found {
		client.Send(`OK use=item.beer`)
		RemoveItemInventory(client, "item.beer")
	}
	client.debuff = "poison"
	client.time_debuff = time - 1
	client.Send(`EVT BEER : -5 player hp during %d turns`, client.time_debuff)
	return 5
}

func UseSword(client *Client, time int) int {
	UseBuff(client, "item.sword", "strengh", time)
	client.Send(`EVT SWORD : +5 player strengh during %d turns`, client.time_buff)
	return 5
}

func UsePoisonPotion(client *Client, time int) int {
	UseBuff(client, "item.poison_potion", "poison", time)
	client.Send(`EVT POISON POTION : -5 npc hp during %d turns`, client.time_buff)
	return 5
}

func UseStrenghPotion(client *Client, time int) int {
	UseBuff(client, "item.strengh_potion", "strengh", time)
	client.Send(`EVT STRENGH POTION : +5 player strengh during %d turns`, client.time_buff)
	return 5
}

func UseHealthPotion(client *Client, time int) int {
	UseBuff(client, "item.health_potion", "health", time)
	client.Send(`EVT HEALTH POTION : +5 player hp during %d turns`, client.time_buff)
	return 5
}
