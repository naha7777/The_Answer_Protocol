package server

import (
	"encoding/json"
	"maps"
	"slices"
	"strconv"
	"strings"
)

type groupInfo struct {
	Id      string   `json:"Id"`
	Leader  string   `json:"Leader"`
	Members []string `json:"Members"`
	Invited []string `json:"Invited"`
}

type Group struct {
	id      string
	leader  string
	members map[string]*Client
	invited map[string]bool
}

func (d *Dispatcher) handleGroupCreate(client *Client, args []string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	// if the player is not in a group, he cans create one
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
			client.Send("ERR 011 GROUP_NAME_IN_USE")
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
	// a player can join a group if he received an invitation from the leader
	// of the group
	if len(args) == 0 {
		client.Send("ERR 003 CANT_FIND_LEADER")
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
		client.Send("ERR 003 CANT_FIND_LEADER")
		return
	}
	if d.groups[leader] == nil {
		client.Send("ERR 004 THIS_USER_IS_NOT_A_LEADER")
		return
	}
	l_group := d.groups[leader]
	if l_group.leader != leader {
		client.Send("ERR 004 THIS_USER_IS_NOT_A_LEADER")
		return
	}
	if !l_group.invited[client.name] {
		client.Send("ERR 005 REQUIRE_INVITATION")
		return
	}
	SendGroup(d, client, "EVT GROUP JOIN %s", client.name)
	l_group.members[client.name] = client
	d.groups[client.name] = l_group
	client.group = l_group
	client.Send("OK group_joined=%s", l_group.id)
	SendGroup(d, client, "EVT GROUP JOIN %s", client.name)
	delete(l_group.invited, client.name)
}

func (d *Dispatcher) handleGroupLeave(client *Client, args []string) {
	// if a player is in a group, he cans leave it
	d.mu.Lock()
	defer d.mu.Unlock()
	if client.group == nil {
		client.Send("ERR 401 NOT_IN_GROUP")
		return
	}
	current_group := client.group
	delete(current_group.members, client.name)
	client.group = nil
	client.Send("OK group_leaved=%s", current_group.id)
	SendGroup(d, client, "EVT GROUP LEAVE %s", client.name)
	if len(current_group.members) > 0 {
		// if this player is the leader of the group, the new leader will be
		// the next who enterred in the group
		for _, m := range current_group.members {
			current_group.leader = m.name
			break
		}
	} else {
		// if the leader was alone in his group, the group has no leader anymore
		// and no one would be able to join it anymore
		current_group.leader = ""
	}
}

func (d *Dispatcher) handleGroupInvite(client *Client, args []string) {
	// a group leader can send an invitation to a user to join his group
	// he needs to be on a group, and a valid username of a player not in the
	// group
	if len(args) == 0 {
		client.Send("ERR 003 CANT_FIND_USER")
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if client.group == nil {
		client.Send("ERR 401 NOT_IN_GROUP")
		return
	}
	if client.name != client.group.leader {
		client.Send("ERR 022 NOT_LEADER")
		return
	}
	username := strings.Join(args, " ")
	if _, exists := d.clients[username]; !exists {
		client.Send("ERR 003 CANT_FIND_USER")
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
			c.Send("EVT GROUP INVITE %s %s", client.name, client.group.id)
		}
	}
}

func (d *Dispatcher) handleInfoGroup(client *Client, args []string) {
	// if a player is in a group, he can access to all informations about it
	// Informations : the group id, the leader, the members and the invited
	// users
	d.mu.Lock()
	defer d.mu.Unlock()
	if client.group == nil {
		client.Send("ERR 401 NOT_IN_GROUP")
		return
	}
	g := client.group
	b, err := json.Marshal(map[string]groupInfo{
		"group": {
			Id:      g.id,
			Leader:  g.leader,
			Members: slices.Sorted(maps.Keys(g.members)),
			Invited: slices.Sorted(maps.Keys(g.invited)),
		},
	})
	if err != nil {
		client.Send("ERR 500 INTERNAL")
		return
	}
	client.Send("OK group_infos=%s", b)
}
