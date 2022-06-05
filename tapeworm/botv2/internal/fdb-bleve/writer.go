package fdbbleve

import (
	"fmt"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	store "github.com/blevesearch/upsidedown_store_api"
)

type Writer struct {
	store *Store
}

func (w *Writer) NewBatch() store.KVBatch {
	return store.NewEmulatedBatch(w.store.mo)
}

func (w *Writer) NewBatchEx(options store.KVBatchOptions) ([]byte, store.KVBatch, error) {
	return make([]byte, options.TotalBytes), w.NewBatch(), nil
}

func (w *Writer) ExecuteBatch(batch store.KVBatch) error {
	eb, ok := batch.(*store.EmulatedBatch)
	if !ok {
		return fmt.Errorf("wrong batch type")
	}

	_, err := w.store.fdb.Transact(func(t fdb.Transaction) (interface{}, error) {
		for k, mergeOps := range eb.Merger.Merges {
			fk := fdb.Key(k)
			existingValue, err := t.Get(fk).Get()
			if err != nil {
				return nil, err
			}

			mergedVal, fullMergeOk := w.store.mo.FullMerge(fk, existingValue, mergeOps)
			if !fullMergeOk {
				return nil, fmt.Errorf("merge operator returned failures")
			}

			t.Set(fk, mergedVal)
		}

		for _, op := range eb.Ops {
			fk := fdb.Key(op.K)
			if op.V != nil {
				t.Set(fk, op.V)
			} else {
				t.Clear(fk)
			}
		}

		return nil, nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (w *Writer) Close() error {
	return nil
}

var _ store.KVWriter = &Writer{}
