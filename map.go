package gomap

import (
	"iter"
	"slices"
)

const defaultNBuckets = 100

type HashFunc[K any] func(K) int
type EqualFunc[K any] func(K, K) bool

type Entry[K, V any] struct {
	Key K
	Val V
}

type Map[K, V any] struct {
	NBuckets  int
	buckets   [][]Entry[K, V]
	size      int
	hash  HashFunc[K]
	equal EqualFunc[K]
}

// New creates a new Map with the specified number of buckets, a hash function, and an equality function.
func New[K, V any](hash HashFunc[K], equal EqualFunc[K], nBuckets ...int) *Map[K, V] {
	var nBkts int
	if len(nBuckets) != 0 {
		nBkts = nBuckets[0]
	} else {
		nBkts = defaultNBuckets
	}
	m := &Map[K, V]{
		NBuckets:  nBkts,
		buckets:   make([][]Entry[K, V], nBkts),
		size:      0,
		hash:  hash,
		equal: equal,
	}
	for i := range m.buckets {
		m.buckets[i] = make([]Entry[K, V], 0, 20)
	}
	return m
}

// Set adds or updates a key-value pair in the map.
// If the key already exists, its value is updated; otherwise, a new entry is created.
func (me *Map[K, V]) Set(k K, v V) {
	bktIdx := me.hash(k) % me.NBuckets

	idx := slices.IndexFunc(me.buckets[bktIdx], func(e Entry[K, V]) bool {
		return me.equal(e.Key, k)
	})
	if idx != -1 { // found
		me.buckets[bktIdx][idx].Val = v
		return
	}

	me.buckets[bktIdx] = append(me.buckets[bktIdx], Entry[K, V]{Key: k, Val: v})
	me.size++
}

// Get retrieves the value associated with the given key.
// Returns the value and true if the key exists, or the zero value and false if it does not.
func (me *Map[K, V]) Get(k K) (V, bool) {
	bktIdx := me.hash(k) % me.NBuckets

	for _, e := range me.buckets[bktIdx] {
		if me.equal(e.Key, k) {
			return e.Val, true
		}
	}

	var zero V
	return zero, false
}

// Delete removes the key-value pair associated with the given key from the map.
// If the key does not exist, the map remains unchanged.
func (me *Map[K, V]) Delete(k K) {
	bktIdx := me.hash(k) % me.NBuckets

	idx := slices.IndexFunc(me.buckets[bktIdx], func(e Entry[K, V]) bool {
		return me.equal(e.Key, k)
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

// Iterator returns a sequence that allows iteration over all the values in the map.
func (me *Map[K, V]) Iterator() iter.Seq[Entry[K, V]] {
	return func(yield func(Entry[K, V]) bool) {
		for _, bkt := range me.buckets {
			for _, e := range bkt {
				if !yield(e) {
					return
				}
			}
		}
	}
}
