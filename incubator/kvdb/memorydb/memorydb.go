package memorydb

import (
	"errors"
	"fmt"
	"sync"

	"github.com/hello2mao/go-common/incubator/kvdb"
)

var (
	// errMemorydbClosed is returned if a memory database was already closed at the
	// invocation of a data access operation.
	errMemorydbClosed = errors.New("database closed")

	// errMemorydbNotFound is returned if a key is requested that is not found in
	// the provided memory database.
	errMemorydbNotFound = errors.New("not found")
)

// Database is an ephemeral key-value store. Apart from basic data storage
// functionality it also supports batch writes and iterating over the keyspace in
// binary-alphabetical order.
type Database struct {
	db   map[string][]byte
	lock sync.RWMutex
}

// Open returns target db instance
func Open(options map[string]interface{}) (kvdb.Database, error) {
	if value, exist := options["size"]; exist {
		if size, ok := value.(int); !ok {
			return nil, fmt.Errorf("open database err: size must be int")
		} else {
			return &Database{
				db: make(map[string][]byte, size),
			}, nil
		}
	} else {
		return &Database{
			db: make(map[string][]byte),
		}, nil
	}
}

// Has retrieves if a key is present in the key-value data store.
func (db *Database) Has(key []byte) (bool, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if db.db == nil {
		return false, errMemorydbClosed
	}
	_, ok := db.db[string(key)]
	return ok, nil
}

// Get retrieves the given key if it's present in the key-value data store.
func (db *Database) Get(key []byte) ([]byte, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if db.db == nil {
		return nil, errMemorydbClosed
	}
	var copiedBytes []byte
	if entry, ok := db.db[string(key)]; ok {
		copiedBytes = make([]byte, len(entry))
		copy(copiedBytes, entry)
		return copiedBytes, nil
	}
	return nil, errMemorydbNotFound
}

// Put inserts the given value into the key-value data store.
func (db *Database) Put(key []byte, value []byte) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.db == nil {
		return errMemorydbClosed
	}
	var copiedBytes []byte
	copiedBytes = make([]byte, len(value))
	copy(copiedBytes, value)
	db.db[string(key)] = copiedBytes
	return nil
}

// Delete removes the key from the key-value data store.
func (db *Database) Delete(key []byte) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.db == nil {
		return errMemorydbClosed
	}
	delete(db.db, string(key))
	return nil
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

// Close close the target db
func (db *Database) Close() error {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.db = nil
	return nil
}
