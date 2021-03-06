package mtoi

import (
	"testing"
	"time"

	"github.com/WindomZ/testify/assert"
)

var cache *Cache

func TestCache_NewCache(t *testing.T) {
	cache = NewCache(0, 0)
	cache = NewCache(20, time.Second)
}

func TestCache_Put(t *testing.T) {
	for i := 0; i < 10; i++ {
		s := string(demo[i])
		cache.Put(s, s, time.Second)
	}
}

func TestCache_Get(t *testing.T) {
	time.Sleep(time.Millisecond * 100)

	for i := 0; i < 10; i++ {
		k := string(demo[i])
		v, ok := cache.Get(k)
		if assert.True(t, ok) {
			s, ok := v.(string)
			assert.True(t, ok)
			assert.NotEmpty(t, s)
			assert.Equal(t, 1, len(s))
		}
	}
}

func TestCache_Contain(t *testing.T) {
	assert.True(t, cache.Contain("a"))
	assert.False(t, cache.Contain("z"))
}

func TestCache_Get2(t *testing.T) {
	time.Sleep(time.Second * 2)
	for i := 0; i < 10; i++ {
		k := string(demo[i])
		_, ok := cache.Get(k)
		assert.False(t, ok)
	}
}

func TestCache_Close(t *testing.T) {
	cache.Close()
}
