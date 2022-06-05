package fdbbleve

import (
	"testing"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	store "github.com/blevesearch/upsidedown_store_api"
	"github.com/blevesearch/upsidedown_store_api/test"
	"github.com/go-logr/zapr"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/logging"
	"go.uber.org/zap"
)

func open(t *testing.T, mo store.MergeOperator) store.KVStore {
	rv, err := New(mo, map[string]interface{}{
		"fdbAPIVersion": 710,
		"clusterFile":   "/tmp/fdb.cluster",
	})
	if err != nil {
		t.Fatal(err)
	}
	return rv
}

func cleanup(t *testing.T) {
	fdb.MustAPIVersion(710)
	db := fdb.MustOpenDatabase("/tmp/fdb.cluster")
	db.Transact(func(t fdb.Transaction) (interface{}, error) {
		t.ClearRange(fdb.KeyRange{
			Begin: fdb.Key("\x00"),
			End:   fdb.Key("\xff"),
		})

		return nil, nil
	})
}

func TestFDBKVCrud(t *testing.T) {
	zapl, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("init logger err=%v", err)
	}
	l := zapr.NewLogger(zapl)
	logging.Logger = l

	s := open(t, nil)
	defer cleanup(t)
	test.CommonTestKVCrud(t, s)
}

func TestFDBReaderIsolation(t *testing.T) {
	s := open(t, nil)
	defer cleanup(t)
	test.CommonTestReaderIsolation(t, s)
}

func TestFDBReaderOwnsGetBytes(t *testing.T) {
	s := open(t, nil)
	defer cleanup(t)
	test.CommonTestReaderOwnsGetBytes(t, s)
}

func TestFDBWriterOwnsBytes(t *testing.T) {
	s := open(t, nil)
	defer cleanup(t)
	test.CommonTestWriterOwnsBytes(t, s)
}

func TestFDBPrefixIterator(t *testing.T) {
	s := open(t, nil)
	defer cleanup(t)
	test.CommonTestPrefixIterator(t, s)
}

func TestFDBPrefixIteratorSeek(t *testing.T) {
	s := open(t, nil)
	defer cleanup(t)
	test.CommonTestPrefixIteratorSeek(t, s)
}

func TestFDBRangeIterator(t *testing.T) {
	s := open(t, nil)
	defer cleanup(t)
	test.CommonTestRangeIterator(t, s)
}

func TestFDBRangeIteratorSeek(t *testing.T) {
	s := open(t, nil)
	defer cleanup(t)
	test.CommonTestRangeIteratorSeek(t, s)
}

func TestFDBMerge(t *testing.T) {
	s := open(t, &test.TestMergeCounter{})
	defer cleanup(t)
	test.CommonTestMerge(t, s)
}
