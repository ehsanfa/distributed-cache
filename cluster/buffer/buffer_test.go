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
	v := cacher.NewVersionBasedCacheValue("hooshang", 1)
	b.Add("hasan", v)
}

func TestReset(t *testing.T) {
	b := CreateInMemoryBuffer()
	v := cacher.NewVersionBasedCacheValue("hooshang", 1)
	b.Add("hasan", v)
	b.Reset()
	if !b.IsEmpty() {
		t.Error("Expected empty buffer")
	}
}

func TestAll(t *testing.T) {
	b := CreateInMemoryBuffer()
	v := cacher.NewVersionBasedCacheValue("hooshang", 1)
	b.Add("hasan", v)
	if reflect.DeepEqual(b.All(), map[string]cacher.CacheValue{"hooshang": v}) {
		t.Error("Expected maps to be equal")
	}
}

func TestSize(t *testing.T) {
	b := CreateInMemoryBuffer()
	v := cacher.NewVersionBasedCacheValue("hooshang", 1)
	b.Add("hasan", v)
	if b.Size() != 1 {
		t.Error("Expected size to be 1")
	}
}

func TestMerge(t *testing.T) {
	b1 := CreateInMemoryBuffer()
	b2 := CreateInMemoryBuffer()
	b1.Add("hasan", cacher.NewVersionBasedCacheValue("hooshang", 1))
	b1.Add("hasan2", cacher.NewVersionBasedCacheValue("hooshang", 1))
	b2.Add("hasan3", cacher.NewVersionBasedCacheValue("hooshang", 1))
	b1.Merge(b2)
	if b1.Size() != 3 {
		t.Error("Expected size to be 3", b1.Size())
	}
}
