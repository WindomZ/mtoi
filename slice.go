package mtoi

import "sync"

type itemSlice struct {
	Key   string
	Value interface{}
}

// Slice
type Slice struct {
	cap    int
	data   []interface{}
	tag    map[string][]int
	stream chan *itemSlice
	lock   *sync.RWMutex
}

func NewSlice(cap int) *Slice {
	if cap <= 2 {
		cap = 2
	}
	c := &Slice{
		cap:    cap,
		data:   make([]interface{}, 0, cap),
		tag:    make(map[string][]int, cap),
		stream: make(chan *itemSlice, cap),
		lock:   new(sync.RWMutex),
	}
	c.start()
	return c
}

func (c *Slice) start() {
	go func() {
		for v, ok := <-c.stream; ok; v, ok = <-c.stream {
			if v != nil {
				c.lock.Lock()
				for ; v != nil && ok; v, ok = <-c.stream {
					c.data = append(c.data, v.Value)
					idx, ok := c.tag[v.Key]
					if !ok {
						idx = make([]int, 0, 2)
					}
					c.tag[v.Key] = append(idx, len(c.data)-1)
					if len(c.data) > c.cap {
						c.cap = len(c.data) + 2
					}
				}
				c.lock.Unlock()
			}
		}
	}()
}

func (c *Slice) Close() {
	close(c.stream)
}

func (c Slice) Size() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.data)
}

func (c *Slice) Clean() {
	c.lock.Lock()
	c.data = make([]interface{}, 0, c.cap)
	c.tag = make(map[string][]int, c.cap)
	c.lock.Unlock()
}

func (c *Slice) Put(k string, v interface{}) {
	c.stream <- &itemSlice{k, v}
	c.stream <- nil
}

func (c *Slice) MulPut() (func(k string, v interface{}), func()) {
	return func(k string, v interface{}) { c.stream <- &itemSlice{k, v} },
		func() { c.stream <- nil }
}

func (c Slice) Get(k string) (res []interface{}, ok bool) {
	c.lock.RLock()
	index, ok := c.tag[k]
	if ok && len(index) != 0 {
		res = make([]interface{}, len(index))
		for i, idx := range index {
			res[i] = c.data[idx]
		}
	}
	c.lock.RUnlock()
	return
}

func (c Slice) MulGet(k string, f func(interface{})) bool {
	c.lock.RLock()
	index, ok := c.tag[k]
	if ok && len(index) != 0 {
		for _, idx := range index {
			f(c.data[idx])
		}
	}
	c.lock.RUnlock()
	return ok
}

func (c Slice) Contain(k string) bool {
	c.lock.RLock()
	_, ok := c.tag[k]
	c.lock.RUnlock()
	return ok
}

func (c Slice) Array() []interface{} {
	return c.data
}
