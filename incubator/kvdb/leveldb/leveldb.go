package leveldb

import (
	"github.com/hello2mao/go-common/incubator/kvdb"
)

type Database struct {

}

func Open(options map[string]interface{}) (kvdb.Database, error) {
	panic("implement me")
}

func (db *Database) Close() error {
	panic("implement me")
}

func (db *Database) Has(key []byte) (bool, error) {
	panic("implement me")
}

func (db *Database) Get(key []byte) ([]byte, error) {
	panic("implement me")
}

func (db *Database) Put(key []byte, value []byte) error {
	panic("implement me")
}

func (db *Database) Delete(key []byte) error {
	panic("implement me")
}

func (db *Database) NewBatch() kvdb.Batch {
	panic("implement me")
}

func (db *Database) NewIterator(prefix []byte, start []byte) kvdb.Iterator {
	panic("implement me")
}

func (db *Database) Stat(property string) (string, error) {
	panic("implement me")
}



