package server

import (
	"strings"
)

type Group struct {
	id		int
	leader	string
	members	map[string]*Client
	invited	map[string]bool
}

func (d *Dispatcher) handleGroupCreate(client *Client, args []string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if client.group != nil {
		client.Send("ERR 402 ALREADY_IN_GROUP")
		return
	}
	g := &Group{
		id:			len(d.groups) + 1,
		leader:		client.name,
		members:	map[string]*Client{client.name: client},
		invited:	map[string]bool{},
	}
	d.groups[client.name] = g
	client.group = g
	client.Send("OK group=%d", g.id)
}

func (d *Dispatcher) handleGroupJoin(client *Client, args []string) {
	if len(args) == 0 {
		client.Send("ERR 104 MISSING_ARGUMENT=LEADER_NAME")
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if client.group != nil {
		client.Send("ERR 402 ALREADY_IN_GROUP")
		return
	}
	leader := strings.Join(args, " ")
	if _, exists := d.clients[leader]; !exists {
		client.Send("ERR 107 CANT_FIND_LEADER")
		return
	}
	l_group := d.groups[leader]
	if l_group.leader != leader {
		client.Send("ERR 108 THIS_USER_IS_NOT_A_LEADER")
		return
	}
	if !l_group.invited[client.name] {
		client.Send("ERR 110 REQUIRE_INVITATION")
		return
	}
	for _, c := range l_group.members {
		c.Send("EVT GROUP JOIN")
	}
	l_group.members[client.name] = client
	d.groups[client.name] = l_group
	client.group = l_group
	client.Send("OK group=%d", l_group.id)
	delete(l_group.invited, client.name)
}

func (d *Dispatcher) handleGroupLeave(client *Client, args []string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if client.group == nil {
		client.Send("ERR 401 NOT_IN_GROUP")
		return
	}
	current_group := client.group
	delete(current_group.members, client.name)
	client.group = nil
	client.Send("OK")
	for _, c := range current_group.members {
		c.Send("EVT GROUP LEAVE")
	}
	if len(current_group.members) > 0 {
		for _, m := range(current_group.members) {
			current_group.leader = m.name
			break
		}
	} else {
		current_group.leader = ""
	}
}

func (d *Dispatcher) handleGroupInvite(client *Client, args []string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	username := strings.Join(args, " ")
	if _, exists := d.clients[username]; !exists {
		client.Send("ERR 107 CANT_FIND_USER")
		return
	}
	for _, m := range(client.group.members) {
		if m.name == username {
			client.Send("ERR 402 ALREADY_IN_GROUP")
			return
		}
	}
	client.group.invited[username] = true

	client.Send("OK")
	for _, c := range d.clients {
		if c.name == username {
			c.Send("EVT GROUP INVITE")
		}
	}
}
