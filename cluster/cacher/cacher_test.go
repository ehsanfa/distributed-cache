package cacher

import (
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	c := CreateInMemoryCache()
	v := NewVersionBasedCacheValue("hooshang", 1)
	c.cache["hasan"] = v
	val, ok := c.Get("hasan")
	if val.GetValue() != "hooshang" || !c.Exists("hasan") {
		t.Error("invalid cached value", val.GetValue())
	}
	if ok != true {
		t.Error("key doesnt exist", val.GetValue())
	}
	c = CreateInMemoryCache()
	val, ok = c.Get("hasan")
	if val != nil {
		t.Error("invalid cached value", val.GetValue())
	}
	if ok || c.Exists("hasan") {
		t.Error("key exists", val.GetValue())
	}
}

func TestSet(t *testing.T) {
	c := CreateInMemoryCache()
	v := NewVersionBasedCacheValue("hooshang", 1)
	c.Set("hasan", v)
	val, _ := c.Get("hasan")
	if val.GetValue() != "hooshang" {
		t.Error("invalid cached value")
	}
}

func TestAll(t *testing.T) {
	c := CreateInMemoryCache()
	v := NewVersionBasedCacheValue("hooshang", 1)
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
	v := NewVersionBasedCacheValue("hooshang", 1)
	c.Set("hasan", v)
	c.Touch("hasan")
	if c.Version("hasan") != 2 {
		t.Error("versions dont match", c.Version("hasan"), 2)
	}
}

func TestDelete(t *testing.T) {
	c := CreateInMemoryCache()
	v := NewVersionBasedCacheValue("hooshang", 1)
	c.Set("hasan", v)
	c.Touch("hasan")
	c.Delete("hasan")
	if c.Version("hasan") != 3 {
		t.Error("versions dont match", c.Version("hasan"), 3)
	}
	newValue, _ := c.Get("hasan")
	if newValue.GetValue() != "" {
		t.Error("Expected empty value, provided this instead", newValue)
	}
}
