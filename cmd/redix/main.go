package main

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log/slog"
	"net"

	"github.com/AlexisPerdomo/redix/internal/config"
	"github.com/AlexisPerdomo/redix/internal/protocol"
	"github.com/AlexisPerdomo/redix/internal/resp"
	"github.com/AlexisPerdomo/redix/internal/server"
)

// “What I cannot create, I do not understand” - Richard Feynman.

func main() {
	ctx := context.Background()
	cfg := server.ServerCfg{
		Port:                  config.GetPort(),
		ConnectionIdleTimeout: config.GetConnectionIdleTimeout(),
	}

	s, err := server.StartServer(cfg)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	s.Serve(ctx, handleConnection)
}

func handleConnection(ctx context.Context, c *server.Connection) {
	if c.State() != server.ConnStateConnected {
		slog.WarnContext(ctx, "connection is not connected")
		return
	}

	// TODO: consider concurrency safety here
	go func() {
		<-ctx.Done()
		c.Close()
	}()

	// TODO: evaluate if needed specifict dinamic buf reader / writer implementation
	// since bufio.Reader / Writer can be memory consuming as default implementation
	r := bufio.NewReader(c)
	w := bufio.NewWriter(c)
	for {

		val, err := protocol.Parse(r)
		if err != nil {
			if err == io.EOF || errors.Is(err, net.ErrClosed) {
				slog.DebugContext(ctx, "connection closed")
				return
			}

			slog.WarnContext(ctx, "error parsing", "err", err)
			return
		}

		if err := resp.Handle(val, w); err != nil {
			slog.WarnContext(ctx, "error handling", "err", err)
			return
		}

		if err := w.Flush(); err != nil {
			slog.WarnContext(ctx, "error flushing", "err", err)
			return
		}
	}
}
