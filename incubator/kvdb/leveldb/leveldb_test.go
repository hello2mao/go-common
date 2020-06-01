package leveldb

import (
	"testing"

	"github.com/hello2mao/go-common/incubator/kvdb"
	"github.com/hello2mao/go-common/incubator/kvdb/dbtest"
)

func TestLeveldb(t *testing.T) {
	t.Run("TestLeveldb", func(t *testing.T) {
		dbtest.TestDatabaseSuite(t, func(options map[string]interface{}) (kvdb.Database, error) {
			return Open(options)
		})
	})
}
