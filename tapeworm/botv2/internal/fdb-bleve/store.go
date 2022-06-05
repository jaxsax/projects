package fdbbleve

import (
	"fmt"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/blevesearch/bleve/v2/registry"
	store "github.com/blevesearch/upsidedown_store_api"
)

const (
	Name = "foundationdb"
)

type Store struct {
	fdb *fdb.Database
	mo  store.MergeOperator
}

func New(mo store.MergeOperator, config map[string]interface{}) (store.KVStore, error) {
	fdbAPIVersion, ok := config["fdbAPIVersion"].(int)
	if !ok {
		return nil, fmt.Errorf("fdbAPIVersion must be specified")
	}

	if err := fdb.APIVersion(fdbAPIVersion); err != nil {
		return nil, err
	}

	fdbClusterConfig, ok := config["clusterFile"].(string)
	if !ok {
		return nil, fmt.Errorf("clusterFile must be specified")
	}

	db, err := fdb.OpenDatabase(fdbClusterConfig)
	if err != nil {
		return nil, err
	}

	return &Store{
		mo:  mo,
		fdb: &db,
	}, nil
}

func (s *Store) Writer() (store.KVWriter, error) {
	return &Writer{
		store: s,
	}, nil
}

func (s *Store) Reader() (store.KVReader, error) {
	tx, err := s.fdb.CreateTransaction()
	if err != nil {
		return nil, err
	}
	return &Reader{
		store: s,
		tx:    &tx,
	}, nil
}

func (s *Store) Close() error {
	return nil
}

func init() {
	registry.RegisterKVStore(Name, New)
}

var _ store.KVStore = &Store{}
