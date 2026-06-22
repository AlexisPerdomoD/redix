package store

import (
	"maps"
	"sync"
	"time"
)

type MemStore struct {
	mtx  sync.RWMutex
	data map[string]storeVal
}

func (st *MemStore) set(key string, val string, expireIn int64) {
	// 0 means no expiration
	var ttl int64
	if expireIn > 0 {
		ttl = time.Now().Unix() + expireIn
	}

	st.data[key] = storeVal{val, ttl}
}

// Get gets the value of a key in the store.
// key is the key to get the value for.
// Returns the value and an error if there was an error.
// Returns nil and nil if the key does not exist or has expired.
func (st *MemStore) get(key string) *storeVal {
	stVal, exists := st.data[key]

	if exists && !stVal.isExpired() {
		return &stVal
	}

	return nil
}

// Set sets the value of a key in the store.
// key is the key to set the value for.
// val is the value to set.
// ttl is the time-to-live of the key in seconds (0 or negative values mean no expiration).
func (st *MemStore) Set(key, val string, expireIn int64) {
	st.mtx.Lock()
	defer st.mtx.Unlock()
	st.set(key, val, expireIn)
}

// HSet sets the value of a key in the store.
// hashKey is the key to set the value for.
// key is the key to set the value for.
// val is the value to set.
// expration is set to 0 as no expiration for new hash keys.
func (st *MemStore) HSet(hashKey, key, val string) error {

	st.mtx.Lock()
	defer st.mtx.Unlock()
	stVal := st.get(hashKey)
	if stVal == nil {
		st.data[hashKey] = storeVal{
			Val: map[string]string{key: val},
		}
		return nil
	}

	hashMap, ok := stVal.Val.(map[string]string)
	if !ok {
		return ErrInvalidHashKey
	}

	hashMap[key] = val
	return nil
}

// Get gets the value of a key in the store as a string.
// key is the key to get the value for.
// Returns the value and an error if there was an error.
// Returns nil and nil if the key does not exist or has expired.
func (st *MemStore) Get(key string) (*string, error) {
	st.mtx.RLock()
	defer st.mtx.RUnlock()
	stVal := st.get(key)

	if stVal == nil {
		return nil, nil
	}

	val, ok := stVal.Val.(string)
	if !ok {
		return nil, ErrWrongType
	}

	return &val, nil
}

// HGet gets the value of a key in the store.
// hashKey is the key to get the value for.
// key is the key to get the value for.
// Returns the value and an error if there was an error.
// Returns nil and nil if the key does not exist or has expired.
func (st *MemStore) HGet(hashKey, key string) (*string, error) {

	st.mtx.RLock()
	defer st.mtx.RUnlock()

	stVal := st.get(hashKey)

	if stVal == nil {
		return nil, nil
	}

	hashMap, ok := stVal.Val.(map[string]string)
	if !ok {
		return nil, ErrWrongType
	}

	val, ok := hashMap[key]
	if !ok {
		return nil, nil
	}

	return &val, nil
}

// HGetAll gets all the values of a key in the store.
// hashKey is the key to get the value for.
// Returns the value and an error if there was an error.
// Returns nil and nil if the key does not exist or has expired.
func (st *MemStore) HGetAll(hashKey string) (map[string]string, error) {
	st.mtx.RLock()
	defer st.mtx.RUnlock()

	stVal := st.get(hashKey)

	if stVal == nil {
		return nil, nil
	}

	val, ok := stVal.Val.(map[string]string)
	if !ok {
		return nil, ErrWrongType
	}

	cp := make(map[string]string, len(val))
	maps.Copy(cp, val)
	return cp, nil
}

// Del deletes the values of a key in the store.
// keys is the key to delete the value for.
func (st *MemStore) Del(keys ...string) {
	if len(keys) == 0 {
		return
	}

	st.mtx.Lock()
	defer st.mtx.Unlock()

	for _, k := range keys {
		if err := ensureValidKey(k); err != nil {
			continue
		}

		delete(st.data, k)
	}
}

