package buffer

import (
	"dbcache/cluster/cacher"
)

type Buffer interface {
	IsEmpty() bool
	Add(string, cacher.CacheValue)
	All() map[string]cacher.CacheValue
	Reset()
	Size() int
	Merge(Buffer)
}
