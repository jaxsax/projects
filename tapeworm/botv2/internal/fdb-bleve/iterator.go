package fdbbleve

import (
	"bytes"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	store "github.com/blevesearch/upsidedown_store_api"
)

type Iterator struct {
	store     *Store
	keyvalues []fdb.KeyValue
	curr      int
	err       error
}

// Seek will advance the iterator to the specified key
func (i *Iterator) Seek(key []byte) {
	for idx, kv := range i.keyvalues {
		if bytes.Compare(kv.Key, key) >= 0 {
			i.curr = idx
			return
		}
	}
}

// Next will advance the iterator to the next key
func (i *Iterator) Next() {
	i.curr++
}

// Key returns the key pointed to by the iterator
// The bytes returned are **ONLY** valid until the next call to Seek/Next/Close
// Continued use after that requires that they be copied.
func (i *Iterator) Key() []byte {
	k, _, valid := i.Current()
	if !valid {
		return nil
	}
	return k
}

// Value returns the value pointed to by the iterator
// The bytes returned are **ONLY** valid until the next call to Seek/Next/Close
// Continued use after that requires that they be copied.
func (i *Iterator) Value() []byte {
	_, v, valid := i.Current()
	if !valid {
		return nil
	}

	return v
}

// Valid returns whether or not the iterator is in a valid state
func (i *Iterator) Valid() bool {
	if i.keyvalues == nil {
		return false
	}

	if i.err != nil {
		return false
	}

	return i.curr < len(i.keyvalues)
}

// Current returns Key(),Value(),Valid() in a single operation
func (i *Iterator) Current() ([]byte, []byte, bool) {
	if !i.Valid() {
		return nil, nil, false
	}

	curr := i.keyvalues[i.curr]
	return curr.Key, curr.Value, true
}

// Close closes the iterator
func (i *Iterator) Close() error {

	return nil
}

var _ store.KVIterator = &Iterator{}
