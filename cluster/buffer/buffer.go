package buffer

import (
	"dbcache/cluster/cacher"
)

type Buffer interface {
	IsEmpty() bool
	Add(c cacher.Cachable)
	All() map[string]cacher.CacheValue
	Reset()
	Size() int
}
