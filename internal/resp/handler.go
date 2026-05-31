package resp

import (
	"io"
	"strings"

	"github.com/AlexisPerdomo/redix/internal/protocol"
)

func handleSimpleStr(s string, w io.Writer) error {
	if s == "PING" {
		return protocol.WrSimpleStr("PONG", w)
	}

	return protocol.WrErr("ERR unknown command", w)
}

func handleBulkStr(_ string, w io.Writer) error {
	return protocol.WrErr("ERR unexpected bulk string", w)
}

func handleErr(_ string, w io.Writer) error {
	return protocol.WrErr("ERR received error", w)
}

func handleInt(_ int64, w io.Writer) error {
	return protocol.WrErr("ERR clients cannot send integer type", w)
}

func handleArray(s []any, w io.Writer) error {
	if len(s) == 0 {
		return protocol.WrErr("ERR empty array", w)
	}

	cmdVal, ok := s[0].(*protocol.RESPVal)
	if !ok {
		return protocol.WrErr("ERR invalid command format", w)
	}

	cmd, ok := cmdVal.Val.(string)
	if !ok {
		return protocol.WrErr("ERR invalid command type value", w)
	}

	switch RESPCommand(strings.ToUpper(cmd)) {
	case RESPCommandGet:
		return protocol.WrErr("ERR not implemented", w)
	case RESPCommandSet:
		return protocol.WrErr("ERR not implemented", w)
	case RESPCommandDel:
		return protocol.WrErr("ERR not implemented", w)
	case RESPCommandKeys:
		return protocol.WrErr("ERR not implemented", w)
	case RESPCommandExists:
		return protocol.WrErr("ERR not implemented", w)
	case RESPCommandPing:
		return protocol.WrSimpleStr("PONG", w)
	default:
		return protocol.WrErr("ERR unknown command", w)

	}
}
