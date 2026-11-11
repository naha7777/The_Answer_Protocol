package server

func (d *Dispatcher) removeClient(client *Client) {
	if client.name == "" {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	client.Send("OK bye")
	delete(d.clients, client.name)
	client.name = ""
	if client.group != nil {
		delete(client.group.members, client.name)
		client.group = nil
	}
	if client.inventory != nil {
		for _, item := range client.inventory {
			client.room.Items = append(client.room.Items, item.Id)
		}
		client.inventory = nil
	}
	if client.room != nil {
		remove(client.name, client.room.Players)
		client.room = nil
	}
}

func remove(name string, place []string) []string {
	i := 0
	for idx, item := range place {
		if item != name {
			place[i] = place[idx]
			i++
		}
	}
	return place[:i]
}
