// +build !go1.9

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
	data      map[string]*itemCache
	stream    chan *itemCache
	interrupt chan bool
	lock      *sync.RWMutex
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
		data:      make(map[string]*itemCache, cap),
		stream:    make(chan *itemCache, cap),
		interrupt: make(chan bool),
		lock:      new(sync.RWMutex),
	}
	c.start()
	return c
}

func (c *Cache) start() {
	go func() {
		for v, ok := <-c.stream; ok; v, ok = <-c.stream {
			if v != nil && len(c.data) < c.cap {
				c.lock.Lock()
				for ; v != nil && ok; v, ok = <-c.stream {
					c.data[v.Key] = v
					if len(c.data) >= c.cap {
						break
					}
				}
				c.lock.Unlock()
			}
		}
	}()
	go func() {
		for ok := true; ok; {
			select {
			case <-c.interrupt:
				ok = false
			case <-time.After(c.interval):
				c.lock.Lock()
				now := time.Now().Unix()
				for _, v := range c.data {
					if v.ExpireTime < now {
						delete(c.data, v.Key)
					}
				}
				c.lock.Unlock()
			}
		}
	}()
}

func (c *Cache) Close() {
	close(c.stream)
	c.interrupt <- true
	close(c.interrupt)
}

func (c Cache) Size() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.data)
}

func (c *Cache) Clean() {
	c.lock.Lock()
	c.data = make(map[string]*itemCache, c.cap)
	c.lock.Unlock()
}

func (c *Cache) Put(k string, v interface{}, expire time.Duration) {
	if k != "" && expire > 0 {
		c.stream <- &itemCache{k, v,
			time.Now().Add(expire).Unix()}
		c.stream <- nil
	}
}

func (c *Cache) Get(k string) (interface{}, bool) {
	c.lock.RLock()
	v, ok := c.data[k]
	c.lock.RUnlock()
	if ok && v.ExpireTime >= time.Now().Unix() {
		return v.Value, true
	}
	return nil, false
}

func (c *Cache) Contain(k string) bool {
	c.lock.RLock()
	_, ok := c.Get(k)
	c.lock.RUnlock()
	return ok
}
