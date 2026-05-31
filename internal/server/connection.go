package server

import (
	"errors"
	"net"
	"time"
)

type ConnState string

const (
	ConnStateConnected ConnState = "connected"
	ConnStateClosing   ConnState = "closing"
	ConnStateClosed    ConnState = "closed"
	ConnStateError     ConnState = "error"
)

var (
	ErrConnIsNotConnected = errors.New("connection is not connected")
)

type Connection struct {
	conn                net.Conn
	state               ConnState
	idleTimeoutDuration *time.Duration
	// TODO: evaluate concurrency safety mechanism here
}

func (c *Connection) State() ConnState {
	return c.state
}

func (c *Connection) Read(b []byte) (int, error) {
	if c.state != ConnStateConnected {
		return 0, ErrConnIsNotConnected
	}

	if len(b) == 0 {
		return 0, nil
	}

	c.refreshDeadline()
	return c.conn.Read(b)
}

func (c *Connection) Write(b []byte) (int, error) {
	if c.state != ConnStateConnected {
		return 0, ErrConnIsNotConnected
	}

	return c.conn.Write(b)
}

func (c *Connection) SetIdleTimeout(d time.Duration) {
	c.idleTimeoutDuration = &d
}

func (c *Connection) refreshDeadline() {
	if c.idleTimeoutDuration == nil {
		return
	}

	c.conn.SetDeadline(time.Now().Add(*c.idleTimeoutDuration))
}

func (c *Connection) Close() error {
	// TODO: consider concurrency safety here
	c.state = ConnStateClosing
	if err := c.conn.Close(); err != nil {
		// TODO: consider concurrency safety here
		c.state = ConnStateError
		return err
	}

	// TODO: consider concurrency safety here
	c.state = ConnStateClosed
	return nil
}
