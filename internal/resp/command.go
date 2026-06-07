package resp

import (
	"io"

	"github.com/AlexisPerdomoD/redix/internal/protocol"
	"github.com/AlexisPerdomoD/redix/internal/store"
)

// pingCmd is a simple command that writes "PONG" to the writer
// or the first argument if it is provided.
func pingCmd(w io.Writer, b []*protocol.RESPVal) error {
	if len(b) > 0 {
		return protocol.FormatWrite(b[0], w)
	}

	return protocol.FormatWrite(
		&protocol.RESPVal{
			Type: protocol.RESPTypeSimpleStr,
			Val:  "PONG",
		}, w)
}

func setCmd(w io.Writer, b []*protocol.RESPVal) error {
	if len(b) != 2 {
		return protocol.FormatWrite(
			&protocol.RESPVal{
				Type: protocol.RESPTypeErr,
				Val:  "ERR wrong number of arguments for 'set' command",
			}, w)
	}

	key, ok := b[0].Val.(string)
	if !ok {
		return protocol.FormatWrite(
			&protocol.RESPVal{
				Type: protocol.RESPTypeErr,
				Val:  "ERR invalid key type",
			}, w)
	}

	val, ok := b[1].Val.(string)
	if !ok {
		return protocol.FormatWrite(
			&protocol.RESPVal{
				Type: protocol.RESPTypeErr,
				Val:  "ERR invalid val type",
			}, w)
	}

	// TODO: hardcoded by now
	var expireIn int64 = 10 * 60
	if err := store.Set(key, val, expireIn); err != nil {
		return err
	}

	return protocol.FormatWrite(&protocol.RESPVal{
		Type: protocol.RESPTypeSimpleStr,
		Val:  "OK",
	}, w)
}

func hsetCmd(w io.Writer, b []*protocol.RESPVal) error {
	if len(b) != 3 {
		return protocol.FormatWrite(
			&protocol.RESPVal{
				Type: protocol.RESPTypeErr,
				Val:  "ERR wrong number of arguments for 'hset' command",
			}, w)
	}

	hashKey, ok := b[0].Val.(string)
	if !ok {
		return protocol.FormatWrite(
			&protocol.RESPVal{
				Type: protocol.RESPTypeErr,
				Val:  "ERR invalid hash key type",
			}, w)
	}

	key, ok := b[1].Val.(string)
	if !ok {
		return protocol.FormatWrite(
			&protocol.RESPVal{
				Type: protocol.RESPTypeErr,
				Val:  "ERR invalid key type",
			}, w)
	}

	val, ok := b[2].Val.(string)
	if !ok {
		return protocol.FormatWrite(
			&protocol.RESPVal{
				Type: protocol.RESPTypeErr,
				Val:  "ERR invalid val type",
			}, w)
	}

	if err := store.HSet(hashKey, key, val); err != nil {
		return err
	}

	return protocol.FormatWrite(&protocol.RESPVal{
		Type: protocol.RESPTypeSimpleStr,
		Val:  "OK",
	}, w)
}

func getCmd(w io.Writer, b []*protocol.RESPVal) error {
	if len(b) != 1 {
		return protocol.FormatWrite(
			&protocol.RESPVal{
				Type: protocol.RESPTypeErr,
				Val:  "ERR wrong number of arguments for 'get' command",
			}, w)
	}

	key, ok := b[0].Val.(string)
	if !ok {
		return protocol.FormatWrite(
			&protocol.RESPVal{
				Type: protocol.RESPTypeErr,
				Val:  "ERR invalid key type",
			}, w)
	}

	res, err := store.Get(key)
	if err != nil {
		return err
	}

	if res == nil {
		return protocol.FormatWrite(protocol.NilBulkStr(), w)
	}

	return protocol.FormatWrite(&protocol.RESPVal{
		Type: protocol.RESPTypeBulkStr,
		Val:  *res,
	}, w)
}

func hgetCmd(w io.Writer, b []*protocol.RESPVal) error {
	if len(b) != 2 {
		return protocol.FormatWrite(
			&protocol.RESPVal{
				Type: protocol.RESPTypeErr,
				Val:  "ERR wrong number of arguments for 'hget' command",
			}, w)
	}

	hashKey, ok := b[0].Val.(string)
	if !ok {
		return protocol.FormatWrite(
			&protocol.RESPVal{
				Type: protocol.RESPTypeErr,
				Val:  "ERR invalid hash key type",
			}, w)
	}

	key, ok := b[1].Val.(string)
	if !ok {
		return protocol.FormatWrite(
			&protocol.RESPVal{
				Type: protocol.RESPTypeErr,
				Val:  "ERR invalid key type",
			}, w)
	}

	res, err := store.HGet(hashKey, key)
	if err != nil {
		return err
	}

	if res == nil {
		return protocol.FormatWrite(protocol.NilBulkStr(), w)
	}

	return protocol.FormatWrite(&protocol.RESPVal{
		Type: protocol.RESPTypeBulkStr,
		Val:  *res,
	}, w)
}

