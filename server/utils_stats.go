package server

import "the_answer_protocol/game"

func (d *Dispatcher) displayStats() {
	// display the number of player connected
	d.mu.Lock()
	defer d.mu.Unlock()
	var targets []*Client
	for _, c := range d.clients {
		if c.name != "" && !c.dead.Load() {
			targets = append(targets, c)
		}
	}
	broadcast(targets, "EVT STATS players=%d", len(targets))
}

func (d *Dispatcher) displayRoomStats() {
	// display the number of player connected in the same room as the client
	d.mu.Lock()
	defer d.mu.Unlock()
	byRoom := map[*game.Room][]*Client{}
	for _, c := range d.clients {
		if c.name != "" && !c.dead.Load() {
			byRoom[c.room] = append(byRoom[c.room], c)
		}
	}
	for _, targets := range byRoom {
		broadcast(targets, "EVT STATS players_in_room=%d", len(targets))
	}
}
