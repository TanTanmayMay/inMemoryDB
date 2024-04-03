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

func (d *DB[K, V]) MoveTo(db *DB[K, V]) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for k, v := range d.m {
		db.m[k] = v
	}
	d.m = make(map[K]V)
}

func (d *DB[K, V]) CopyTo(db *DB[K, V]) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for k, v := range d.m {
		db.m[k] = v
	}
}

func (d *DB[K, V]) Keys() []K {
	d.mu.RLock()
	defer d.mu.RUnlock()
	keys := make([]K, 0, len(d.m))
	for k, _ := range d.m {
		keys = append(keys, k)
	}
	return keys
}
