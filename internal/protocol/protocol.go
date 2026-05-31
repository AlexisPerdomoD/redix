// Package protocol provides the protocol layer for parsing and serializing.
package protocol

import (
	"bytes"
	"errors"
	"io"
	"strconv"
)

type RESPType byte

const CR = '\r'
const LF = '\n'
const CRLF = "\r\n"

const MaxBulkStrLen = 512 << 20
const MaxArrayLen = 1 << 32

const (
	RESPTypeSimpleString RESPType = '+'
	RESPTypeError        RESPType = '-'
	RESPTypeInteger      RESPType = ':'
	RESPTypeBulkString   RESPType = '$'
	RESPTypeArray        RESPType = '*'
)

var (
	ErrMalformedInput = errors.New("malformed input")
	ErrInvalidLenVal  = errors.New("invalid len val")
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

	// ReadBytes reads until the first occurrence of delim in the input. Returns a slice containing the data up to and including the delimiter.
	ReadBytes(delim byte) ([]byte, error)

	// ReadString reads until the first occurrence of delim in the input. Returns a string containing the data up to and including the delimiter.
	ReadString(delim byte) (string, error)
}

func parseSimpleStr(r LineReader) (string, error) {
	b, err := r.ReadBytes(LF)
	if err != nil {
		return "", err
	}
	return string(bytes.TrimRight(b, CRLF)), nil
}

func parseInt(r LineReader) (int64, error) {
	b, err := r.ReadBytes(LF)
	if err != nil {
		return 0, err
	}

	strVal := string(bytes.TrimRight(b, CRLF))
	val, err := strconv.ParseInt(strVal, 10, 64)
	if err != nil {
		return 0, err
	}

	return val, nil
}

func parseBulkStr(r LineReader) (string, error) {
	b, err := r.ReadBytes(LF)
	if err != nil {
		return "", err
	}

	strLen := string(bytes.TrimRight(b, CRLF))
	ln, err := strconv.ParseInt(strLen, 10, 64)
	if err != nil {
		return "", err
	}

	if ln < -1 || ln > MaxBulkStrLen {
		return "", ErrInvalidLenVal
	}

	// no second line to read case
	if ln == -1 {
		return "", nil
	}

	strBuf := make([]byte, ln+2)
	_, err = io.ReadFull(r, strBuf)
	if err != nil {
		return "", err
	}

	return string(strBuf[:ln]), nil
}

func parseArray(r LineReader) ([]*RESPVal, error) {
	b, err := r.ReadBytes(LF)
	if err != nil {
		return nil, err
	}

	ln, err := strconv.ParseInt(string(bytes.TrimRight(b, CRLF)), 10, 64)
	if err != nil {
		return nil, err
	}

	if ln < -1 || ln > MaxArrayLen {
		return nil, ErrInvalidLenVal
	}

	if ln == -1 {
		return nil, nil
	}

	// minimal initial len
	res := make([]*RESPVal, min(64, ln))

	for range ln {
		val, err := Parse(r)
		if err != nil {
			return nil, err
		}

		res = append(res, val)
	}

	return res, nil
}

func parseErr(r LineReader) (string, error) {
	return parseSimpleStr(r)
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
		val, err = parseSimpleStr(r)

	case RESPTypeError:
		val, err = parseErr(r)

	case RESPTypeInteger:
		val, err = parseInt(r)

	case RESPTypeBulkString:
		val, err = parseBulkStr(r)

	case RESPTypeArray:
		val, err = parseArray(r)
	default:
		return nil, ErrMalformedInput
	}

	if err != nil {
		return nil, err
	}

	return &RESPVal{
		Type: rt,
		Val:  val,
	}, nil

}
