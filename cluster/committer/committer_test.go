package committer

import (
	"dbcache/cluster/buffer"
	"dbcache/cluster/cacher"
	"testing"
)

func TestCommit(t *testing.T) {
	cache := cacher.CreateInMemoryCache()
	cache.Set("hasan", cacher.NewVersionBasedCacheValue("hooshang", 0))

	buff := buffer.CreateInMemoryBuffer()
	buff.Add("hasan", cacher.NewVersionBasedCacheValue("hooshang", 2))
	buff.Add("hasan_new", cacher.NewVersionBasedCacheValue("hooshang_new", 1))

	vc := CreateVersionCompareCommitter()
	vc.Commit(cache, buff)

	if v, _ := cache.Get("hasan"); v.Version() != 2 {
		t.Error("version is not updated")
	}

	if v, _ := cache.Get("hasan_new"); v.Version() != 1 {
		t.Error("new value version is not correct")
	} else if v.GetValue() != "hooshang_new" {
		t.Error("new value is not correctly set")
	}
}
