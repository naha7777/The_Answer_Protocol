package server

import (
	"fmt"
	"net"
	"the_answer_protocol/game"
)

type Client struct {
	conn      net.Conn
	name      string
	group     *Group
	room      *game.Room
	inventory []*game.Item
	hp        int
	hp_max    int
	status    string
}

func (c *Client) Send(format string, a ...any) {
	fmt.Fprintf(c.conn, format+"\n", a...)
}

func (c *Client) Name() string {
	return c.name
}