// HDel deletes the values of a key in the store.
// hashKey is the key to delete the value for.
// keys is the key to delete the value for.
func (st *MemStore) HDel(hashKey string, keys ...string) {
	if len(keys) == 0 {
		return
	}

	st.mtx.Lock()
	defer st.mtx.Unlock()

	stVal := st.get(hashKey)
	if stVal == nil {
		return
	}

	hashMap, ok := stVal.Val.(map[string]string)
	if !ok {
		return
	}

	for _, k := range keys {
		if err := ensureValidKey(k); err != nil {
			continue
		}

		delete(hashMap, k)
	}
}

// Exists checks if a key exists in the store.
func (st *MemStore) Exists(keys ...string) int64 {
	st.mtx.RLock()
	defer st.mtx.RUnlock()
	var res int64
	for _, k := range keys {
		if st.get(k) != nil {
			res++
		}
	}

	return res

}

// HExists checks if a key exists in a hash key in the store.
// hashKey is the key to check for.
// key is the key to check for.
func (st *MemStore) HExists(hashKey, key string) bool {
	st.mtx.RLock()
	defer st.mtx.RUnlock()

	stVal := st.get(hashKey)
	if stVal == nil {
		return false
	}

	hashMap, ok := stVal.Val.(map[string]string)
	if !ok {
		return false
	}

	_, exists := hashMap[key]
	return exists
}

// TTL returns the seconds until the key expires.
// -2 means the key does not exist.
// -1 means the key is not expirable.
// 0 means the key has expired.
// >0 means the key has not expired.
func (st *MemStore) TTL(key string) int64 {
	st.mtx.RLock()
	defer st.mtx.RUnlock()

	stVal := st.get(key)
	if stVal == nil {
		return -2
	}

	if stVal.isNotExpirable() {
		return -1
	}

	if stVal.isExpired() {
		return 0
	}

	return stVal.TTL - time.Now().Unix()

}

// Expire sets the time-to-live(unix time) of a key in the store.
// key is the key to set the time-to-live for.
// expireIn is the time-to-live of the key in seconds (0 or negative values mean no expiration).
func (st *MemStore) Expire(key string, expireIn int64) error {
	st.mtx.Lock()
	defer st.mtx.Unlock()

	stVal := st.get(key)

	if stVal == nil {
		return nil
	}

	// 0 means no expiration
	var ttl int64
	if expireIn > 0 {
		ttl = time.Now().Unix() + expireIn
	}
	stVal.TTL = ttl
	st.data[key] = *stVal
	return nil
}

// FlushAll flushes all the keys in the store.
func (st *MemStore) FlushAll() {
	st.mtx.Lock()
	defer st.mtx.Unlock()
	st.data = make(map[string]storeVal)
}

// cleaningLoopMem periodically removes expired keys from the store.
//
// Expiration cleanup is performed in two phases:
//
//  1. A full keyspace scan under RLock collects candidate expired keys.
//  2. Candidates are revalidated and removed under Lock.
//
// This approach minimizes exclusive lock contention while preserving
// correctness in the presence of concurrent writes. Keys may be modified,
// renewed, or deleted between phases, therefore expiration state is
// revalidated before removal.
//
// Cleanup complexity remains O(total_keys) because the entire keyspace
// must be scanned on each iteration. Future optimizations could maintain
// a dedicated expiration index to avoid full scans and make cleanup
// proportional to the number of expiring keys.
func cleaningLoopMem(s *MemStore, expireTimeLoop time.Duration) {
	if expireTimeLoop <= 0 {
		expireTimeLoop = defaultExpireTimeLoop
	}

	ticker := time.NewTicker(expireTimeLoop)
	for range ticker.C {
		expiredKeys := make([]string, 0)
		s.mtx.RLock()
		for k, v := range s.data {
			if v.isExpired() {
				expiredKeys = append(expiredKeys, k)
			}
		}
		s.mtx.RUnlock()

		if len(expiredKeys) == 0 {
			continue
		}

		s.mtx.Lock()
		for _, k := range expiredKeys {
			v, ok := s.data[k]
			if !ok || !v.isExpired() {
				continue
			}

			delete(s.data, k)
		}

		s.mtx.Unlock()
	}
}

func NewMemStore() *MemStore {
	s := &MemStore{
		mtx:  sync.RWMutex{},
		data: make(map[string]storeVal),
	}

	go cleaningLoopMem(s, defaultExpireTimeLoop)
	return s
}
