package gomap

import (
	"iter"
	"slices"
)

type Entry[K, V any] struct {
	key K
	val V
}

type Map[K, V any] struct {
	NBuckets  int
	buckets   [][]Entry[K, V]
	size      int
	hashFunc  func(K) int
	equalFunc func(K, K) bool
}

// TODO: maybe i will make nBuckets as optional param

// New creates a new Map with the specified number of buckets, a hash function, and an equality function.
func New[K, V any](nBuckets int, hashf func(K) int, equalf func(K, K) bool) *Map[K, V] {
	return &Map[K, V]{
		NBuckets:  nBuckets,
		buckets:   make([][]Entry[K, V], nBuckets),
		size:      0,
		hashFunc:  hashf,
		equalFunc: equalf,
	}
}

// Set adds or updates a key-value pair in the map.
// If the key already exists, its value is updated; otherwise, a new entry is created.
func (me *Map[K, V]) Set(k K, v V) {
	bktIdx := me.hashFunc(k) % me.NBuckets

	idx := slices.IndexFunc(me.buckets[bktIdx], func(e Entry[K, V]) bool {
		return me.equalFunc(e.key, k)
	})
	if idx != -1 { // found
		me.buckets[bktIdx][idx].val = v
		return
	}

	me.buckets[bktIdx] = append(me.buckets[bktIdx], Entry[K, V]{key: k, val: v})
	me.size++
}

// Get retrieves the value associated with the given key.
// Returns the value and true if the key exists, or the zero value and false if it does not.
func (me *Map[K, V]) Get(k K) (V, bool) {
	bktIdx := me.hashFunc(k) % me.NBuckets

	for _, e := range me.buckets[bktIdx] {
		if me.equalFunc(e.key, k) {
			return e.val, true
		}
	}

	var zero V
	return zero, false
}

// Delete removes the key-value pair associated with the given key from the map.
// If the key does not exist, the map remains unchanged.
func (me *Map[K, V]) Delete(k K) {
	bktIdx := me.hashFunc(k) % me.NBuckets

	idx := slices.IndexFunc(me.buckets[bktIdx], func(e Entry[K, V]) bool {
		return me.equalFunc(e.key, k)
	})
	if idx != -1 { // found
		me.buckets[bktIdx] = slices.Delete(me.buckets[bktIdx], idx, idx+1)
		me.size--
	}
}

// Size returns the number of key-value pairs currently in the map.
func (me *Map[K, V]) Size() int {
	return me.size
}

// TODO: comment and test
func (me *Map[K, V]) Iterator() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, bkt := range me.buckets {
			for _, e := range bkt {
				if !yield(e.val) {
					return
				}
			}
		}
	}
}
