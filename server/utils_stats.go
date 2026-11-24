package server

func (d *Dispatcher) displayStats() {
	// Send an EVT with the number of players connected to the server
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, c := range d.clients {
		c.Send("EVT STATS players=%d", len(d.clients))
	}
}

func (d *Dispatcher) displayRoomStats() {
	// // Send an EVT with the number of players in the client room
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, c := range d.clients {
		for _, str := range c.room.Players {
			if c.name == str {
				c.Send("EVT STATS players_in_room=%d", len(c.room.Players))
			}
		}
	}
}
