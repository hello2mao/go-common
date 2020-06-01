package memorydb

import (
	"testing"

	"github.com/hello2mao/go-common/incubator/kvdb"
	"github.com/hello2mao/go-common/incubator/kvdb/dbtest"
)

func TestMemorydb(t *testing.T) {
	t.Run("TestMemorydb", func(t *testing.T) {
		dbtest.TestDatabaseSuite(t, func(options map[string]interface{}) (kvdb.Database, error) {
			return Open(options)
		})
	})
}