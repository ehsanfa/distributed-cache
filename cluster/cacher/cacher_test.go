package cacher

import (
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	c := CreateInMemoryCache()
	v := NewCacheValue()
	v.Value = "hooshang"
	c.cache["hasan"] = v
	val, ok := c.Get("hasan")
	if val.Value != "hooshang" || !c.Exists("hasan") {
		t.Error("invalid cached value", val.Value)
	}
	if ok != true {
		t.Error("key doesnt exist", val.Value)
	}
	c = CreateInMemoryCache()
	val, ok = c.Get("hasan")
	if val.Value != "" {
		t.Error("invalid cached value", val.Value)
	}
	if ok || c.Exists("hasan") {
		t.Error("key exists", val.Value)
	}
}

func TestSet(t *testing.T) {
	c := CreateInMemoryCache()
	v := NewCacheValue()
	v.Value = "hooshang"
	c.Set("hasan", v)
	val, _ := c.Get("hasan")
	if val.Value != "hooshang" {
		t.Error("invalid cached value")
	}
}

func TestAll(t *testing.T) {
	c := CreateInMemoryCache()
	v := NewCacheValue()
	v.Value = "hooshang"
	c.Set("hasan", v)
	all := c.All()
	m := make(map[string]CacheValue)
	m["hasan"] = v
	if !reflect.DeepEqual(all, m) {
		t.Error("maps don't match", all, m)
	}
	c.Replace(m)
	all = c.All()
	if !reflect.DeepEqual(all, m) {
		t.Error("maps don't match", all, m)
	}
}

func TestVersion(t *testing.T) {
	c := CreateInMemoryCache()
	v := NewCacheValue()
	v.Value = "hooshang"
	c.Set("hasan", v)
	c.Touch("hasan")
	if c.Version("hasan") != 1 {
		t.Error("versions dont match", c.Version("hasan"), 1)
	}
}

func TestDelete(t *testing.T) {
	c := CreateInMemoryCache()
	v := NewCacheValue()
	v.Value = "hooshang"
	c.Set("hasan", v)
	c.Touch("hasan")
	c.Delete("hasan")
	if c.Version("hasan") != 2 {
		t.Error("versions dont match", c.Version("hasan"), 2)
	}
	newValue, _ := c.Get("hasan")
	if newValue.Value != "" {
		t.Error("Expected empty value, provided this instead", newValue)
	}
}
