package buffer

import (
	"dbcache/cluster/cacher"
	"reflect"
	"testing"
)

type cacheValue struct {
	val string
	ver int
}

func (c *cacheValue) GetValue() string {
	return c.val
}
func (c *cacheValue) SetValue(val string) {
	c.val = val
}
func (c *cacheValue) Version() int {
	return c.ver
}
func (c *cacheValue) IncrementVersion() {
	c.ver += 1
}

func TestIsEmpty(t *testing.T) {
	b := CreateInMemoryBuffer()
	if !b.IsEmpty() {
		t.Error("Expected empty buffer")
	}
}

func TestAdd(t *testing.T) {
	b := CreateInMemoryBuffer()
	v := &cacheValue{"hooshang", 1}
	c := &cachable{key: "hasan", value: v}
	b.Add(c)
}

func TestReset(t *testing.T) {
	b := CreateInMemoryBuffer()
	v := &cacheValue{"hooshang", 1}
	c := &cachable{key: "hasan", value: v}
	b.Add(c)
	b.Reset()
	if !b.IsEmpty() {
		t.Error("Expected empty buffer")
	}
}

func TestAll(t *testing.T) {
	b := CreateInMemoryBuffer()
	v := &cacheValue{"hooshang", 1}
	c := &cachable{key: "hasan", value: v}
	b.Add(c)
	if reflect.DeepEqual(b.All(), map[string]cacher.CacheValue{"hooshang": v}) {
		t.Error("Expected maps to be equal")
	}
}

func TestSize(t *testing.T) {
	b := CreateInMemoryBuffer()
	v := &cacheValue{"hooshang", 1}
	c := &cachable{key: "hasan", value: v}
	b.Add(c)
	if b.Size() != 1 {
		t.Error("Expected size to be 1")
	}
}

func TestMerge(t *testing.T) {
	b1 := CreateInMemoryBuffer()
	b2 := CreateInMemoryBuffer()
	b1.Add(&cachable{key: "hasan", value: &cacheValue{"hooshang", 1}})
	b1.Add(&cachable{key: "hasan2", value: &cacheValue{"hooshang", 1}})
	b2.Add(&cachable{key: "hasan3", value: &cacheValue{"hooshang", 1}})
	b1.Merge(b2)
	if b1.Size() != 3 {
		t.Error("Expected size to be 3", b1.Size())
	}
}
