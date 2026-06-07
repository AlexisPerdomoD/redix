// Package store provides storage code for the redix server.
package store

import "errors"

var (
	ErrInvalidHashKey      = errors.New("invalid hash key")
	ErrInvalidKeyFormat    = errors.New("invalid key format")
	ErrExpectedExistingKey = errors.New("expected existing key and found none")
	ErrWrongType           = errors.New("operation against a key holding the wrong kind of value")
)

type storeVal struct {
	Val any
	TTL int64
}

func ensureValidKey(key string) error {
	if len(key) == 0 {
		return ErrInvalidKeyFormat
	}

	return nil
}

var mem = NewMemStore()

// Set sets the value of a key in the store.
// key is the key to set the value for.
// val is the value to set.
// ttl is the time-to-live of the key in seconds (0 or negative values mean no expiration).
func Set(key, val string, expireIn int64) error {
	if err := ensureValidKey(key); err != nil {
		return err
	}

	mem.Set(key, val, expireIn)
	return nil
}

// HSet sets the value of a key in the store.
// hashKey is the key to set the value for.
// key is the key to set the value for.
// val is the value to set.
// expration is set to 0 as no expiration for new hash keys.
func HSet(hashKey, key, val string) error {
	if err := ensureValidKey(hashKey); err != nil {
		return ErrInvalidHashKey
	}

	if err := ensureValidKey(key); err != nil {
		return err
	}

	return mem.HSet(hashKey, key, val)
}

// Get gets the value of a key in the store as a string.
// key is the key to get the value for.
// Returns the value and an error if there was an error.
// Returns nil and nil if the key does not exist or has expired.
func Get(key string) (*string, error) {
	if err := ensureValidKey(key); err != nil {
		return nil, err
	}

	return mem.Get(key)
}

// HGet gets the value of a key in the store.
// hashKey is the key to get the value for.
// key is the key to get the value for.
// Returns the value and an error if there was an error.
// Returns nil and nil if the key does not exist or has expired.
func HGet(hashKey, key string) (*string, error) {
	if err := ensureValidKey(hashKey); err != nil {
		return nil, err
	}

	if err := ensureValidKey(key); err != nil {
		return nil, err
	}

	return mem.HGet(hashKey, key)
}

// HGetAll gets all the values of a key in the store.
// hashKey is the key to get the value for.
// Returns the value and an error if there was an error.
// Returns nil and nil if the key does not exist or has expired.
func HGetAll(hashKey string) (map[string]string, error) {
	if err := ensureValidKey(hashKey); err != nil {
		return nil, err
	}

	return mem.HGetAll(hashKey)
}

// Del deletes the values of a key in the store.
// keys is the key to delete the value for.
func Del(keys ...string) {
	if len(keys) == 0 {
		return
	}

	mem.Del(keys...)
}

// HDel deletes the values of a key in the store.
// hashKey is the key to delete the value for.
// keys is the key to delete the value for.
func HDel(hashKey string, keys ...string) {
	if err := ensureValidKey(hashKey); err != nil {
		return
	}

	if len(keys) == 0 {
		return
	}

	mem.HDel(hashKey, keys...)
}

// Exists checks if a key exists in the store.
func Exists(keys ...string) int64 {
	return mem.Exists(keys...)
}

// HExists checks if a key exists in a hash key in the store.
// hashKey is the key to check for.
// key is the key to check for.
// Returns true if the key exists, false otherwise.
func HExists(hashKey, key string) bool {
	if err := ensureValidKey(hashKey); err != nil {
		return false
	}

	if err := ensureValidKey(key); err != nil {
		return false
	}

	return mem.HExists(hashKey, key)
}

// TTL returns the time-to-live of a key in the store in unix time or nil if not expirable or key does not exist.
// key is the key to get the time-to-live for.
func TTL(key string) (*int64, error) {
	if err := ensureValidKey(key); err != nil {
		return nil, err
	}

	return mem.TTL(key)
}

// Expire sets the time-to-live(unix time) of a key in the store.
// key is the key to set the time-to-live for.
// expireIn is the time-to-live of the key in seconds (0 or negative values mean no expiration).
func Expire(key string, expireIn int64) error {
	if err := ensureValidKey(key); err != nil {
		return err
	}

	return mem.Expire(key, expireIn)
}

// FlushAll flushes all the keys in the store.
func FlushAll() {
	mem.FlushAll()
}
