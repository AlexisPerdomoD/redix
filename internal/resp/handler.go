package resp

import (
	"io"
	"strings"

	"github.com/AlexisPerdomoD/redix/internal/protocol"
)

func handleSimpleStr(s string, w io.Writer) error {
	if s == "PING" {
		return protocol.Write(&protocol.RESPVal{
			Type: protocol.RESPTypeSimpleStr,
			Val:  "PONG",
		}, w)
	}

	return protocol.Write(&protocol.RESPVal{
		Type: protocol.RESPTypeErr,
		Val:  "ERR unknown command",
	}, w)
}

func handleBulkStr(_ string, w io.Writer) error {
	return protocol.Write(&protocol.RESPVal{
		Type: protocol.RESPTypeErr,
		Val:  "ERR unexpected bulk string",
	}, w)
}

func handleErr(_ string, w io.Writer) error {
	return protocol.Write(&protocol.RESPVal{
		Type: protocol.RESPTypeErr,
		Val:  "ERR received error",
	}, w)
}

func handleInt(_ int64, w io.Writer) error {
	return protocol.Write(&protocol.RESPVal{
		Type: protocol.RESPTypeErr,
		Val:  "ERR clients cannot send integer type",
	}, w)
}

func handleArray(s []*protocol.RESPVal, w io.Writer) error {
	if len(s) == 0 {
		return protocol.Write(&protocol.RESPVal{
			Type: protocol.RESPTypeErr,
			Val:  "ERR empty array",
		}, w)
	}

	cmdVal := s[0]
	if cmdVal == nil {
		return protocol.Write(&protocol.RESPVal{
			Type: protocol.RESPTypeErr,
			Val:  "ERR nil command",
		}, w)
	}

	cmd, ok := cmdVal.Val.(string)
	if !ok {
		return protocol.Write(&protocol.RESPVal{
			Type: protocol.RESPTypeErr,
			Val:  "ERR invalid command type value",
		}, w)
	}

	noImpl := &protocol.RESPVal{
		Type: protocol.RESPTypeErr,
		Val:  "ERR not implemented",
	}
	switch RESPCommand(strings.ToUpper(cmd)) {
	case RESPCommandSet:
		return protocol.Write(noImpl, w)
	case RESPCommandGet:
		return protocol.Write(noImpl, w)
	case RESPCommandDel:
		return protocol.Write(noImpl, w)
	case RESPCommandKeys:
		return protocol.Write(noImpl, w)
	case RESPCommandExists:
		return protocol.Write(noImpl, w)
	case RESPCommandExpire:
		return protocol.Write(noImpl, w)
	case RESPCommandTTL:
		return protocol.Write(noImpl, w)
	case RESPCommandPing:
		return pingCmd(w, s[1:])
	case RESPCommandCommand:
		return protocol.Write(noImpl, w)
	case RESPCommandCommandDocs:
		return protocol.Write(noImpl, w)
	case RESPCommandCommandInfo:
		return protocol.Write(noImpl, w)
	case RESPCommandInfoServer:
		return protocol.Write(noImpl, w)
	default:
		return protocol.Write(&protocol.RESPVal{
			Type: protocol.RESPTypeErr,
			Val:  "ERR invalid command",
		}, w)
	}
}
