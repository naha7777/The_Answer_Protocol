package server

func (d *Dispatcher) removeClient(client *Client) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if client.name == "" {
		return
	}
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
		client.room.Players = remove(client.name, client.room.Players)
		client.room = nil
	}
	delete(d.clients, client.name)
	client.name = ""
	client.Send("OK bye")
}

func remove(name string, place []string) []string {
	var new []string
	for _, item := range place {
		if item != name {
			new = append(new, item)
		}
	}
	return new
}
