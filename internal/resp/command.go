package resp

import (
	"io"

	"github.com/AlexisPerdomo/redix/internal/protocol"
)

func pingCmd(w io.Writer) error {
	return protocol.Write(&protocol.RESPVal{
		Type: protocol.RESPTypeSimpleStr,
		Val:  "PONG",
	}, w)
}
