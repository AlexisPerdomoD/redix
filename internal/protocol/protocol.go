// Package protocol provides the protocol layer for parsing and serializing.
package protocol

import (
	"bytes"
	"errors"
	"io"
)

type RESPType byte

const (
	RESPTypeSimpleString RESPType = '+'
	RESPTypeError        RESPType = '-'
	RESPTypeInteger      RESPType = ':'
	RESPTypeBulkString   RESPType = '$'
	RESPTypeArray        RESPType = '*'
)

type RESPVal struct {
	Type RESPType
	Val  any
}

type LineReader interface {
	io.Reader

	// ReadByte reads and returns a single byte.
	// If no byte is available, returns an error.
	ReadByte() (byte, error)

	// ReadBytes reads until the first occurrence of delim in the input. Does not return the delimiter.
	ReadBytes(delim byte) ([]byte, error)

	// ReadString reads until the first occurrence of delim in the input. Returns the data read before the delimiter.
	ReadString(delim byte) (string, error)
}

type LineWriter interface {
	io.Writer

	// Flush writes any buffered data to the underlying io.Writer.
	Flush() error
}

func ParseString(r LineReader) (string, error) {
	b, err := r.ReadBytes('\n')
	if err != nil {
		return "", err
	}
	return string(bytes.TrimRight(b, "\r")), nil
}

func ParseInteger(r LineReader) (int64, error) {
	return 0, errors.New("not implemented")
}

func ParseBulkString(r LineReader) (string, error) {
	return "", errors.New("not implemented")
}

func ParseArray(r LineReader) ([]any, error) {
	return nil, errors.New("not implemented")
}

func ParseError(r LineReader) (string, error) {
	return "", errors.New("not implemented")
}

func Parse(r LineReader) (*RESPVal, error) {
	t, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	rt := RESPType(t)
	var val any

	switch rt {
	case RESPTypeSimpleString:
		val, err = ParseString(r)

	case RESPTypeError:
		val, err = ParseError(r)

	case RESPTypeInteger:
		val, err = ParseInteger(r)

	case RESPTypeBulkString:
		val, err = ParseBulkString(r)

	case RESPTypeArray:
		val, err = ParseArray(r)
	default:
		return nil, errors.New("unknown type")
	}

	if err != nil {
		return nil, err
	}

	return &RESPVal{
		Type: rt,
		Val:  val,
	}, nil

}
