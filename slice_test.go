package mtoi

import (
	"testing"
	"time"

	"github.com/WindomZ/testify/assert"
)

var array *Slice

func TestSlice_NewSlice(t *testing.T) {
	array = NewSlice(2)
}

func TestSlice_Put(t *testing.T) {
	for i := 0; i < 10; i++ {
		s := string(demo[i])
		array.Put(s, s)
	}
}

func TestSlice_Size(t *testing.T) {
	assert.Equal(t, 10, array.Size())
}

func TestSlice_Get(t *testing.T) {
	for i := 0; i < 10; i++ {
		k := string(demo[i])
		vs, ok := array.Get(k)
		if assert.True(t, ok) && assert.Equal(t, 1, len(vs)) {
			for _, v := range vs {
				s, ok := v.(string)
				assert.True(t, ok)
				assert.NotEmpty(t, s)
				assert.Equal(t, 1, len(s))
			}
		}
	}
}

func TestSlice_Contain(t *testing.T) {
	assert.True(t, array.Contain("a"))
	assert.False(t, array.Contain("z"))
}

func TestSlice_Array(t *testing.T) {
	assert.NotEmpty(t, array.Array())
}

func TestSlice_Clean(t *testing.T) {
	array.Clean()
}

func TestSlice_Size2(t *testing.T) {
	assert.Equal(t, 0, array.Size())
}

func TestSlice_MulPut(t *testing.T) {
	f, stop := array.MulPut()
	for i := 10; i < 20; i++ {
		s := string(demo[i%5+10])
		f(s, s)
	}
	stop()
}

func TestSlice_Size3(t *testing.T) {
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, 10, array.Size())
}

func TestSlice_Get2(t *testing.T) {
	for i := 10; i < 20; i++ {
		k := string(demo[i%5+10])
		vs, ok := array.Get(k)
		if assert.True(t, ok) && assert.Equal(t, 2, len(vs)) {
			for _, v := range vs {
				s, ok := v.(string)
				assert.True(t, ok)
				assert.NotEmpty(t, s)
				assert.Equal(t, 1, len(s))
			}
		}
	}
}

func TestSlice_MulGet(t *testing.T) {
	for i := 10; i < 20; i++ {
		k := string(demo[i%5+10])
		ok := array.MulGet(k, func(v interface{}) {
			s, ok := v.(string)
			assert.True(t, ok)
			assert.NotEmpty(t, s)
			assert.Equal(t, 1, len(s))
		})
		assert.True(t, ok)
	}
}

func TestSlice_Close(t *testing.T) {
	array.Close()
}
