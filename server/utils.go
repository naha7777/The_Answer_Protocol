package server

func (d *Dispatcher) removeClient(client *Client) {
	if client.name == "" {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	client.Send("OK bye")
	for _, c := range(d.clients) {
		if c.name == client.name {
			delete(d.clients, c.name)
		}
	}
	client.name = ""
	if client.group != nil {
		delete(client.group.members, client.name)
	}
	client.group = nil
	delete(client.room.Visitors, client.name)
	client.room = nil
}
