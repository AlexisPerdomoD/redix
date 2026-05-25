package main

import (
	"context"
	"io"
	"log/slog"

	"github.com/AlexisPerdomo/redix/internal/server"
)

// “What I cannot create, I do not understand” - Richard Feynman.

func main() {
	ctx := context.Background()

	cfg := server.ServerCfg{
		Port: "6379",
		Ctx:  ctx,
	}

	s, err := server.StartServer(cfg)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	for {
		c, err := s.Accept()
		if err != nil {
			slog.ErrorContext(s.Ctx(), "error accepting connection", "err", err)
			continue
		}

		buf := make([]byte, 1024)
		_, err = c.Read(buf)
		if err != nil {
			if err == io.EOF {
				slog.InfoContext(s.Ctx(), "connection closed")
				continue
			}

			slog.ErrorContext(s.Ctx(), "error reading from connection", "err", err)
			continue
		}

		c.Write(buf)
	}
}
