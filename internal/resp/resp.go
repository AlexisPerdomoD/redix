// Package resp provides the compatibility layer for the redis protocol (RESP).
package resp

import (
	"errors"
	"io"

	"github.com/AlexisPerdomo/redix/internal/protocol"
)

type RESPCommand string

const (
	RESPCommandGet    RESPCommand = "GET"
	RESPCommandSet    RESPCommand = "SET"
	RESPCommandDel    RESPCommand = "DEL"
	RESPCommandKeys   RESPCommand = "KEYS"
	RESPCommandExists RESPCommand = "EXISTS"
	RESPCommandPing   RESPCommand = "PING"
)

var (
	ErrUnknownType    = errors.New("unknown type")
	ErrInvalidType    = errors.New("invalid type assertion")
	ErrUnknownCommand = errors.New("unknown command")
	ErrInvalidCommand = errors.New("invalid command")
)

func Handle(val *protocol.RESPVal, w io.Writer) error {
	if val == nil {
		return protocol.WrNil(w)
	}

	switch val.Type {
	case protocol.RESPTypeSimpleString:
		v, ok := val.Val.(string)
		if !ok {
			return ErrInvalidType
		}

		return handleSimpleStr(v, w)

	case protocol.RESPTypeError:
		v, ok := val.Val.(string)
		if !ok {
			return ErrInvalidType
		}
		return handleErr(v, w)

	case protocol.RESPTypeInteger:
		v, ok := val.Val.(int64)
		if !ok {
			return ErrInvalidType
		}

		return handleInt(v, w)
	case protocol.RESPTypeBulkString:
		v, ok := val.Val.(string)
		if !ok {
			return ErrInvalidType
		}

		return handleBulkStr(v, w)
	case protocol.RESPTypeArray:
		v, ok := val.Val.([]any)
		if !ok {
			return ErrInvalidType
		}

		return handleArray(v, w)
	default:
		return ErrUnknownType
	}
}
