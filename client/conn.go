package client

import (
	"bufio"
	"net"
	"strings"
	"sync"
	"time"
	"fmt"
	"os"
	"errors"
	"sync/atomic"
)

type Message struct {
	Kind	string
	Raw		string
	Rest	string
}

type Conn struct {
	conn      net.Conn
	mu        sync.Mutex
	messages  chan Message
	done      chan struct{}
	closeOnce sync.Once
	closed    atomic.Bool
}

func Dial(addr string) (*Conn, error) {
	conn, err := net.DialTimeout("tcp", addr, 5 * time.Second)
	if err != nil {
		return nil, err
	}

	c := &Conn{conn: conn, messages: make(chan Message, 64), done: make(chan struct{}),}

	go scan(c)
	return c, nil
}

func scan(c *Conn) {
	defer close(c.messages)
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		raw := scanner.Text()
		kind, rest, _ := strings.Cut(raw, " ")
		msg := Message{Kind: kind, Raw: raw, Rest: rest}
		select {
		case c.messages <- msg:
		case <-c.done:
			return
		}
	}
	if err := scanner.Err(); err != nil &&
		!c.closed.Load() && !errors.Is(err, net.ErrClosed) {
		fmt.Fprintln(os.Stderr, err)
	}
}

func (c *Conn) Send(line string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	line += "\n"
	msg := []byte(line)
	if _, err := c.conn.Write(msg); err != nil {
		return err
	}
	return nil
}

func (c *Conn) Messages() <- chan Message{
	return c.messages
}

func (c *Conn) Done() <-chan struct{} {
	return c.done
}

func (c *Conn) Close() error {
	var err error
	c.closeOnce.Do(func() {
		c.closed.Store(true)
		close(c.done)
		err = c.conn.Close()
	})
	return err
}
