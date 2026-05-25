package server

import (
	"context"
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
	conn  net.Conn
	state ConnState
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

	return c.conn.Read(b)
}

// TODO: evaluate if conceptualy appropriate
func (c *Connection) ReadCtx(ctx context.Context, b []byte) (int, error) {
	deadline, ok := ctx.Deadline()
	if ok {
		if err := c.conn.SetReadDeadline(deadline); err != nil {
			return 0, err
		}

		defer c.conn.SetReadDeadline(time.Time{})
	}

	return c.Read(b)
}

func (c *Connection) Write(b []byte) (int, error) {
	if c.state != ConnStateConnected {
		return 0, ErrConnIsNotConnected
	}

	return c.conn.Write(b)
}

// TODO: evaluate if conceptualy appropriate
func (c *Connection) WriteCtx(ctx context.Context, b []byte) (int, error) {
	return 0, errors.New("not implemented")
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
