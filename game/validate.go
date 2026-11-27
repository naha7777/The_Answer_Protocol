package game

import (
	"fmt"
	"sort"
	"strings"
)

// check the consistency of the references in the loaded world
func (w *World) Validate() (errs, warns []string) {
	bad := func(f string, a ...any) {
		errs = append(errs, fmt.Sprintf(f, a...))
	}
	warn := func(f string, a ...any) {
		warns = append(warns, fmt.Sprintf(f, a...))
	}

	// item names
	itemNames := map[string]bool{}
	for key, it := range w.Items {
		if it == nil {
			bad("item %q is empty", key)
			continue
		}
		if it.Id != key {
			bad("item key %q does not match its id %q", key, it.Id)
		}
		itemNames[strings.ToLower(it.Name)] = true
	}

	// entry points
	if w.Rooms[w.Start] == nil {
		bad("start room %q does not exist", w.Start)
	}
	if w.Rooms[w.Respawn] == nil {
		bad("respawn room %q does not exist", w.Respawn)
	}

	// rooms
	for key, r := range w.Rooms {
		if r == nil {
			bad("room %q is empty", key)
			continue
		}
		if r.Id != key {
			bad("room key %q does not match its id %q", key, r.Id)
		}
		for dir, dest := range r.Exits {
			if w.Rooms[dest] == nil {
				bad("room %q: exit %q leads to unknown room %q", key, dir, dest)
			}
		}
		for _, id := range r.Items {
			if w.Items[id] == nil {
				bad("room %q: unknown item %q", key, id)
			}
		}
		for _, id := range r.NPCs {
			if w.NPCs[id] == nil {
				bad("room %q: unknown npc %q", key, id)
			}
		}
	}

	// NPC
	for id, n := range w.NPCs {
		if n == nil {
			bad("npc %q is empty", id)
			continue
		}
		if n.Hostile && (n.Hp <= 0 || n.Max_hp < n.Hp) {
			bad("npc %q: invalid hp (%d/%d)", id, n.Hp, n.Max_hp)
		}
		if len(n.Dialogue) == 0 {
			warn("npc %q has no dialogue", id)
		}
		for _, qid := range n.Quest {
			q := w.Quests[qid]
			if q == nil {
				bad("npc %q: unknown quest %q", id, qid)
			} else if q.NPC != id {
				warn("npc %q lists quest %q, which belongs to %q", id, qid,
					q.NPC)
			}
		}
	}

	// quests
	for key, q := range w.Quests {
		if q == nil {
			bad("quest %q is empty", key)
			continue
		}
		if q.Id != key {
			bad("quest key %q does not match its id %q", key, q.Id)
		}
		if w.NPCs[q.NPC] == nil {
			bad("quest %q: unknown giver npc %q", key, q.NPC)
		}
		if q.Status != "available" {
			warn("quest %q starts with status %q (expected \"available\")", key,
				q.Status)
		}
		if len(q.Objective) == 0 {
			bad("quest %q has no objective", key)
		}
		for _, o := range q.Objective {
			if w.Items[o] == nil && w.NPCs[o] == nil &&
				!itemNames[strings.ToLower(o)] {
				bad("quest %q: objective %q is neither an item nor an npc", key,
					o)
			}
		}
		for _, rw := range q.Reward {
			if w.Items[rw] == nil {
				bad("quest %q: unknown reward item %q", key, rw)
			}
		}
	}

	// start room accessibility
	seen := map[string]bool{}
	queue := []string{w.Start}
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		r := w.Rooms[id]
		if seen[id] || r == nil {
			continue
		}
		seen[id] = true
		for _, dest := range r.Exits {
			queue = append(queue, dest)
		}
	}
	for id := range w.Rooms {
		if !seen[id] {
			warn("room %q is not reachable from the start room", id)
		}
	}

	sort.Strings(errs)
	sort.Strings(warns)
	return errs, warns
}
