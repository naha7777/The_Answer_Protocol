package server

import (
	"strings"
	"strconv"
)

type Group struct {
	id      string
	leader  string
	members map[string]*Client
	invited map[string]bool
}

func (d *Dispatcher) handleGroupCreate(client *Client, args []string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if client.group != nil {
		client.Send("ERR 402 ALREADY_IN_GROUP")
		return
	}
	id := strings.Join(args, " ")
	if id == "" {
		id = strconv.Itoa(len(d.groups) + 1)
	}
	for _, g := range d.groups {
		if g.id == id && len(g.members) > 0 {
			client.Send("ERR 11 GROUP_NAME_IN_USE")
			return
		}
	}
	g := &Group{
		id:      id,
		leader:  client.name,
		members: map[string]*Client{client.name: client},
		invited: map[string]bool{},
	}
	d.groups[client.name] = g
	client.group = g
	client.Send("OK group_created=%s", g.id)
}

func (d *Dispatcher) handleGroupJoin(client *Client, args []string) {
	if len(args) == 0 {
		client.Send("ERR 03 CANT_FIND_LEADER")
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
		client.Send("ERR 03 CANT_FIND_LEADER")
		return
	}
	if d.groups[leader] == nil {
		client.Send("ERR 04 THIS_USER_IS_NOT_A_LEADER")
		return
	}
	l_group := d.groups[leader]
	if l_group.leader != leader {
		client.Send("ERR 04 THIS_USER_IS_NOT_A_LEADER")
		return
	}
	if !l_group.invited[client.name] {
		client.Send("ERR 05 REQUIRE_INVITATION")
		return
	}
	for _, c := range l_group.members {
		c.Send("EVT GROUP JOIN %s", client.name)
	}
	l_group.members[client.name] = client
	d.groups[client.name] = l_group
	client.group = l_group
	client.Send("OK group_joined=%s", l_group.id)
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
	client.Send("OK group_leaved=%s", client.group.id)
	client.group = nil
	for _, c := range current_group.members {
		c.Send("EVT GROUP LEAVE %s", client.name)
	}
	if len(current_group.members) > 0 {
		for _, m := range current_group.members {
			current_group.leader = m.name
			break
		}
	} else {
		current_group.leader = ""
	}
}

func (d *Dispatcher) handleGroupInvite(client *Client, args []string) {
	if len(args) == 0 {
		client.Send("ERR 03 CANT_FIND_USER")
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if client.group == nil {
		client.Send("ERR 401 NOT_IN_GROUP")
		return
	}
	username := strings.Join(args, " ")
	if _, exists := d.clients[username]; !exists {
		client.Send("ERR 03 CANT_FIND_USER")
		return
	}
	for _, m := range client.group.members {
		if m.name == username {
			client.Send("ERR 402 ALREADY_IN_GROUP")
			return
		}
	}
	client.group.invited[username] = true

	client.Send("OK group_invited=%s", username)
	for _, c := range d.clients {
		if c.name == username {
			c.Send("EVT GROUP INVITE %s", client.name)
		}
	}
}
