// +build go1.9

package mtoi

import (
	"sync"
	"time"
)

type itemCache struct {
	Key        string
	Value      interface{}
	ExpireTime int64
}

type Cache struct {
	cap       int
	interval  time.Duration
	data      sync.Map
	stream    chan *itemCache
	interrupt chan bool
}

func NewCache(cap int, interval time.Duration) *Cache {
	if cap <= 2 {
		cap = 2
	}
	if interval <= 0 {
		interval = time.Minute * 5
	}
	c := &Cache{
		cap:       cap,
		interval:  interval,
		stream:    make(chan *itemCache, cap),
		interrupt: make(chan bool),
	}
	c.start()
	return c
}

func (c *Cache) start() {
	go func() {
		for v, ok := <-c.stream; ok; v, ok = <-c.stream {
			if v != nil {
				for ; v != nil && ok; v, ok = <-c.stream {
					c.data.Store(v.Key, v)
				}
			}
		}
	}()
	go func() {
		for ok := true; ok; {
			select {
			case <-c.interrupt:
				ok = false
			case <-time.After(c.interval):
				now := time.Now().Unix()
				c.data.Range(func(key, value interface{}) bool {
					if v, ok := value.(*itemCache); !ok || v.ExpireTime <= now {
						c.data.Delete(key)
					}
					return true
				})
			}
		}
	}()
}

func (c *Cache) Close() {
	close(c.stream)
	c.interrupt <- true
	close(c.interrupt)
}

func (c *Cache) Put(k string, v interface{}, expire time.Duration) {
	if k != "" && expire > 0 {
		c.stream <- &itemCache{k, v,
			time.Now().Add(expire).Unix()}
		c.stream <- nil
	}
}

func (c *Cache) Get(k string) (interface{}, bool) {
	if value, ok := c.data.Load(k); ok {
		v, ok := value.(*itemCache)
		if ok && v.ExpireTime > time.Now().Unix() {
			return v.Value, true
		}
		c.data.Delete(k)
	}
	return nil, false
}

func (c *Cache) Contain(k string) bool {
	_, ok := c.Get(k)
	return ok
}
