// Package protocol provides the protocol layer for parsing and serializing.
package protocol

import (
	"bytes"
	"errors"
	"fmt"
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
	RESPTypeSimpleStr RESPType = '+'
	RESPTypeBulkStr   RESPType = '$'
	RESPTypeErr       RESPType = '-'
	RESPTypeInt       RESPType = ':'
	RESPTypeArray     RESPType = '*'
)

var (
	ErrMalformedInput  = errors.New("malformed input")
	ErrInvalidLenVal   = errors.New("invalid len val")
	ErrInvalidValCast  = errors.New("invalid val cast")
	ErrInvalidTypeCast = errors.New("invalid val type cast")
)

type RESPVal struct {
	Type RESPType
	Val  any
}

func (r *RESPVal) String() string {
	return fmt.Sprintf("&{Type:%d Val:%v}", r.Type, r.Val)
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
	res := make([]*RESPVal, 0, min(64, ln))

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
	case RESPTypeSimpleStr:
		val, err = parseSimpleStr(r)

	case RESPTypeErr:
		val, err = parseErr(r)

	case RESPTypeInt:
		val, err = parseInt(r)

	case RESPTypeBulkStr:
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

// WRITERS

func writeSimpleStr(s string, w io.Writer) error {
	_, err := w.Write([]byte("+" + s + "\r\n"))
	return err
}

func writeBulkStr(s string, w io.Writer) error {
	_, err := w.Write([]byte("$" + strconv.Itoa(len(s)) + "\r\n" + s + "\r\n"))
	return err
}

func writeErr(s string, w io.Writer) error {
	_, err := w.Write([]byte("-" + s + "\r\n"))
	return err
}

func writeInt(n int64, w io.Writer) error {
	_, err := w.Write([]byte(":" + strconv.FormatInt(n, 10) + "\r\n"))
	return err
}

func writeNil(w io.Writer) error {
	_, err := w.Write([]byte("$-1\r\n"))
	return err
}

func writeArray(s []*RESPVal, w io.Writer) error {
	if s == nil {
		return writeNil(w)
	}

	_, err := fmt.Fprintf(w, "*%d\r\n", len(s))
	if err != nil {
		return err
	}

	for _, val := range s {
		if err = Write(val, w); err != nil {
			return err
		}

	}

	return nil
}

func Write(val *RESPVal, w io.Writer) error {
	if val == nil {
		return writeNil(w)
	}

	switch val.Type {
	case RESPTypeSimpleStr:
		str, ok := val.Val.(string)
		if !ok {
			return ErrInvalidValCast
		}

		return writeSimpleStr(str, w)

	case RESPTypeBulkStr:
		str, ok := val.Val.(string)
		if !ok {
			return ErrInvalidValCast
		}

		return writeBulkStr(str, w)

	case RESPTypeInt:
		n, ok := val.Val.(int64)
		if !ok {
			return ErrInvalidValCast
		}

		return writeInt(n, w)

	case RESPTypeErr:
		str, ok := val.Val.(string)
		if !ok {
			return ErrInvalidValCast
		}

		return writeErr(str, w)

	case RESPTypeArray:
		arr, ok := val.Val.([]*RESPVal)
		if !ok {
			return ErrInvalidValCast
		}

		return writeArray(arr, w)

	default:
		return ErrInvalidTypeCast
	}
}
