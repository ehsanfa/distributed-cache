package buffer

import (
	"dbcache/cluster/cacher"
	"reflect"
	"testing"
)

func TestIsEmpty(t *testing.T) {
	b := CreateInMemoryBuffer()
	if !b.IsEmpty() {
		t.Error("Expected empty buffer")
	}
}

func TestAdd(t *testing.T) {
	b := CreateInMemoryBuffer()
	v := cacher.NewCacheValue()
	v.Value = "hooshang"
	c := &cacher.CacheEntity{Key: "hasan", Value: v}
	b.Add(c)
}

func TestReset(t *testing.T) {
	b := CreateInMemoryBuffer()
	v := cacher.NewCacheValue()
	v.Value = "hooshang"
	c := &cacher.CacheEntity{Key: "hasan", Value: v}
	b.Add(c)
	b.Reset()
	if !b.IsEmpty() {
		t.Error("Expected empty buffer")
	}
}

func TestAll(t *testing.T) {
	b := CreateInMemoryBuffer()
	v := cacher.NewCacheValue()
	v.Value = "hooshang"
	c := &cacher.CacheEntity{Key: "hasan", Value: v}
	b.Add(c)
	if reflect.DeepEqual(b.All(), map[string]cacher.CacheValue{"hooshang": v}) {
		t.Error("Expected maps to be equal")
	}
}

func TestSize(t *testing.T) {
	b := CreateInMemoryBuffer()
	v := cacher.NewCacheValue()
	v.Value = "hooshang"
	c := &cacher.CacheEntity{Key: "hasan", Value: v}
	b.Add(c)
	if b.Size() != 1 {
		t.Error("Expected size to be 1")
	}
}
