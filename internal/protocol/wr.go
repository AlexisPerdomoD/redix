package protocol

import (
	"io"
	"strconv"
)

func WrSimpleStr(s string, w io.Writer) error {
	_, err := w.Write([]byte("+" + s + "\r\n"))
	return err
}

func WrBulkStr(s string, w io.Writer) error {
	_, err := w.Write([]byte("$" + strconv.Itoa(len(s)) + "\r\n" + s + "\r\n"))
	return err
}

func WrErr(s string, w io.Writer) error {
	_, err := w.Write([]byte("-" + s + "\r\n"))
	return err
}

func WrInt(n int64, w io.Writer) error {
	_, err := w.Write([]byte(":" + strconv.FormatInt(n, 10) + "\r\n"))
	return err
}

func WrNil(w io.Writer) error {
	_, err := w.Write([]byte("$-1\r\n"))
	return err
}
