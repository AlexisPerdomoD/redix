package server

import (
	"context"
	"errors"
	"net"
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

type S struct {
	listener net.Listener
	ctx      context.Context
	state    ServerState
}

func (s *S) State() ServerState {
	return s.state
}

func (s *S) Ctx() context.Context {
	return s.ctx
}

func (s *S) Close() error {
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

func (s *S) Accept() (*Connection, error) {
	conn, err := s.listener.Accept()
	if err != nil {
		return nil, err
	}

	return &Connection{
		conn:  conn,
		state: ConnStateConnected,
	}, nil
}

type ServerCfg struct {
	Port string
	Ctx  context.Context
}

// StartServer starts a server on the given configuration.
func StartServer(cfg ServerCfg) (*S, error) {
	l, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		return nil, err
	}

	return &S{
		listener: l,
		state:    ServerStateRunning,
		ctx:      cfg.Ctx,
	}, nil
}
