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
	// try TCP connection, stop after waiting for 5 secondes
	conn, err := net.DialTimeout("tcp", addr, 5 * time.Second)
	if err != nil {
		return nil, err
	}

	c := &Conn{
		conn: 		conn, 
		messages: 	make(chan Message, 64), // until 64 messages in wait
		done: 		make(chan struct{}),
	}

	// reading goroutine
	go scan(c)
	return c, nil
}

func scan(c *Conn) {
	// when scan is finished, we close the channel
	defer close(c.messages)
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		raw := scanner.Text()
		// cut the line and stock it into Message struct
		kind, rest, _ := strings.Cut(raw, " ")
		msg := Message{Kind: kind, Raw: raw, Rest: rest}
		// send the channel
		select {
		// normal case : send the msg
		case c.messages <- msg:
		// if Close() : we leave
		case <-c.done:
			return
		}
	}
	// if it closes because of an error, we print it on stderr
	if err := scanner.Err(); err != nil &&
		!c.closed.Load() && !errors.Is(err, net.ErrClosed) {
		fmt.Fprintln(os.Stderr, err)
	}
}

const writeTimeout = 5 * time.Second

func (c *Conn) Send(line string) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    if err := c.conn.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
        return err
    }

    if _, err := c.conn.Write([]byte(line + "\n")); err != nil {
        c.Close()
        return err
    }
    return nil
}

func (c *Conn) Messages() <- chan Message{
	// return the channel in a reading only mode
	return c.messages
}

func (c *Conn) Done() <-chan struct{} {
	// can wait the closure with <-c.Done()
	return c.done
}

func (c *Conn) Close() error {
	var err error
	// Once : only one execution even if Close() is called more times
	c.closeOnce.Do(func() {
		// close
		c.closed.Store(true)
		// wake up scan() if blocked
		close(c.done)
		// close the socket
		err = c.conn.Close()
	})
	return err
}
