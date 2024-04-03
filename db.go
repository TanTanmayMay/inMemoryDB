package inmemorydb

import "sync"

type DB[K comparable, V any] struct {
	m  map[K]V
	mu sync.RWMutex
}

func New[K comparable, V any]() *DB[K, V] {
	return &DB[K, V]{
		m:  make(map[K]V),
		mu: sync.RWMutex{},
	}
}

func (d *DB[K, V]) Set(k K, v V) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.m[k] = v
}

func (d *DB[K, V]) Get(k K) (V, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	val, ok := d.m[k]
	return val, ok
}

func (d *DB[K, V]) Delete(k K) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.m, k)
}

// TransferTo Function is used to Transfer data from Source DB to Destination DB
// Source -> Locked for the entire duration
// Destination -> Locked for the function call duration
func (src *DB[K, V]) TransferTo(dest *DB[K, V]) {
	src.mu.RLock()
	dest.mu.Lock()
	defer src.mu.RUnlock()
	defer dest.mu.Unlock()
	for k, v := range src.m {
		dest.m[k] = v
	}
	src.m = make(map[K]V)
}

// CopyTo Function is used to copy data from Source DB to Destination DB
// Source -> Locked for the entire duration
// Destination -> Locked for the function call duration
func (src *DB[K, V]) CopyTo(dest *DB[K, V]) {
	src.mu.RLock()
	dest.mu.Lock()
	defer src.mu.RUnlock()
	defer dest.mu.Unlock()
	for k, v := range src.m {
		dest.m[k] = v
	}
}

// Returns the slice of keys present in the Database
func (d *DB[K, V]) Keys() []K {
	d.mu.RLock()
	defer d.mu.RUnlock()
	keys := make([]K, 0, len(d.m))
	for k, _ := range d.m {
		keys = append(keys, k)
	}
	return keys
}
