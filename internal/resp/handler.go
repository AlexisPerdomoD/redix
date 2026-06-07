package resp

import (
	"io"
	"strings"

	"github.com/AlexisPerdomoD/redix/internal/protocol"
)

func handleSimpleStr(s string, w io.Writer) error {
	if s == "PING" {
		return protocol.FormatWrite(&protocol.RESPVal{
			Type: protocol.RESPTypeSimpleStr,
			Val:  "PONG",
		}, w)
	}

	return protocol.FormatWrite(&protocol.RESPVal{
		Type: protocol.RESPTypeErr,
		Val:  "ERR unknown command",
	}, w)
}

func handleBulkStr(_ string, w io.Writer) error {
	return protocol.FormatWrite(&protocol.RESPVal{
		Type: protocol.RESPTypeErr,
		Val:  "ERR unexpected bulk string",
	}, w)
}

func handleErr(_ string, w io.Writer) error {
	return protocol.FormatWrite(&protocol.RESPVal{
		Type: protocol.RESPTypeErr,
		Val:  "ERR received error",
	}, w)
}

func handleInt(_ int64, w io.Writer) error {
	return protocol.FormatWrite(&protocol.RESPVal{
		Type: protocol.RESPTypeErr,
		Val:  "ERR clients cannot send integer type",
	}, w)
}

func handleArray(s []*protocol.RESPVal, w io.Writer) error {
	if len(s) == 0 {
		return protocol.FormatWrite(&protocol.RESPVal{
			Type: protocol.RESPTypeErr,
			Val:  "ERR empty array",
		}, w)
	}

	cmdVal := s[0]
	if cmdVal == nil {
		return protocol.FormatWrite(&protocol.RESPVal{
			Type: protocol.RESPTypeErr,
			Val:  "ERR nil command",
		}, w)
	}

	cmd, ok := cmdVal.Val.(string)
	if !ok {
		return protocol.FormatWrite(&protocol.RESPVal{
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
		return setCmd(w, s[1:])
	case RESPCommandHSet:
		return hsetCmd(w, s[1:])
	case RESPCommandGet:
		return getCmd(w, s[1:])
	case RESPCommandHGet:
		return hgetCmd(w, s[1:])
	case RESPCommandDel:
		return delCmd(w, s[1:])
	case RESPCommandHDel:
		return hdelCmd(w, s[1:])
	case RESPCommandKeys:
		return protocol.FormatWrite(noImpl, w)
	case RESPCommandExists:
		return existsCmd(w, s[1:])
	case RESPCommandExpire:
		return expireCmd(w, s[1:])
	case RESPCommandTTL:
		return ttlCmd(w, s[1:])
	case RESPCommandPing:
		return pingCmd(w, s[1:])
	case RESPCommandCommand:
		return protocol.FormatWrite(noImpl, w)
	case RESPCommandCommandDocs:
		return protocol.FormatWrite(noImpl, w)
	case RESPCommandCommandInfo:
		return protocol.FormatWrite(noImpl, w)
	case RESPCommandInfoServer:
		return protocol.FormatWrite(noImpl, w)
	default:
		return protocol.FormatWrite(&protocol.RESPVal{
			Type: protocol.RESPTypeErr,
			Val:  "ERR invalid command",
		}, w)
	}
}
