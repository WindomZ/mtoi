// +build !go1.9

package mtoi

import "sync"

type itemKV struct {
	Key   string
	Value interface{}
}

type KV struct {
	cap    int
	data   map[string]interface{}
	stream chan *itemKV
	lock   *sync.RWMutex
}

func NewKV(cap int) *KV {
	if cap <= 2 {
		cap = 2
	}
	c := &KV{
		cap:    cap,
		data:   make(map[string]interface{}, cap),
		stream: make(chan *itemKV, cap),
		lock:   new(sync.RWMutex),
	}
	c.start()
	return c
}

func (c *KV) start() {
	go func() {
		for v, ok := <-c.stream; ok; v, ok = <-c.stream {
			if v != nil && len(c.data) < c.cap {
				c.lock.Lock()
				for ; v != nil && ok; v, ok = <-c.stream {
					c.data[v.Key] = v.Value
					if len(c.data) >= c.cap {
						break
					}
				}
				c.lock.Unlock()
			}
		}
	}()
}

func (c *KV) Close() {
	close(c.stream)
}

func (c KV) Size() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.data)
}

func (c *KV) Clean() {
	c.lock.Lock()
	c.data = make(map[string]interface{}, c.cap)
	c.lock.Unlock()
}

func (c *KV) Put(k string, v interface{}) {
	c.stream <- &itemKV{k, v}
	c.stream <- nil
}

func (c *KV) MulPut() (func(k string, v interface{}), func()) {
	return func(k string, v interface{}) { c.stream <- &itemKV{k, v} },
		func() { c.stream <- nil }
}

func (c KV) Get(k string) (interface{}, bool) {
	c.lock.RLock()
	v, ok := c.data[k]
	c.lock.RUnlock()
	return v, ok
}

func (c KV) Contain(k string) bool {
	c.lock.RLock()
	_, ok := c.data[k]
	c.lock.RUnlock()
	return ok
}
