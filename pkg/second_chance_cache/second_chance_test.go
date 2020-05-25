package cache

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecondChanceCache(t *testing.T) {
	cache := NewSecondChanceCache(2)
	assert.NotNil(t, cache)

	cache.Add("a", "xyz")

	cache.Add("b", "123")
	// get b, b referenced bit is set to true
	obj, ok := cache.Get("b")
	assert.True(t, ok)
	assert.Equal(t, "123", obj.(string))

	// add c. victim scan: delete a and set b as the next candidate of a victim
	cache.Add("c", "777")

	// check a is deleted
	_, ok = cache.Get("a")
	assert.False(t, ok)

	// add d. victim scan: b referenced bit is set to false and delete c
	cache.Add("d", "555")

	// check c is deleted
	_, ok = cache.Get("c")
	assert.False(t, ok)

	// check b and d
	obj, ok = cache.Get("b")
	assert.True(t, ok)
	assert.Equal(t, "123", obj.(string))
	obj, ok = cache.Get("d")
	assert.True(t, ok)
	assert.Equal(t, "555", obj.(string))
}

func TestSecondChanceCacheConcurrent(t *testing.T) {
	cache := NewSecondChanceCache(25)

	workers := 16
	wg := sync.WaitGroup{}
	wg.Add(workers)

	key1 := fmt.Sprintf("key1")
	val1 := key1

	for i := 0; i < workers; i++ {
		id := i
		key2 := fmt.Sprintf("key2-%d", i)
		val2 := key2

		go func() {
			for j := 0; j < 10000; j++ {
				key3 := fmt.Sprintf("key3-%d-%d", id, j)
				val3 := key3
				cache.Add(key3, val3)

				val, ok := cache.Get(key1)
				if ok {
					assert.Equal(t, val1, val.(string))
				}
				cache.Add(key1, val1)

				val, ok = cache.Get(key2)
				if ok {
					assert.Equal(t, val2, val.(string))
				}
				cache.Add(key2, val2)

				key4 := fmt.Sprintf("key4-%d", j)
				val4 := key4
				val, ok = cache.Get(key4)
				if ok {
					assert.Equal(t, val4, val.(string))
				}
				cache.Add(key4, val4)

				val, ok = cache.Get(key3)
				if ok {
					assert.Equal(t, val3, val.(string))
				}
			}

			wg.Done()
		}()
	}
	wg.Wait()
}
