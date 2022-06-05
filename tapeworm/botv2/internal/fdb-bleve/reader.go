package fdbbleve

import (
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	store "github.com/blevesearch/upsidedown_store_api"
)

type Reader struct {
	store *Store
	tx    *fdb.Transaction
}

// Get returns the value associated with the key
// If the key does not exist, nil is returned.
// The caller owns the bytes returned.
func (r *Reader) Get(key []byte) ([]byte, error) {
	future := r.tx.Get(fdb.Key(key))
	value, err := future.Get()
	if err != nil {
		return nil, err
	}

	return value, nil
}

// MultiGet retrieves multiple values in one call.
func (r *Reader) MultiGet(keys [][]byte) ([][]byte, error) {
	values := make([][]byte, len(keys))
	for _, key := range keys {
		future := r.tx.Get(fdb.Key(key))
		value, err := future.Get()
		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}

	return values, nil
}

// PrefixIterator returns a KVIterator that will
// visit all K/V pairs with the provided prefix
func (r *Reader) PrefixIterator(prefix []byte) store.KVIterator {
	iterator := &Iterator{
		store: r.store,
	}

	pr, err := fdb.PrefixRange(prefix)
	if err != nil {
		return &Iterator{
			err: err,
		}
	}

	rangeResult := r.tx.GetRange(pr, fdb.RangeOptions{})
	items, err := rangeResult.GetSliceWithError()
	if err != nil {
		return &Iterator{
			err: err,
		}
	}

	iterator.keyvalues = items

	return iterator
}

// RangeIterator returns a KVIterator that will
// visit all K/V pairs >= start AND < end
func (r *Reader) RangeIterator(start, end []byte) store.KVIterator {
	if end == nil {
		end = []byte("\xff")
	}

	rangeKey := fdb.KeyRange{
		Begin: fdb.Key(start),
		End:   fdb.Key(end),
	}

	iterator := &Iterator{
		store: r.store,
	}

	rangeResult := r.tx.GetRange(rangeKey, fdb.RangeOptions{})
	items, err := rangeResult.GetSliceWithError()
	if err != nil {
		return &Iterator{
			err: err,
		}
	}

	iterator.keyvalues = items

	return iterator
}

// Close closes the iterator
func (r *Reader) Close() error {
	return r.tx.Commit().Get()
}

var _ store.KVReader = &Reader{}
