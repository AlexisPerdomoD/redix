// Package protocol provides the protocol layer for parsing and serializing.
package protocol

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
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

func NilBulkStr() *RESPVal {
	return &RESPVal{Type: RESPTypeBulkStr, Val: nil}
}

func NilArray() *RESPVal {
	return &RESPVal{Type: RESPTypeArray, Val: nil}
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

func formatSimpleStr(s string) string {
	return "+" + s + CRLF
}

func formatBulkStr(s string) string {
	return "$" + strconv.Itoa(len(s)) + CRLF + s + CRLF
}

func formatErr(s string) string {
	return "-" + s + CRLF
}

func formatInt(n int64) string {
	return ":" + strconv.FormatInt(n, 10) + CRLF
}

func formatBulkStrNil() string {
	return "$-1" + CRLF
}

func formatArrayNil() string {
	return "*-1" + CRLF
}

func Format(val *RESPVal) (string, error) {
	if val == nil {
		return "", ErrInvalidValCast
	}

	switch val.Type {
	case RESPTypeSimpleStr:
		str, ok := val.Val.(string)
		if !ok {
			return "", ErrInvalidValCast
		}

		return formatSimpleStr(str), nil

	case RESPTypeBulkStr:
		if val.Val == nil {
			return formatBulkStrNil(), nil
		}

		str, ok := val.Val.(string)
		if !ok {
			return "", ErrInvalidValCast
		}

		return formatBulkStr(str), nil

	case RESPTypeInt:
		n, ok := val.Val.(int64)
		if !ok {
			return "", ErrInvalidValCast
		}

		return formatInt(n), nil

	case RESPTypeErr:
		str, ok := val.Val.(string)
		if !ok {
			return "", ErrInvalidValCast
		}

		return formatErr(str), nil

	case RESPTypeArray:
		if val.Val == nil {
			return formatArrayNil(), nil
		}

		arr, ok := val.Val.([]*RESPVal)
		if !ok {
			return "", ErrInvalidValCast
		}

		if len(arr) == 0 {
			return "*0\r\n", nil
		}

		sb := strings.Builder{}
		fmt.Fprintf(&sb, "*%d\r\n", len(arr))

		for _, val := range arr {
			el, err := Format(val)
			if err != nil {
				return "", err
			}
			sb.WriteString(el)
		}

		return sb.String(), nil

	default:
		return "", ErrInvalidTypeCast
	}
}

// WRITER
func FormatWrite(val *RESPVal, w io.Writer) error {
	res, err := Format(val)
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, res)
	return err
}
