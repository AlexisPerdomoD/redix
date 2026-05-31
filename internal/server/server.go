package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"time"
)

var (
	ErrServerHasBeenClosed = errors.New("serverHasBeenClosed")
	ErrServerIsNotRunning  = errors.New("serverIsNotRunning")
)

type ServerState string

const (
	ServerStateRunning  ServerState = "running"
	ServerStateStopping ServerState = "stopping"
	ServerStateError    ServerState = "error"
	ServerStateStopped  ServerState = "stopped"
)

type ServerCfg struct {
	Port                  string
	ConnectionIdleTimeout *time.Duration
}

type Server struct {
	listener net.Listener
	state    ServerState
	cfg      ServerCfg
}

func (s *Server) State() ServerState {
	return s.state
}

func (s *Server) Close() error {
	// TODO: consider concurrency safety here
	s.state = ServerStateStopping

	err := s.listener.Close()
	if err != nil {
		// TODO: consider concurrency safety here
		s.state = ServerStateError
		return err
	}

	// TODO: consider concurrency safety here
	s.state = ServerStateStopped
	return nil
}

func (s *Server) accept() (*Connection, error) {
	if s.state != ServerStateRunning {
		return nil, ErrServerIsNotRunning
	}

	conn, err := s.listener.Accept()
	if err != nil {
		return nil, err
	}

	return &Connection{
		conn:  conn,
		state: ConnStateConnected,
	}, nil
}

func (s *Server) Serve(ctx context.Context, h func(context.Context, *Connection)) {
	go func() {
		<-ctx.Done()
		s.Close()
	}()

	for {
		c, err := s.accept()
		if err != nil {
			if err == ErrServerHasBeenClosed || err == ErrServerIsNotRunning {
				return
			}

			slog.ErrorContext(ctx, "error accepting connection", "err", err)
			continue
		}

		// TODO: evaluate needed concurrency safety mechanism here
		// go h(ctx, c)
		if s.cfg.ConnectionIdleTimeout != nil {
			c.SetIdleTimeout(*s.cfg.ConnectionIdleTimeout)
		}

		h(ctx, c)
	}
}

// StartServer starts a server on the given configuration.
func StartServer(cfg ServerCfg) (*Server, error) {
	l, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		return nil, err
	}

	return &Server{
		listener: l,
		state:    ServerStateRunning,
	}, nil
}