func delCmd(w io.Writer, b []*protocol.RESPVal) error {
	if len(b) == 0 {
		return protocol.FormatWrite(
			&protocol.RESPVal{
				Type: protocol.RESPTypeErr,
				Val:  "ERR wrong number of arguments for 'del' command",
			}, w)
	}
	keys := make([]string, 0, len(b))
	for _, v := range b {
		key, ok := v.Val.(string)
		if !ok {
			return protocol.FormatWrite(
				&protocol.RESPVal{
					Type: protocol.RESPTypeErr,
					Val:  "ERR invalid key type",
				}, w)
		}
		keys = append(keys, key)
	}

	store.Del(keys...)
	return protocol.FormatWrite(&protocol.RESPVal{
		Type: protocol.RESPTypeSimpleStr,
		Val:  "OK",
	}, w)
}

func hdelCmd(w io.Writer, b []*protocol.RESPVal) error {
	if len(b) < 2 {
		return protocol.FormatWrite(
			&protocol.RESPVal{
				Type: protocol.RESPTypeErr,
				Val:  "ERR wrong number of arguments for 'hdel' command",
			}, w)
	}

	hashKey, ok := b[0].Val.(string)
	if !ok {
		return protocol.FormatWrite(
			&protocol.RESPVal{
				Type: protocol.RESPTypeErr,
				Val:  "ERR invalid hash key type",
			}, w)
	}

	keys := make([]string, 0, len(b)-1)
	for i := 1; i < len(b); i++ {
		key, ok := b[i].Val.(string)
		if !ok {
			return protocol.FormatWrite(
				&protocol.RESPVal{
					Type: protocol.RESPTypeErr,
					Val:  "ERR invalid key type",
				}, w)
		}

		keys = append(keys, key)
	}

	store.HDel(hashKey, keys...)
	return protocol.FormatWrite(&protocol.RESPVal{
		Type: protocol.RESPTypeSimpleStr,
		Val:  "OK",
	}, w)
}

func existsCmd(w io.Writer, b []*protocol.RESPVal) error {
	if len(b) == 0 {
		return protocol.FormatWrite(
			&protocol.RESPVal{
				Type: protocol.RESPTypeErr,
				Val:  "ERR wrong number of arguments for 'exists' command",
			}, w)
	}

	keys := make([]string, 0, len(b))
	for _, v := range b {
		key, ok := v.Val.(string)
		if !ok {
			return protocol.FormatWrite(
				&protocol.RESPVal{
					Type: protocol.RESPTypeErr,
					Val:  "ERR invalid key type",
				}, w)
		}

		keys = append(keys, key)
	}

	count := store.Exists(keys...)

	return protocol.FormatWrite(&protocol.RESPVal{
		Type: protocol.RESPTypeInt,
		Val:  count,
	}, w)
}

func expireCmd(w io.Writer, b []*protocol.RESPVal) error {
	if len(b) != 2 {
		return protocol.FormatWrite(
			&protocol.RESPVal{
				Type: protocol.RESPTypeErr,
				Val:  "ERR wrong number of arguments for 'expire' command",
			}, w)
	}

	key, ok := b[0].Val.(string)
	if !ok {
		return protocol.FormatWrite(
			&protocol.RESPVal{
				Type: protocol.RESPTypeErr,
				Val:  "ERR invalid key type",
			}, w)
	}

	expireIn, ok := b[1].Val.(int64)
	if !ok {
		return protocol.FormatWrite(
			&protocol.RESPVal{
				Type: protocol.RESPTypeErr,
				Val:  "ERR invalid expireIn type",
			}, w)
	}

	if err := store.Expire(key, expireIn); err != nil {
		return err
	}

	return protocol.FormatWrite(&protocol.RESPVal{
		Type: protocol.RESPTypeSimpleStr,
		Val:  "OK",
	}, w)
}

func ttlCmd(w io.Writer, b []*protocol.RESPVal) error {
	if len(b) != 1 {
		return protocol.FormatWrite(
			&protocol.RESPVal{
				Type: protocol.RESPTypeErr,
				Val:  "ERR wrong number of arguments for 'ttl' command",
			}, w)
	}

	key, ok := b[0].Val.(string)
	if !ok {
		return protocol.FormatWrite(
			&protocol.RESPVal{
				Type: protocol.RESPTypeErr,
				Val:  "ERR invalid key type",
			}, w)
	}

	ttl, err := store.TTL(key)
	if err != nil {
		return err
	}

	if ttl == nil {
		return protocol.FormatWrite(protocol.NilBulkStr(), w)
	}

	return protocol.FormatWrite(&protocol.RESPVal{
		Type: protocol.RESPTypeInt,
		Val:  *ttl,
	}, w)
}
