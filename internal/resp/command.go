package resp

import (
	"io"

	"github.com/AlexisPerdomoD/redix/internal/protocol"
)

// pingCmd is a simple command that writes "PONG" to the writer
// or the first argument if it is provided.
func pingCmd(w io.Writer, b []*protocol.RESPVal) error {
	if len(b) > 0 {
		return protocol.Write(b[0], w)
	}

	return protocol.Write(&protocol.RESPVal{
		Type: protocol.RESPTypeSimpleStr,
		Val:  "PONG",
	}, w)
}
