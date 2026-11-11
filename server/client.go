package server

import (
	"net"
	"the_answer_protocol/game"
	"fmt"
)

type Client struct {
	conn		net.Conn
	name		string
	group		*Group
	room		*game.Room
	inventory	[]*game.Item
}

func (c *Client) Send(format string, a ...any) {
	fmt.Fprintf(c.conn, format+"\n", a...)
}

func (c *Client) Name() string {
	return c.name
}
