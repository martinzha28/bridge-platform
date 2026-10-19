package ws

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

const defaultLeaveGrace = 5 * time.Second

// member is a connected client's presence at one table.
type member struct {
	id   int    // per-table, in join order
	name string // unique at this table (may carry a " (2)" suffix)
}

// PersonView is one connected client as shown to the table.
type PersonView struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Guest bool   `json:"guest"`
	Seat  string `json:"seat"` // "" while watching
}

// AddObserver registers a client as watching this table, cancels any
// pending reap, and announces the arrival unless it is the same person
// reconnecting. Safe to call more than once for the same client.
func (t *Table) AddObserver(c *Client) PersonView {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.reapTimer != nil {
		t.reapTimer.Stop()
		t.reapTimer = nil
	}
	if m, ok := t.observers[c]; ok {
		return t.personView(c, m)
	}

	// The creator's arrival is not news; nobody else is here yet.
	first := t.nextMemberID == 0
	returning := t.identityPresent(c.identity)
	if timer, ok := t.pendingLeaves[c.identity]; ok {
		timer.Stop()
		delete(t.pendingLeaves, c.identity)
		returning = true
	}

	t.nextMemberID++
	m := &member{id: t.nextMemberID, name: t.uniqueName(c)}
	t.observers[c] = m

	// The newcomer gets this line in the chat history sent after joining.
	if !first && !returning && !t.closed {
		t.broadcastChatExcept(t.appendSystem(m.name+" joined"), c)
	}
	return t.personView(c, m)
}

// RemoveObserver drops a client and, unless the same person is still
// connected, posts "<name> left" after leaveGrace — cancelled if they
// come back first.
func (t *Table) RemoveObserver(c *Client) {
	t.mu.Lock()
	defer t.mu.Unlock()
	m, ok := t.observers[c]
	if !ok {
		return
	}
	delete(t.observers, c)
	if t.closed || t.identityPresent(c.identity) {
		return
	}

	id, name := c.identity, m.name
	if old, ok := t.pendingLeaves[id]; ok {
		old.Stop()
	}
	var timer *time.Timer
	timer = time.AfterFunc(t.leaveGrace, func() {
		t.mu.Lock()
		defer t.mu.Unlock()
		if t.closed || t.pendingLeaves[id] != timer {
			return
		}
		delete(t.pendingLeaves, id)
		t.broadcastChat(t.appendSystem(name + " left"))
	})
	t.pendingLeaves[id] = timer
}

// identityPresent reports whether a client with this identity is
// connected. Must be called with t.mu held.
func (t *Table) identityPresent(identity string) bool {
	for other := range t.observers {
		if other.identity == identity {
			return true
		}
	}
	return false
}

// uniqueName returns c.name, with a " (n)" suffix if another person at
// this table holds it. The same person's other sockets don't count, so a
// reconnect keeps its name. Must be called with t.mu held.
func (t *Table) uniqueName(c *Client) string {
	taken := func(name string) bool {
		for other, m := range t.observers {
			if other.identity != c.identity && strings.EqualFold(m.name, name) {
				return true
			}
		}
		return false
	}
	name := c.name
	for n := 2; taken(name); n++ {
		name = fmt.Sprintf("%s (%d)", c.name, n)
	}
	return name
}

// seatOf returns the seat c holds, or "". Must be called with t.mu held.
func (t *Table) seatOf(c *Client) string {
	for d, p := range t.players {
		if p == Player(c) {
			return d.String()
		}
	}
	return ""
}

// personView must be called with t.mu held.
func (t *Table) personView(c *Client, m *member) PersonView {
	return PersonView{ID: m.id, Name: m.name, Guest: c.guest, Seat: t.seatOf(c)}
}

// people lists every connected client in join order. Must be called
// with t.mu held.
func (t *Table) people() []PersonView {
	out := make([]PersonView, 0, len(t.observers))
	for c, m := range t.observers {
		out = append(out, t.personView(c, m))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
